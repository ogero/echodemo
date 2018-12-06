package components

import (
	"bitbucket.org/ogero/echodemo/storage/types"
	"github.com/Unknwon/i18n"
)

type Menu struct {
	Label   string
	Url     string
	Visible bool
	Menu    *Menu
}

func GetMenu(c *EchoContext) []Menu {
	u := c.Echo().Reverse
	p := func(permission string) bool {
		return c.UserIsGranted(permission)
	}
	return []Menu{
		{Label: i18n.Tr(c.GetLang(), "menu.users"), Url: u("users.list"), Visible: p(types.Permission_Users_List)},
		{Label: i18n.Tr(c.GetLang(), "menu.settings"), Url: u("settings.list"), Visible: p(types.Permission_Settings_List)},
	}
}
