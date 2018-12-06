package components

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime/debug"
)

// Setups logrus component, using a filesystem log hook.
// It's recommended to call this function like this:
//
//  defer components.SetupLogrus(settings.LogFile, settings.LogLevel, settings.LogAsJSON)()
func SetupLogrus(LogFile string, LogLevel string, LogAsJSON bool) func() {
	os.MkdirAll(path.Dir(LogFile), os.ModePerm)
	pathMap := lfshook.PathMap{
		logrus.DebugLevel: LogFile,
		logrus.InfoLevel:  LogFile,
		logrus.WarnLevel:  LogFile,
		logrus.ErrorLevel: LogFile,
		logrus.FatalLevel: LogFile,
		logrus.PanicLevel: LogFile,
	}
	var formatter logrus.Formatter
	if LogAsJSON {
		formatter = &logrus.JSONFormatter{}
		logrus.SetFormatter(formatter)
	}
	fsHook := lfshook.NewHook(pathMap, formatter)

	logrus.AddHook(fsHook)

	level, err := logrus.ParseLevel(LogLevel)
	if err != nil {
		level = logrus.DebugLevel
	}
	logrus.SetLevel(level)
	if err != nil {
		logrus.WithError(err).Error("Using default logrus level debug instead")
	}

	return func() {
		err := recover()
		if err != nil {
			debug.PrintStack()
			logrus.Fatal(err)
		} else {
			logrus.Info("Goodbye!")
		}
	}
}
