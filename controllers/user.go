package controllers

import (
	"bitbucket.org/ogero/echodemo/components"
	"bitbucket.org/ogero/echodemo/storage"
	"bitbucket.org/ogero/echodemo/storage/types"
	"fmt"
	"github.com/Unknwon/i18n"
	"github.com/astaxie/beego/utils/pagination"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"time"
)

const USERS_PER_PAGE = 8

func obtainFilteredUserQuery(c *components.EchoContext) (query *gorm.DB) {
	query = c.Storage().Table("users")
	if dateFrom := c.QueryParam("created_date_from"); len(dateFrom) > 0 {
		if t, err := time.Parse("02-01-2006", dateFrom); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if dateTo := c.QueryParam("created_date_to"); len(dateTo) > 0 {
		if t, err := time.Parse("02-01-2006", dateTo); err == nil {
			query = query.Where("created_at <= ?", t.Add(24*time.Hour-1*time.Millisecond))
		}
	}
	if terms := c.QueryParam("terms"); len(terms) > 0 {
		query = query.Where("email LIKE ?", fmt.Sprintf("%%%s%%", terms))
	}
	return
}

func GetUsers_Login(c echo.Context) error {
	return c.Render(http.StatusOK, "default/users.login.html", map[string]interface{}{})
}

func PostUsers_Logout(c echo.Context) error {
	cc := c.(*components.EchoContext)
	if cc.GetLoggedUserID() > 0 {
		cc.SetLoggedUserID(0)
	}
	return c.Redirect(http.StatusSeeOther, cc.Echo().Reverse("home.index"))
}

func PostUsers_Login(c echo.Context) error {
	cc := c.(*components.EchoContext)
	email := c.FormValue("email")
	password := c.FormValue("password")
	user := &storage.User{}
	if cc.Storage().Select("id, password").Where("email = ? AND deleted_at IS NULL", email).First(user).RecordNotFound() ||
		!user.CheckPasswordMatch(password) {
		cc.AddFlash("Wrong email or password")
		return cc.Redirect(http.StatusSeeOther, cc.Echo().Reverse("users.login"))
	} else {
		cc.SetLoggedUserID(user.ID)
		return cc.Redirect(http.StatusSeeOther, cc.Echo().Reverse("users.list"))
	}
}

func GetUsers_List(c echo.Context) error {
	cc := c.(*components.EchoContext)
	if !cc.UserIsGranted(types.Permission_Users_List) {
		return c.Render(http.StatusForbidden, "default/error.html", map[string]interface{}{"code": http.StatusForbidden})
	}
	var users []storage.User
	query := obtainFilteredUserQuery(cc)
	count := 0
	query.Where("deleted_at IS NULL").Count(&count)
	paginator := pagination.NewPaginator(c.Request(), USERS_PER_PAGE, count)
	query.Offset(paginator.Offset()).Limit(USERS_PER_PAGE).Order("email asc").Find(&users)
	return c.Render(http.StatusOK, "default/users.list.html", map[string]interface{}{
		"users":             users,
		"created_date_from": c.QueryParam("created_date_from"),
		"created_date_to":   c.QueryParam("created_date_to"),
		"terms":             c.QueryParam("terms"),
		"filter":            c.QueryParam("filter"),
		"paginator":         paginator,
		"canCreate":         cc.UserIsGranted(types.Permission_Users_Create),
		"canUpdate":         cc.UserIsGranted(types.Permission_Users_Update),
		"canDelete":         cc.UserIsGranted(types.Permission_Users_Delete),
	})
}

func AnyUsers_CreateUpdate(c echo.Context) error {
	cc := c.(*components.EchoContext)
	id, _ := strconv.Atoi(c.Param("id"))
	user := &storage.User{}
	if (id == 0 && !cc.UserIsGranted(types.Permission_Users_Create)) ||
		(id != 0 && !cc.UserIsGranted(types.Permission_Users_Update)) {
		return c.Render(http.StatusForbidden, "default/error.html", map[string]interface{}{"code": http.StatusForbidden})
	}
	if id != 0 {
		if err := cc.Storage().First(user, id).Error; err != nil {
			cc.Logrus.WithError(err).Warn("Failed to find user %d", id)
			return c.Render(http.StatusNotFound, "default/error.html", map[string]interface{}{"code": http.StatusNotFound})
		}
	}
	if cc.Request().Method == http.MethodPost {
		oldPassword := user.Password
		if err := c.Bind(user); err == nil {
			if len(user.Password) == 0 {
				user.Password = oldPassword
			}
			// TODO: Validate
			if len(user.Email) != 0 && len(user.Password) != 0 {
				if err := cc.Storage().Save(user).Error; err == nil {
					cc.AddFlash(i18n.Tr(cc.GetLang(), "users.user_saved"))
					return cc.Redirect(http.StatusSeeOther, cc.Echo().Reverse("users.list"))
				} else {
					cc.AddFlash(err.Error())
				}
			} else {
				cc.AddFlash(i18n.Tr(cc.GetLang(), "users.bad_input"))
			}
		}
	}
	user.Password = ""
	return c.Render(http.StatusOK, "default/users.create-update.html", map[string]interface{}{
		"user":  user,
		"roles": types.Role.List(0),
	})
}

func PostUsers_Delete(c echo.Context) error {
	cc := c.(*components.EchoContext)
	if !cc.UserIsGranted(types.Permission_Users_Delete) {
		return c.Render(http.StatusForbidden, "default/error.html", map[string]interface{}{"code": http.StatusForbidden})
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 0 {
		return cc.String(http.StatusBadRequest, "User not specified")
	}
	if err := cc.Storage().Delete(&storage.User{Model: gorm.Model{ID: uint(id)}}).Error; err == nil {
		cc.AddFlash(i18n.Tr(cc.GetLang(), "users.user_deleted"))
	} else {
		cc.Logrus.WithError(err).Warn("Failed to delete user %d", id)
	}
	return cc.NoContent(http.StatusOK)
}
