package token

import (
	"errors"
)

var (
	TokenNotFoundErr = errors.New("token not found")
	ExpiredErr       = errors.New("token is expired")
	NotValidYetErr   = errors.New("token not active yet")
	MalformedErr     = errors.New("that's not even a token")
	InvalidErr       = errors.New("couldn't handle this token")
)

const (
	userTokenHeader = "Authorization"
	userTokenParam  = "_user_token"

	expireSeconds = 3600
	issuer        = "bbs-go"
	secret        = ""
)
