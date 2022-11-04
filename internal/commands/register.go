package commands

import (
	"bot/pkg/talks"
	"bot/pkg/websocket"
	"context"
	"fmt"
)

// 注册一个机器，绑定别名
func Register(ctx context.Context) {
	// registerCode := -1
	target := ctx.Value(talks.TargetNameKey).(string)
	registerName := ""

	hosts := ""
	hub := websocket.NewHub()
	for index, clientIdent := range hub.CollectClientsIdentity() {
		hosts += fmt.Sprintf("%d %s\n", index, clientIdent)
	}

	w := newWorkFlow(ctx)
	w.add("start", "请回复别名，30秒有效", func(ctx context.Context, replyMsg string) string {
		registerName = replyMsg
		return ""
	})
	w.start("start")
	eT := w.getCostTime()

	fmt.Printf("register %v = %v\n", target, registerName)
	clientAlias := websocket.NewClientAlias()
	clientAlias.Add(target, registerName)

	talks.MakeTalkEnd(ctx, fmt.Sprintf("注册完成，耗时：%v，本次服务结束", eT))
}

func init() {
	registerTaskCommand("Register", Register)
}
