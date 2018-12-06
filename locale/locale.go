package locale

import (
	"bitbucket.org/ogero/echodemo/embed"
	"github.com/Unknwon/i18n"
	"github.com/sirupsen/logrus"
)

var availableLocales = map[string]bool{}

// InitLocales attempts to initialize all defined locales from embed, returning a map of locales with its load result
func InitLocales(logger *logrus.Logger) *map[string]bool {
	log := logger.WithField("entity", "Locale")
	firstRunLocalesCount := len(availableLocales)

	//For a better approach use embed.ReadFile("locale/locale_xx.ini")
	if _, ok := availableLocales["en"]; !ok {
		if err := i18n.SetMessage("en", embed.FileLocaleLocaleEnIni); err != nil {
			availableLocales["en"] = false
			log.WithError(err).WithField("lang", "en").Error("Failed when calling i18n.SetMessage")
		} else {
			availableLocales["en"] = true
		}
	}

	if firstRunLocalesCount == 0 {
		log.WithField("locales", availableLocales).Info("Locales loading finished.")
	}

	return &availableLocales
}
