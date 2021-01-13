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

func BuildUserDefaultIfNull(id int64) *model.UserInfo {
	user := cache.UserCache.Get(id)
	if user == nil {
		user = &model.User{}
		user.Id = id
		user.Username = simple.SqlNullString(strconv.FormatInt(id, 10))
		user.CreateTime = date.NowTimestamp()
	}
	return BuildUser(user)
}

func BuildUserById(id int64) *model.UserInfo {
	user := cache.UserCache.Get(id)
	return BuildUser(user)
}

func BuildUser(user *model.User) *model.UserInfo {
	if user == nil {
		return nil
	}
	roles := strings.Split(user.Roles, ",")
	ret := &model.UserInfo{
		Id:                   user.Id,
		Username:             user.Username.String,
		Nickname:             user.Nickname,
		Avatar:               user.Avatar,
		SmallAvatar:          HandleOssImageStyleAvatar(user.Avatar),
		BackgroundImage:      user.BackgroundImage,
		SmallBackgroundImage: HandleOssImageStyleSmall(user.BackgroundImage),
		Email:                user.Email.String,
		EmailVerified:        user.EmailVerified,
		Type:                 user.Type,
		Roles:                roles,
		HomePage:             user.HomePage,
		Description:          user.Description,
		Score:                user.Score,
		TopicCount:           user.TopicCount,
		CommentCount:         user.CommentCount,
		PasswordSet:          len(user.Password) > 0,
		Forbidden:            user.IsForbidden(),
		Status:               user.Status,
		CreateTime:           user.CreateTime,
	}
	if len(ret.Description) == 0 {
		ret.Description = "这家伙很懒，什么都没留下"
	}
	if user.Status == constants.StatusDeleted {
		ret.Username = "blacklist"
		ret.Nickname = "黑名单用户"
		ret.Email = ""
		ret.HomePage = ""
		ret.Description = ""
		ret.Score = 0
		ret.Forbidden = true
	}
	return ret
}

func BuildUsers(users []model.User) []model.UserInfo {
	if len(users) == 0 {
		return nil
	}
	var responses []model.UserInfo
	for _, user := range users {
		item := BuildUser(&user)
		if item != nil {
			responses = append(responses, *item)
		}
	}
	return responses
}
