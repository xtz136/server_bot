package handler

import (
	"bot/internal/commands"
	"bot/pkg/config"
	"bot/pkg/logging"
	"bot/pkg/talks"
	"bot/pkg/targets"
	"context"
	"strings"
	"time"
)

func Handler(talk talks.TalkInterface) {
	sender := talk.GetSender()
	reply := talk.GetReply()
	command := talk.GetCommand()
	senderName := talk.GetSenderName()

	// 会话中，直接将命令通过管道发送给正在执行的工作者
	if !talk.IsFirstTalk() {
		if command == "取消" {
			sender <- "已取消任务"
			close(sender)
		} else {
			reply <- command
		}
		return
	}

	log := logging.Log.With().
		Str("app", "handler").
		Str("command", command).
		Str("sender", senderName).
		Logger()

	ctx := context.Background()
	ctx = context.WithValue(ctx, talks.SenderKey, sender)
	ctx = context.WithValue(ctx, talks.ReplyKey, reply)
	ctx = context.WithValue(ctx, talks.LoggerKey, log)
	ctx, cancel := context.WithTimeout(ctx, 5*60*time.Second)

	// 监听 sender，转发消息给客户。并监听/处理结束会话
	// 当 sender 被关闭，说明任务已经结束
	go func() {
		defer cancel()
		for msg := range sender {
			talk.ReplyMessage(ctx, msg)
		}
		senderName := talk.GetSenderName()
		command := talk.GetCommand()
		talks.DestoryTalkSession(senderName, command)
		log.Info().Msg("talk end")
	}()

	if command == "帮助" {
		go commands.TaskCommands["ShowHelper"](ctx)
		return
	}

	// 解析命令和目标
	// 例如 command = "重启阿里云"，则 task = "重启"，target = "阿里云"
	for taskName, task := range config.C.Tasks {
		if !strings.HasPrefix(command, taskName) {
			continue
		}
		targetName := command[len(taskName):]
		if target, err := targets.GetTarget(targetName); err != nil {
			talks.MakeTalkEnd(ctx, targetName+" 这个机器没有配置，请联系管理员")
			return
		} else {
			ctx = context.WithValue(ctx, talks.TargetKey, target)
		}
		ctx = context.WithValue(ctx, talks.TaskKey, &task)
		ctx = context.WithValue(ctx, talks.TargetNameKey, targetName)
		go commands.TaskCommands[task.Name](ctx)
		return
	}
	// 没有找到对应的命令，显示帮助
	go commands.TaskCommands["ShowHelper"](ctx)
}
