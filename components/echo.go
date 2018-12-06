package components

import (
	"github.com/fgrosse/goldi"
	"github.com/ipfans/echo-session"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func NewEcho(di *goldi.Container, config *Config, logger *logrus.Logger) *echo.Echo {
	log := logger.WithField("entity", "App")
	e := echo.New()
	e.HideBanner = true
	e.Debug = config.LogLevel == "debug"
	if len(config.SessionName) == 0 {
		log.Panic("SessionName can't be empty")
	}
	e.Use(session.Sessions(
		config.SessionName,
		session.NewCookieStore([][]byte{[]byte(config.SessionAuthKey), []byte(config.SessionEncKey)}...),
	))
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			return h(NewEchoContext(context, di, log))
		}
	})
	var err error
	if e.Renderer, err = NewTemplateRenderer(); err != nil {
		log.Panic(err)
	}
	e.Logger = NewLogrusWrapper(logger)
	return e
}
