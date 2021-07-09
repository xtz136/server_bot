package commands

import (
	"fmt"
)

func Dummy(ctx Context) {
	w := newWorkFlow(ctx)
	w.add(func(ctx Context) bool {
		ctx.Sender <- "有什么需要帮助吗"
		return true
	})
	w.add(func(ctx Context) bool {
		msg, ok := <-ctx.Reply
		if ok {
			ctx.Log.Info().Msg("got repl")
			ctx.State["msg"] = msg
			return true
		}
		return false
	})
	w.add(func(ctx Context) bool {
		ctx.Log.Info().Msg("to done")
		ctx.Sender <- fmt.Sprintf("你回复了：%s", ctx.State["msg"])
		return true
	})

	if w.start() {
		eT := w.getCostTime()
		MakeTalkEnd(ctx.Sender, fmt.Sprintf("测试结束，耗时：%v，本次服务结束", eT))
	}
}
