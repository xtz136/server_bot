package logging

import (
	"io"
	"os"
	"path"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log zerolog.Logger

func GetLog() zerolog.Logger {
	writers := []io.Writer{}

	if viper.GetBool("log.enable_console") {
		writers = append(
			writers,
			zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339},
		)
	}

	if viper.GetBool("log.enable_file") {
		dir := viper.GetString("log.log_file_dir")
		filename := viper.GetString("log.log_file_name")
		writers = append(
			writers,
			rollingFile(dir, filename),
		)
	}

	switch viper.GetInt("log.level") {
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

func InitLog() {
	Log = GetLog()
}
