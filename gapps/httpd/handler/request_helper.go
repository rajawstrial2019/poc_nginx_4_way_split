package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type RequestHelper struct {

}

func (requestHelper *RequestHelper) GetHeader(ctx *gin.Context, headerName string) string {
	multi := ctx.Request.Header[http.CanonicalHeaderKey(headerName)]
	if len(multi) > 0 {
		return multi[0]
	}
	return ""
}

