package logger

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

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
	// logDir := "./dir/temp.txt"
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
