package handler

import (
	"bot/internal/commands"
	"bot/pkg/talks"
	"context"
	"sync"
	"testing"
	"time"
)

type FakeTalk struct {
	command string
	sender  chan string
	replay  chan string
	talkNum int
}

func (t *FakeTalk) setCommand(command string) {
	t.talkNum += 1
	t.command = command
}

func (t *FakeTalk) ReplyMessage(ctx context.Context, msg string) {
}

func (t *FakeTalk) GetSender() chan string {
	return t.sender
}

func (t *FakeTalk) GetSenderName() string {
	return "Tester"
}

func (t *FakeTalk) GetReply() chan string {
	return t.replay
}

func (t *FakeTalk) GetCommand() string {
	return t.command
}

func (t *FakeTalk) IsFirstTalk() bool {
	return t.talkNum == 1
}

func NewFakeTalk() *FakeTalk {
	return &FakeTalk{
		command: "",
		sender:  make(chan string),
		replay:  make(chan string),
		talkNum: 0,
	}
}

func TestHandler(t *testing.T) {
	tests := []struct {
		name         string
		commands     []string
		callFuncName string
		callNum      int
	}{
		{"target not exists", []string{"重启XXX"}, "", 0},
		{"call help when command is empty", []string{""}, "ShowHelper", 1},
		{"call help when command is not found", []string{"XXXSWEE"}, "ShowHelper", 1},
		{"call help", []string{"帮助"}, "ShowHelper", 1},
		// 下面测试依赖配置文件，确保 ~/server_bot.yaml 存在
		{"call normal func", []string{"重启本地"}, "Restart", 1},
		{"cancel talk", []string{"重启本地", "取消"}, "Restart", 1},
		{"unlockip talk", []string{"解锁ip本地"}, "UnlockIP", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			talk := NewFakeTalk()
			tt := tt

			wg := &sync.WaitGroup{}
			wg.Add(tt.callNum)
			commands.TaskCommands = make(map[string]func(context.Context))
			commands.TaskCommands[tt.callFuncName] = func(ctx context.Context) {
				wg.Done()
			}

			for _, command := range tt.commands {
				talk.setCommand(command)
				Handler(talk)
			}

			done := make(chan struct{})
			go func() {
				wg.Wait()
				done <- struct{}{}
			}()

			select {
			case <-done:
				// ensure session is empty
				if len(talks.Sessions) != 0 {
					t.Errorf("%s: session is not empty", tt.name)
				}
			case <-time.After(1 * time.Second):
				t.Error("test is timeout")
			}
		})
	}
}
