package storage

import (
	"bitbucket.org/ogero/echodemo/storage/types"
	"github.com/elithrar/simple-scrypt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
	"time"
)

type GormStore struct {
	*gorm.DB
	logger *logrus.Entry
}

type GormConfig struct {
	GormConnDialect  string
	GormConnArgs     string
	GormDebugQueries bool
	GormMustUsers    []GormConfigMustUsers
}

type GormConfigMustUsers struct {
	Email    string
	Password string
	Role     int
}

// NewCormStore creates a GormStore instance based on provided gorm configuration
func NewGormStore(config *GormConfig, logger *logrus.Logger) (gormStore *GormStore) {
	log := logger.WithField("entity", "Renderer")
	if len(config.GormConnDialect) == 0 || len(config.GormConnArgs) == 0 {
		log.Panic("GormConnDialect and GormConnArgs must be specified")
	}
	db, err := gorm.Open(config.GormConnDialect, config.GormConnArgs)
	if err != nil {
		log.WithError(err).Panic("Failed when calling gorm.Open")
	}
	db.LogMode(config.GormDebugQueries)
	scryptParams, err := scrypt.Calibrate(500*time.Millisecond, 64, scrypt.Params{})
	if err != nil {
		log.Error("Failed when calling scrypt.Calibrate. Falling back to scrypt.DefaultParams")
		scryptParams = scrypt.DefaultParams
	} else {
		scrypt.DefaultParams = scryptParams
	}
	gormStore = &GormStore{
		DB:     db,
		logger: log,
	}
	gormStore.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8 auto_increment=1")
	gormStore.AutoMigrate(&User{}, &Setting{})
	var c int
	for i := 0; i < len(config.GormMustUsers); i++ {
		if gormStore.Model(&User{}).Where("email = ?", config.GormMustUsers[i].Email).Count(&c); c == 0 && len(config.GormMustUsers[i].Email) > 0 {
			role := types.Role(config.GormMustUsers[i].Role)
			gormStore.Create(&User{
				Email:    config.GormMustUsers[i].Email,
				Password: config.GormMustUsers[i].Password,
				Role:     role,
			})
			log.Debugf("Created %s %s", role.FriendlyValue(), config.GormMustUsers[i].Email)
		}
	}

	for name, value := range DefaultSettings {
		if gormStore.Model(&Setting{}).Where("name = ?", name).Count(&c); c == 0 && len(name) > 0 {
			gormStore.Create(&Setting{
				Name:  name,
				Value: value,
			})
			log.Debugf("Created setting %s", name)
		}
	}
	return
}
