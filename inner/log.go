package inner

import (
	stdlog "log"

	"github.com/huskyrobotdog/toolbox-go/log"
	"go.uber.org/zap"
)

func Debug(msg string) {
	if log.Instance != nil {
		log.Instance.WithOptions(zap.AddCallerSkip(1)).Debug(msg)
	} else {
		stdlog.Print(msg)
	}
}

func Info(msg string) {
	if log.Instance != nil {
		log.Instance.WithOptions(zap.AddCallerSkip(1)).Info(msg)
	} else {
		stdlog.Print(msg)
	}
}

func Warn(msg string) {
	if log.Instance != nil {
		log.Instance.WithOptions(zap.AddCallerSkip(1)).Warn(msg)
	} else {
		stdlog.Print(msg)
	}
}

func Error(msg string) {
	if log.Instance != nil {
		log.Instance.WithOptions(zap.AddCallerSkip(1)).Error(msg)
	} else {
		stdlog.Print(msg)
	}
}

func Fatal(msg string) {
	if log.Instance != nil {
		log.Instance.WithOptions(zap.AddCallerSkip(1)).Fatal(msg)
	} else {
		stdlog.Fatal(msg)
	}
}

func Panic(msg string) {
	if log.Instance != nil {
		log.Instance.WithOptions(zap.AddCallerSkip(1)).Panic(msg)
	} else {
		stdlog.Panic(msg)
	}
}
