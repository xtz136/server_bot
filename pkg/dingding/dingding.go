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
	"github.com/rs/zerolog"
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
	NotifyUrl string
	Sender    string
	Command   string
	appSecret string
}

func (dd *DingDingAPP) Parse(c *gin.Context) {

	rq := DingDingRequest{}
	c.BindJSON(&rq)

	dd.NotifyUrl = rq.Sessionwebhook
	dd.Sender = rq.Sendernick
	dd.Command = strings.TrimSpace(rq.Text.Content)
}

func (dd DingDingAPP) getSenderName() string {
	return dd.Sender
}

func (dd DingDingAPP) getCommand() string {
	return dd.Command
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
func (dd DingDingAPP) Notify(text string) {
	data := dd.Response(text)
	dataJson, _ := json.Marshal(data)

	logging.Log.Debug().
		Str("app", "dingding").
		Str("text", text).
		Msg("notify")

	dhc := http_client.NewDumbHttpClient(10)
	_, err := http_client.PostJson(dhc, dd.NotifyUrl, bytes.NewBuffer(dataJson))
	if err != nil {
		logging.Log.Error().AnErr("err", err).Msg("notify error")
	}
}

func DingDing(handler func(string, chan string, chan string, zerolog.Logger)) gin.HandlerFunc {
	ddapp := DingDingAPP{appSecret: config.C.DingDing.AppSecret}

	return func(c *gin.Context) {
		ddapp.Parse(c)
		senderName := ddapp.getSenderName()
		command := ddapp.getCommand()

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
			log.Debug().Str("timestamp", timestamp).Str("sign", sign).Int("err", err).Msg("非法操作")
			ddapp.Notify("非法操作")
			return
		}

		isFirst, sender, reply := talk.ContinueTaskSession(senderName, command)
		log.Info().Bool("isFirst", isFirst).Msg("got request")

		// 会话中，直接将命令通过管道发送给正在执行的机器人
		if !isFirst {
			reply <- command
			return
		}

		// 开始会话，将机器人的回复/结果，发送给用户。并处理结束会话
		go func(sender chan string, reply chan string, senderName string, command string) {
			for {
				select {
				case msg, ok := <-sender:
					if ok {
						ddapp.Notify(msg)
					} else {
						talk.CloseTaskSession(senderName, command)
						log.Info().Msg("talk end")
						return
					}
				}
			}
		}(sender, reply, senderName, command)

		// 开始进入任务工作
		handler(command, sender, reply, log)
	}
}
