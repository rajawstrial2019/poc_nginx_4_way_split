package handler

import (
	"net/http"
	"time"
	"fmt"
	"CASB-Securlet/securlets/gapps/notification-be/authserver/model"
	"github.com/gin-gonic/gin"
)

func GetHealthCheck(appCtx *model.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appCtx.Logger.Infof("GAAPS - HEALTH CHECK - Got request.")
		currentTime := time.Now()
		ctx.Writer.WriteString(fmt.Sprintf("Health check response at %s.", currentTime.Format(DATE_FORMAT)))
		ctx.Status(http.StatusOK)
	}
}
