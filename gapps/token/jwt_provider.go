package token

import (
	"errors"
	"time"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"CASB-Securlet/securlets/gapps/notification-be/authserver/model"
	"encoding/base64"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

type JWTProvider struct {
	Logger    *logrus.Entry
}

func (jwtProvider *JWTProvider) GetSHA26HashAsHex(data string) (string){
	h := sha256.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	keyHash := hex.EncodeToString(bs)
	return keyHash
}

func (jwtProvider *JWTProvider) convertToken(token string) (string, error) {
	tokenParts := strings.Split(token, ".")
	var jwt64 string
	if len(tokenParts) == 3 {
		header, err := jwtProvider.convertStdBase64ToRawBase64(tokenParts[0])
		if err != nil {
			return jwt64, err
		}
		data, err := jwtProvider.convertStdBase64ToRawBase64(tokenParts[1])
		if err != nil {
			return jwt64, err
		}
		sign, err := jwtProvider.convertStdBase64ToRawBase64(tokenParts[2])
		if err != nil {
			return jwt64, err
		}

		jwt64 = header + "." + data + "." + sign
		return jwt64, nil
	} else {
		return jwt64, errors.New("invalid JWT Token")
	}
}

func (jwtProvider *JWTProvider) VerifyToken(token string, key string, skipSignVerify bool) (*model.Claims, error) {

	keyfunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			jwtProvider.Logger.Debug("Signature verification failed.")
			return nil, ErrInvalidToken
		}
		return []byte(key), nil
	}

	if skipSignVerify {
		keyfunc = nil
		jwtProvider.Logger.Warn("!!!Signature validation is skipped !!!")
	}

	jwt64, _ := jwtProvider.convertToken(token)
	jwtProvider.Logger.Debugf("Converted the Token in compatible format.")

	jwtToken, err := jwt.ParseWithClaims(jwt64, &model.Claims{}, keyfunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			jwtProvider.Logger.Debugf("Token is Expired. %s", err)
			return nil, ErrExpiredToken
		}
		if skipSignVerify && (verr.Errors == jwt.ValidationErrorUnverifiable || verr.Errors == jwt.ValidationErrorSignatureInvalid) {
			jwtProvider.Logger.Debugf("Signature Verification was skipped. %s", verr)
		} else {
			jwtProvider.Logger.Debugf("Token is Invalid. %s", err)
			return nil, ErrInvalidToken
		}
	}
	claims, ok := jwtToken.Claims.(*model.Claims)
	if !ok {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func (jwtProvider *JWTProvider) convertStdBase64ToRawBase64(input string) (string, error) {
	var temp []byte
	temp, err := base64.URLEncoding.DecodeString(input)
	if err != nil {
		return input, nil
	}
	return base64.RawURLEncoding.EncodeToString(temp), nil
}

func (jwtProvider *JWTProvider) CreateDummyToken(key string) (string, error) {
	var claims model.Claims
	claims.Domain = "dcidevtest.org"
	claims.ExpireAt = time.Now().Add(10 * 24 * time.Hour).Unix()
	claims.Version = 2
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	dummy, _:= jwtToken.SignedString([]byte(key))
	jwt64, _ := jwtProvider.convertToken(dummy)
	return jwt64, nil
}