package main

import (
	"strings"

	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"

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
	notifyUrl string
	sender    string
}

func (dd *DingDingAPP) request(c *gin.Context) string {

	// by test
	// jsonData, _ := ioutil.ReadAll(c.Request.Body)
	// fmt.Printf("reqeust %v\n", string(jsonData))

	rq := DingDingRequest{}
	c.BindJSON(&rq)

	dd.notifyUrl = rq.Sessionwebhook
	dd.sender = rq.Sendernick
	return strings.TrimSpace(rq.Text.Content)
}

func (dd DingDingAPP) response(text string) DingDingResponse {
	rs := DingDingResponse{}
	rs.Msgtype = "text"
	rs.Text.Content = text
	// rs.At.Atmobiles = []string{"13427692994"}
	return rs
}

func (dd DingDingAPP) notify(text string) {
	data := dd.response(text)
	dataJson, _ := json.Marshal(data)

	Log.Debug().
		Str("app", "dingding").
		Str("text", text).
		Msg("notify")

	var PTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}
	client := http.Client{
		Transport: PTransport,
	}

	resp, err := client.Post(dd.notifyUrl, "application/json", bytes.NewBuffer(dataJson))
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
}
