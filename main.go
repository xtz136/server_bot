package main

import (
	"bot/internal/commands"
	"bot/pkg/config"
	"bot/pkg/dingding"
	"bot/pkg/health"
	"bot/pkg/talk"
	"bot/pkg/task"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Talk(command string, sender chan string, reply chan string, log zerolog.Logger) {
	ctx := commands.Context{
		"",
		sender,
		reply,
		nil,
		nil,
		log,
		make(commands.State),
		talk.MakeTalkEnd,
	}

	// 解析命令和目标
	// 例如 command = "重启阿里云"，则 command = "重启"，target = "阿里云"
	for taskName, taskC := range config.C.Tasks {
		if !strings.HasPrefix(command, taskName) {
			continue
		}
		targetName := command[len(taskName):]
		if target, err := task.GetTarget(targetName); err != nil {
			talk.MakeTalkEnd(sender, "这个系统没有配置，请联系管理员")
			return
		} else {
			ctx.Target = target
		}
		ctx.Task = &taskC
		ctx.TargetName = targetName
		go commands.TaskCommands[taskC.Name](ctx)
		return
	}
	// 没有找到对应的命令，显示帮助
	go commands.ShowHelper(ctx)
}

func main() {
	r := gin.Default()
	r.POST("/bot/dingding/talk", dingding.DingDing(Talk))

	if len(config.C.Beat) > 0 {
		go health.BeatCheckHealth()
	}
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
