package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Instance *zap.Logger

type Config struct {
	Level string `json:"level"`
}

func Initialization(c *Config) {
	level := zap.InfoLevel
	if c.Level != "" {
		l, err := zapcore.ParseLevel(c.Level)
		if err != nil {
			panic(err)
		}
		level = l
	}
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	Instance = logger
	Instance.Debug("log initialization complete")
}
