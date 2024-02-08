package common

import "github.com/sirupsen/logrus"

var Logger *logrus.Logger = logrus.New()

type LogModule interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	Trace(args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Tracef(format string, args ...interface{})
}

func NewLogModule(name string) LogModule {
	return Logger.WithField("module", name)
}
