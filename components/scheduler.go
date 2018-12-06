package components

import (
	"bitbucket.org/ogero/echodemo/components/schedulertasks"
	"github.com/bamzi/jobrunner"
	"github.com/sirupsen/logrus"
	"time"
)

type Scheduler struct {
	DataProvider SchedulerDataProvider
	logger       *logrus.Entry
}

type SchedulerDataProvider interface {
	SendMailsCronExpression() (string, error)
}

// NewScheduler creates a scheduler manager and spawn all schedules after 5 seconds
func NewScheduler(dataProvider SchedulerDataProvider, logger *logrus.Logger) (scheduler *Scheduler) {
	scheduler = &Scheduler{
		DataProvider: dataProvider,
		logger:       logger.WithField("entity", "Scheduler"),
	}
	go time.AfterFunc(5*time.Second, func() {
		scheduler.logger.Info("Starting jobrunner and scheduling tasks")
		jobrunner.Start()
		scheduler.RescheduleWelcomeMailTask()
	})
	return scheduler
}

func (o *Scheduler) RemoveJobByName(JobName string) {
	entries := jobrunner.Entries()
	for _, entry := range entries {
		if j, ok := entry.Job.(*jobrunner.Job); ok && j.Name == JobName {
			o.logger.Debugf("Removing Job %s", j.Name)
			jobrunner.Remove(entry.ID)
		}
	}
}
func (o *Scheduler) RescheduleWelcomeMailTask() (err error) {
	taskName := "WelcomeMailTask"
	o.RemoveJobByName(taskName)
	if scheduleExpression, err := o.DataProvider.SendMailsCronExpression(); err == nil {
		if len(scheduleExpression) > 0 {
			taskJob := scheduling.NewWelcomeMailTask(o.logger)
			err = jobrunner.Schedule(scheduleExpression, taskJob)
			if err != nil {
				o.logger.Errorf("Failed to schedule task %s. %f.", taskName, err)
				return err
			} else {
				o.logger.Debugf("Task %s scheduled at %s.", taskName, scheduleExpression)
			}
		} else {
			o.logger.Debugf("Task %s schedule is empty, it was disabled.", taskName)
		}
	} else {
		o.logger.Errorf("Failed to obtain schedule for task %s. %f", taskName, err)
		return err
	}
	return nil
}
