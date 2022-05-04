package main

import (
	"bot/pkg/beat"
	"bot/pkg/config"
	"bot/pkg/dingding"
	"bot/pkg/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/bot/dingding/talk", dingding.DingDing(handler.Handler))

	if len(config.C.Beat) > 0 {
		go beat.StartBeat()
	}
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
