package migrations

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/repositories"
	"strings"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

func migrate_smtp_config_to_sys_config() error {
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		now := dates.NowTimestamp()
		existing := repositories.SysConfigRepository.GetByKey(tx, constants.SysConfigSmtpConfig)
		if existing != nil && !isEmptySmtpConfigValue(existing.Value) {
			return nil
		}

		smtpCfg := dto.SmtpConfig{
			Host:     strings.TrimSpace(config.Instance.Smtp.Host),
			Port:     strings.TrimSpace(config.Instance.Smtp.Port),
			Username: strings.TrimSpace(config.Instance.Smtp.Username),
			Password: strings.TrimSpace(config.Instance.Smtp.Password),
			SSL:      config.Instance.Smtp.SSL,
		}
		if smtpCfg.Host == "" || smtpCfg.Port == "" || smtpCfg.Username == "" || smtpCfg.Password == "" {
			return nil
		}

		value := toConfigValue(smtpCfg)
		if existing == nil {
			name, description := smtpConfigMetaByLanguage()
			return repositories.SysConfigRepository.Create(tx, &models.SysConfig{
				Key:         constants.SysConfigSmtpConfig,
				Value:       value,
				Name:        name,
				Description: description,
				CreateTime:  now,
				UpdateTime:  now,
			})
		}
		existing.Value = value
		existing.UpdateTime = now
		return repositories.SysConfigRepository.Update(tx, existing)
	})
}

func isEmptySmtpConfigValue(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return true
	}
	var cfg dto.SmtpConfig
	if err := jsons.Parse(value, &cfg); err != nil {
		return false
	}
	return strings.TrimSpace(cfg.Host) == "" &&
		strings.TrimSpace(cfg.Port) == "" &&
		strings.TrimSpace(cfg.Username) == "" &&
		strings.TrimSpace(cfg.Password) == ""
}

func smtpConfigMetaByLanguage() (name, description string) {
	if config.Instance.Language == config.LanguageEnUS {
		return "SMTP Config", "SMTP Config"
	}
	return "SMTP配置", "SMTP配置"
}
