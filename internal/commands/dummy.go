package commands

import (
	"bot/pkg/talks"
	"context"
	"fmt"
)

// 测试任务
func Dummy(ctx context.Context) {
	w := newWorkFlow(ctx)
	w.add("start", "有什么需要帮助吗？", func(ctx context.Context, replyMsg string) string {
		err := talks.ReplyMsg(ctx, fmt.Sprintf("你回复了：%s", replyMsg))
		if err != nil {
			panic(err)
		}
		return "pingjia"
	})
	w.add("pingjia", "请为我的服务评价，回复数字1为好评，数字2为中等，数字3为差评。", func(ctx context.Context, replyMsg string) string {
		if replyMsg == "1" {
			talks.ReplyMsg(ctx, "本次服务好评！")
		} else if replyMsg == "2" {
			talks.ReplyMsg(ctx, "本次服务中等！")
		} else if replyMsg == "3" {
			talks.ReplyMsg(ctx, "本次服务差评！")
		} else {
			talks.ReplyMsg(ctx, "回复错误，请重新回复！")
			return "pingjia"
		}
		return ""
	})
	w.start("start")
	eT := w.getCostTime()
	talks.MakeTalkEnd(ctx, fmt.Sprintf("测试结束，耗时：%v，本次服务结束", eT))
}

func init() {
	registerTaskCommand("Dummy", Dummy)
}
