package components

import (
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"io"
)

type LogrusWrapper struct {
	logger *logrus.Logger
}

func NewLogrusWrapper(logger *logrus.Logger) *LogrusWrapper {
	return &LogrusWrapper{
		logger: logger,
	}
}

func (l *LogrusWrapper) Output() io.Writer {
	return logrus.StandardLogger().Out
}
func (l *LogrusWrapper) SetOutput(w io.Writer) {
	l.logger.Error("LogrusWrapper::SetOutput Not implemented")
}
func (l *LogrusWrapper) Prefix() string {
	return ""
}
func (l *LogrusWrapper) SetPrefix(p string) {
}
func (l *LogrusWrapper) Level() log.Lvl {
	switch logrus.GetLevel() {
	case logrus.ErrorLevel:
		return log.ERROR
	case logrus.WarnLevel:
		return log.WARN
	case logrus.InfoLevel:
		return log.INFO
	default:
		return log.DEBUG
	}
}
func (l *LogrusWrapper) SetLevel(v log.Lvl) {
	switch v {
	case log.ERROR:
		logrus.SetLevel(logrus.ErrorLevel)
	case log.WARN:
		logrus.SetLevel(logrus.WarnLevel)
	case log.INFO:
		logrus.SetLevel(logrus.InfoLevel)
	default:
		logrus.SetLevel(logrus.DebugLevel)
	}
}
func (l *LogrusWrapper) SetHeader(h string) {
}
func (l *LogrusWrapper) Print(i ...interface{}) {
	l.logger.Print(i)
}
func (l *LogrusWrapper) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args)
}
func (l *LogrusWrapper) Printj(j log.JSON) {
	l.logger.Print(j)
}
func (l *LogrusWrapper) Debug(i ...interface{}) {
	l.logger.Debug(i)
}
func (l *LogrusWrapper) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args)
}
func (l *LogrusWrapper) Debugj(j log.JSON) {
	l.logger.Debug(j)
}
func (l *LogrusWrapper) Info(i ...interface{}) {
	l.logger.Info(i)
}
func (l *LogrusWrapper) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args)
}
func (l *LogrusWrapper) Infoj(j log.JSON) {
	l.logger.Info(j)
}
func (l *LogrusWrapper) Warn(i ...interface{}) {
	l.logger.Warn(i)
}
func (l *LogrusWrapper) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args)
}
func (l *LogrusWrapper) Warnj(j log.JSON) {
	l.logger.Warn(j)
}
func (l *LogrusWrapper) Error(i ...interface{}) {
	l.logger.Error(i)
}
func (l *LogrusWrapper) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args)
}
func (l *LogrusWrapper) Errorj(j log.JSON) {
	l.logger.Error(j)
}
func (l *LogrusWrapper) Fatal(i ...interface{}) {
	l.logger.Fatal(i)
}
func (l *LogrusWrapper) Fatalj(j log.JSON) {
	l.logger.Fatal(j)
}
func (l *LogrusWrapper) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args)
}
func (l *LogrusWrapper) Panic(i ...interface{}) {
	l.logger.Panic(i)
}
func (l *LogrusWrapper) Panicj(j log.JSON) {
	l.logger.Panic(j)
}
func (l *LogrusWrapper) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args)
}
