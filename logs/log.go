package logger

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger = NewLogger()

var logDir = "./logs/logs.txt"

func NewLogger() *zap.SugaredLogger {

	if _, err := os.Stat(logDir); errors.Is(err, os.ErrNotExist) {
		path := strings.Split(logDir, "/")
		os.MkdirAll(filepath.Join(path[:len(path)-1]...), os.ModePerm)
		// os.Create(logDir)
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
