package controllers

import (
	"bitbucket.org/ogero/echodemo/components"
	"bitbucket.org/ogero/echodemo/dist"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"net/http"
)

func GetHome_Index(c echo.Context) error {
	return c.Render(http.StatusOK, "default/home.index.html", map[string]interface{}{})
}

func GetHome_Version(c echo.Context) error {
	return c.JSON(http.StatusOK, struct {
		Tag    string
		Commit string
		Date   string
	}{Tag: dist.GitTag, Commit: dist.CommitHash, Date: dist.Timestamp})
}

func GetHome_PrivacyPolicy(c echo.Context) error {
	cc := c.(*components.EchoContext)
	cc.Logrus.WithError(errors.New("dummy error example")).Error()
	cc.AddFlash("Dummy flash message")
	return c.Render(http.StatusInternalServerError, "default/error.html", map[string]interface{}{"code": http.StatusInternalServerError})
}
