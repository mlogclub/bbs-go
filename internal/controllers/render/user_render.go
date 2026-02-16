package render

import (
	"math"

	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/resp"
	"bbs-go/internal/pkg/bbsurls"
	"bbs-go/internal/pkg/idcodec"
	"bbs-go/internal/pkg/locales"
	"strconv"
	"strings"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/spf13/cast"
)

func BuildUserInfoDefaultIfNull(id int64) *resp.UserInfo {
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

func BuildUserInfo(user *models.User) *resp.UserInfo {
	if user == nil {
		return nil
	}
	ret := &resp.UserInfo{
		Id:           idcodec.Encode(user.Id),
		Type:         user.Type,
		Nickname:     user.Nickname,
		Gender:       user.Gender,
		Birthday:     user.Birthday,
		TopicCount:   user.TopicCount,
		CommentCount: user.CommentCount,
		FansCount:    user.FansCount,
		FollowCount:  user.FollowCount,
		Score:        user.Score,
		Exp:          user.Exp,
		Level:        user.Level,
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

	if levelConfig := cache.LevelConfigCache.GetByLevel(user.Level); levelConfig != nil {
		ret.LevelTitle = levelConfig.Title
	}

	if len(ret.Description) == 0 {
		ret.Description = locales.Get("user.default_description")
	}
	if user.Status == constants.StatusDeleted {
		ret.Nickname = locales.Get("user.blacklist")
		ret.Description = ""
		ret.Score = 0
		ret.Exp = 0
		ret.Level = 0
		ret.Forbidden = true
	} else {
		ret.ExpProgress = buildExpProgress(user)
	}
	return ret
}

// buildExpProgress 根据用户当前经验与等级配置，计算当前等级内经验进度（用于进度条与文案展示）
func buildExpProgress(user *models.User) *resp.ExpProgressResponse {
	if user == nil {
		return nil
	}
	current := cache.LevelConfigCache.GetByLevel(user.Level)
	if current == nil {
		return nil
	}
	needCurrent := current.NeedExp
	expInCurrent := user.Exp - needCurrent
	if expInCurrent < 0 {
		expInCurrent = 0
	}

	next := cache.LevelConfigCache.GetByLevel(user.Level + 1)
	var expNeedForNext int
	var isMaxLevel bool
	if next == nil {
		expNeedForNext = 0
		isMaxLevel = true
	} else {
		expNeedForNext = next.NeedExp - current.NeedExp
		isMaxLevel = false
	}

	var percent int
	if isMaxLevel || expNeedForNext <= 0 {
		percent = 100
	} else {
		if expInCurrent > expNeedForNext {
			expInCurrent = expNeedForNext
		}
		percent = int(math.Round(float64(expInCurrent) * 100 / float64(expNeedForNext)))
		if percent > 100 {
			percent = 100
		}
	}

	return &resp.ExpProgressResponse{
		CurrentExp:          user.Exp,
		Level:               user.Level,
		LevelTitle:          current.Title,
		ExpInCurrentLevel:   expInCurrent,
		ExpNeedForNextLevel: expNeedForNext,
		ExpProgressPercent:  percent,
		IsMaxLevel:          isMaxLevel,
	}
}

func BuildUserDetail(user *models.User) *resp.UserDetail {
	if user == nil {
		return nil
	}
	ret := &resp.UserDetail{
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

func BuildUserProfile(user *models.User) *resp.UserProfile {
	if user == nil {
		return nil
	}
	ret := &resp.UserProfile{
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
	return bbsurls.AbsUrl("/res/images/avatars/" + cast.ToString(avatarIndex) + ".png")
}
