package commands

import (
	"time"

	"github.com/rs/zerolog"
)

type State map[string]interface{}

type Context = struct {
	Name   string
	Sender chan string
	Reply  chan string
	Log    zerolog.Logger
	State  State
	MakeTalkEnd func(chan string, string)
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


type stepFunc func(ctx Context) bool

type WorkFlow struct {
	ctx   Context
	funcs []stepFunc
	startTime time.Time
	costTime time.Duration
}

func (this *WorkFlow) add(f stepFunc) {
	this.funcs = append(this.funcs, f)
}

func (this *WorkFlow) start() bool {
	success := true
	for _, f := range this.funcs {
		next := f(this.ctx)
		if !next {
			success = false
			break
		}
	}
	this.costTime = time.Since(this.startTime)
	return success
}

func (this *WorkFlow) getCostTime() time.Duration {
	return this.costTime
}

func newWorkFlow(ctx Context) WorkFlow {
	w := WorkFlow{ctx: ctx}
	w.startTime = time.Now()
	return w
}