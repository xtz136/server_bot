package commands

import (
	"bot/pkg/config"
)

func ShowHelper(ctx Context) {

	msg := "命令列表：\n"

	for taskName, task := range config.C.Tasks {
		if !task.Hidden {
			msg += "* " + taskName + "\n"
		}
	}

	msg += "机器列表：\n"

	for targetName, target := range config.C.Targets {
		if !target.Hidden {
			msg += "* " + targetName + "\n"
		}
	}

	msg += "使用命令+机器名称，如：\n"
	msg += "* 重启阿里云\n"

	ctx.MakeTalkEnd(ctx.Sender, msg)
}

func init() {
	registerTaskCommand("ShowHelper", ShowHelper)
}
