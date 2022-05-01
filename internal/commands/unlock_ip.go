package commands

import (
	"bot/pkg/config"
	"bot/pkg/http_client"
	"bot/pkg/task"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

func UnlockIP(ctx Context) {
	hc := http_client.NewHttpClient(3)
	dhc := http_client.NewDumbHttpClient(3)
	w := newWorkFlow(ctx)
	targetTask := task.ListTargetTask(ctx.Target, ctx.Task, &config.C.Variables)
	lockIP := targetTask[0].Command

	// 获取所有被封的IP
	w.add(func(ctx Context) bool {
		respBody, err := http_client.Get(hc, lockIP)
		if err != nil {
			ctx.Log.Error().Err(err).Str("url", lockIP).Msg("get lock ip")
			return false
		}

		lresp := UnlockIPResponse{}
		json.Unmarshal(respBody, &lresp)

		keys := lresp.Data
		// 不需要处理，结束
		if len(keys) == 0 {
			ctx.Sender <- "没有需要解封的IP"
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

			// FIXME 最多重试次数，和超时时间
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
		if _, err := http_client.Delete(dhc, lockIP+"&key="+url.QueryEscape(keys[choose])); err != nil {
			ctx.Log.Error().Err(err).Msg("remove lock ip")
			return false
		}
		return true
	})

	if w.start() {
		ctx.MakeTalkEnd(ctx.Sender, "解除限制成功，本次服务结束")
	} else {
		ctx.MakeTalkEnd(ctx.Sender, "解除限制失败，如有疑问请联系管理员，本次服务结束")
	}
}

func init() {
	registerTaskCommand("UnlockIP", UnlockIP)
}
