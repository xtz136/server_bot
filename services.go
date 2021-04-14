package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

var ctx = context.Background()

type System struct {
	Name           string `json:"name"`
	RedisURL       string `json:"redis_url"`
	LockIP         string `json:"lock_ip"`
	ListUserSystem string `json:"list_user_system"`
	MakeToken      string `json:"make_token"`
	Restart        []struct {
		Command string `json:"command"`
		Check   string `json:"check"`
	} `json:"restart"`
}

type LockIPResponse struct {
	Status  bool     `json:"status"`
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
	ReqID   string   `json:"req_id"`
}

func makeTalkEnd(sender chan string, lastMsg string) {
	if lastMsg != "" {
		sender <- lastMsg
	}

	time.Sleep(time.Duration(1) * time.Second)
	sender <- ""
}

func dummy(systemName string, sender chan string, reply chan string, logger zerolog.Logger) {
	bT := time.Now()
	time.Sleep(time.Duration(3) * time.Second)
	eT := time.Since(bT)

	makeTalkEnd(sender, fmt.Sprintf("测试 %s 完成，耗时：%v，本次服务结束", systemName, eT))
}

func restart(systemName string, sender chan string, reply chan string, logger zerolog.Logger) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	sender <- "开始重启阿里云，耐心等待"

	b := System{}
	viper.UnmarshalKey("systems."+systemName, &b)

	if len(b.Restart) == 0 {
		logger.Info().Msg("system name is invalid")
		return
	}

	bT := time.Now()
	client := http.Client{
		Timeout: 1 * time.Second,
	}

	raiseError := func() {
		eT := time.Since(bT)
		makeTalkEnd(sender, fmt.Sprintf("汪，重启%s失败，耗时: %v, 请联系管理员，本次服务结束", b.Name, eT))
		return
	}

	for _, item := range b.Restart {
		command := item.Command
		check := item.Check

		_, err := client.Get(command)
		if err != nil {
			logger.Error().Err(err).Msg("")
			raiseError()
			return
		}
		logger.Info().Str("request", command).Msg("run command")

		// 循环100次，每次3秒，一共5分钟
		has_error := true
		for i := 0; i < 100; i++ {
			_, err := client.Get(check)
			if err != nil {
				time.Sleep(time.Duration(3) * time.Second)
				continue
			}
			has_error = false
			break
		}

		if has_error {
			logger.Warn().Msg("check failed")
			raiseError()
			return
		}
	}

	eT := time.Since(bT)
	makeTalkEnd(sender, fmt.Sprintf("汪，重启%s完成，耗时：%v，本次服务结束", b.Name, eT))
}

func enableUser(systemName string, sender chan string, reply chan string, logger zerolog.Logger) {
	// sender <- "请输入用户账号"
	// account := <-reply
	// fmt.Printf("enable %s user\n", account)
	// sender <- fmt.Sprintf("用户%s已启用", account)
	// sender <- ""
	makeTalkEnd(sender, "这个功能还在建设中")
}

func allowIP(systemName string, sender chan string, reply chan string, logger zerolog.Logger) {
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
		makeTalkEnd(sender, "没有需要解封的IP，本次服务结束")
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
			// rdb.Del(ctx, "cache:user:validerr:").Result()
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
	makeTalkEnd(sender, fmt.Sprintf("解除 %s 限制成功，本次服务结束", keys[index]))
}

type UserSystem struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		UserAlias           string `json:"user_alias"`
		SystemID            string `json:"system_id"`
		SystemAlias         string `json:"system_alias"`
		DanweiID            string `json:"danwei_id"`
		DanweiAlias         string `json:"danwei_alias"`
		DanweiParentAlias   string `json:"danwei_parent_alias"`
		QuanxianID          string `json:"quanxian_id"`
		QuanxianAlias       string `json:"quanxian_alias"`
		QuanxianParentAlias string `json:"quanxian_parent_alias"`
	} `json:"data"`
	ReqID string `json:"req_id"`
}

type Token struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
	ReqID   string `json:"req_id"`
}

func loginUser(systemName string, sender chan string, reply chan string, logger zerolog.Logger) {
	tokenList := viper.GetString("systems." + systemName + ".list_user_system")
	tokenMake := viper.GetString("systems." + systemName + ".make_token")

	sender <- "请回复登录账号，同时回复申请原因(用空格隔开)。例如：\n\nxiaoming 火箭升空演示，需要测试推进系统"
	accountAndWhy := strings.SplitN(<-reply, " ", 2)
	account := accountAndWhy[0]
	why := accountAndWhy[1]

	if account == "" || why == "" {
		makeTalkEnd(sender, "回复有误，本次服务结束")
		return
	}

	// 找到这个用户绑定的所有系统
	resp, err := http.Get(tokenList + "&account=" + account)
	if err != nil {
		logger.Error().Err(err).Str("url", tokenList).Msg("list user system")
		makeTalkEnd(sender, "")
		return
	}
	defer resp.Body.Close()

	// 让客服选一个
	userSystem := UserSystem{}
	json.NewDecoder(resp.Body).Decode(&userSystem)

	// 不需要处理，结束
	if len(userSystem.Data) == 0 {
		makeTalkEnd(sender, "这个用户没有绑定系统，本次服务结束")
		return
	}

	// 让客服选择要处理的IP
	query := "请选择，回复序号，格式[用户名][系统名][部门名]：\n\n"
	for i, item := range userSystem.Data {
		query += fmt.Sprintf("%d. [%s][%s][%s]\n", i, item.UserAlias, item.SystemAlias, item.DanweiParentAlias)
	}
	sender <- query
	var index int

	for {
		answer := <-reply
		i, err := strconv.Atoi(answer)
		if err != nil || i >= len(userSystem.Data) {
			sender <- "选择错误，请重新选择!"
		} else {
			index = i
			break
		}
	}

	// 请求token
	{
		params := url.Values{
			"account":     {account},
			"system_id":   {userSystem.Data[index].SystemID},
			"danwei_id":   {userSystem.Data[index].DanweiID},
			"quanxian_id": {userSystem.Data[index].QuanxianID},
			"why":         {why},
		}
		resp, err := http.Get(tokenMake + "&" + params.Encode())
		if err != nil {
			logger.Error().Err(err).Str("url", tokenMake).Msg("make token")
			makeTalkEnd(sender, "")
			return
		}
		token := Token{}
		json.NewDecoder(resp.Body).Decode(&token)
		makeTalkEnd(sender, fmt.Sprintf("访问下方链接(%s)即可登录,本次服务结束\n%s", token.Message, token.Data))
		defer resp.Body.Close()

		logger.Info().
			Str("account", params.Get("account")).
			Str("danwei_id", params.Get("danwei_id")).
			Str("quanxian_id", params.Get("quanxian_id")).
			Str("why", params.Get("why")).
			Msg("make token")
	}
}
