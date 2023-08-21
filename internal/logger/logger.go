package logger

import (
	"context"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *otelzap.Logger

func Init(debug bool) Syncer {
	var err error
	config := zap.NewProductionConfig()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = "" // to hide stacktrace info
	config.EncoderConfig = encoderConfig
	if debug {
		config.Level.SetLevel(zap.DebugLevel)
	}

	l, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	Log = otelzap.New(l, otelzap.WithTraceIDField(true), otelzap.WithMinLevel(zap.DebugLevel))
	log = Log.Sugar()
	return log
}

type Logger interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

type Syncer interface {
	Sync() error
}

func FromContext(ctx context.Context) Logger {
	return log.Ctx(ctx)
}

func Infof(message string, fields ...interface{}) {
	log.Infof(message, fields...)
}

func Debugf(message string, fields ...interface{}) {
	log.Debugf(message, fields...)
}

func Errorf(message string, fields ...interface{}) {
	log.Errorf(message, fields...)
}

func Fatalf(message string, fields ...interface{}) {
	log.Fatalf(message, fields...)
}

var log *otelzap.SugaredLogger

func init() {
	Init(false)
}
