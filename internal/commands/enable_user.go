package commands

import "github.com/rs/zerolog"

func EnableUser(systemName string, sender chan string, reply chan string, logger zerolog.Logger) {
	// sender <- "请输入用户账号"
	// account := <-reply
	// fmt.Printf("enable %s user\n", account)
	// sender <- fmt.Sprintf("用户%s已启用", account)
	// sender <- ""
	MakeTalkEnd(sender, "这个功能还在建设中")
}
