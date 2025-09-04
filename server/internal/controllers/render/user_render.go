package render

import (
	"strconv"
	"strings"

	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/bbsurls"
	"bbs-go/internal/pkg/locales"

	"github.com/spf13/cast"

	"bbs-go/internal/pkg/simple/common/dates"
	"bbs-go/internal/pkg/simple/common/strs"
	"bbs-go/internal/pkg/simple/sqls"
)

func BuildUserInfoDefaultIfNull(id int64) *models.UserInfo {
	user := cache.UserCache.Get(id)
	if user == nil {
		user = &models.User{}
		user.Id = id
		user.Type = constants.UserTypeNormal
		user.Username = sqls.SqlNullString(strconv.FormatInt(id, 10))
		user.Nickname = locales.Getf("user.anonymous", id)
		user.CreateTime = dates.NowTimestamp()
	}
	return BuildUserInfo(user)
}

func BuildUserInfo(user *models.User) *models.UserInfo {
	if user == nil {
		return nil
	}
	ret := &models.UserInfo{
		Id:           user.Id,
		Type:         user.Type,
		Nickname:     user.Nickname,
		Gender:       user.Gender,
		Birthday:     user.Birthday,
		TopicCount:   user.TopicCount,
		CommentCount: user.CommentCount,
		FansCount:    user.FansCount,
		FollowCount:  user.FollowCount,
		Score:        user.Score,
		Description:  user.Description,
		CreateTime:   user.CreateTime,
		Forbidden:    user.IsForbidden(),
	}
	if strs.IsNotBlank(user.Avatar) {
		ret.Avatar = user.Avatar
		ret.SmallAvatar = HandleOssImageStyleAvatar(user.Avatar)
	} else {
		// avatar := RandomAvatar(user.Id)
		// ret.Avatar = avatar
		// ret.SmallAvatar = avatar
	}
	if len(ret.Description) == 0 {
		ret.Description = locales.Get("user.default_description")
	}
	if user.Status == constants.StatusDeleted {
		ret.Nickname = locales.Get("user.blacklist")
		ret.Description = ""
		ret.Score = 0
		ret.Forbidden = true
	}
	return ret
}

func BuildUserDetail(user *models.User) *models.UserDetail {
	if user == nil {
		return nil
	}
	ret := &models.UserDetail{
		UserInfo:             *BuildUserInfo(user),
		Username:             user.Username.String,
		BackgroundImage:      user.BackgroundImage,
		SmallBackgroundImage: HandleOssImageStyleSmall(user.BackgroundImage),
		HomePage:             user.HomePage,
		Status:               user.Status,
	}
	if user.Status == constants.StatusDeleted {
		ret.Username = "blacklist"
		ret.HomePage = ""
	}
	return ret
}

func BuildUserProfile(user *models.User) *models.UserProfile {
	if user == nil {
		return nil
	}
	ret := &models.UserProfile{
		UserDetail:    *BuildUserDetail(user),
		Email:         user.Email.String,
		EmailVerified: user.EmailVerified,
		PasswordSet:   len(user.Password) > 0,
	}

	if strs.IsNotBlank(user.Roles) {
		ret.Roles = strings.Split(user.Roles, ",")
	}
	return ret
}

func RandomAvatar(userId int64) string {
	avatarCount := 128
	avatarIndex := userId % int64(avatarCount)
	return bbsurls.AbsUrl("/images/avatars/" + cast.ToString(avatarIndex) + ".png")
}
