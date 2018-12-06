package providers

import "bitbucket.org/ogero/echodemo/storage"

type SchedulerProvider struct {
	storage *storage.GormStore
}

func NewSchedulerProvider(storage *storage.GormStore) *SchedulerProvider {
	return &SchedulerProvider{
		storage: storage,
	}
}

func (o *SchedulerProvider) SendMailsCronExpression() (string, error) {
	setting := &storage.Setting{}
	if err := o.storage.Where("name = ?", storage.SettingsKeyWelcomeMailSchedule).First(setting).Error; err == nil {
		return setting.Value, nil
	} else {
		return "", err
	}
}
