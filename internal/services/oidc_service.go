package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/oidc"
	"database/sql"
	"errors"
	"strings"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

func OIDCProvider(key string) (dto.OIDCProviderConfig, bool) {
	for _, provider := range SysConfigService.GetLoginConfig().OIDCProviders {
		if provider.Key == key {
			return provider, true
		}
	}
	return dto.OIDCProviderConfig{}, false
}

func OIDCLogin(provider dto.OIDCProviderConfig, claims *oidc.Claims) (*models.User, error) {
	db := sqls.DB()
	var identity models.OIDCIdentity
	err := db.Where("issuer = ? AND subject = ?", provider.Issuer, claims.Subject).First(&identity).Error
	if err == nil {
		return UserService.Get(identity.UserId), nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	nickname := strings.TrimSpace(claims.Name)
	if nickname == "" {
		nickname = strings.TrimSpace(claims.PreferredName)
	}
	if nickname == "" {
		nickname = "OIDC User"
	}
	if len([]rune(nickname)) > 16 {
		nickname = string([]rune(nickname)[:16])
	}
	user := &models.User{Nickname: nickname, Avatar: claims.Picture, Status: constants.StatusOk, CreateTime: dates.NowTimestamp(), UpdateTime: dates.NowTimestamp()}
	// A matching email is intentionally not linked to an existing account. It is
	// only stored for newly-created users when it does not conflict locally.
	if claims.EmailVerified && claims.Email != "" && UserService.GetByEmail(claims.Email) == nil {
		user.Email = sql.NullString{String: claims.Email, Valid: true}
		user.EmailVerified = true
	}
	return user, db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		return tx.Create(&models.OIDCIdentity{ProviderKey: provider.Key, Issuer: provider.Issuer, Subject: claims.Subject, UserId: user.Id, Email: claims.Email, Nickname: nickname, Avatar: claims.Picture, CreateTime: dates.NowTimestamp(), UpdateTime: dates.NowTimestamp()}).Error
	})
}
