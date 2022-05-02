package main

import (
	"bot/pkg/config"
	"bot/pkg/dingding"
	"bot/pkg/handler"
	"bot/pkg/health"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/bot/dingding/talk", dingding.DingDing(handler.Handler))

	if len(config.C.Beat) > 0 {
		go health.BeatCheckHealth()
	}
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
