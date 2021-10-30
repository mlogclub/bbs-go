package render

import (
	"bbs-go/cache"
	"bbs-go/model"
	"bbs-go/model/constants"
	"github.com/mlogclub/simple"
	"github.com/mlogclub/simple/date"
	"strconv"
	"strings"
)

func BuildUserSimpleInfoDefaultIfNull(id int64) *model.UserSimpleInfo {
	user := cache.UserCache.Get(id)
	if user == nil {
		user = &model.User{}
		user.Id = id
		user.Username = simple.SqlNullString(strconv.FormatInt(id, 10))
		user.Nickname = "匿名用户" + strconv.FormatInt(id, 10)
		user.CreateTime = date.NowTimestamp()
	}
	return BuildUserSimpleInfo(user)
}

func BuildUserSimpleInfo(user *model.User) *model.UserSimpleInfo {
	if user == nil {
		return nil
	}
	ret := &model.UserSimpleInfo{
		Id:           user.Id,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		SmallAvatar:  HandleOssImageStyleAvatar(user.Avatar),
		TopicCount:   user.TopicCount,
		CommentCount: user.CommentCount,
		Score:        user.Score,
		Description:  user.Description,
		CreateTime:   user.CreateTime,
	}
	if user.Status == constants.StatusDeleted {
		ret.Nickname = "黑名单用户"
		ret.Description = ""
	}
	return ret
}

func BuildUser(user *model.User) *model.UserInfo {
	if user == nil {
		return nil
	}
	roles := strings.Split(user.Roles, ",")
	ret := &model.UserInfo{
		UserSimpleInfo:       *BuildUserSimpleInfo(user),
		Username:             user.Username.String,
		BackgroundImage:      user.BackgroundImage,
		SmallBackgroundImage: HandleOssImageStyleSmall(user.BackgroundImage),
		Type:                 user.Type,
		Roles:                roles,
		HomePage:             user.HomePage,
		Forbidden:            user.IsForbidden(),
		Status:               user.Status,
	}
	if len(ret.Description) == 0 {
		ret.Description = "这家伙很懒，什么都没留下"
	}
	if user.Status == constants.StatusDeleted {
		ret.Username = "blacklist"
		ret.HomePage = ""
		ret.Score = 0
		ret.Forbidden = true
	}
	return ret
}

func BuildUserProfile(user *model.User) *model.UserProfile {
	if user == nil {
		return nil
	}
	ret := &model.UserProfile{
		UserInfo:      *BuildUser(user),
		Email:         user.Email.String,
		EmailVerified: user.EmailVerified,
		PasswordSet:   len(user.Password) > 0,
	}
	return ret
}
