package log

import (
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	loggerOpts = []zap.Option{zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel)}
	encoder    = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     RFC3339NanoEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	output   zapcore.WriteSyncer
	logLevel zapcore.Level
	log      *zap.Logger
	logMu    sync.Mutex
)

var PotentialInconsistentData = zap.Bool("potential_inconsistent_data", true)

func init() {
	defaultLog()
}

func createLogger(o zapcore.WriteSyncer, lvl zapcore.Level) {
	output = o
	core := zapcore.NewCore(encoder, output, lvl)
	log = zap.New(core).WithOptions(loggerOpts...)
}

func defaultLog() {
	logMu.Lock()
	defer logMu.Unlock()

	output = os.Stderr
	logLevel = zap.InfoLevel
	createLogger(output, logLevel)
}

func SetLogLevel(lvl string) {
	logMu.Lock()
	defer logMu.Unlock()

	logLevel = zap.InfoLevel

	switch lvl {
	case "debug":
		logLevel = zap.DebugLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "error":
		logLevel = zap.ErrorLevel
	}

	createLogger(output, logLevel)
}

// Debug add log entry with or without fields to debug level
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

// Info add log entry with or without fields to info level
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Warn add log entry with or without fields to warn level
func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

// Error add log entry with or without fields to error level
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// Fatal add log entry with or without fields to fatal level
func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

// Panic add log entry with or without fields to panic level
func Panic(msg string, fields ...zap.Field) {
	log.Panic(msg, fields...)
}

type Logger struct {
	log *zap.Logger
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.log.Info(msg, zap.Any("data", keysAndValues))
}
func (l *Logger) Error(err error, msg string, keysAndValues ...interface{}) {
	l.log.Error(msg, zap.Error(err), zap.Any("data", keysAndValues))
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log.Error(fmt.Sprintf(format, args...))
}

func GetLogger() *Logger {
	return &Logger{
		log: log,
	}
}
