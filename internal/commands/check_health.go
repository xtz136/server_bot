package commands

import (
	"bot/pkg/config"
	"bot/pkg/http_client"
	"bot/pkg/talks"
	"bot/pkg/targets"
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type healthUnit struct {
	log           zerolog.Logger
	checkClient   http_client.HttpClientInterface
	commandClient http_client.HttpClientInterface
}

// 返回一个整形和错误
// 整形：0：正常，1：健康检查错误，2：重启错误，3：完成重启
func checkHealth(ctx context.Context, hu *healthUnit, check string, command string) (int, error) {
	var err error
	bT := time.Now()
	defer func() {
		eT := time.Now()
		hu.log.Debug().TimeDiff("cost time", eT, bT).Msg("check health done and wait group -1")
	}()

	// 检查服务是否超时
	_, err = http_client.Get(ctx, hu.checkClient, check)
	if err == nil {
		// 服务正常
		return 0, nil
	}
	if !os.IsTimeout(err) {
		// 服务不正常，但是不是响应超时问题
		hu.log.Error().Err(err).Msg("check health")
		return 1, errors.New("check health failed")
	}

	hu.log.Info().Str("check", check).Msg("get timeout")

	// 如果请求超时了，就重启服务
	_, err = http_client.Get(ctx, hu.commandClient, command)
	if err != nil {
		// 重启服务失败
		hu.log.Error().Err(err).Msg("restart")
		return 2, errors.New("restart failed")
	}
	hu.log.Info().Str("service", command).Msg("restart")
	return 3, nil
}

func CheckHealthGroup(ctx context.Context) {
	target := ctx.Value(talks.TargetKey).([]config.Target)
	task := ctx.Value(talks.TaskKey).(*config.Task)
	logger := ctx.Value(talks.LoggerKey).(zerolog.Logger)
	targetTask := targets.ListTargetTask(target, task, &config.C.Variables)
	beatTasksLen := len(targetTask)
	logger.Debug().Int("count", beatTasksLen).Msg("add wait group")

	hu := &healthUnit{
		checkClient:   http_client.NewDumbHttpClient(5 * time.Second),
		commandClient: http_client.NewDumbHttpClient(5 * time.Second),
		log:           logger,
	}
	wg := &sync.WaitGroup{}
	wg.Add(beatTasksLen)

	for _, item := range targetTask {
		go func(hu *healthUnit, item *targets.TargetTask) {
			defer wg.Done()
			checkHealth(ctx, hu, item.Check, item.Command)
		}(hu, &item)
	}

	wg.Wait()
}

func init() {
	registerTaskCommand("CheckHealth", CheckHealthGroup)
}
