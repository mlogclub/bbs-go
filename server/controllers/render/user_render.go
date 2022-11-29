package render

import (
	"bbs-go/cache"
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/pkg/bbsurls"
	"strconv"
	"strings"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/spf13/cast"
)

func BuildUserInfoDefaultIfNull(id int64) *model.UserInfo {
	user := cache.UserCache.Get(id)
	if user == nil {
		user = &model.User{}
		user.Id = id
		user.Username = sqls.SqlNullString(strconv.FormatInt(id, 10))
		user.Nickname = "匿名用户" + strconv.FormatInt(id, 10)
		user.CreateTime = dates.NowTimestamp()
	}
	return BuildUserInfo(user)
}

func BuildUserInfo(user *model.User) *model.UserInfo {
	if user == nil {
		return nil
	}
	ret := &model.UserInfo{
		Id:           user.Id,
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

func BuildUserDetail(user *model.User) *model.UserDetail {
	if user == nil {
		return nil
	}
	ret := &model.UserDetail{
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

func BuildUserProfile(user *model.User) *model.UserProfile {
	if user == nil {
		return nil
	}
	ret := &model.UserProfile{
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
