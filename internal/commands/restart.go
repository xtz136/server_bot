package commands

import (
	"bot/pkg/config"
	"bot/pkg/http_client"
	"bot/pkg/task"
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

type restartUnit struct {
	log           zerolog.Logger
	checkClient   http_client.HttpClientInterface
	commandClient http_client.HttpClientInterface
}

func restart(ru *restartUnit, check string, command string, checkInterval time.Duration) int {
	// 调用重启链接
	_, err := http_client.Get(ru.commandClient, command)
	if err != nil {
		// 调用失败，打印错误，然后下一个
		ru.log.Error().Err(err).Msg("")
		return 1
	}
	ru.log.Info().Str("request", command).Msg("run command")

	// 不需要检查
	if check == "" {
		return 0
	}

	// 服务已经重启了，这里会每3秒检查一次，直到检查链接响应
	// 循环100次，每次3秒，一共5分钟
	maxTry := 100
	for i := 1; i <= maxTry; i++ {
		if i == maxTry {
			// 超过了5分钟，没有检查到服务响应，则认为服务重启失败
			ru.log.Warn().Str("request", check).Msg("check failed")
			return 2
		}

		_, err := http_client.Get(ru.checkClient, check)
		if err != nil {
			time.Sleep(checkInterval)
			continue
		} else {
			// 检查到服务响应，则认为服务重启成功
			ru.log.Info().Str("request", check).Msg("check success")
			return 0
		}
	}

	return 0
}

// 重启进程/服务，需要预先配置检查链接和重启链接
// 服务都是一个集群的，不应该同时重启，会被客户感知。
// 需要按顺序重启多个服务，中间如果有错误，则会跳过错误的服务
func RestartGroup(ctx Context) {

	ctx.Sender <- fmt.Sprintf("开始重启%s，耐心等待", ctx.TargetName)

	bT := time.Now()
	client := http_client.NewDumbHttpClient(5)

	targetTask := task.ListTargetTask(ctx.Target, ctx.Task, &config.C.Variables)
	ru := &restartUnit{
		log:           ctx.Log,
		checkClient:   client,
		commandClient: client,
	}

	for _, item := range targetTask {
		result_code := restart(ru, item.Check, item.Command, 3*time.Second)
		if result_code != 0 {
			eT := time.Since(bT)
			ctx.MakeTalkEnd(ctx.Sender, fmt.Sprintf("汪，重启%s失败，耗时: %v, 请联系管理员，本次服务结束", ctx.TargetName, eT))
			return
		}
	}

	eT := time.Since(bT)
	ctx.MakeTalkEnd(ctx.Sender, fmt.Sprintf("汪，重启%s完成，耗时：%v，本次服务结束", ctx.TargetName, eT))

}

func init() {
	registerTaskCommand("Restart", RestartGroup)
}
