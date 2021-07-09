package talk

import (
	"time"
)

var Talks = map[string][]chan string{}

func CreateTaskSession() (chan string, chan string) {
	sender := make(chan string)
	reply := make(chan string)
	return sender, reply
}

func ContinueTaskSession(sender string) (bool, chan string, chan string) {
	for name := range Talks {
		if name == sender {
			return false, Talks[name][0], Talks[name][1]
		}
	}

	a, b := CreateTaskSession()
	Talks[sender] = []chan string{
		a, b,
	}
	return true, a, b
}

func CloseTaskSession(sender string) {
	for name := range Talks {
		if name == sender {
			delete(Talks, name)
		}
	}
}

func MakeTalkEnd(sender chan string, lastMsg string) {
	if lastMsg != "" {
		sender <- lastMsg
	}

	time.Sleep(time.Duration(1) * time.Second)
	close(sender)
}
