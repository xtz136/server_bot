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
	unlockIPUrl := targetTask[0].Command

	// 获取所有被封的IP
	respBody, err := http_client.Get(hc, unlockIPUrl)
	if err != nil {
		ctx.Log.Error().Err(err).Str("url", unlockIPUrl).Msg("get lock ip")
		return
	}

	resp := UnlockIPResponse{}
	json.Unmarshal(respBody, &resp)
	keys := resp.Data
	// 不需要处理，结束
	if len(keys) == 0 {
		ctx.Sender <- "没有需要解封的IP，本次服务结束"
		return
	}
	query := "请选择，回复序号：\n"
	for i, key := range keys {
		query += fmt.Sprintf("%d. %s\n", i, key)
	}

	// 让客服选择要处理的IP
	w.add("choose_release", query, func(ctx Context, replayMsg string) string {
		i, err := strconv.Atoi(replayMsg)
		if err != nil || i >= len(keys) {
			ctx.Sender <- "选择错误，请重新选择!"
			return "choose_release"
		} else {
			// 根据客服的选项，解封IP
			if _, err := http_client.Delete(dhc, unlockIPUrl+"&key="+url.QueryEscape(keys[i])); err != nil {
				ctx.Log.Error().Err(err).Msg("remove lock ip")
			}
		}
		return ""
	})

	if w.start("choose_release") == 0 {
		ctx.MakeTalkEnd(ctx.Sender, "解除限制成功，本次服务结束")
	} else {
		ctx.MakeTalkEnd(ctx.Sender, "解除限制失败，如有疑问请联系管理员，本次服务结束")
	}
}

func init() {
	registerTaskCommand("UnlockIP", UnlockIP)
}
