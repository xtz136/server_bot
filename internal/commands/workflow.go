package commands

import (
	"bot/pkg/talks"
	"context"
	"time"

	"github.com/rs/zerolog"
)

type stepFunc func(context.Context, string) string
type workStep struct {
	name     string
	msg      string
	stepFunc stepFunc
}

type WorkFlow struct {
	ctx       context.Context
	funcs     []workStep
	startTime time.Time
}

func (wf *WorkFlow) add(name string, msg string, f stepFunc) {
	wf.funcs = append(wf.funcs, workStep{name, msg, f})
}

func (wf *WorkFlow) getNext(name string) *workStep {
	for _, ws := range wf.funcs {
		if ws.name == name {
			return &ws
		}
	}
	return nil
}

// 开始这个流程，参数name代表第一个步骤的名称，示例：
//  w := newWorkFlow(ctx)
//  w.add("start", "你好，有什么可以帮到你？", func(ctx Context, replyMsg string) string {
//    return "next_step_name" // 根据返回值，跳到某个步骤
//  })
//  w.add("next_step_name", "还有什么可以帮到你？", func(ctx Context, replyMsg string) string {
//    return "" // 返回空字符串，代表流程结束
//  })
//  resultStatus := w.start("start")
//
// resultStatus 可能是以下几种：
//  0：成功
//  1：步骤没找到
//  2：sender channel 被关闭
//  3：reply channel 被关闭
//  4：进入了死循环，同一个步骤连续执行了10次
func (wf *WorkFlow) start(name string) int {
	lastName := name
	defaultMaxCall := 10
	maxCall := defaultMaxCall
	logger := wf.ctx.Value(talks.LoggerKey).(zerolog.Logger)

	for {
		ws := wf.getNext(name)
		if ws == nil {
			logger.Error().Str("name", name).Msg("workstep not found")
			return 1
		}
		// 发送消息给客户
		err := talks.ReplyMsg(wf.ctx, ws.msg)
		if err != nil {
			logger.Error().Err(err).Msg("reply message")
			return 2
		}
		// 等待客户回复
		msg, err := talks.ReceiveMsg(wf.ctx)
		if err != nil {
			logger.Error().Err(err).Msg("receive message")
			return 3
		}
		// 执行流程步骤，name 指明了下一个步骤的名称
		name = ws.stepFunc(wf.ctx, msg)
		// 流程完成
		if name == "" {
			return 0
		}
		// 避免死循环，同一个步骤不能连续调用10次
		if lastName == name {
			if maxCall--; maxCall == 0 {
				logger.Error().Str("name", name).Msg("circular call")
				return 4
			}
		} else {
			lastName = name
			maxCall = defaultMaxCall
		}

	}
}

func (wf *WorkFlow) getCostTime() time.Duration {
	return time.Since(wf.startTime)
}

func newWorkFlow(ctx context.Context) *WorkFlow {
	w := WorkFlow{ctx: ctx}
	w.startTime = time.Now()
	return &w
}
