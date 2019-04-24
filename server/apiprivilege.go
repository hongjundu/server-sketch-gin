package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hongjundu/go-level-logger"
	"github.com/hongjundu/go-rest-api-helper"
	"net/http"
	"strings"
)

const (
	privilegeAll   = "*"
	privilegeAdmin = "admin"
	privilegeView  = "view"
)

type apiPrivilege struct {
	allowPrivileges []string
}

func newApiPrivilege(privileges ...string) *apiPrivilege {
	privelege := &apiPrivilege{make([]string, 0, 0)}
	for _, priv := range privileges {
		privelege.allowPrivileges = append(privelege.allowPrivileges, priv)
	}
	return privelege
}

func (privelege *apiPrivilege) Handler(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIp := c.ClientIP()
		token := c.GetHeader("token")

		logger.Infof("client ip: %v", clientIp)
		logger.Infof("token: %s", token)

		if token != "123456" {
			c.JSON(http.StatusUnauthorized, apihelper.NewErrorResponse(apihelper.NewError("unauthorized", "unauthorized")))
			c.Abort()
			return
		}

		// TODO: should parse from token
		privileges := make([]string, 0, 0)
		strPriveleges := c.GetHeader("privileges")
		if len(strPriveleges) > 0 {
			privileges = strings.Split(strPriveleges, ",")
		}

		if !privelege.allow(privileges) {
			c.JSON(http.StatusUnauthorized, apihelper.NewErrorResponse(apihelper.NewError("no_privilege", "no privilege")))
			c.Abort()
		} else {
			handler(c)
		}
	}

}

func (privelege *apiPrivilege) allow(userPrivileges []string) bool {
	if stringSlice(privelege.allowPrivileges).Contains(privilegeAll) ||
		stringSlice(userPrivileges).Contains(privilegeAll) ||
		stringSlice(userPrivileges).Contains(privilegeAdmin) {
		return true
	}

	for _, userPriv := range userPrivileges {
		if stringSlice(privelege.allowPrivileges).Contains(userPriv) {
			return true
		}
	}

	return false
}

type stringSlice []string

func (array stringSlice) Contains(str string) bool {
	for _, val := range array {
		if strings.Compare(val, str) == 0 {
			return true
		}
	}
	return false
}
