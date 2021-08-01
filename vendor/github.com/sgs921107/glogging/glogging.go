package glogging

import (
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

// 起别名
type (
	// Fields logrus fields
	Fields = logrus.Fields
	// Logger logrus logger
	Logger = logrus.Logger
)

// Logging logging
type Logging interface{
	GetLogger()	*logrus.Logger
}

// Log log
type log struct {
	options	*Options
	logger			*logrus.Logger
}

// getLogger 获取一个logger
func (l *log) GetLogger() *logrus.Logger {
	if l.logger == nil {
		l.logger = logrus.New()
		l.initLogger(l.logger)
	}
	return l.logger
}

// setLoggerLevel set log level
func (l *log) setLoggerLevel(logger *logrus.Logger) {
	switch strings.ToUpper(l.options.Level) {
	case "DEBUG":
		logger.SetLevel(logrus.DebugLevel)
	case "INFO":
		logger.SetLevel(logrus.InfoLevel)
	case "WARNING":
		logger.SetLevel(logrus.WarnLevel)
	case "ERROR":
		logger.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		logger.SetLevel(logrus.FatalLevel)
	case "PANIC":
		logger.SetLevel(logrus.PanicLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
}

// setFormater set log level
func (l *log) setFormater(logger *logrus.Logger) {
	switch strings.ToUpper(l.options.Formater) {
	case "TEXT":
		logger.SetFormatter(&logrus.TextFormatter{})
	default:
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
}

// cratePattern create filename pattern
func (l *log) createPattern() string {
	var p string
	duration := l.options.RotationTime
	switch {
	case duration < time.Hour:
		p = ".%Y%m%d%H%M"
	case duration < time.Hour * 24:
		p = ".%Y%m%d%H"
	case duration >= time.Hour * 24:
		p = ".%Y%m%d"
	default:
		p = ".%Y%m%d%H"
	}
	return p
}

// setOutPut set log output
func (l *log) setOutPut(logger *logrus.Logger) {
	filePath := l.options.FilePath
	if filePath == "" {
		logger.SetOutput(os.Stdout)
		return
	}
	if l.options.RotationMaxAge != 0 && l.options.RotationMaxAge < l.options.RotationTime {
		l.options.RotationMaxAge = l.options.RotationTime
	}
	writer, err := rotatelogs.New(
		filePath+l.createPattern(),
		rotatelogs.WithLinkName(filePath),
		rotatelogs.WithRotationTime(l.options.RotationTime),
		rotatelogs.WithMaxAge(l.options.RotationMaxAge),
		// rotatelogs.ForceNewFile(),
	)
	if err != nil {
		logger.Panicf("set log file failed: %s", filePath)
	} else {
		logger.SetOutput(writer)
	}
}

// 对logger进行配置
func (l *log) initLogger(logger *Logger) {
	// 配置日志等级
	l.setLoggerLevel(logger)
	// 配置日志格式
	l.setFormater(logger)
	l.setOutPut(logger)
	if l.options.NoLock {
		logger.SetNoLock()
	}
	logger.SetReportCaller(true)
}

// NewLogging 生成一个logging
func NewLogging(options *Options) Logging {
	var logging Logging
	logging = &log{
		options: options,
	}
	return logging
}
