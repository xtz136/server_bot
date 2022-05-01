package commands

import (
	"bot/pkg/config"
	"time"

	"github.com/rs/zerolog"
)

type State map[string]interface{}

type Context = struct {
	TargetName string
	// 通道，发送消息给用户
	Sender chan string
	// 通道，接收用户回复，并传给具体的任务工作者
	Reply chan string
	// 目标
	Target *config.Target
	// 任务
	Task *config.Task
	// 日志对象
	Log zerolog.Logger
	// 在 workflow 中的状态，可以在每个 step 中访问
	State State
	// 结束会话的函数
	MakeTalkEnd func(chan string, string)
}

type UnlockIPResponse struct {
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
	ctx       Context
	funcs     []stepFunc
	startTime time.Time
	costTime  time.Duration
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

// 保存任务别名和任务执行函数的关系
var TaskCommands = map[string]func(Context){}

func registerTaskCommand(name string, f func(Context)) {
	TaskCommands[name] = f
}
