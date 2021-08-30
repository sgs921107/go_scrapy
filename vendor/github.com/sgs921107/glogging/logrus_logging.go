package glogging

import (
	"sync"
	"strings"

	"github.com/sirupsen/logrus"
)

// 起别名
type (
	// LogrusFields logrus fields
	LogrusFields = logrus.Fields
	// LogrusLogger logrus logger
	LogrusLogger = logrus.Logger
)

// LogrusLogging logging
type LogrusLogging interface{
	BaseLogging
	GetLogger()	*LogrusLogger
}

// Log log
type logrusLog struct {
	baseLog
	logger	*LogrusLogger
	once	sync.Once
}

// getLogger 获取一个logger
func (l *logrusLog) GetLogger() *logrus.Logger {
	l.once.Do(func() {
		l.logger = logrus.New()
		l.initLogger()
	})
	return l.logger
}

// setLoggerLevel set log level
func (l *logrusLog) setLoggerLevel() {
	switch strings.ToUpper(l.options.Level) {
	case "DEBUG":
		l.logger.SetLevel(logrus.DebugLevel)
	case "INFO":
		l.logger.SetLevel(logrus.InfoLevel)
	case "WARNING":
		l.logger.SetLevel(logrus.WarnLevel)
	case "ERROR":
		l.logger.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		l.logger.SetLevel(logrus.FatalLevel)
	case "PANIC":
		l.logger.SetLevel(logrus.PanicLevel)
	case "TRACE":
		l.logger.SetLevel(logrus.TraceLevel)
	default:
		l.logger.SetLevel(logrus.DebugLevel)
	}
}

// setFormater set log level
func (l *logrusLog) setFormatter() {
	switch strings.ToUpper(l.options.Formatter) {
	case "TEXT":
		l.logger.SetFormatter(&logrus.TextFormatter{})
	default:
		l.logger.SetFormatter(&logrus.JSONFormatter{})
	}
}

// 对logger进行配置
func (l *logrusLog) initLogger() {
	// 配置日志等级
	l.setLoggerLevel()
	// 配置日志格式
	l.setFormatter()
	// set output
	l.logger.SetOutput(l.Output())
	if l.options.NoLock {
		l.logger.SetNoLock()
	}
	l.logger.SetReportCaller(true)
}

// NewLogrusLogging 生成一个logging
func NewLogrusLogging(options Options) LogrusLogging {
	return &logrusLog{
		baseLog: baseLog{
			options: options,
		},
	}
}
