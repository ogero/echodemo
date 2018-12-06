package components

import (
	"bitbucket.org/ogero/echodemo/storage"
	"bytes"
	"github.com/fgrosse/goldi"
	"github.com/ipfans/echo-session"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

const ECHO_SESSION_USER_ID_KEY = "user-id"

type EchoContext struct {
	DI     *goldi.Container
	Logrus *logrus.Entry
	user   *storage.User
	echo.Context
}

func NewEchoContext(context echo.Context, di *goldi.Container, logger *logrus.Entry) *EchoContext {
	return &EchoContext{
		Context: context,
		DI:      di,
		Logrus:  logger,
	}
}

func (c *EchoContext) Render(code int, name string, data interface{}) (err error) {
	if c.Echo().Renderer == nil {
		return echo.ErrRendererNotRegistered
	}
	buf := new(bytes.Buffer)
	if err = c.Echo().Renderer.Render(buf, name, data, c); err != nil {
		return
	}
	return c.HTMLBlob(code, buf.Bytes())
}

func (c *EchoContext) Storage() *storage.GormStore {
	return c.DI.MustGet("storage").(*storage.GormStore)
}

func (c *EchoContext) RBAC() *RBAC {
	return c.DI.MustGet("rbac").(*RBAC)
}

// GetLoggedUserID obtains current user id. Returns 0 for anonymous and >0 for logged users.
func (c *EchoContext) GetLoggedUserID() uint {
	if user, ok := session.Default(c).Get(ECHO_SESSION_USER_ID_KEY).(uint); ok {
		return user
	}
	return 0
}

// SetLoggedUser sets current logged user.
func (c *EchoContext) SetLoggedUserID(UserID uint) {
	ses := session.Default(c)
	if UserID > 0 {
		c.Logrus.Debugf("User %d logged in", UserID)
		ses.Set(ECHO_SESSION_USER_ID_KEY, UserID)
	} else {
		userID := c.GetLoggedUserID()
		c.Logrus.Debugf("User %d logged out", userID)
		ses.Delete(ECHO_SESSION_USER_ID_KEY)
	}
	// Set curent user to nil, so forcing a refresh on calls to GetLoggedUser
	c.user = nil
	ses.Save()
}

func (c *EchoContext) AddFlash(flash string) {
	ses := session.Default(c)
	ses.AddFlash(flash)
	ses.Save()
}

func (c *EchoContext) Flashes() (flashes []string) {
	ses := session.Default(c)
	v := ses.Flashes()
	for i := 0; i < len(v); i++ {
		if f, ok := v[i].(string); ok {
			flashes = append(flashes, f)
		}
	}
	ses.Save()
	return
}

func (c *EchoContext) GetLang() string {
	return "en"
}

// GetLoggedUser obtains current user. Returns nil for anonymous, and User for logged users.
func (c *EchoContext) GetLoggedUser(refresh ...bool) (user *storage.User) {
	userID := c.GetLoggedUserID()
	if userID > 0 {
		if c.user == nil || c.user.ID == 0 || (len(refresh) > 0 && refresh[0] == true) {
			user = &storage.User{}
			if !c.Storage().First(user, userID).RecordNotFound() {
				c.user = user
			} else {
				c.Logrus.Error(gorm.ErrRecordNotFound)
				user = nil
			}
		} else {
			user = c.user
		}
	}
	return
}

// CheckRole checks if user is logged in, and if its role is allowed to perform the requested permission
func (c *EchoContext) UserIsGranted(permission string) bool {
	if user := c.GetLoggedUser(); user != nil {
		return c.RBAC().ExistsAndIsGranted(permission, user.Role.RBACValue())
	}
	return false
}
