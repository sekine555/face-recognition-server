package logger

import (
	"face-recognition/config"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log *zap.Logger
)

type envVals struct {
	filePath string
	stdout   bool
	level    zapcore.Level
}

// ログに関する環境変数を設定
func getEnv() (*envVals, error) {
	res := envVals{}
	// ログ出力先
	res.filePath = config.Config.LoggerFilePath
	// 標準出力する
	res.stdout = true
	// ログレベル
	level := config.Config.LoggerLevel
	switch level {
	case "debug":
		res.level = zapcore.DebugLevel
	case "error":
		res.level = zapcore.ErrorLevel
	default:
		res.level = zapcore.InfoLevel
	}
	return &res, nil
}

// 初期処理
func init() {
	envVals, err := getEnv()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// ログ出力先の設定
	var outputPaths []string
	if envVals.filePath != "" {
		outputPaths = append(outputPaths, envVals.filePath)
	}
	// 標準出力可否
	if envVals.stdout {
		outputPaths = append(outputPaths, "stdout")
	}
	// ログの設定
	logConfig := zap.Config{
		OutputPaths: outputPaths,
		Level: zap.NewAtomicLevelAt(envVals.level),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "msg",
			CallerKey:    "caller",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	if Log, err = logConfig.Build(); err != nil {
		panic(err)
	}
}