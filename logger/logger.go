package logger

import (
	"io"
	"os"
	"time"

	"github.com/mahesh-yadav/go-recipes-api/config"
	"github.com/rs/zerolog"
)

func SetupLogger(config *config.Config) zerolog.Logger {
	// lumberjackLogger := &lumberjack.Logger{
	// 	Filename: config.LogFile,
	// 	MaxSize:  config.LogMaxSizeInMB, // megabytes
	// 	MaxAge:   config.LogMaxAge,
	// 	Compress: config.LogCompress,
	// }

	multiWriter := io.MultiWriter(os.Stdout)
	zerolog.TimeFieldFormat = "2006/01/02 15:04:05"
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}

	logger := zerolog.New(multiWriter).
		Level(config.GetLogLevel()).
		With().
		Str("service_name", "recipes-api"). // Global context
		Int("pid", os.Getpid()).
		Timestamp().
		Caller().
		Logger()

	return logger
}
