package commands

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/spf13/viper"
)

func AllowIP(ctx Context) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	lockIP := viper.GetString("systems." + ctx.Name + ".lock_ip")
	w := newWorkFlow(ctx)

	// 获取所有被封的IP
	w.add(func(ctx Context) bool {
		resp, err := http.Get(lockIP)
		defer resp.Body.Close()
		if err != nil {
			ctx.Log.Error().Err(err).Str("url", lockIP).Msg("get lock ip")
			ctx.Sender <- ""
			return false
		}
		lresp := LockIPResponse{}
		json.NewDecoder(resp.Body).Decode(&lresp)
		keys := lresp.Data
		// 不需要处理，结束
		if len(keys) == 0 {
			ctx.MakeTalkEnd(ctx.Sender, "没有需要解封的IP，本次服务结束")
			return false
		}
		ctx.State["keys"] = keys
		return true
	})

	// 让客服选择要处理的IP
	w.add(func(ctx Context) bool {
		keys := ctx.State["keys"].([]string)
		query := "请选择，回复序号：\n"
		for i, key := range keys {
			query += fmt.Sprintf("%d. %s\n", i, key)
		}
		ctx.Sender <- query

		for {
			answer, ok := <-ctx.Reply
			if !ok {
				return false
			}

			i, err := strconv.Atoi(answer)
			if err != nil || i >= len(keys) {
				ctx.Sender <- "选择错误，请重新选择!"
			} else {
				ctx.State["choose"] = i
				break
			}
		}

		return true
	})

	// 根据客服的选项，解封IP
	w.add(func(ctx Context) bool {
		keys := ctx.State["keys"].([]string)
		choose := ctx.State["choose"].(int)

		req, err := http.NewRequest("DELETE", lockIP+"&key="+url.QueryEscape(keys[choose]), nil)
		if err != nil {
			ctx.Log.Error().Err(err).Msg("remove lock ip")
		}
		http.DefaultClient.Do(req)
		return true
	})

	if w.start() {
		ctx.MakeTalkEnd(ctx.Sender, "解除限制成功，本次服务结束")
	}
}
