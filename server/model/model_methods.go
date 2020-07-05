package model

import (
	"github.com/mlogclub/simple"
	"strings"
	"time"
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
	return u.ForbiddenEndTime > simple.NowTimestamp()
}

// HasRole 是否有指定角色
func (u *User) HasRole(role string) bool {
	roles := strings.Split(u.Roles, ",")
	if len(roles) == 0 {
		return false
	}
	return simple.Contains(role, roles)
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

// GetRoles 获取角色
func (u *User) GetRoles() []string {
	if simple.IsBlank(u.Roles) {
		return nil
	}
	ss := strings.Split(u.Roles, ",")
	if len(ss) == 0 {
		return nil
	}
	var roles []string
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if simple.IsNotBlank(s) {
			roles = append(roles, s)
		}
	}
	return roles
}

// InObservationPeriod 是否在观察期，observeHour观察时长
func (u *User) InObservationPeriod(observeHour int) bool {
	if observeHour <= 0 {
		return false
	}
	return simple.TimeFromTimestamp(u.CreateTime).Add(time.Hour * time.Duration(observeHour)).After(time.Now())
}
