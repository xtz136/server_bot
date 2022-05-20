package main

import (
	"bot/pkg/beat"
	"bot/pkg/config"
	"bot/pkg/dingding"
	"bot/pkg/handler"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/bot/dingding/talk", dingding.DingDing(handler.Handler))

	if len(config.C.Beat) > 0 {
		go beat.StartBeat()
	}
	r.Run(fmt.Sprintf("%s:%d", config.C.App.Host, config.C.App.Port))
}
