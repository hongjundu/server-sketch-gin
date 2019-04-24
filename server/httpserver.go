package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hongjundu/go-level-logger"
	"github.com/rs/cors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpServer struct {
	public *apiPrivilege
	admin  *apiPrivilege
	view   *apiPrivilege
	router *gin.Engine
}

func NewHttpServer() *HttpServer {
	return &HttpServer{
		public: newApiPrivilege(privilegeAll),
		admin:  newApiPrivilege(privilegeAdmin),
		view:   newApiPrivilege(privilegeView),
		router: gin.New(),
	}
}

func (server *HttpServer) configRouter() {
	logger.Debugf("[HttpServer] configRouter")

	server.router.Use(gin.Logger(), gin.Recovery())

	server.router.NoRoute(ginHandlerFunc(server.notFound))

	v1 := server.router.Group("/v1")
	authorized := v1.Group("/", server.Auth())

	authorized.GET("/public", ginHandlerFunc(server.publicApi))
	authorized.GET("/allowall", server.public.Handler(ginHandlerFunc(server.allowAllApi)))
	authorized.GET("/view", server.view.Handler(ginHandlerFunc(server.viewApi)))
	authorized.POST("/admin", server.admin.Handler(ginHandlerFunc(server.adminApi)))
}

func (server *HttpServer) Run(port int) error {
	logger.Debugf("HttpServer Start")

	server.configRouter()

	// CORS
	handler := cors.Default().Handler(server.router)
	c := cors.New(cors.Options{
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"authorization", "Content-Type"},
		OptionsPassthrough: true,
		AllowCredentials:   true,
	})
	handler = c.Handler(handler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("HttpServer: listen error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)

	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Infof("shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("server shutdown error: %v", err)
	}

	logger.Infof("server exit.")

	return nil
}
