package logger

import (
	"os"
	"sync"

	"go.opentelemetry.io/contrib/bridges/otelzap"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerInterface interface {
	Info(message string, fields ...zap.Field)
	Fatal(message string, fields ...zap.Field)
	Debug(message string, fields ...zap.Field)
	Error(message string, fields ...zap.Field)
	Warn(message string, fields ...zap.Field)
}

type Logger struct {
	Log *zap.Logger
}

var once sync.Once
var instance LoggerInterface

func NewLogger(service string, loggerProvider *sdklog.LoggerProvider) (LoggerInterface, error) {
	var setupErr error

	once.Do(func() {
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		stdoutCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		)

		otelCore := otelzap.NewCore(
			service,
			otelzap.WithLoggerProvider(loggerProvider),
		)

		core := zapcore.NewTee(stdoutCore, otelCore)

		logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		instance = &Logger{Log: logger}
	})

	return instance, setupErr
}

func (l *Logger) Info(message string, fields ...zap.Field)  { l.Log.Info(message, fields...) }
func (l *Logger) Fatal(message string, fields ...zap.Field) { l.Log.Fatal(message, fields...) }
func (l *Logger) Debug(message string, fields ...zap.Field) { l.Log.Debug(message, fields...) }
func (l *Logger) Error(message string, fields ...zap.Field) { l.Log.Error(message, fields...) }
func (l *Logger) Warn(message string, fields ...zap.Field)  { l.Log.Warn(message, fields...) }
