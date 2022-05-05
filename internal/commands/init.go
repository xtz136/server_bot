package commands

import (
	"bot/pkg/config"

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

// 保存任务别名和任务执行函数的关系
var TaskCommands = map[string]func(Context){}

func registerTaskCommand(name string, f func(Context)) {
	TaskCommands[name] = f
}
