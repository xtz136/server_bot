package services

import (
	"fmt"
	"time"
)

func Dummy(ctx Task) {
	steps := struct {
		start int
		reply int
		done  int
	}{
		1, 2, 3,
	}

	var msg string
	status := steps.start
	bT := time.Now()

loop:
	for {
		switch status {
		case steps.start:
			ctx.Sender <- "有什么需要帮助吗"
			status = steps.reply

		case steps.reply:
			msg = <-ctx.Reply
			ctx.Log.Info().Msg("got repl")
			status = steps.done

		case steps.done:
			ctx.Log.Info().Msg("to done")
			ctx.Sender <- fmt.Sprintf("你回复了：%s", msg)
			break loop
		}
	}

	eT := time.Since(bT)

	MakeTalkEnd(ctx.Sender, fmt.Sprintf("测试结束，耗时：%v，本次服务结束", eT))
}
