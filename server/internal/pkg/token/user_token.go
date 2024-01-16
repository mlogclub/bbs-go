package token

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/web/params"
	"github.com/spf13/cast"
)

type UserClaims struct {
	*jwt.RegisteredClaims

	UserId   int64  `json:"userId"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

func CreateToken(userId int64, nickname, avatar string) (string, error) {
	var (
		expiredAt = time.Now().Add(time.Duration(expireSeconds) * time.Second)
	)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &UserClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        cast.ToString(userId),
		},
		UserId:   userId,
		Nickname: nickname,
		Avatar:   avatar,
	})
	return claims.SignedString([]byte(secret))
}

func GetUser(c iris.Context) (user *UserClaims) {
	token := getToken(c)
	if strs.IsNotBlank(token) {
		user, _ = parseToken(token)
	}
	return
}

func getToken(c iris.Context) string {
	token := c.Request().Header.Get(userTokenHeader)
	if strs.IsNotBlank(token) {
		if strings.HasPrefix(token, "Bearer ") {
			return token[7:]
		}
		return token
	}
	token, _ = params.Get(c, userTokenParam)
	return token
}

func parseToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(secret), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, MalformedErr
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ExpiredErr
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, NotValidYetErr
			} else {
				return nil, InvalidErr
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, InvalidErr
	} else {
		return nil, InvalidErr
	}
}
