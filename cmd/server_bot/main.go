package main

import (
	"bot/pkg/beat"
	"bot/pkg/config"
	"bot/pkg/handler"
	"bot/pkg/websocket"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	hub := websocket.NewHub()
	go hub.Run()

	r := gin.Default()
	r.POST("/bot/dingding/talk", handler.DingDing(handler.Handler))

	r.GET("/bot/client", gin.WrapF(func(w http.ResponseWriter, req *http.Request) {
		websocket.ServeWs(hub, w, req)
	}))

	r.GET("/", gin.WrapF(func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "home.html")
	}))

	r.GET("/bot/info", gin.WrapF(func(w http.ResponseWriter, req *http.Request) {
		msg := ""
		hub := websocket.NewHub()
		for clientName, clientIdents := range hub.CollectClientsName() {
			msg += fmt.Sprintf("* %s[%s]\n", clientName, clientIdents)
		}
		w.Write([]byte(msg))
	}))

	if len(config.C.Beat) > 0 {
		go beat.StartBeat()
	}
	r.Run(fmt.Sprintf("%s:%d", config.C.App.Host, config.C.App.Port))
}
