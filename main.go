package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var EnableConfig = initConfig()
var Log = getLog()
var ddapp = DingDingAPP{}

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
	return true, a, b
}

func DingDing(h func(t string, s chan string, r chan string, l zerolog.Logger)) gin.HandlerFunc {
	return func(c *gin.Context) {
		command := ddapp.request(c)
		log := Log.With().
			Str("app", "dingding").
			Str("module", "services").
			Str("command", command).
			Str("sender", ddapp.sender).
			Logger()

		isFirst, sender, reply := continueTaskSession(ddapp.sender)

		if isFirst {
			// select {
			// case msg, ok := (<-sender):
			// 	if ok {
			// 		ddapp.notify(msg)
			// 	}else{
			// 		log.Info().Msg("talk end")
			// 		break
			// 	}
			// }
			go func() {
				for {
					msg := <-sender
					if msg == "" {
						log.Info().Msg("talk end")
						break
					} else {
						ddapp.notify(msg)
					}
				}
			}()
			h(command, sender, reply, log)
		} else {
			if command == "取消" {
				// TODO 协程需要增加超时退出机制
				sender <- ""
			} else {
				reply <- command
			}
		}
	}
}

func Talk(command string, sender chan string, reply chan string, log zerolog.Logger) {

	switch command {
	case "重启阿里云":
		go restart("阿里云", sender, reply, log)
	case "重启长沙":
		go restart("长沙", sender, reply, log)
	case "启用阿里云用户":
		go enableUser("阿里云", sender, reply, log)
	case "解封阿里云IP限制":
		go allowIP("阿里云", sender, reply, log)
	case "使用阿里云用户登录":
		go loginUser("阿里云", sender, reply, log)
	case "测试":
		go dummy("测试", sender, reply, log)
	default:
		helpMsg := `命令列表：
		*. 帮助
		*. 重启阿里云
		*. 重启长沙
		*. 启用阿里云用户
		*. 解封阿里云IP限制
		`
		makeTalkEnd(sender, helpMsg)
	}
}

func initConfig() bool {
	// 这个方法有点奇怪，为了实现在 getLog 调用前，初始化 viper
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	return true
}

func main() {
	r := gin.Default()
	r.POST("/bot/dingding/talk", DingDing(Talk))
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
