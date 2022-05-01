package talk

import (
	"time"
)

type talk struct {
	senderChan chan string
	replyChan  chan string
	senderName string
	command    string
}

var talks = []talk{}

func addTalk(sender chan string, reply chan string, senderName string, command string) {
	talks = append(talks, talk{sender, reply, senderName, command})
}

func removeTalk(senderName string, command string) bool {
	for i, talk := range talks {
		if talk.senderName == senderName && talk.command == command {
			close(talk.senderChan)
			close(talk.replyChan)

			lastLen := len(talks) - 1
			talks[i] = talks[lastLen]
			talks[lastLen] = talk
			talks = talks[:lastLen]
			return true
		}
	}
	return false
}

func findTalk(command string, senderName string) *talk {
	for _, talk := range talks {
		if talk.senderName == senderName && talk.command == command {
			return &talk
		}
	}
	return nil
}

// 创建会话session或者找到已存在的会话session。
//
// 第一个返回值为true表示创建成功，false表示已存在。
func ContinueTaskSession(senderName string, command string) (bool, chan string, chan string) {
	talk := findTalk(command, senderName)
	if talk != nil {
		return false, talk.senderChan, talk.replyChan
	}

	senderChan := make(chan string)
	replyChan := make(chan string)
	addTalk(senderChan, replyChan, senderName, command)
	return true, senderChan, replyChan
}

// 关闭会话session
func CloseTaskSession(senderName string, command string) {
	removeTalk(senderName, command)
}

func MakeTalkEnd(sender chan string, lastMsg string) {
	if lastMsg != "" {
		sender <- lastMsg
	}

	time.Sleep(time.Duration(1) * time.Second)
	close(sender)
}
