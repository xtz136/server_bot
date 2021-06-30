package services

import (
	"time"

	"github.com/rs/zerolog"
)

type Task = struct {
	Name   string
	Sender chan string
	Reply  chan string
	Log    zerolog.Logger
}

type System struct {
	Name           string `json:"name"`
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

func MakeTalkEnd(sender chan string, lastMsg string) {
	if lastMsg != "" {
		sender <- lastMsg
	}

	time.Sleep(time.Duration(1) * time.Second)
	close(sender)
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
