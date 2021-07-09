package dingding

import (
	"bot/pkg/logging"
	"bot/pkg/talk"
	"strings"

	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"net/http"
)

type DingDingRequest struct {
	Conversationid string `json:"conversationId"`
	Atusers        []struct {
		Dingtalkid string `json:"dingtalkId"`
	} `json:"atUsers"`
	Chatbotcorpid             string `json:"chatbotCorpId"`
	Chatbotuserid             string `json:"chatbotUserId"`
	Msgid                     string `json:"msgId"`
	Sendernick                string `json:"senderNick"`
	Isadmin                   bool   `json:"isAdmin"`
	Senderstaffid             string `json:"senderStaffId"`
	Sessionwebhookexpiredtime int64  `json:"sessionWebhookExpiredTime"`
	Createat                  int64  `json:"createAt"`
	Sendercorpid              string `json:"senderCorpId"`
	Conversationtype          string `json:"conversationType"`
	Senderid                  string `json:"senderId"`
	Conversationtitle         string `json:"conversationTitle"`
	Isinatlist                bool   `json:"isInAtList"`
	Sessionwebhook            string `json:"sessionWebhook"`
	Text                      struct {
		Content string `json:"content"`
	} `json:"text"`
	Msgtype string `json:"msgtype"`
}

type DingDingResponse struct {
	Msgtype string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	At struct {
		Atmobiles []string `json:"atMobiles"`
	} `json:"at"`
}

type DingDingAPP struct {
	NotifyUrl string
	Sender    string
}

func (dd *DingDingAPP) Request(c *gin.Context) string {

	// by test
	// jsonData, _ := ioutil.ReadAll(c.Request.Body)
	// fmt.Printf("reqeust %v\n", string(jsonData))

	rq := DingDingRequest{}
	c.BindJSON(&rq)

	dd.NotifyUrl = rq.Sessionwebhook
	dd.Sender = rq.Sendernick
	return strings.TrimSpace(rq.Text.Content)
}

func (dd DingDingAPP) Response(text string) DingDingResponse {
	rs := DingDingResponse{}
	rs.Msgtype = "text"
	rs.Text.Content = text
	// rs.At.Atmobiles = []string{"13427692994"}
	return rs
}

func (dd DingDingAPP) Notify(text string) {
	data := dd.Response(text)
	dataJson, _ := json.Marshal(data)

	logging.Log.Debug().
		Str("app", "dingding").
		Str("text", text).
		Msg("notify")

	var PTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}
	client := http.Client{
		Transport: PTransport,
	}

	resp, err := client.Post(dd.NotifyUrl, "application/json", bytes.NewBuffer(dataJson))
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
}

func DingDing(h func(t string, s chan string, r chan string, l zerolog.Logger)) gin.HandlerFunc {
	ddapp := DingDingAPP{}

	return func(c *gin.Context) {
		command := ddapp.Request(c)
		log := logging.Log.With().
			Str("app", "dingding").
			Str("module", "command").
			Str("command", command).
			Str("sender", ddapp.Sender).
			Logger()

		isFirst, sender, reply := talk.ContinueTaskSession(ddapp.Sender)
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
						talk.CloseTaskSession(ddapp.Sender)
						log.Info().Msg("talk end")
						return
					}
				}
			}
		}(sender, reply)

		h(command, sender, reply, log)
	}
}
