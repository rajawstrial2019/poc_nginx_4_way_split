package handler

import (
	"CASB-Securlet/securlets/gapps/notification-be/authserver/model"
	"CASB-Securlet/securlets/gapps/notification-be/authserver/token"
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
)

//HEADER Constants
const (
	HEADER_CHANNEL_TOKEN  = "X-Goog-Channel-Token"
	HEADER_CHANNEL_ID     = "X-Goog-Channel-Id"
	HEADER_RESOURCE_ID    = "X-Goog-Resource-Id"
	HEADER_RESOURCE_STATE = "X-Goog-Resource-State"
)

//ERROR Constants
const (
	ERROR_INVALID_TOKEN_URL_ENCODING  = "invalid Token (a061)"
	ERROR_INVALID_TOKEN_DECOMPRESSION = "invalid Token (a062)"
	ERROR_INVALID_TOKEN_CONTENT       = "invalid Token (a063)"

	ERROR_BAD_REQUEST_MISSING_HEADER    = "bad request, missing required header (a051)"
	ERROR_BAD_REQUEST_MISSING_PATH_ARGS = "bad request, missing required path args (a052)"
)

//ERROR Constants
const (
	DATE_FORMAT = "Jan 2, 2006 at 3:04pm (MST)"
)


type EventPathParam struct {
	TenantName       string `uri:"tenant_name" binding:"required,min=1"`
	SecurletInstance string `uri:"securlet_instance" binding:"required,min=1"`
	InstanceId       string `uri:"instance_id" binding:"required,min=1"`
}

type EventQueryParam struct {
	UserEmail string `form:"user_email"`
	Organizer string `form:"organizer"`
	IsTeam    bool   `form:"is_team"`
}

