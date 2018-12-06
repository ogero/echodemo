package main

import (
	"bitbucket.org/ogero/echodemo/components"
	"bitbucket.org/ogero/echodemo/controllers"
	"bitbucket.org/ogero/echodemo/dist"
	"bitbucket.org/ogero/echodemo/embed"
	"bitbucket.org/ogero/echodemo/locale"
	"bitbucket.org/ogero/echodemo/storage"
	"bitbucket.org/ogero/echodemo/storage/providers"
	"flag"
	"github.com/fgrosse/goldi"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"log"
)

// DefaultConfig provides default values for not required settings
func DefaultConfig() *components.Config {
	return &components.Config{
		HTTPServerListenAddr: ":2332",
		SessionName:          "echodemo",
		SessionAuthKey:       "MwaS3BURJ15JBaEYPz5N7eN3ZehXYS6P",
		SessionEncKey:        "EVMH7vX2Pc56wiDSoq8wjDm9lVm8HuIO",
		RBACFile:             "rbac.json",
		LogFile:              "echodemo.log",
		LogLevel:             "debug",
		LogAsJSON:            true,
	}
}

func main() {

	tomlConfigFile := flag.String("config", "echodemo.ini", "Path to the settings file")
	version := flag.Bool("version", false, "Prints version")
	flag.Parse()

	if *version {
		log.Println("App:", dist.GitTag)
		log.Println("Build Date:", dist.Timestamp)
		log.Println("Commit:", dist.CommitHash)
		return
	}

	config := DefaultConfig()
	components.ReadConfig(*tomlConfigFile, config, true)

	defer components.SetupLogrus(config.LogFile, config.LogLevel, config.LogAsJSON)()
	logger := logrus.StandardLogger()
	logger.WithField("entity", "App").WithFields(logrus.Fields{
		"GitTag":     dist.GitTag,
		"Build Date": dist.Timestamp,
		"Commit":     dist.CommitHash,
	}).Info("Starting")

	registry := goldi.NewTypeRegistry()
	container := goldi.NewContainer(registry, map[string]interface{}{})

	container.InjectInstance("config", config)
	container.InjectInstance("logger", logger)
	container.InjectInstance("locale.availableLocales", locale.InitLocales(logger))
	container.RegisterType("storage", storage.NewGormStore, &config.GormConfig, "@logger")
	container.RegisterType("scheduler.provider", providers.NewSchedulerProvider, "@storage")
	container.RegisterType("scheduler", components.NewScheduler, "@scheduler.provider", "@logger")
	container.RegisterType("rbac", components.NewRBAC, "@config", "@logger", false)
	container.RegisterType("echoserver", components.NewEcho, container, "@config", "@logger")

	container.MustGet("scheduler")
	echoServer := container.MustGet("echoserver").(*echo.Echo)
	controllers.SetRoutes(echoServer, embed.Handler)
	echoServer.Logger.Fatal(echoServer.Start(config.HTTPServerListenAddr))
}
