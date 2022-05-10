package beat

import (
	"bot/internal/commands"
	"bot/pkg/config"
	"bot/pkg/logging"
	"bot/pkg/talks"
	"bot/pkg/tasks"
	"context"
	"time"
)

func eachBeatTasks() {

	log := logging.Log.With().
		Str("module", "beat").
		Logger()

	for _, b := range config.C.Beat {
		target, err := tasks.GetTarget(b.TargetName)
		if err != nil {
			log.Error().Err(err).Msg("get target")
		}
		task, err := tasks.GetTask(b.TaskName)
		if err != nil {
			log.Error().Err(err).Msg("get task")
		}

		ctx := context.Background()
		ctx = context.WithValue(ctx, talks.LoggerKey, log)
		ctx = context.WithValue(ctx, talks.TargetKey, target)
		ctx = context.WithValue(ctx, talks.TargetNameKey, b.TargetName)
		ctx = context.WithValue(ctx, talks.TaskKey, task)
		ctx, cancel := context.WithTimeout(ctx, 5*60*time.Second)

		go func(ctx context.Context, cancel context.CancelFunc) {
			defer cancel()
			commands.TaskCommands[task.Name](ctx)
		}(ctx, cancel)
	}

}

// 定期执行任务
func StartBeat() {

	for {
		eachBeatTasks()
		time.Sleep(3 * time.Second)
	}

}
