package commands

import (
	"time"
)

type stepFunc func(Context, string) string
type workStep struct {
	name     string
	msg      string
	stepFunc stepFunc
}

type WorkFlow struct {
	ctx       Context
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
//  w.add("start", "你好，有什么可以帮到你？", func(ctx Context, replayMsg string) string {
//    return "next_step_name" // 根据返回值，跳到某个步骤
//  })
//  w.add("next_step_name", "还有什么可以帮到你？", func(ctx Context, replayMsg string) string {
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

	for {
		ws := wf.getNext(name)
		if ws == nil {
			wf.ctx.Log.Error().Str("name", name).Msg("workstep not found")
			return 1
		}
		// 发送消息给客户
		select {
		case _, ok := <-wf.ctx.Sender:
			// 防止sender提前被关闭
			if !ok {
				wf.ctx.Log.Error().Msg("sender channel closed")
				return 2
			}
			wf.ctx.Sender <- ws.msg
		default:
			wf.ctx.Sender <- ws.msg
		}
		// 等待客户回复
		replyMsg, ok := <-wf.ctx.Reply
		if !ok {
			wf.ctx.Log.Error().Msg("reply channel closed")
			return 3
		}
		// 执行流程步骤，name 指明了下一个步骤的名称
		name = ws.stepFunc(wf.ctx, replyMsg)
		// 流程完成
		if name == "" {
			return 0
		}
		// 避免死循环，同一个步骤不能连续调用10次
		if lastName == name {
			if maxCall--; maxCall == 0 {
				wf.ctx.Log.Error().Str("name", name).Msg("circular call")
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

func newWorkFlow(ctx Context) *WorkFlow {
	w := WorkFlow{ctx: ctx}
	w.startTime = time.Now()
	return &w
}
