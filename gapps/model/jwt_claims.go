package model

import (
	"errors"
	"time"
)

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

type Claims struct {
	Domain   string		`json:"dom"`
	ExpireAt int64		`json:"exp"`
	Version  int		`json:"ver"`
}

func (claims *Claims) Valid() error {
	exp := time.Unix(claims.ExpireAt, 0)

	//FIXME - Implement Token Expire Time Validation
	exp = exp.AddDate(1, 0, 0)

	if time.Now().After(exp) {
		return ErrExpiredToken
	}
	return nil
}
