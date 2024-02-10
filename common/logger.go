package common

import "github.com/sirupsen/logrus"

var Logger *logrus.Logger = logrus.New()

func NewLogModule(name string) *logrus.Entry {
	return Logger.WithField("module", name)
}

func NewLogModules(name string, subs ...string) []*logrus.Entry {
	logs := make([]*logrus.Entry, len(subs)+1)
	logs[0] = NewLogModule(name)
	for i, sub := range subs {
		logs[i] = NewLogModule(name + "." + sub)
	}
	return logs
}
