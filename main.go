package main

import (
	"bot/internal/commands"
	"bot/pkg/dingding"
	"bot/pkg/talk"
	"fmt"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Notify func(t string)

func Talk(command string, sender chan string, reply chan string, log zerolog.Logger) {

	ctx := commands.Context{
		"",
		sender,
		reply,
		log,
		make(commands.State),
		talk.MakeTalkEnd,
	}

	if strings.HasPrefix(command, "重启") {
		ctx.Name = command[6:]
		go commands.Restart(ctx)
		return
	}

	switch command {
	case "启用阿里云用户":
		ctx.Name = "阿里云"
		go commands.EnableUser(ctx)
	case "解封阿里云IP限制":
		ctx.Name = "阿里云"
		go commands.AllowIP(ctx)
	case "使用阿里云用户登录":
		ctx.Name = "阿里云"
		go commands.LoginUser(ctx)
	case "测试":
		ctx.Name = "测试"
		go commands.Dummy(ctx)
	case "获取会话数":
		talk.MakeTalkEnd(sender, fmt.Sprintf("%d", runtime.NumGoroutine()))
	default:
		helpMsg := "命令列表：\n"
		helpMsg += "* 帮助\n"

		// 重启服务器动态读取配置文件
		for _, k := range commands.ListServices() {
			helpMsg += "* 重启" + k + "\n"
		}

		helpMsg += "* 启用阿里云用户\n"
		helpMsg += "* 解封阿里云IP限制"

		talk.MakeTalkEnd(sender, helpMsg)
	}
}

func main() {
	r := gin.Default()
	r.POST("/bot/dingding/talk", dingding.DingDing(Talk))
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
