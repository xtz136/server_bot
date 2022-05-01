package logging

import (
	"bot/pkg/config"
	"io"
	"os"
	"path"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log zerolog.Logger

func GetLog() zerolog.Logger {
	writers := []io.Writer{}

	if config.C.Log.EnableConsole {
		writers = append(
			writers,
			zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339},
		)
	}

	if config.C.Log.EnableFile {
		dir := config.C.Log.LogFileDir
		filename := config.C.Log.LogFileName
		writers = append(
			writers,
			rollingFile(dir, filename),
		)
	}

	switch config.C.Log.Level {
	case 0:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	multi := zerolog.MultiLevelWriter(writers...)
	log := zerolog.New(multi).With().Timestamp().Logger()
	return log
}

func rollingFile(dir string, filename string) io.Writer {
	if err := os.MkdirAll(dir, 0744); err != nil {
		log.Error().Err(err).Str("path", dir).Msg("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename: path.Join(dir, filename),
		MaxAge:   30,
	}
}

func init() {
	Log = GetLog()
}
