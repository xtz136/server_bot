package dingding

import (
	"bot/pkg/config"
	"bot/pkg/http_client"
	"bot/pkg/logging"
	"bot/pkg/talk"
	"crypto/hmac"
	"crypto/sha256"
	"strconv"
	"strings"
	"time"

	"bytes"
	"encoding/base64"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func hmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	sha := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(sha)
}

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
	notifyUrl   string
	appSecret   string
	senderName  string
	sender      chan string
	reply       chan string
	command     string
	isFirstTalk bool
}

func (dd *DingDingAPP) Parse(c *gin.Context) {

	rq := DingDingRequest{}
	c.BindJSON(&rq)

	dd.notifyUrl = rq.Sessionwebhook
	dd.senderName = rq.Sendernick
	dd.command = strings.TrimSpace(rq.Text.Content)
}

func (dd DingDingAPP) GetSenderName() string {
	return dd.senderName
}

func (dd DingDingAPP) GetCommand() string {
	return dd.command
}

func (dd DingDingAPP) GetSender() chan string {
	return dd.sender
}

func (dd DingDingAPP) GetReply() chan string {
	return dd.reply
}

func (dd DingDingAPP) IsFirstTalk() bool {
	return dd.isFirstTalk
}

func (dd DingDingAPP) Response(text string) DingDingResponse {
	rs := DingDingResponse{}
	rs.Msgtype = "text"
	rs.Text.Content = text
	return rs
}

// 验证钉钉接口是否可信，根据钉钉规范验证
func (dd DingDingAPP) check(timestamp string, sign string) int {
	t, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return 1
	}

	// 检查 timestamp
	gap := time.Now().Unix() - t
	if gap < 0 && gap > 86400 {
		return 2
	}

	// 检查 sign
	appSecret := dd.appSecret
	string_to_sign := timestamp + "\n" + appSecret

	if hmacSha256(string_to_sign, appSecret) != sign {
		return 3
	}

	return 0
}

// 给钉钉发送消息，目前只能发文本信息
func (dd DingDingAPP) ReplyMessage(text string) {
	data := dd.Response(text)
	dataJson, _ := json.Marshal(data)

	logging.Log.Debug().
		Str("app", "dingding").
		Str("text", text).
		Msg("notify")

	dhc := http_client.NewDumbHttpClient(10)
	_, err := http_client.PostJson(dhc, dd.notifyUrl, bytes.NewBuffer(dataJson))
	if err != nil {
		logging.Log.Error().AnErr("err", err).Msg("notify error")
	}
}

func DingDing(handler func(talk.TalkInterface)) gin.HandlerFunc {

	return func(c *gin.Context) {
		ddapp := DingDingAPP{appSecret: config.C.DingDing.AppSecret}
		ddapp.Parse(c)
		senderName := ddapp.GetSenderName()
		command := ddapp.GetCommand()

		log := logging.Log.With().
			Caller().
			Str("app", "dingding").
			Str("command", command).
			Str("sender", senderName).
			Logger()

		// 确认是钉钉服务器发送的请求
		timestamp := c.Request.Header.Get("timestamp")
		sign := c.Request.Header.Get("sign")
		if err := ddapp.check(timestamp, sign); err != 0 {
			log.Warn().Str("timestamp", timestamp).Str("sign", sign).Int("err", err).Msg("非法操作")
			ddapp.ReplyMessage("非法操作")
			return
		}

		isFirst, sender, reply := talk.ContinueTaskSession(senderName, command)
		log.Info().Bool("isFirst", isFirst).Msg("got request")

		// 开始进入任务工作
		ddapp.sender = sender
		ddapp.reply = reply
		ddapp.isFirstTalk = isFirst
		handler(&ddapp)
	}
}
