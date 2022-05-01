package health

import (
	"bot/pkg/config"
	"bot/pkg/http_client"
	"bot/pkg/logging"
	"bot/pkg/task"
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
func onceCheckHealth(h *healthUnit, command string, check string, wg *sync.WaitGroup) (int, error) {
	var err error
	bT := time.Now()
	defer func() {
		eT := time.Now()
		h.log.Debug().TimeDiff("cost time", eT, bT).Msg("check health done and wait group -1")
		wg.Done()
	}()

	// 检查服务是否超时
	_, err = http_client.Get(h.checkClient, check)
	if err == nil {
		// 服务正常
		return 0, nil
	}
	if !os.IsTimeout(err) {
		// 服务不正常，但是不是响应超时问题
		h.log.Error().Err(err).Msg("check health")
		return 1, errors.New("check health failed")
	}

	h.log.Info().Str("check", check).Msg("get timeout")

	// 如果请求超时了，就重启服务
	_, err = http_client.Get(h.commandClient, command)
	if err != nil {
		// 重启服务失败
		h.log.Error().Err(err).Msg("restart")
		return 2, errors.New("restart failed")
	}
	h.log.Info().Str("service", command).Msg("restart")
	return 3, nil
}

func eachCheckHealths(h *healthUnit) {

	for _, b := range config.C.Beat {
		targetC, err := task.GetTarget(b.TargetName)
		if err != nil {
			h.log.Error().Err(err).Msg("get target")
		}
		taskC, err := task.GetTask(b.TaskName)
		if err != nil {
			h.log.Error().Err(err).Msg("get task")
		}
		targetTask := task.ListTargetTask(targetC, taskC, &config.C.Variables)
		healthItemsLen := len(targetTask)
		wg := &sync.WaitGroup{}
		wg.Add(healthItemsLen)
		h.log.Debug().Int("count", healthItemsLen).Msg("add wait group")

		for _, item := range targetTask {
			go onceCheckHealth(h, item.Check, item.Command, wg)
		}

		wg.Wait()
	}

}

// 定时检查服务的情况，当服务停止响应/超时时，都会重启对应的服务进程
func BeatCheckHealth() {

	log := logging.Log.With().
		Caller().
		Str("module", "beat").
		Logger()

	h := &healthUnit{
		checkClient:   http_client.NewDumbHttpClient(5),
		commandClient: http_client.NewDumbHttpClient(5),
		log:           log,
	}
	h.log.Info().Msg("start check health")

	for {
		eachCheckHealths(h)
		time.Sleep(3 * time.Second)
	}

}
