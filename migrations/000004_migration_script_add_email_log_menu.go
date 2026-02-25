package migrations

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

const emailLogMenuID int64 = 201

func migrate_add_email_log_menu() error {
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		now := dates.NowTimestamp()
		title := "邮件记录"
		if config.Instance.Language == config.LanguageEnUS {
			title = "Email Logs"
		}

		menu := repositories.MenuRepository.Get(tx, emailLogMenuID)
		if menu == nil {
			menu = repositories.MenuRepository.Take(tx, "path = ?", "/system/email-log")
		}

		if menu == nil {
			menu = &models.Menu{
				Model:      models.Model{Id: emailLogMenuID},
				ParentId:   18,
				Type:       string(constants.MenuTypeMenu),
				Name:       "EmailLog",
				Title:      title,
				Icon:       "",
				Path:       "/system/email-log",
				Component:  "system/email-log/index",
				SortNo:     25,
				Status:     constants.StatusOk,
				CreateTime: now,
				UpdateTime: now,
			}
			if err := repositories.MenuRepository.Create(tx, menu); err != nil {
				return err
			}
		} else {
			menu.ParentId = 18
			menu.Type = string(constants.MenuTypeMenu)
			menu.Name = "EmailLog"
			menu.Title = title
			menu.Icon = ""
			menu.Path = "/system/email-log"
			menu.Component = "system/email-log/index"
			menu.SortNo = 25
			menu.Status = constants.StatusOk
			menu.UpdateTime = now
			if err := repositories.MenuRepository.Update(tx, menu); err != nil {
				return err
			}
		}

		for _, roleCode := range []string{constants.RoleOwner, constants.RoleAdmin} {
			role := repositories.RoleRepository.Take(tx, "code = ?", roleCode)
			if role == nil {
				continue
			}
			if repositories.RoleMenuRepository.Take(tx, "role_id = ? AND menu_id = ?", role.Id, menu.Id) != nil {
				continue
			}
			if err := repositories.RoleMenuRepository.Create(tx, &models.RoleMenu{
				RoleId:     role.Id,
				MenuId:     menu.Id,
				CreateTime: now,
			}); err != nil {
				return err
			}
		}

		return nil
	})
}
