package logger

import (
	"io"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	DefaultLogger Logger
)

type Logger interface {
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
	Writer() io.Writer
}

type SugaredLogger struct {
	*zap.SugaredLogger
	writer zapcore.WriteSyncer
}

func (s *SugaredLogger) Writer() io.Writer {
	return s.writer
}

func InitDefault(logfile string, level string) {
	lv := zapcore.InfoLevel
	switch strings.ToLower(level) {
	case "debug":
		lv = zapcore.DebugLevel
	case "error":
		lv = zapcore.ErrorLevel
	case "warn":
		lv = zapcore.WarnLevel
	}
	DefaultLogger = Init(logfile, lv)
}

func Debugf(format string, v ...interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.Debugf(format, v...)
	}

}

func Infof(format string, v ...interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.Infof(format, v...)
	}

}

func Warnf(format string, v ...interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.Warnf(format, v...)
	}

}

func Errorf(format string, v ...interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.Errorf(format, v...)
	}

}

func Fatalf(format string, v ...interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.Fatalf(format, v...)
	}

}

var (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
)

// Init 初始化日志
func Init(logfile string, level zapcore.Level) *SugaredLogger {
	writer := logWriter(logfile, 500, false)
	log := zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(logEncoder()),
		// logStdOut(),
		writer,
		zapcore.Level(level),
	))
	return &SugaredLogger{
		SugaredLogger: log.Sugar(),
		writer:        writer,
	}
}

func logStdOut() zapcore.WriteSyncer {
	return zapcore.AddSync(os.Stdout)
}

func logWriter(logFile string, size int, compress bool) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:  logFile,
		MaxSize:   size,
		LocalTime: true,
		Compress:  compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func logEncoder() zapcore.EncoderConfig {
	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder.EncodeLevel = zapcore.CapitalLevelEncoder
	return encoder
}
