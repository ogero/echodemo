package scheduling

import (
	"github.com/sirupsen/logrus"
)

type WelcomeMailTask struct {
	logger *logrus.Entry
}

func NewWelcomeMailTask(logger *logrus.Entry) WelcomeMailTask {
	task := WelcomeMailTask{
		logger: logger.WithField("task", "WelcomeMail"),
	}
	return task
}

func (o WelcomeMailTask) Run() {
	o.logger.Debugf("Sending welcome emails")
}