func PostEvent(appCtx *model.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appCtx.Logger.Info("EVENT - Got request.")

		var pathParam EventPathParam
		if err := ctx.ShouldBindUri(&pathParam); err != nil {
			appCtx.Logger.Errorf("Required Path arguments missing. %s", err)
			ctx.Writer.WriteString(fmt.Sprintf("Required Path arguments missing. %s", err))
			ctx.Status(http.StatusBadRequest)
			return
		}

		var queryParam EventQueryParam
		if err := ctx.ShouldBindQuery(&queryParam); err != nil {
			appCtx.Logger.Errorf("Required Query arguments missing. %s", err)
			ctx.Writer.WriteString(fmt.Sprintf("Required Query arguments missing. %s", err))
			ctx.Status(http.StatusBadRequest)
			return
		}

		var tokenList []string
		var notificationVersion int
		var primaryDomain string
		var tokenExpire time.Time
		skipSignVerify := true

		_, err := verifyRequest(appCtx, ctx)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			ctx.Writer.WriteString(fmt.Sprintf("Required Query arguments missing. %s", err))
			return
		}
		requestHelper := RequestHelper{}

		resourceState := requestHelper.GetHeader(ctx, HEADER_RESOURCE_STATE)
		channelId := requestHelper.GetHeader(ctx, HEADER_CHANNEL_ID)
		resourceId := requestHelper.GetHeader(ctx, HEADER_RESOURCE_ID)

		authToken, userEmail, organizer, instanceId, isTeamDrive, err := parseRequest(appCtx, ctx, pathParam, queryParam)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		if resourceState == "sync" {
			appCtx.Logger.Infof("Received sync notification for user %s on channel %s", userEmail, channelId)
			ctx.Status(http.StatusOK)
			return
		}

		if len(pathParam.SecurletInstance) == 0 {
			securletInstance, err := lookupSecurletInstance(pathParam.TenantName, userEmail)
			if err != nil {
				ctx.Status(http.StatusUnauthorized)
				return
			}

			_ = securletInstance
		}

		key := channelId + userEmail + instanceId

		jwtProvider := token.JWTProvider{Logger: appCtx.Logger}

		tokenParts := strings.Split(authToken, ".")
		if len(tokenParts) == 3 {
			keyHash := jwtProvider.GetSHA26HashAsHex(key)

			//dummy, _ := jwtProvider.CreateDummyToken(keyHash)
			//authToken = dummy
			claims, err := jwtProvider.VerifyToken(authToken, keyHash, skipSignVerify)
			if err != nil {
				ctx.Status(http.StatusUnauthorized)
				return
			} else {
				tokenList = append(tokenList, instanceId)
				if claims != nil {
					notificationVersion = claims.Version
					primaryDomain = claims.Domain
					tokenExpire = time.Unix(claims.ExpireAt, 0)
				}
			}
		} else {
			appCtx.Logger.Warn("To be implemented")
		}

		_ = authToken
		_ = userEmail
		_ = organizer
		_ = isTeamDrive
		_ = instanceId
		_ = channelId
		_ = resourceId
		_ = resourceState
		_ = notificationVersion
		_ = primaryDomain

		nonce := requestHelper.GetHeader(ctx, "nonce")
		currentTime := time.Now()
		ctx.Writer.WriteString(fmt.Sprintf("Hit the EVENT for %s\n\n\n", appCtx.Name))
		ctx.Writer.WriteString(fmt.Sprintf("- Response from %s running at port %s.\n", appCtx.Name, appCtx.Port))
		ctx.Writer.WriteString(fmt.Sprintf("- Host: %s\n", ctx.Request.URL.Host))
		ctx.Writer.WriteString(fmt.Sprintf("- URL:%s\n", ctx.Request.URL.Path))

		ctx.Writer.WriteString("\n\nJWT DATA:\n")
		
		ctx.Writer.WriteString(fmt.Sprintf("- SkipSignVerify: %v\n", skipSignVerify))
		ctx.Writer.WriteString(fmt.Sprintf("- NotificationVersion: %d\n", notificationVersion))
		ctx.Writer.WriteString(fmt.Sprintf("- PrimaryDomain: %s\n", primaryDomain))
		ctx.Writer.WriteString(fmt.Sprintf("- ExpiresAt: %s\n", tokenExpire.Format(DATE_FORMAT)))

		ctx.Writer.WriteString("\n\nHEADER DATA:\n")
		ctx.Writer.WriteString(fmt.Sprintf("- ResourceId: %s\n", resourceId))
		ctx.Writer.WriteString(fmt.Sprintf("- channelId: %s\n", channelId))

		ctx.Writer.WriteString("\n\nPATH DATA:\n")
		ctx.Writer.WriteString(fmt.Sprintf("- TenantName: %s\n", pathParam.TenantName))
		ctx.Writer.WriteString(fmt.Sprintf("- InstanceId: %s\n", pathParam.InstanceId))
		ctx.Writer.WriteString(fmt.Sprintf("- SecurletInstance: %s\n", pathParam.SecurletInstance))

		ctx.Writer.WriteString("\n\nQUERY DATA:\n")
		ctx.Writer.WriteString(fmt.Sprintf("- UserEmail: %s\n", userEmail))
		ctx.Writer.WriteString(fmt.Sprintf("- Organizer: %s\n", queryParam.Organizer))
		ctx.Writer.WriteString(fmt.Sprintf("- IsTeam: %v\n", queryParam.IsTeam))
		
		ctx.Writer.WriteString(fmt.Sprintf("- Nonce: %s\n", nonce))
		ctx.Writer.WriteString(fmt.Sprintf("- Current Time is: %s.\n", currentTime.Format(DATE_FORMAT)))

		ctx.Status(http.StatusOK)
	}
}

func lookupSecurletInstance(tenantName string, userEmail string) (string, error) {
	/*
				# Notification came from older endpoints for older account
		        if not securlet_instance:
				    securlet_instance = lookup_securlet_instance(mongo_client[tenant_name], user_email=user_email)
	*/
	return "TDB", nil
}

func verifyRequest(appCtx *model.AppContext, ctx *gin.Context) (bool, error) {
	return verifyRequestHeaders(appCtx, ctx)
}

