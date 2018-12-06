package components

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mikespook/gorbac"
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

type RBAC struct {
	*gorbac.RBAC
	permissions gorbac.Permissions
	logger      *logrus.Entry
}

// NewRBAC creates an RBAC from the provided configuration
func NewRBAC(config *Config, logger *logrus.Logger, panicOnFail ...bool) (rbac *RBAC) {
	log := logger.WithField("entity", "RBAC")
	shouldPanicOnFail := len(panicOnFail) == 1 && panicOnFail[0] == true

	file, err := os.Open(config.RBACFile)
	if err != nil {
		expectedRBACPath := config.RBACFile
		if !path.IsAbs(expectedRBACPath) {
			if cwd, err := os.Getwd(); err == nil {
				expectedRBACPath = path.Join(cwd, config.RBACFile)
			}
		}
		err = errors.New(fmt.Sprintf("Error when opening RBAC file %s: %v", expectedRBACPath, err))
		failLog := log.WithError(err).WithField("RBACFile", config.RBACFile)
		if shouldPanicOnFail {
			failLog.Panic("Setup RBAC failed")
		} else {
			failLog.Error("Setup RBAC failed")
		}
		return nil
	}
	defer file.Close()

	var rbacRules map[string]map[string][]string
	if err := json.NewDecoder(file).Decode(&rbacRules); err != nil {
		failLog := log.WithError(err).WithField("RBACFile", config.RBACFile)
		if shouldPanicOnFail {
			failLog.Panic("Setup RBAC failed")
		} else {
			failLog.Error("Setup RBAC failed")
		}
		return nil
	}

	rbac = &RBAC{
		RBAC:        gorbac.New(),
		logger:      log,
		permissions: make(gorbac.Permissions),
	}

	rolesRulesCount, permissionsRulesCount := 0, 0

	jsonRoles := rbacRules["roles"]
	for rid, pids := range jsonRoles {
		role := gorbac.NewStdRole(rid)
		rolesRulesCount++
		for _, pid := range pids {
			_, ok := rbac.permissions[pid]
			if !ok {
				rbac.permissions[pid] = gorbac.NewStdPermission(pid)
				permissionsRulesCount++
			}
			role.Assign(rbac.permissions[pid])
		}
		rbac.Add(role)
	}

	jsonInheritance := rbacRules["inheritance"]
	for rid, parents := range jsonInheritance {
		if err := rbac.SetParents(rid, parents); err != nil {
			if shouldPanicOnFail {
				rbac.logger.WithError(err).Panic("Setup RBAC failed")
			}
			return nil
		}
	}

	rbac.logger.Infof("Loaded %d permissions on %d roles", permissionsRulesCount, rolesRulesCount)

	return rbac
}

func (o *RBAC) ExistsAndIsGranted(permission string, role string) bool {
	if p, ok := o.permissions[permission]; ok {
		return o.IsGranted(role, p, nil)
	} else {
		o.logger.Infof("Unknown permission %s", permission)
		return false
	}
}
