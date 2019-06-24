package model

import "time"

const (
	OauthClientStatusOk       = 0
	OauthClientStatusDisabled = 1
)

type OauthClient struct {
	Model
	ClientId     string    `gorm:"unique;not null;size:64" json:"clientId" form:"clientId"`
	ClientSecret string    `gorm:"size:128" json:"clientSecret" form:"clientSecret"`
	Domain       string    `gorm:"size:1024" json:"domain" form:"domain"`
	CallbackUrl  string    `gorm:"size:1024" json:"callbackUrl" form:"callbackUrl"`
	Status       int       `gorm:"index:idx_status" json:"status" form:"status"`
	CreateTime   time.Time `json:"createTime" form:"createTime"`
}

type OauthToken struct {
	Model
	ExpiredAt    int64  `gorm:"index:idx_expired_at" json:"expiredAt"`
	Code         string `gorm:"index:idx_code" json:"code"`
	AccessToken  string `gorm:"index:idx_access_token" json:"accessToken"`
	RefreshToken string `gorm:"index:idx_refresh_token" json:"refreshToken"`
	Data         string `gorm:"type:text" json:"data"`
}
