package commands

import (
	"bot/pkg/config"
	"bot/pkg/http_client"
	"bot/pkg/talks"
	"bot/pkg/targets"
	"bot/pkg/tasks"
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

type restartJob struct {
	checkClient   http_client.HttpClientInterface
	commandClient http_client.HttpClientInterface
	checkInterval time.Duration
	check         string
	command       string
	result        int
}

func (r *restartJob) Execute(ctx context.Context) int {
	return restart(ctx, r.check, r.command, r)
}

func (r *restartJob) SetResult(result int) {
	r.result = result
}

func (r *restartJob) GetResult() int {
	return r.result
}

func restart(ctx context.Context, check string, command string, ros *restartJob) int {
	logger := ctx.Value(talks.LoggerKey).(zerolog.Logger)
	// 调用重启链接
	_, err := http_client.Get(ctx, ros.commandClient, command)
	if err != nil {
		// 调用失败，打印错误
		logger.Error().Err(err).Msg("")
		return 1
	}
	logger.Info().Str("request", command).Msg("run command")

	// 不需要检查
	if check == "" {
		return 0
	}

	// 服务已经重启了，这里会每隔一段时间检查一次，直到检查链接响应
	// 最多循环100次
	maxTry := 100
	for i := 1; i <= maxTry; i++ {
		if i == maxTry {
			// 超过了5分钟，没有检查到服务响应，则认为服务重启失败
			logger.Warn().Str("request", check).Msg("check failed")
			return 2
		}

		_, err := http_client.Get(ctx, ros.checkClient, check)
		if err != nil {
			time.Sleep(ros.checkInterval)
			continue
		} else {
			// 检查到服务响应，则认为服务重启成功
			logger.Info().Str("request", check).Msg("check success")
			return 0
		}
	}

	return 0
}

// 重启进程/服务，需要预先配置检查链接和重启链接
// 服务都是一个集群的，不应该同时重启，会被客户感知。
// 需要按顺序重启多个服务，中间如果有错误，则会跳过错误的服务
func RestartTarget(ctx context.Context) {
	targetName := ctx.Value(talks.TargetNameKey).(string)
	target := ctx.Value(talks.TargetKey).([]config.Target)
	task := ctx.Value(talks.TaskKey).(*config.Task)
	talks.ReplyMsg(ctx, fmt.Sprintf("开始重启%s，耐心等待", targetName))

	start := time.Now()
	client := http_client.NewDumbHttpClient(5 * time.Second)
	targetTask := targets.ListTargetTask(target, task, &config.C.Variables)

	jobs := []tasks.Job{}
	for _, item := range targetTask {
		ros := &restartJob{
			checkClient:   client,
			commandClient: client,
			checkInterval: 2 * time.Second,
			check:         item.Check,
			command:       item.Command,
		}
		jobs = append(jobs, ros)
	}
	if tasks.NewParallelTasks(ctx, jobs) {
		elapsed := time.Since(start)
		talks.MakeTalkEnd(ctx, fmt.Sprintf("重启%s完成，耗时：%v，本次服务结束", targetName, elapsed))
	} else {
		elapsed := time.Since(start)
		talks.MakeTalkEnd(ctx, fmt.Sprintf("重启%s失败，耗时: %v, 请联系管理员，本次服务结束", targetName, elapsed))
	}

}

func init() {
	registerTaskCommand("Restart", RestartTarget)
}
