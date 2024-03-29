package talks

import (
	"context"
	"errors"
	"sync"
	"time"
)

type key int

const (
	SenderKey key = iota
	ReplyKey
	TargetNameKey
	TargetKey
	TaskKey
	LoggerKey
	Cancel
)

type TalkInterface interface {
	// 给客户发送消息
	ReplyMessage(context.Context, string)
	// 获取 sender
	GetSender() chan string
	// 获取客户名称
	GetSenderName() string
	// 获取 reply
	GetReply() chan string
	// 获取客户下达的指令
	GetCommand() string
	// 是否是第一次会话
	IsFirstTalk() bool
}

// 保存会话session
type Session struct {
	SenderChan chan string
	ReplyChan  chan string
	SenderName string
	Command    string
}

var Sessions = []Session{}

func addTalk(sender chan string, reply chan string, senderName string, command string) {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	Sessions = append(Sessions, Session{sender, reply, senderName, command})
}

func removeTalk(senderName string, command string) bool {
	for i, talk := range Sessions {
		if talk.SenderName == senderName && talk.Command == command {
			var mutex sync.Mutex
			mutex.Lock()
			defer mutex.Unlock()

			// 关闭 reply
			close(talk.ReplyChan)
			// 把talk从Sessions中删掉
			lastLen := len(Sessions) - 1
			Sessions[i] = Sessions[lastLen]
			Sessions[lastLen] = talk
			Sessions = Sessions[:lastLen]

			return true
		}
	}
	return false
}

func findTalk(command string, senderName string) *Session {
	for _, talk := range Sessions {
		if talk.SenderName == senderName {
			return &talk
		}
	}
	return nil
}

// 给客户发送消息
func ReplyMsg(ctx context.Context, msg string) error {
	sender := ctx.Value(SenderKey).(chan string)

	select {
	case <-ctx.Done():
		return errors.New("talk is closed")
	case _, ok := <-sender:
		// 防止sender提前被关闭
		if !ok {
			return errors.New("sender is closed")
		}
		sender <- msg
	default:
		sender <- msg
	}
	return nil
}

// 接受客户的回复
func ReceiveMsg(ctx context.Context) (string, error) {
	reply := ctx.Value(ReplyKey).(chan string)

	select {
	case <-ctx.Done():
		return "", errors.New("talk is closed")
	case msg, ok := <-reply:
		if !ok {
			return "", errors.New("reply is closed")
		}
		return msg, nil
	case <-time.After(30 * time.Second):
		return "", errors.New("talk is timeout")
	}
}

// 创建会话session或者找到已存在的会话session。
//
// 第一个返回值为true表示创建成功，false表示已存在。
func ContinueTaskSession(senderName string, command string) (bool, chan string, chan string) {
	talk := findTalk(command, senderName)
	if talk != nil {
		return false, talk.SenderChan, talk.ReplyChan
	}

	senderChan := make(chan string)
	replyChan := make(chan string)
	addTalk(senderChan, replyChan, senderName, command)
	return true, senderChan, replyChan
}

// 消除会话session
func DestoryTalkSession(senderName string, command string) {
	removeTalk(senderName, command)
}

// 关闭整个会话/任务，允许发最后一个消息
func MakeTalkEnd(ctx context.Context, lastMsg string) {
	if lastMsg != "" {
		ReplyMsg(ctx, lastMsg)
		// 等待消息发送完毕
		// time.Sleep(time.Duration(1) * time.Second)
	}

	sender := ctx.Value(SenderKey).(chan string)
	close(sender)
}