func verifyRequestHeaders(appCtx *model.AppContext, ctx *gin.Context) (bool, error) {
	requestHelper := RequestHelper{}
	channelTokenHeader := requestHelper.GetHeader(ctx, HEADER_CHANNEL_TOKEN)
	if len(channelTokenHeader) == 0 {
		err := errors.New(ERROR_BAD_REQUEST_MISSING_HEADER)
		appCtx.Logger.Errorf("%s - %s", err, HEADER_CHANNEL_TOKEN)
		return false, err
	}
	channleId := requestHelper.GetHeader(ctx, HEADER_CHANNEL_ID)
	if len(channleId) == 0 {
		err := errors.New(ERROR_BAD_REQUEST_MISSING_HEADER)
		appCtx.Logger.Errorf("%s - %s", err, HEADER_CHANNEL_ID)
		return false, err
	}
	resourceId := requestHelper.GetHeader(ctx, HEADER_RESOURCE_ID)
	if len(resourceId) == 0 {
		err := errors.New(ERROR_BAD_REQUEST_MISSING_HEADER)
		appCtx.Logger.Errorf("%s - %s", err, HEADER_RESOURCE_ID)
		return false, err
	}
	return true, nil
}

// Token, user_email, organizer, isTeam, instanceId
func parseRequest(appCtx *model.AppContext, ctx *gin.Context, pathParam EventPathParam, queryParam EventQueryParam) (string, string, string, string, bool, error) {
	var isTeam bool = false
	var instanceId string = pathParam.InstanceId
	var organizer string = "unknown"
	var authToken string
	var userEmail string
	var uncompressedData string
	var err error
	requestHelper := RequestHelper{}

	for {
		channelToken := requestHelper.GetHeader(ctx, HEADER_CHANNEL_TOKEN)

		// Since go by default does the application/x-www-form-urlencoded, so we don't need unquote_plus
		channelToken, encodingErr := url.QueryUnescape(channelToken)
		if encodingErr != nil {
			err := errors.New(ERROR_INVALID_TOKEN_URL_ENCODING)
			appCtx.Logger.Errorf("%s - Invalid url encoding(%s)", err, encodingErr)
			break
		}

		var tokenParts []string = strings.SplitN(channelToken, ":", 2)

		compressionIndicator := tokenParts[0]

		if compressionIndicator == "t" {
			authToken = tokenParts[1]

			userEmail = queryParam.UserEmail
			organizer = queryParam.Organizer
			isTeam = queryParam.IsTeam
			/*
				isTeamArg := queryParam.IsTeam
				if len(isTeamArg) > 0 {
					isTeam = isTeamArg == "True"
				}
			*/
			break
		} else {
			var uncompressedBytes []byte
			if compressionIndicator == "c" {
				reader, callErr := zlib.NewReader(bytes.NewReader([]byte(tokenParts[1])))
				if callErr != nil {
					err := errors.New(ERROR_INVALID_TOKEN_DECOMPRESSION)
					appCtx.Logger.Errorf("%s - Decompression failed (NewReader) (%s)", err, callErr)
					break
				}
				defer reader.Close()

				uncompressedBytes, callErr = ioutil.ReadAll(reader)
				if callErr != nil {
					err := errors.New(ERROR_INVALID_TOKEN_DECOMPRESSION)
					appCtx.Logger.Errorf("%s - Decompression failed (ReadAll) (%s)", err, callErr)
					break
				}
				uncompressedData = string(uncompressedBytes)
			} else {
				uncompressedData = channelToken
			}

			parts := strings.Split(string(uncompressedData), ":")
			if len(parts) < 3 {
				err := errors.New(ERROR_INVALID_TOKEN_CONTENT)
				appCtx.Logger.Errorf("%s - not valid content", err)
				break
			}
			if len(parts) == 3 {
				authToken = parts[0]
				userEmail = parts[1]
				instanceId = parts[2]
				break
			} else {
				parts = strings.SplitN(string(uncompressedData), ":", 4)
				authToken = parts[0]
				userEmail = parts[1]
				instanceId = parts[2]
				organizer = parts[3]
				break
			}
		}
	}
	return authToken, userEmail, organizer, instanceId, isTeam, err
}
