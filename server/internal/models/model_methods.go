package models

import (
	"bbs-go/internal/models/constants"
	"strings"
	"time"

	"github.com/mlogclub/simple/common/arrays"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/strs"
)

// IsForbidden 是否禁言
func (u *User) IsForbidden() bool {
	if u.ForbiddenEndTime == 0 {
		return false
	}
	// 永久禁言
	if u.ForbiddenEndTime == -1 {
		return true
	}
	// 判断禁言时间
	return u.ForbiddenEndTime > dates.NowTimestamp()
}

// HasRole 是否有指定角色
func (u *User) HasRole(role string) bool {
	roles := strings.Split(u.Roles, ",")
	if len(roles) == 0 {
		return false
	}
	return arrays.Contains(role, roles)
}

// HasAnyRole 是否有指定的任意角色
func (u *User) HasAnyRole(roles ...string) bool {
	if len(roles) == 0 {
		return false
	}
	for _, role := range roles {
		if u.HasRole(role) {
			return true
		}
	}
	return false
}

// IsOwnerOrAdmin 是否是管理员
func (u *User) IsOwnerOrAdmin() bool {
	return u.HasAnyRole(constants.RoleOwner, constants.RoleAdmin)
}

// GetRoles 获取角色
func (u *User) GetRoles() []string {
	if strs.IsBlank(u.Roles) {
		return nil
	}
	ss := strings.Split(u.Roles, ",")
	if len(ss) == 0 {
		return nil
	}
	var roles []string
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if strs.IsNotBlank(s) {
			roles = append(roles, s)
		}
	}
	return roles
}

// InObservationPeriod 是否在观察期
// observeSeconds 观察时长
func (u *User) InObservationPeriod(observeSeconds int) bool {
	if observeSeconds <= 0 {
		return false
	}
	return dates.FromTimestamp(u.CreateTime).Add(time.Second * time.Duration(observeSeconds)).After(time.Now())
}

// GetTitle 获取帖子的标题
func (t *Topic) GetTitle() string {
	if t.Type == constants.TopicTypeTweet {
		if strs.IsNotBlank(t.Content) {
			return t.Content
		} else {
			return "分享图片"
		}
	} else {
		return t.Title
	}
}
