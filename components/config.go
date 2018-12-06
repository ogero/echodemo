package components

import (
	"bitbucket.org/ogero/echodemo/storage"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

type Config struct {
	HTTPServerListenAddr string
	SessionName          string
	SessionAuthKey       string
	SessionEncKey        string
	RBACFile             string
	LogLevel             string
	LogFile              string
	LogAsJSON            bool
	storage.GormConfig
}

// Reads a TOML file into the provided interface.
// It's recommended to call this function like this:
//
//  config := flag.String("config", "service-config.ini", "Path to the settings file")
//  flag.Parse()
//
//  settings := DefaultConfig()
//  components.ReadConfig(*tomlConfigFile, config, true)
func ReadConfig(tomlConfigFile string, config interface{}, panicOnFail ...bool) error {

	shouldPanicOnFail := len(panicOnFail) == 1 && panicOnFail[0] == true

	_, err := os.Stat(tomlConfigFile)
	if err != nil {
		expectedConfigPath := tomlConfigFile
		if !path.IsAbs(expectedConfigPath) {
			if cwd, err := os.Getwd(); err == nil {
				expectedConfigPath = path.Join(cwd, tomlConfigFile)
			}
		}
		err = errors.New(fmt.Sprintf("Error when poking config file %s: %v", expectedConfigPath, err))
		if shouldPanicOnFail {
			logrus.WithError(err).Panic("ReadConfig failed")
		}
		return err
	}

	if _, err := toml.DecodeFile(tomlConfigFile, config); err != nil {
		err = errors.New(fmt.Sprintf("Error when decoding config file: %v", err))
		if shouldPanicOnFail {
			logrus.WithError(err).Panic("ReadConfig failed")
		}
		return err
	}

	return nil
}
