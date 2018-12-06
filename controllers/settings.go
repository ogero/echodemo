package controllers

import (
	"bitbucket.org/ogero/echodemo/components"
	"bitbucket.org/ogero/echodemo/storage"
	"bitbucket.org/ogero/echodemo/storage/types"
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/Unknwon/i18n"
	"github.com/astaxie/beego/utils/pagination"
	"github.com/bamzi/jobrunner"
	"github.com/labstack/echo"
	"net/http"
	"os"
	"time"
)

func GetSettings_JobRunner(c echo.Context) (err error) {
	cc := c.(*components.EchoContext)
	if !cc.UserIsGranted(types.Permission_JobRunner) {
		return c.Render(http.StatusForbidden, "default/error.html", map[string]interface{}{"code": http.StatusForbidden})
	}
	return c.Render(http.StatusOK, "default/settings.jobrunner.html", map[string]interface{}{
		"now":       time.Now(),
		"jobrunner": jobrunner.StatusPage(),
	})
}

func GetSettings_List(c echo.Context) error {
	cc := c.(*components.EchoContext)
	if !cc.UserIsGranted(types.Permission_Settings_List) {
		return c.Render(http.StatusForbidden, "default/error.html", map[string]interface{}{"code": http.StatusForbidden})
	}
	var settings []storage.Setting
	cc.Storage().Order("name asc").Find(&settings)
	return c.Render(http.StatusOK, "default/settings.list.html", map[string]interface{}{
		"settings":     settings,
		"canJobrunner": cc.UserIsGranted(types.Permission_JobRunner),
		"canUpdate":    cc.UserIsGranted(types.Permission_Settings_Update),
	})
}

func AnySettings_CreateUpdate(c echo.Context) error {
	cc := c.(*components.EchoContext)
	if !cc.UserIsGranted(types.Permission_Settings_Update) {
		return c.Render(http.StatusForbidden, "default/error.html", map[string]interface{}{"code": http.StatusForbidden})
	}
	setting := &storage.Setting{}
	name := c.Param("name")
	if err := cc.Storage().Where("name = ?", name).First(setting).Error; err != nil {
		cc.Logrus.WithError(err).Warn("Failed to find setting %d", name)
		return c.Render(http.StatusNotFound, "default/error.html", map[string]interface{}{"code": http.StatusNotFound})
	}
	if cc.Request().Method == http.MethodPost {
		if err := c.Bind(setting); err == nil {
			// TODO: Validate
			if err := cc.Storage().Save(setting).Error; err == nil {
				switch name {
				case storage.SettingsKeyWelcomeMailSchedule:
					err := cc.DI.MustGet("scheduler").(*components.Scheduler).RescheduleWelcomeMailTask()
					if err != nil && len(setting.Value) != 0 {
						cc.Logrus.WithError(err).Error("Failed when rescheduling Batch Email")
						return c.Render(http.StatusInternalServerError, "default/error.html", map[string]interface{}{"code": http.StatusInternalServerError})
					}
				}
				cc.AddFlash(i18n.Tr(cc.GetLang(), "settings.settings_saved"))
				return cc.Redirect(http.StatusSeeOther, cc.Echo().Reverse("settings.list"))
			} else {
				cc.AddFlash(err.Error())
			}
		}
	}
	return c.Render(http.StatusOK, "default/settings.update.html", map[string]interface{}{
		"setting": setting,
	})
}

func GetSettings_LogRead(c echo.Context) error {
	cc := c.(*components.EchoContext)
	const BUFSIZE = 1024 * 3

	logFile := cc.DI.MustGet("config").(*components.Config).LogFile

	stat, err := os.Stat(logFile)
	if err != nil {
		return err
	}
	file, err := os.Open(logFile)
	if err != nil {
		return err
	}
	defer file.Close()

	paginator := pagination.NewPaginator(c.Request(), BUFSIZE, stat.Size())

	start := stat.Size() - int64(paginator.Page())*BUFSIZE
	if start < 0 {
		start = 0
	}
	buf := make([]byte, BUFSIZE)
	n, err := file.ReadAt(buf, start)

	reader := bufio.NewReader(bytes.NewReader(buf[:n]))
	logs := make([]map[string]interface{}, 0)
	var line []byte
	for {
		var lineMap map[string]interface{}
		line, err = reader.ReadSlice('\n')
		if err := json.Unmarshal(line, &lineMap); err == nil {
			if v, ok := lineMap["time"]; ok {
				lineMap["time"], _ = time.Parse(time.RFC3339, v.(string))
			}
			j, _ := json.MarshalIndent(lineMap, "", "  ")
			lineMap["json"] = string(j)
			logs = append(logs, lineMap)
		}
		if err != nil {
			break
		}
	}
	for i := len(logs)/2 - 1; i >= 0; i-- {
		opp := len(logs) - 1 - i
		logs[i], logs[opp] = logs[opp], logs[i]
	}

	return c.Render(http.StatusOK, "default/settings.readlog.html", map[string]interface{}{
		"logs":      logs,
		"paginator": paginator,
	})
}
