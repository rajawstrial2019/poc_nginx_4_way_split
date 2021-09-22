package api

import (
	"CASB-Securlet/securlets/gapps/notification-be/authserver/httpd/handler"
	"CASB-Securlet/securlets/gapps/notification-be/authserver/model"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func NewServer(appCtx *model.AppContext) (*Server, error) {
	server := Server{}

	defaultRouter := gin.Default()

	defaultRouter.GET("/", handler.GetRoot(appCtx))
	defaultRouter.GET("/api/v1/healthcheck/", handler.GetHealthCheck(appCtx))
	defaultRouter.GET("/auth/", handler.GetAuth(appCtx))
	defaultRouter.POST("/:tenant_name/:securlet_instance/:instance_id/drive/notification/", handler.PostEvent(appCtx))
	server.router = defaultRouter

	return &server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
