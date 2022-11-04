package commands

import (
	"bot/pkg/config"
	"bot/pkg/talks"
	"bot/pkg/websocket"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func restartClients(clientsName []string) {
	commandJson, _ := json.Marshal(map[string]string{
		"type":    "command",
		"command": "restart",
		"timeout": "10.0",
		"clients": strings.Join(clientsName, ","),
	})

	hub := websocket.NewHub()
	hub.SendMessageToAll(commandJson)
}

// 重启进程/服务
func RestartTarget2(ctx context.Context) {
	// 兼容旧的重启服务
	target := ctx.Value(talks.TargetKey).([]config.Target)
	if target != nil {
		RestartTarget(ctx)
		return
	}

	logger := ctx.Value(talks.LoggerKey).(zerolog.Logger)
	targetName := ctx.Value(talks.TargetNameKey).(string)
	talks.ReplyMsg(ctx, fmt.Sprintf("开始重启%s，耐心等待", targetName))
	start := time.Now()

	// if tasks.NewParallelTasks(ctx, jobs) {
	// 	elapsed := time.Since(start)
	// 	talks.MakeTalkEnd(ctx, fmt.Sprintf("重启%s完成，耗时：%v，本次服务结束", targetName, elapsed))
	// } else {
	// 	elapsed := time.Since(start)
	// 	talks.MakeTalkEnd(ctx, fmt.Sprintf("重启%s失败，耗时: %v, 请联系管理员，本次服务结束", targetName, elapsed))
	// }

	hub := websocket.NewHub()
	ca := websocket.NewClientAlias()
	// 找到所有属于这个别名的客户端
	clientIdents := []string{}
	for client := range hub.ListClients() {
		curClientIdent := client.GetIdentity()
		for _, clientIdent := range ca.Alias[targetName] {
			if curClientIdent == clientIdent {
				clientIdents = append(clientIdents, client.GetName())
			}
		}
	}
	clientsLenHalf := len(clientIdents) / 2

	// 先重启一半完成，在重启另外一半。与负载搭配可以不让客户感知正在重启
	// 先重启后面的，这样在重启数量是单数的情况下可以更好
	shutdownClients := clientIdents[clientsLenHalf:]
	restartClients(shutdownClients)
	logger.Debug().Str(
		"clients", fmt.Sprintf("%v", shutdownClients),
	).Str(
		"all_clients", fmt.Sprintf("%v", clientIdents),
	).Msg("pre restart clients")

	notifyerChan := make(chan *websocket.Client)
	doneChan := make(chan int, 1)
	timeout := time.Second * 100
	defer close(doneChan)

	// 检查是否重启完成：有多少个程序接收到了重启指令，如果正常，那么就有多少个程序重新注册
	go hub.ObserveRegister(notifyerChan, doneChan)
	raiseClientsNum := 0
	for {
		select {
		case <-time.After(timeout):
			doneChan <- 1
			elapsed := time.Since(start)
			talks.MakeTalkEnd(ctx, fmt.Sprintf("重启%s失败，耗时：%v，本次服务结束", targetName, elapsed))
			return
		case <-notifyerChan:
			raiseClientsNum += 1
			if raiseClientsNum == len(shutdownClients) {
				doneChan <- 1
				goto FINISH
			}
		}
	}

FINISH:
	// 重启剩下的一半，理论上前面一半重启正常，那么剩下一半也应该正常重启，所以不检查了
	restartClients(clientIdents[:clientsLenHalf])
	logger.Debug().Str(
		"clients", fmt.Sprintf("%v", clientIdents[:clientsLenHalf]),
	).Msg("post restart clients")
	elapsed := time.Since(start)
	talks.MakeTalkEnd(ctx, fmt.Sprintf("重启%s完成，耗时：%v，本次服务结束", targetName, elapsed))
}

func init() {
	registerTaskCommand("Restart2", RestartTarget2)
}
