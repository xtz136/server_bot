package commands

import (
	"reflect"
	"testing"
	"time"
)

func echoStep(ctx Context, replyMsg string) string {
	// 根据回复，跳到某个步骤
	return replyMsg
}

func make_workflow() *WorkFlow {
	return &WorkFlow{
		ctx: Context{Sender: make(chan string), Reply: make(chan string)},
		funcs: []workStep{
			{"start", "start?", echoStep},
			{"circle", "circle?", echoStep},
			{"done", "done?", echoStep},
		},
		startTime: time.Now(),
	}
}

func TestWorkFlow_getNext(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name     string
		workflow *WorkFlow
		args     args
		want     *workStep
	}{
		{"not found", make_workflow(), args{"not exists"}, nil},
		{
			"find first",
			make_workflow(),
			args{"start"},
			&workStep{"start", "start?", echoStep},
		},
		{
			"find not first",
			make_workflow(),
			args{"done"},
			&workStep{"done", "done?", echoStep},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wf := tt.workflow
			defer close(wf.ctx.Sender)
			defer close(wf.ctx.Reply)

			got := wf.getNext(tt.args.name)
			// stepFunc 不重要
			if got != nil {
				got.stepFunc = nil
			}
			if tt.want != nil {
				tt.want.stepFunc = nil
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WorkFlow.getNext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkFlow_start(t *testing.T) {
	tests := []struct {
		name          string
		workflow      *WorkFlow
		senderMsgList []string
		replyMsgList  []string
		args          string
		want          int
	}{
		{
			"good",
			make_workflow(),
			[]string{"start?", "done?"},
			[]string{"done", ""},
			"start",
			0,
		},
		{
			"quit early",
			make_workflow(),
			[]string{"start?", "done"},
			[]string{""},
			"start",
			0,
		},
		{
			"invalid step name",
			make_workflow(),
			[]string{},
			[]string{},
			"not exists",
			1,
		},
		{
			"close sender early",
			make_workflow(),
			[]string{"start?"},
			[]string{"done", ""},
			"start",
			2,
		},
		{
			"close replay early",
			make_workflow(),
			[]string{"start?", "done?"},
			[]string{"done"},
			"start",
			3,
		},
		{
			"circular call",
			make_workflow(),
			[]string{"start?", "start?", "start?", "start?", "start?", "start?", "start?", "start?", "start?", "start?"},
			[]string{"start", "start", "start", "start", "start", "start", "start", "start", "start", "start"},
			"start",
			4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wf := tt.workflow
			// 模拟客户接收消息
			go func() {
				defer close(wf.ctx.Sender)
				for _, expect := range tt.senderMsgList {
					msg := <-wf.ctx.Sender
					if msg != expect {
						t.Errorf("WorkFlow.start() sender msg = %v, want %v", msg, expect)
					}
				}
			}()
			// 模拟服务器回复
			go func() {
				defer close(wf.ctx.Reply)
				for _, msg := range tt.replyMsgList {
					wf.ctx.Reply <- msg
				}
			}()
			if got := wf.start(tt.args); got != tt.want {
				t.Errorf("WorkFlow.start() = %v, want %v", got, tt.want)
			}
		})
	}
}
