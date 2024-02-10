package common

import "github.com/sirupsen/logrus"

var Logger *logrus.Logger = logrus.New()

func NewLogModule(name string) *logrus.Entry {
	return Logger.WithField("module", name)
}
