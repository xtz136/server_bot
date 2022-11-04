package commands

import (
	"bot/pkg/config"
	"bot/pkg/talks"
	"bot/pkg/websocket"
	"context"
	"fmt"
	"strings"
)

func ShowHelper(ctx context.Context) {

	msg := "命令列表：\n"

	for taskName := range config.C.Tasks {
		if !strings.HasSuffix(taskName, "*") {
			msg += "* " + taskName + "\n"
		}
	}

	msg += "机器列表：\n"

	for targetName := range config.C.Targets {
		if !strings.HasSuffix(targetName, "*") {
			msg += "* " + targetName + "\n"
		}
	}

	// 加上websocket的客户端
	hub := websocket.NewHub()
	for clientName, clientIdents := range hub.CollectClientsName() {
		// msg += "* " + clientName + "\n"
		msg += fmt.Sprintf("* %s[%s]\n", clientName, clientIdents)
	}

	msg += "使用命令+机器名称，如：\n"
	msg += "重启阿里云\n"

	talks.MakeTalkEnd(ctx, msg)
}

func init() {
	registerTaskCommand("ShowHelper", ShowHelper)
}
