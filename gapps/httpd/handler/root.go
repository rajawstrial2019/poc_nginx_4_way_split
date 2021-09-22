package handler

import (
	"net/http"
	"time"
	"fmt"
	"CASB-Securlet/securlets/gapps/notification-be/authserver/model"
	"github.com/gin-gonic/gin"
)

func GetRoot(appCtx *model.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appCtx.Logger.Infof("ROOT - Got request.")
		requestHelper := RequestHelper{}
		nonce := requestHelper.GetHeader(ctx, "nonce")
		currentTime := time.Now()
		ctx.Writer.WriteString(fmt.Sprintf("Hit the root for %s\n\n\n", appCtx.Name))
		ctx.Writer.WriteString(fmt.Sprintf("- Response from %s running at port %s.\n", appCtx.Name, appCtx.Port))
		ctx.Writer.WriteString(fmt.Sprintf("- Host: %s URL:%s\n", ctx.Request.URL.Host, ctx.Request.URL.Path))
		ctx.Writer.WriteString(fmt.Sprintf("- Nonce: %s\n", nonce))
		ctx.Writer.WriteString(fmt.Sprintf("- Current Time is: %s.", currentTime.Format(DATE_FORMAT)))
		ctx.Status(http.StatusOK)
	}
}
