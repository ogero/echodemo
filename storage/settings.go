package storage

import "github.com/jinzhu/gorm"

var SettingsKeyWelcomeMailSchedule = "job.welcome_mail"
var DefaultSettings = make(map[string]string)

func init() {
	DefaultSettings[SettingsKeyWelcomeMailSchedule] = "@midnight"
}

type Setting struct {
	gorm.Model
	Name  string
	Value string `form:"value"`
}
