package commands

import (
	"bot/pkg/talks"
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
)

func echoStep(ctx context.Context, replyMsg string) string {
	// 根据回复，跳到某个步骤
	return replyMsg
}

func make_workflow() *WorkFlow {

	ctx := context.Background()
	ctx = context.WithValue(ctx, talks.SenderKey, make(chan string))
	ctx = context.WithValue(ctx, talks.ReplyKey, make(chan string))
	ctx = context.WithValue(ctx, talks.LoggerKey, log.With().Logger())
	ctx, _ = context.WithTimeout(ctx, 100*time.Millisecond)

	return &WorkFlow{
		ctx: ctx,
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
			defer func() {
				sender := wf.ctx.Value(talks.SenderKey).(chan string)
				close(sender)
			}()
			defer func() {
				reply := wf.ctx.Value(talks.ReplyKey).(chan string)
				close(reply)
			}()

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
				sender := wf.ctx.Value(talks.SenderKey).(chan string)
				defer close(sender)
				for _, expect := range tt.senderMsgList {
					msg := <-sender
					if msg != expect {
						t.Errorf("WorkFlow.start() sender msg = %v, want %v", msg, expect)
					}
				}
			}()
			// 模拟服务器回复
			go func() {
				reply := wf.ctx.Value(talks.ReplyKey).(chan string)
				defer close(reply)
				for _, msg := range tt.replyMsgList {
					reply <- msg
				}
			}()
			if got := wf.start(tt.args); got != tt.want {
				t.Errorf("WorkFlow.start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkFlow_Timeout_start(t *testing.T) {
	tests := []struct {
		name          string
		workflow      *WorkFlow
		senderMsgList []string
		replyMsgList  []string
		args          string
		want          int
	}{
		{
			"it timeout by blocked reply",
			make_workflow(),
			[]string{"start?", "done?"},
			[]string{"done"},
			"start",
			3,
		},
		{
			"it timeout by blocked sender",
			make_workflow(),
			[]string{"start?"},
			[]string{"done", ""},
			"start",
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wf := tt.workflow
			// 模拟客户接收消息
			go func() {
				sender := wf.ctx.Value(talks.SenderKey).(chan string)
				for _, expect := range tt.senderMsgList {
					msg := <-sender
					if msg != expect {
						t.Errorf("WorkFlow.start() sender msg = %v, want %v", msg, expect)
					}
				}
			}()
			// 模拟服务器回复
			go func() {
				reply := wf.ctx.Value(talks.ReplyKey).(chan string)
				for _, msg := range tt.replyMsgList {
					reply <- msg
				}
			}()
			if got := wf.start(tt.args); got != tt.want {
				t.Errorf("WorkFlow.start() = %v, want %v", got, tt.want)
			}
		})
	}
}
