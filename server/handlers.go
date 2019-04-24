package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hongjundu/go-rest-api-helper"
	"net/http"
)

func ginHandlerFunc(f func(c *gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		rsp, err := f(c)
		if err == nil {
			c.JSON(http.StatusOK, rsp)
		} else {
			c.JSON(http.StatusInternalServerError, apihelper.NewErrorResponse(err))
		}
	}
}

func (server *HttpServer) notFound(c *gin.Context) (response interface{}, err error) {
	err = apihelper.NewError("not_found", "not found")
	return
}

func (server *HttpServer) publicApi(c *gin.Context) (response interface{}, err error) {
	response = apihelper.NewOKResponse("public api")
	return
}

func (server *HttpServer) allowAllApi(c *gin.Context) (response interface{}, err error) {
	response = apihelper.NewOKResponse("allow all api")
	return
}

func (server *HttpServer) adminApi(c *gin.Context) (response interface{}, err error) {
	response = apihelper.NewOKResponse("admin api")
	return
}

func (server *HttpServer) viewApi(c *gin.Context) (response interface{}, err error) {
	response = apihelper.NewOKResponse("view api")
	return
}

func (server *HttpServer) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if len(token) == 0 {
			c.JSON(http.StatusUnauthorized, apihelper.NewErrorResponse(apihelper.NewError("unauthorized", "unauthorized")))
			c.Abort()
		}
	}
}
