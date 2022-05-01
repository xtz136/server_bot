package commands

import (
	"bot/pkg/config"
	"bot/pkg/http_client"
	"fmt"
	"time"
)

// 重启进程/服务，需要预先配置检查链接和重启链接
// 服务都是一个集群的，不应该同时重启，会被客户感知。
// 需要按顺序重启多个服务，中间如果有错误，则会跳过错误的服务
func Restart(ctx Context) {

	b := config.C.Systems[ctx.Name]
	if len(b.Restart) == 0 {
		ctx.Sender <- fmt.Sprintln("这个系统没有配置，请联系管理员")
		return
	}

	ctx.Sender <- fmt.Sprintf("开始重启%s，耐心等待", b.Name)

	bT := time.Now()
	client := http_client.NewDumbHttpClient(5)

	raiseError := func() {
		eT := time.Since(bT)
		ctx.MakeTalkEnd(ctx.Sender, fmt.Sprintf("汪，重启%s失败，耗时: %v, 请联系管理员，本次服务结束", b.Name, eT))
	}

	for _, item := range b.Restart {
		command := item.Command
		check := item.Check

		// 调用重启链接
		_, err := http_client.Get(client, command)
		if err != nil {
			// 调用失败，打印错误，然后下一个
			ctx.Log.Error().Err(err).Msg("")
			raiseError()
			return
		}
		ctx.Log.Info().Str("request", command).Msg("run command")

		// 服务已经重启了，这里会每3秒检查一次，直到检查链接响应
		// 循环100次，每次3秒，一共5分钟
		maxTry := 100
		for i := 1; i <= maxTry; i++ {
			_, err := http_client.Get(client, check)
			if err != nil {
				time.Sleep(time.Duration(3) * time.Second)
				continue
			}
			if i == maxTry {
				// 超过了5分钟，没有检查到服务响应，则认为服务重启失败
				ctx.Log.Warn().Msg("check failed")
				raiseError()
				return
			}
		}
	}

	eT := time.Since(bT)
	ctx.MakeTalkEnd(ctx.Sender, fmt.Sprintf("汪，重启%s完成，耗时：%v，本次服务结束", b.Name, eT))
}
