package common

import (
	"bbs-go/internal/models"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/idcodec"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web/params"
)

var validate = validator.New()

const (
	current_user_key = "__current_user"
)

func SetCurrentUser(ctx iris.Context, user *models.User) {
	ctx.Values().Set(current_user_key, user)
}

func GetCurrentUserID(ctx iris.Context) int64 {
	user := GetCurrentUser(ctx)
	if user != nil {
		return user.Id
	}
	return 0
}

func GetCurrentUser(ctx iris.Context) *models.User {
	if v := ctx.Values().Get(current_user_key); v != nil {
		if u, ok := v.(*models.User); ok {
			return u
		}
	}
	return nil
}

func CheckLogin(ctx iris.Context) (*models.User, error) {
	user := GetCurrentUser(ctx)
	if user == nil {
		return nil, errs.NotLogin()
	}
	return user, nil
}

func IsLogin(ctx iris.Context) bool {
	return GetCurrentUser(ctx) != nil
}

func GetID(ctx iris.Context, name string) int64 {
	idStr, _ := params.Get(ctx, name)
	return idcodec.Decode(idStr)
}
