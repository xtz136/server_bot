package commands

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func AllowIP(systemName string, sender chan string, reply chan string, logger zerolog.Logger) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	lockIP := viper.GetString("systems." + systemName + ".lock_ip")

	// 获取所有被封的IP
	resp, err := http.Get(lockIP)
	if err != nil {
		logger.Error().Err(err).Str("url", lockIP).Msg("get lock ip")
		sender <- ""
		return
	}
	defer resp.Body.Close()

	lresp := LockIPResponse{}
	json.NewDecoder(resp.Body).Decode(&lresp)
	keys := lresp.Data

	// 不需要处理，结束
	if len(keys) == 0 {
		MakeTalkEnd(sender, "没有需要解封的IP，本次服务结束")
		return
	}

	// 让客服选择要处理的IP
	query := "请选择，回复序号：\n"
	for i, key := range lresp.Data {
		query += fmt.Sprintf("%d. %s\n", i, key)
	}
	sender <- query
	var index int

	for {
		answer := <-reply
		i, err := strconv.Atoi(answer)
		if err != nil || i >= len(keys) {
			sender <- "选择错误，请重新选择!"
		} else {
			index = i
			break
		}
	}

	// 根据客服的选项，解封IP
	{
		req, err := http.NewRequest("DELETE", lockIP+"&key="+url.QueryEscape(keys[index]), nil)
		if err != nil {
			logger.Error().Err(err).Msg("remove lock ip")
		}
		http.DefaultClient.Do(req)
	}

	// 结束
	MakeTalkEnd(sender, fmt.Sprintf("解除 %s 限制成功，本次服务结束", keys[index]))
}
