package commands

func EnableUser(ctx Context) {
	// sender <- "请输入用户账号"
	// account := <-reply
	// fmt.Printf("enable %s user\n", account)
	// sender <- fmt.Sprintf("用户%s已启用", account)
	// sender <- ""
	ctx.MakeTalkEnd(ctx.Sender, "这个功能还在建设中")
}
