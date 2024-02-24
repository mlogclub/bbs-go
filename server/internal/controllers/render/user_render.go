package render

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/bbsurls"
	"strconv"
	"strings"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/spf13/cast"
)

func BuildUserInfoDefaultIfNull(id int64) *models.UserInfo {
	user := cache.UserCache.Get(id)
	if user == nil {
		user = &models.User{}
		user.Id = id
		user.Type = constants.UserTypeNormal
		user.Username = sqls.SqlNullString(strconv.FormatInt(id, 10))
		user.Nickname = "匿名用户" + strconv.FormatInt(id, 10)
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
		avatar := RandomAvatar(user.Id)
		ret.Avatar = avatar
		ret.SmallAvatar = avatar
	}
	if len(ret.Description) == 0 {
		ret.Description = "这家伙很懒，什么都没留下"
	}
	if user.Status == constants.StatusDeleted {
		ret.Nickname = "黑名单用户"
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
