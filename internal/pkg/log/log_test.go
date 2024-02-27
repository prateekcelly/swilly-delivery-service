package log

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var buf *testWriteSyncer

type testWriteSyncer struct {
	*bytes.Buffer
}

func (tws *testWriteSyncer) Sync() error {
	return nil
}

func TestMain(m *testing.M) {
	buf = &testWriteSyncer{
		new(bytes.Buffer),
	}
	createLogger(buf, logLevel)

	os.Exit(m.Run())
}

const (
	debugMessage  = "debug message"
	infoMessage   = "info message"
	warnMessage   = "warn message"
	errorMessage  = "error message"
	fatalMessage  = "fatal message"
	panicMessage  = "panic message"
	notNilMessage = "Should not nil"
	nilMessage    = "Should nil"
)

func TestSetLogLevelInfo(t *testing.T) {
	assert.Nil(t, log.Check(zap.DebugLevel, debugMessage), nilMessage)
	assert.NotNil(t, log.Check(zap.InfoLevel, infoMessage), notNilMessage)
	assert.NotNil(t, log.Check(zap.WarnLevel, warnMessage), notNilMessage)
	assert.NotNil(t, log.Check(zap.ErrorLevel, errorMessage), notNilMessage)
	assert.NotNil(t, log.Check(zap.FatalLevel, fatalMessage), notNilMessage)
	assert.NotNil(t, log.Check(zap.PanicLevel, panicMessage), notNilMessage)
}

func TestSetLogLevelDebug(t *testing.T) {
	SetLogLevel("debug")
	assert.NotNil(t, log.Check(zap.DebugLevel, debugMessage), notNilMessage)
	assert.NotNil(t, log.Check(zap.InfoLevel, infoMessage), notNilMessage)
	assert.NotNil(t, log.Check(zap.WarnLevel, warnMessage), notNilMessage)
	assert.NotNil(t, log.Check(zap.ErrorLevel, errorMessage), notNilMessage)
	assert.NotNil(t, log.Check(zap.FatalLevel, fatalMessage), notNilMessage)
	assert.NotNil(t, log.Check(zap.PanicLevel, panicMessage), notNilMessage)
}

func TestSetLogLevelWarn(t *testing.T) {
	SetLogLevel("warn")
	assert.Nil(t, log.Check(zap.DebugLevel, debugMessage), nilMessage)
	assert.Nil(t, log.Check(zap.InfoLevel, infoMessage), nilMessage)
	assert.NotNil(t, log.Check(zap.WarnLevel, warnMessage), notNilMessage)
	assert.NotNil(t, log.Check(zap.ErrorLevel, errorMessage), notNilMessage)
	assert.NotNil(t, log.Check(zap.FatalLevel, fatalMessage), notNilMessage)
	assert.NotNil(t, log.Check(zap.PanicLevel, panicMessage), notNilMessage)
}

func TestSetLogLevelError(t *testing.T) {
	SetLogLevel("error")
	assert.Nil(t, log.Check(zap.DebugLevel, debugMessage), nilMessage)
	assert.Nil(t, log.Check(zap.InfoLevel, infoMessage), nilMessage)
	assert.Nil(t, log.Check(zap.WarnLevel, warnMessage), notNilMessage)
	assert.NotNil(t, log.Check(zap.ErrorLevel, errorMessage), notNilMessage)
	assert.NotNil(t, log.Check(zap.FatalLevel, fatalMessage), notNilMessage)
	assert.NotNil(t, log.Check(zap.PanicLevel, panicMessage), notNilMessage)
}

func TestSetLogLevelFatal(t *testing.T) {
	SetLogLevel("fatal")
	assert.Nil(t, log.Check(zap.DebugLevel, debugMessage), nilMessage)
	assert.NotNil(t, log.Check(zap.InfoLevel, infoMessage), nilMessage)
	assert.NotNil(t, log.Check(zap.WarnLevel, warnMessage), nilMessage)
	assert.NotNil(t, log.Check(zap.ErrorLevel, errorMessage), nilMessage)
	assert.NotNil(t, log.Check(zap.FatalLevel, fatalMessage), notNilMessage)
	assert.NotNil(t, log.Check(zap.PanicLevel, panicMessage), notNilMessage)
}

func TestSetLogLevelPanic(t *testing.T) {
	SetLogLevel("panic")
	assert.Nil(t, log.Check(zap.DebugLevel, debugMessage), nilMessage)
	assert.NotNil(t, log.Check(zap.InfoLevel, infoMessage), nilMessage)
	assert.NotNil(t, log.Check(zap.WarnLevel, warnMessage), notNilMessage)
	assert.NotNil(t, log.Check(zap.ErrorLevel, errorMessage), notNilMessage)
	assert.NotNil(t, log.Check(zap.FatalLevel, fatalMessage), notNilMessage)
	assert.NotNil(t, log.Check(zap.PanicLevel, panicMessage), notNilMessage)
}

func TestLogging(t *testing.T) {
	const panicName = "Panic"
	buf.Reset()
	SetLogLevel("debug")

	type tc struct {
		funcToTest func(msg string, fields ...zap.Field)
		level      zapcore.Level
	}

	testCases := map[string]tc{
		"Debug":   {funcToTest: Debug, level: zap.DebugLevel},
		"Info":    {funcToTest: Info, level: zap.InfoLevel},
		"Warn":    {funcToTest: Warn, level: zap.WarnLevel},
		"Error":   {funcToTest: Error, level: zap.ErrorLevel},
		panicName: {funcToTest: Panic, level: zap.PanicLevel},
	}
	msg := "test message"

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			if name == panicName {
				assert.True(t, assert.Panics(t, func() {
					testCase.funcToTest(msg)
				}), "Panic should return true")
			} else {
				testCase.funcToTest(msg)
			}

			logMap := make(map[string]interface{})
			err := json.NewDecoder(buf).Decode(&logMap)
			assert.Nil(t, err, "Error should nil while decoding log to json")

			assert.Equal(t, msg, logMap["message"], "Should be equal")
			assert.Equal(t, testCase.level.String(), logMap["level"], "Should be equal")
		})
	}
}
