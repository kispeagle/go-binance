package logger

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger = NewLogger()

func NewLogger() *zap.SugaredLogger {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	logDir := os.Getenv("LOG_DIR")
	if _, err := os.Stat(logDir); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(logDir, os.ModePerm)
	}

	conf := zap.NewProductionConfig()
	conf.OutputPaths = append(conf.OutputPaths, logDir)
	conf.EncoderConfig.TimeKey = "timestamp"
	conf.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	if os.Getenv("LOG_DEBUG") == "1" {
		conf.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	logger, err := conf.Build()
	if err != nil {
		log.Fatal(err)
	}

	return logger.Sugar()
}
