package commands

import (
	"fmt"
)

// 测试任务
func Dummy(ctx Context) {
	w := newWorkFlow(ctx)
	w.add("start", "有什么需要帮助吗？", func(ctx Context, replyMsg string) string {
		ctx.Sender <- fmt.Sprintf("你回复了：%s", replyMsg)
		return "pingjia"
	})
	w.add("pingjia", "请为我的服务评价，回复数字1为好评，数字2为中等，数字3为差评。", func(ctx Context, replyMsg string) string {
		if replyMsg == "1" {
			ctx.Sender <- "本次服务好评！"
		} else if replyMsg == "2" {
			ctx.Sender <- "本次服务中等！"
		} else if replyMsg == "3" {
			ctx.Sender <- "本次服务差评！"
		} else {
			ctx.Sender <- "回复错误，请重新回复！"
			return "pingjia"
		}
		return ""
	})
	w.start("start")
	eT := w.getCostTime()
	ctx.MakeTalkEnd(ctx.Sender, fmt.Sprintf("测试结束，耗时：%v，本次服务结束", eT))
}

func init() {
	registerTaskCommand("Dummy", Dummy)
}
