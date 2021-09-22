package handler

import (
	"fmt"
	"net/http"
	"time"
	"strings"
	"CASB-Securlet/securlets/gapps/notification-be/authserver/model"	
	"github.com/gin-gonic/gin"
)

const (
	HEADER_TENANT = "X-CASB-TENANT"
	HEADER_AUTHORIZATION = "Authorization"
	HEADER_ORIGINAL_URI = "X-Original-URI"
	HEADER_NONCE = "nonce"
)

func GetAuth(appCtx *model.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appCtx.Logger.Infof("AUTH - Got request.")
		requestHelper := RequestHelper{}

		nonce := requestHelper.GetHeader(ctx, HEADER_NONCE)
		currentTime := time.Now()
		appCtx.Logger.Infof("Header - %s", ctx.Request.Header)

		var original_uri_header string
		var authorization_header string
		original_uri_header = requestHelper.GetHeader(ctx, HEADER_ORIGINAL_URI)
		authorization_header = requestHelper.GetHeader(ctx, HEADER_AUTHORIZATION)
		if len(authorization_header) == 0 {
			appCtx.Logger.Errorf("AUTH - Authorization header is missing.")

			ctx.Writer.WriteString("Authorization header is missing")
			ctx.Writer.WriteString(fmt.Sprintf("- Response from %s running at port %s.\n", appCtx.Name, appCtx.Port))
			ctx.Writer.WriteString(fmt.Sprintf("- Host: %s URL:%s\n", ctx.Request.URL.Host, ctx.Request.URL.Path))
			ctx.Writer.WriteString(fmt.Sprintf("- Nonce: %s\n", nonce))
			ctx.Writer.WriteString(fmt.Sprintf("- Header: %s\n", ctx.Request.Header))
			ctx.Writer.WriteString(fmt.Sprintf("- Current Time is: %s.\n", currentTime.Format(DATE_FORMAT)))
			ctx.Status(http.StatusUnauthorized)
			return
		}

		if len(original_uri_header) == 0 {
			appCtx.Logger.Errorf("AUTH - X-Original-URI header is missing.")

			ctx.Writer.WriteString(fmt.Sprintf("- Response Headers set: %s= %s.\n", HEADER_TENANT, requestHelper.GetHeader(ctx, HEADER_TENANT)))
			ctx.Writer.WriteString("!!! X-Original-URI header is missing - Unauthorized Error - Status code 401 !!!\n\n\n")
			ctx.Writer.WriteString(fmt.Sprintf("- Response from %s running at port %s.\n", appCtx.Name, appCtx.Port))
			ctx.Writer.WriteString(fmt.Sprintf("- Host: %s URL:%s\n", ctx.Request.URL.Host, ctx.Request.URL.Path))
			ctx.Writer.WriteString(fmt.Sprintf("- Nonce: %s\n", nonce))
			ctx.Writer.WriteString(fmt.Sprintf("- Header: %s\n", ctx.Request.Header))
			ctx.Writer.WriteString(fmt.Sprintf("- Current Time is: %s.\n", currentTime.Format(DATE_FORMAT)))
			ctx.Status(http.StatusUnauthorized)
			return
		} else {
			if strings.Contains(original_uri_header, "gapps") {
				appCtx.Logger.Infof("AUTH - Gapps is OK.")

				success := GetTenantFromAuthToken(appCtx, ctx)
				if success {
					ctx.Status(http.StatusOK)
				} else {
					ctx.Status(http.StatusUnauthorized)
				}

				ctx.Writer.WriteString(fmt.Sprintf("- Response from %s running at port %s.\n", appCtx.Name, appCtx.Port))
				ctx.Writer.WriteString(fmt.Sprintf("- Host: %s URL:%s\n", ctx.Request.URL.Host, ctx.Request.URL.Path))
				ctx.Writer.WriteString(fmt.Sprintf("- Nonce: %s\n", nonce))
				ctx.Writer.WriteString(fmt.Sprintf("- Authorization: %s\n", authorization_header))
				ctx.Writer.WriteString(fmt.Sprintf("- Response Headers set: %s= %s.\n", HEADER_TENANT, requestHelper.GetHeader(ctx, HEADER_TENANT)))
				ctx.Writer.WriteString(fmt.Sprintf("- Current Time is: %s.\n", currentTime.Format(DATE_FORMAT)))
			} else {
				appCtx.Logger.Errorf("AUTH - Non-Gapps request is NOT OK.")

				ctx.Status(http.StatusUnauthorized)

				ctx.Writer.WriteString("!!! X-Original-URI header is missing - Unauthorized Error - Status code 401 !!!\n\n\n")
				ctx.Writer.WriteString(fmt.Sprintf("- Response from %s running at port %s.\n", appCtx.Name, appCtx.Port))
				ctx.Writer.WriteString(fmt.Sprintf("- Host: %s URL:%s\n", ctx.Request.URL.Host, ctx.Request.URL.Path))
				ctx.Writer.WriteString(fmt.Sprintf("- Nonce: %s\n", nonce))
				ctx.Writer.WriteString(fmt.Sprintf("- Header: %s\n", ctx.Request.Header))
				ctx.Writer.WriteString(fmt.Sprintf("- Current Time is: %s.\n", currentTime.Format(DATE_FORMAT)))
			}
		}

		ctx.Status(http.StatusOK)
	}
}

func GetTenantFromAuthToken(appCtx *model.AppContext, ctx *gin.Context) bool {
	requestHelper := RequestHelper{}
	auth := requestHelper.GetHeader(ctx, HEADER_AUTHORIZATION)
	appCtx.Logger.Infof("Inside Token parsing")
	tokenSplit := strings.Split(auth, " ")
	if len(tokenSplit) > 1 {
		if tokenSplit[0] != "Bearer" {
			ctx.Status(http.StatusUnauthorized)
			ctx.Writer.WriteString("!!!Invalid gapps token, Bearer key word missing - 'Bearer gapps_tenant=acme-3648'\n\n")
			appCtx.Logger.Errorf("!!!Invalid gapps token, Bearer key word missing - 'Bearer gapps_tenant=acme-3648'\n\n")
		} else {
			tenantSplit := strings.Split(tokenSplit[1], "=")
			if len(tenantSplit) > 1 {
				if tenantSplit[0] == "gapps_tenant" {
					ctx.Writer.Header().Set(HEADER_TENANT, tenantSplit[1])
					ctx.Status(http.StatusOK)
					appCtx.Logger.Infof("Validated box token.\n\n")
					return true
				} else {
					ctx.Status(http.StatusUnauthorized)
					ctx.Writer.WriteString("!!!Invalid gapps token, Bearer key word missing - 'Bearer gapps_tenant=acme-3648'\n\n")
					appCtx.Logger.Errorf("!!!Invalid gapps token, first item has to be tenant - 'Bearer gapps_tenant=acme-3648'\n\n")
				}
			} else {
				ctx.Status(http.StatusUnauthorized)
				ctx.Writer.WriteString("!!!Invalid gapps token, tenant missing - 'Bearer gapps_tenant=acme-3648'\n\n")
				appCtx.Logger.Errorf("!!!Invalid gapps token, tenant missing - 'Bearer gapps_tenant=acme-3648'\n\n")
			}
		}
	} else {
		ctx.Status(http.StatusUnauthorized)
		ctx.Writer.WriteString("!!!Invalid gapps token - 'Bearer gapps_tenant=acme-3648'\n\n")
		appCtx.Logger.Errorf("!!!Invalid gapps token - 'Bearer gapps_tenant=acme-3648'\n\n")
	}
	return false
}