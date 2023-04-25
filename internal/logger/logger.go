package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLog *zap.SugaredLogger

func init() {
	Init(false)
}

func Init(debug bool) *zap.SugaredLogger {
	var err error
	config := zap.NewProductionConfig()
	enccoderConfig := zap.NewProductionEncoderConfig()
	zapcore.TimeEncoderOfLayout("Jan _2 15:04:05.000000000")
	enccoderConfig.StacktraceKey = "" // to hide stacktrace info
	config.EncoderConfig = enccoderConfig
	if debug {
		config.Level.SetLevel(zap.DebugLevel)
	}

	l, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	zapLog = l.Sugar()
	return zapLog
}

func Info(message string, fields ...interface{}) {
	zapLog.Infof(message, fields...)
}

func Debug(message string, fields ...interface{}) {
	zapLog.Debugf(message, fields...)
}

func Error(message string, fields ...interface{}) {
	zapLog.Errorf(message, fields...)
}

func Fatal(message string, fields ...interface{}) {
	zapLog.Fatalf(message, fields...)
}
