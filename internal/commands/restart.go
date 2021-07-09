package commands

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

func Restart(ctx Context) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	b := System{}
	viper.UnmarshalKey("systems."+ctx.Name, &b)

	if len(b.Restart) == 0 {
		ctx.Sender <- fmt.Sprintf("这个系统没有配置，请联系管理员")
		return
	}

	ctx.Sender <- fmt.Sprintf("开始重启%s，耐心等待", b.Name)

	bT := time.Now()
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	raiseError := func() {
		eT := time.Since(bT)
		ctx.MakeTalkEnd(ctx.Sender, fmt.Sprintf("汪，重启%s失败，耗时: %v, 请联系管理员，本次服务结束", b.Name, eT))
		return
	}

	for _, item := range b.Restart {
		command := item.Command
		check := item.Check

		_, err := client.Get(command)
		if err != nil {
			ctx.Log.Error().Err(err).Msg("")
			raiseError()
			return
		}
		ctx.Log.Info().Str("request", command).Msg("run command")

		// 循环100次，每次3秒，一共5分钟
		has_error := true
		for i := 0; i < 100; i++ {
			resp, err := client.Get(check)
			if err != nil {
				time.Sleep(time.Duration(3) * time.Second)
				continue
			}

			if resp.StatusCode == 200 {
				has_error = false
				break
			} else {
				time.Sleep(time.Duration(3) * time.Second)
				continue
			}
		}

		if has_error {
			ctx.Log.Warn().Msg("check failed")
			raiseError()
			return
		}
	}

	eT := time.Since(bT)
	ctx.MakeTalkEnd(ctx.Sender, fmt.Sprintf("汪，重启%s完成，耗时：%v，本次服务结束", b.Name, eT))
}
