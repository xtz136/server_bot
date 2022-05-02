package handler

import (
	"bot/internal/commands"
	"bot/pkg/config"
	"bot/pkg/logging"
	talks "bot/pkg/talk"
	tasks "bot/pkg/task"
	"strings"
)

func Handler(talk talks.TalkInterface) {
	sender := talk.GetSender()
	reply := talk.GetReply()
	command := talk.GetCommand()
	senderName := talk.GetSenderName()

	// 会话中，直接将命令通过管道发送给正在执行的工作者
	if !talk.IsFirstTalk() {
		reply <- command
		return
	}

	log := logging.Log.With().
		Caller().
		Str("app", "handler").
		Str("command", command).
		Str("sender", senderName).
		Logger()

	ctx := commands.Context{
		"",
		sender,
		reply,
		nil,
		nil,
		log,
		make(commands.State),
		talks.MakeTalkEnd,
	}

	// 监听 sender，转发消息给客户。并监听/处理结束会话
	// 当 sender 被关闭，说明任务已经结束
	go func(sender chan string, talk talks.TalkInterface) {
		for {
			select {
			case msg, ok := <-sender:
				if ok {
					go talk.ReplyMessage(msg)
				} else {
					senderName := talk.GetSenderName()
					command := talk.GetCommand()
					talks.DestoryTalkSession(senderName, command)
					log.Info().Msg("talk end")
					return
				}
			}
		}
	}(sender, talk)

	// 解析命令和目标
	// 例如 command = "重启阿里云"，则 task = "重启"，target = "阿里云"
	for taskName, task := range config.C.Tasks {
		if !strings.HasPrefix(command, taskName) {
			continue
		}
		targetName := command[len(taskName):]
		if target, err := tasks.GetTarget(targetName); err != nil {
			talks.MakeTalkEnd(sender, targetName+" 这个机器没有配置，请联系管理员")
			return
		} else {
			ctx.Target = target
		}
		ctx.Task = &task
		ctx.TargetName = targetName
		go commands.TaskCommands[task.Name](ctx)
		return
	}
	// 没有找到对应的命令，显示帮助
	go commands.TaskCommands["ShowHelper"](ctx)
}
