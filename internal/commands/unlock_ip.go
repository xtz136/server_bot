package commands

import (
	"bot/pkg/config"
	"bot/pkg/http_client"
	"bot/pkg/talks"
	"bot/pkg/targets"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

func UnlockIP(ctx context.Context) {
	hc := http_client.NewHttpClient(3 * time.Second)
	dhc := http_client.NewDumbHttpClient(3 * time.Second)
	w := newWorkFlow(ctx)
	target := ctx.Value(talks.TargetKey).([]config.Target)
	task := ctx.Value(talks.TaskKey).(*config.Task)
	logger := ctx.Value(talks.LoggerKey).(zerolog.Logger)
	targetTask := targets.ListTargetTask(target, task, &config.C.Variables)
	unlockIPUrl := targetTask[0].Command

	// 获取所有被封的IP
	respBody, err := http_client.Get(ctx, hc, unlockIPUrl)
	if err != nil {
		logger.Error().Err(err).Str("url", unlockIPUrl).Msg("get lock ip")
		return
	}

	resp := UnlockIPResponse{}
	json.Unmarshal(respBody, &resp)
	keys := resp.Data
	// 不需要处理，结束
	if len(keys) == 0 {
		talks.MakeTalkEnd(ctx, "没有需要解封的IP，本次服务结束")
		return
	}
	query := "请选择，回复序号：\n"
	for i, key := range keys {
		query += fmt.Sprintf("%d. %s\n", i, key)
	}

	// 让客服选择要处理的IP
	w.add("choose_release", query, func(ctx context.Context, replyMsg string) string {
		i, err := strconv.Atoi(replyMsg)
		if err != nil || i >= len(keys) {
			talks.ReplyMsg(ctx, "选择错误，请重新选择!")
			return "choose_release"
		} else {
			// 根据客服的选项，解封IP
			if _, err := http_client.Delete(ctx, dhc, unlockIPUrl+"&key="+url.QueryEscape(keys[i])); err != nil {
				logger.Error().Err(err).Msg("remove lock ip")
			}
		}
		return ""
	})

	if w.start("choose_release") == 0 {
		talks.MakeTalkEnd(ctx, "解除限制成功，本次服务结束")
	} else {
		talks.MakeTalkEnd(ctx, "解除限制失败，如有疑问请联系管理员，本次服务结束")
	}
}

func init() {
	registerTaskCommand("UnlockIP", UnlockIP)
}
