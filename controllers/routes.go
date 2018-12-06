package controllers

import (
	"bitbucket.org/ogero/echodemo/components"
	"github.com/labstack/echo"
	"golang.org/x/net/webdav"
	"net/http"
)

func SetRoutes(s *echo.Echo, staticFilesHandler *webdav.Handler) {

	s.GET("/assets/*", echo.WrapHandler(staticFilesHandler))

	// Attach routes and its handlers
	s.GET("/", GetHome_Index).Name = "home.index"
	s.GET("/version", GetHome_Version).Name = "version"
	s.GET("/privacy", GetHome_PrivacyPolicy).Name = "home.privacy-policy"
	// Users
	s.GET("/users/login", GetUsers_Login, AnonymousOnly).Name = "users.login"
	s.POST("/users/login", PostUsers_Login, AnonymousOnly)
	s.POST("/users/logout", PostUsers_Logout, AuthenticatedOnly).Name = "users.logout"
	s.GET("/users/list", GetUsers_List, AuthenticatedOnly).Name = "users.list"
	s.GET("/users/create", AnyUsers_CreateUpdate, AuthenticatedOnly).Name = "users.create"
	s.POST("/users/create", AnyUsers_CreateUpdate, AuthenticatedOnly)
	s.GET("/users/update/:id", AnyUsers_CreateUpdate, AuthenticatedOnly).Name = "users.update"
	s.POST("/users/update/:id", AnyUsers_CreateUpdate, AuthenticatedOnly)
	s.POST("/users/delete/:id", PostUsers_Delete, AuthenticatedOnly).Name = "users.delete"
	//// Settings
	s.GET("/settings/jobrunner", GetSettings_JobRunner, AuthenticatedOnly).Name = "settings.jobrunner"
	s.GET("/settings/list", GetSettings_List, AuthenticatedOnly).Name = "settings.list"
	s.GET("/settings/update/:name", AnySettings_CreateUpdate, AuthenticatedOnly).Name = "settings.update"
	s.POST("/settings/update/:name", AnySettings_CreateUpdate, AuthenticatedOnly)
	s.GET("/settings/log", GetSettings_LogRead, AuthenticatedOnly).Name = "settings.readlog"
}

// AuthenticatedOnly checks if user is logged in. Redirects to login page if anonymous.
func AuthenticatedOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(*components.EchoContext)
		userID := cc.GetLoggedUserID()
		if userID <= 0 {
			return c.Redirect(http.StatusSeeOther, cc.Echo().Reverse("users.login"))
		}
		return next(c)
	}
}

// AnonymousOnly checks if user is anonymous. Redirects to index page if logged in.
func AnonymousOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(*components.EchoContext)
		userID := cc.GetLoggedUserID()
		if userID > 0 {
			return c.Redirect(http.StatusSeeOther, cc.Echo().Reverse("home.index"))
		}
		return next(c)
	}
}
