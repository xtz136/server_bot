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

var Sender = make(chan string)
var Reply = make(chan string)

type Notify func(t string)

type APP interface {
	Request()
	Response()
	Notify()
}

func DingDing(h func(t string, s chan string, r chan string, l zerolog.Logger)) gin.HandlerFunc {
	var talking bool = false

	return func(c *gin.Context) {
		command := ddapp.request(c)
		log := Log.With().
			Str("app", "dingding").
			Str("module", "services").
			Str("command", command).
			Str("sender", ddapp.sender).
			Logger()

		if talking && command == "取消" {
			// TODO 协程需要增加超时退出机制
			Sender <- ""
		} else if talking {
			Reply <- command
		} else {
			go func() {
				for {
					msg := <-Sender
					if msg == "" {
						talking = false
						log.Info().Msg("talk end")
						break
					} else {
						talking = true
						ddapp.notify(msg)
					}
				}
			}()
			h(command, Sender, Reply, log)
		}
	}
}

func Talk(command string, sender chan string, reply chan string, log zerolog.Logger) {

	switch command {
	case "重启阿里云":
		go restart("aly", sender, reply, log)
	case "启用阿里云用户":
		go enableUser("aly", sender, reply, log)
	case "解封阿里云IP限制":
		go allowIP("aly", sender, reply, log)
	case "使用阿里云用户登录":
		go loginUser("aly", sender, reply, log)
	case "测试":
		go dummy("xx系统", sender, reply, log)
	default:
		helpMsg := `命令列表：
		1. 帮助
		2. 重启阿里云
		3. 启用阿里云用户
		4. 解封阿里云IP限制
		5. 使用阿里云用户登录
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
