package main

import (
	"bot/internal/commands"
	"bot/pkg/dingding"
	"bot/pkg/logging"
	"fmt"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// var Sender = make(chan string)
// var Reply = make(chan string)

type Notify func(t string)

type APP interface {
	Request()
	Response()
	Notify()
}

var Talks = map[string][]chan string{}

func createTaskSession() (chan string, chan string) {
	sender := make(chan string)
	reply := make(chan string)
	return sender, reply
}

func continueTaskSession(sender string) (bool, chan string, chan string) {
	for name := range Talks {
		if name == sender {
			return false, Talks[name][0], Talks[name][1]
		}
	}

	a, b := createTaskSession()
	Talks[sender] = []chan string{
		a, b,
	}
	return true, a, b
}

func closeTaskSession(sender string) {
	for name := range Talks {
		if name == sender {
			delete(Talks, name)
		}
	}
}

func DingDing(h func(t string, s chan string, r chan string, l zerolog.Logger)) gin.HandlerFunc {
	ddapp := dingding.DingDingAPP{}

	return func(c *gin.Context) {
		command := ddapp.Request(c)
		log := logging.Log.With().
			Str("app", "dingding").
			Str("module", "command").
			Str("command", command).
			Str("sender", ddapp.Sender).
			Logger()

		isFirst, sender, reply := continueTaskSession(ddapp.Sender)
		log.Info().Bool("isFirst", isFirst).Msg("got request")

		// 会话中
		if !isFirst {
			if command == "取消" {
				sender <- "取消成功"
				close(sender)
			} else {
				reply <- command
			}
			return
		}

		// 开始会话
		// TODO 协程需要增加超时退出机制
		go func(sender chan string, reply chan string) {
			for {
				select {
				case msg, ok := <-sender:
					if ok {
						ddapp.Notify(msg)
					} else {
						close(reply)
						closeTaskSession(ddapp.Sender)
						log.Info().Msg("talk end")
						return
					}
				}
			}
		}(sender, reply)

		h(command, sender, reply, log)
	}
}

func Talk(command string, sender chan string, reply chan string, log zerolog.Logger) {

	ctx := commands.Context{
		"",
		sender,
		reply,
		log,
		make(commands.State),
	}

	if strings.HasPrefix(command, "重启") {
		go commands.Restart(command[6:], sender, reply, log)
		return
	}

	switch command {
	case "启用阿里云用户":
		go commands.EnableUser("阿里云", sender, reply, log)
	case "解封阿里云IP限制":
		go commands.AllowIP("阿里云", sender, reply, log)
	case "使用阿里云用户登录":
		go commands.LoginUser("阿里云", sender, reply, log)
	case "测试":
		ctx.Name = "测试"
		go commands.Dummy(ctx)
	case "获取会话数":
		commands.MakeTalkEnd(sender, fmt.Sprintf("%d", runtime.NumGoroutine()))
	default:
		helpMsg := "命令列表：\n"
		helpMsg += "* 帮助\n"

		// 重启服务器动态读取配置文件
		for _, k := range commands.ListServices() {
			helpMsg += "* 重启" + k + "\n"
		}

		helpMsg += "* 启用阿里云用户\n"
		helpMsg += "* 解封阿里云IP限制"

		commands.MakeTalkEnd(sender, helpMsg)
	}
}

func initConfig() bool {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	return true
}

func main() {
	initConfig()
	logging.InitLog()

	r := gin.Default()
	r.POST("/bot/dingding/talk", DingDing(Talk))
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
