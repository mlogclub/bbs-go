package migrations

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/repositories"

	"encoding/json"
	"strconv"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

type roleSeed struct {
	ID     int64
	Type   int
	Name   string
	Code   string
	SortNo int
	Remark string
	Status int
}

type categorySeed struct {
	ID          int64
	Name        string
	Description string
	Logo        string
	SortNo      int
	Status      int
}

type sysConfigSeed struct {
	Key         string
	Value       interface{}
	Name        string
	Description string
}

type seedData struct {
	Roles      []roleSeed
	Categories []categorySeed
	SysConfigs []sysConfigSeed
}

func migrate_init() error {
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		now := dates.NowTimestamp()
		seed := seedForLanguage()

		for _, r := range seed.Roles {
			existing := repositories.RoleRepository.Take(tx, "code = ?", r.Code)
			if existing != nil {
				existing.Type = r.Type
				existing.Name = r.Name
				existing.SortNo = r.SortNo
				existing.Remark = r.Remark
				existing.Status = r.Status
				existing.UpdateTime = now
				if err := repositories.RoleRepository.Update(tx, existing); err != nil {
					return err
				}
				continue
			}

			role := &models.Role{
				Model:      models.Model{Id: r.ID},
				Type:       r.Type,
				Name:       r.Name,
				Code:       r.Code,
				SortNo:     r.SortNo,
				Remark:     r.Remark,
				Status:     r.Status,
				CreateTime: now,
				UpdateTime: now,
			}
			if err := repositories.RoleRepository.Create(tx, role); err != nil {
				return err
			}
		}

		categoryIDMap := make(map[int64]int64)
		for _, n := range seed.Categories {
			existing := repositories.CategoryRepository.Take(tx, "name = ?", n.Name)
			if existing != nil {
				existing.Description = n.Description
				existing.Logo = n.Logo
				existing.SortNo = n.SortNo
				existing.Status = n.Status
				if err := repositories.CategoryRepository.Update(tx, existing); err != nil {
					return err
				}
				categoryIDMap[n.ID] = existing.Id
				continue
			}

			category := &models.Category{
				Model:       models.Model{Id: n.ID},
				Name:        n.Name,
				Description: n.Description,
				Logo:        n.Logo,
				SortNo:      n.SortNo,
				Status:      n.Status,
				CreateTime:  now,
			}
			if err := repositories.CategoryRepository.Create(tx, category); err != nil {
				return err
			}
			categoryIDMap[n.ID] = category.Id
		}

		for _, c := range seed.SysConfigs {
			existing := repositories.SysConfigRepository.GetByKey(tx, c.Key)
			if existing != nil {
				existing.Value = toConfigValue(c.Value)
				existing.Name = c.Name
				existing.Description = c.Description
				existing.UpdateTime = now
				if err := repositories.SysConfigRepository.Update(tx, existing); err != nil {
					return err
				}
				continue
			}

			cfg := &models.SysConfig{
				Key:         c.Key,
				Value:       toConfigValue(c.Value),
				Name:        c.Name,
				Description: c.Description,
				CreateTime:  now,
				UpdateTime:  now,
			}
			if err := repositories.SysConfigRepository.Create(tx, cfg); err != nil {
				return err
			}
		}

		// ensure defaultCategoryId sys config points to created category id
		if categoryID := categoryIDMap[1]; categoryID > 0 {
			if cfg := repositories.SysConfigRepository.GetByKey(tx, constants.SysConfigDefaultCategoryId); cfg != nil {
				cfg.Value = strconv.FormatInt(categoryID, 10)
				cfg.UpdateTime = now
				if err := repositories.SysConfigRepository.Update(tx, cfg); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func toConfigValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	default:
		if b, err := json.Marshal(val); err == nil {
			return string(b)
		}
	}
	return ""
}

func seedForLanguage() seedData {
	lang := config.Instance.Language
	if !lang.IsValid() {
		lang = config.DefaultLanguage
	}
	if lang == config.LanguageEnUS {
		return seedData{
			Roles: []roleSeed{
				{ID: 1, Type: constants.RoleTypeSystem, Name: "Owner", Code: constants.RoleOwner, SortNo: 0, Remark: "Owner with highest privileges", Status: constants.StatusOk},
			},
			Categories: []categorySeed{
				{ID: 1, Name: "Default", Description: "", Logo: "", SortNo: 0, Status: constants.StatusOk},
			},
			SysConfigs: []sysConfigSeed{
				{Key: constants.SysConfigSiteTitle, Value: "BBS-GO Demo Site", Name: "Site Title", Description: "Site Title"},
				{Key: constants.SysConfigSiteDescription, Value: "BBS-GO, an open source community system based on Go language", Name: "Site Description", Description: "Site Description"},
				{Key: constants.SysConfigBaseURL, Value: "/", Name: "Site URL", Description: "Site URL"},
				{Key: constants.SysConfigSiteKeywords, Value: []string{"bbs-go"}, Name: "Site Keywords", Description: "Site Keywords"},
				{Key: constants.SysConfigSiteNavs, Value: []map[string]string{
					{"title": "Topics", "url": "/topics"},
					{"title": "Articles", "url": "/articles"},
					// {"title": "Tasks", "url": "/tasks"},
				}, Name: "Site Navigation", Description: "Site Navigation"},
				{Key: constants.SysConfigDefaultCategoryId, Value: "1", Name: "Default Category", Description: "Default Category"},
				{Key: constants.SysConfigTokenExpireDays, Value: "365", Name: "User Login Validity Period (Days)", Description: "User Login Validity Period (Days)"},
				{Key: constants.SysConfigUrlRedirect, Value: "false"},
				{Key: constants.SysConfigEnableHideContent, Value: "false"},
				{Key: constants.SysConfigSiteLogo, Value: "/res/images/logo.png"},
				{Key: constants.SysConfigSiteNotification, Value: ""},
				{Key: constants.SysConfigRecommendTags, Value: ""},
				{Key: constants.SysConfigModules, Value: map[string]bool{"tweet": true, "topic": true, "qa": true, "article": false}},
				{Key: constants.SysConfigSmtpConfig, Value: dto.SmtpConfig{}},
				{Key: constants.SysConfigUploadConfig, Value: dto.UploadConfig{
					EnableUploadMethod: dto.Local,
					AliyunOss:          dto.AliyunOssUploadConfig{},
					TencentCos:         dto.TencentCosUploadConfig{},
					AwsS3:              dto.AwsS3UploadConfig{},
				}},
			},
		}
	}

	chinese := chineseSeedText(lang)
	return seedData{
		Roles: []roleSeed{
			{ID: 1, Type: constants.RoleTypeSystem, Name: chinese("超级管理员", "站長"), Code: constants.RoleOwner, SortNo: 0, Remark: chinese("超级管理员拥有最高权限", "站長擁有最高權限"), Status: constants.StatusOk},
		},
		Categories: []categorySeed{
			{ID: 1, Name: chinese("默认节点", "預設節點"), Description: "", Logo: "", SortNo: 0, Status: constants.StatusOk},
		},
		SysConfigs: []sysConfigSeed{
			{Key: constants.SysConfigSiteTitle, Value: "bbs-go演示站", Name: chinese("站点标题", "網站標題"), Description: chinese("站点标题", "網站標題")},
			{Key: constants.SysConfigSiteDescription, Value: chinese("bbs-go，基于Go语言的开源社区系统", "bbs-go，基於Go語言的開源社群系統"), Name: chinese("站点描述", "網站描述"), Description: chinese("站点描述", "網站描述")},
			{Key: constants.SysConfigBaseURL, Value: "/", Name: chinese("网站URL", "網站URL"), Description: chinese("网站URL", "網站URL")},
			{Key: constants.SysConfigSiteKeywords, Value: []string{"bbs-go"}, Name: chinese("站点关键字", "網站關鍵字"), Description: chinese("站点关键字", "網站關鍵字")},
			{Key: constants.SysConfigSiteNavs, Value: []map[string]string{
				{"title": chinese("话题", "話題"), "url": "/topics"},
				{"title": "文章", "url": "/articles"},
				// {"title": "任务", "url": "/tasks"},
			}, Name: chinese("站点导航", "網站導航"), Description: chinese("站点导航", "網站導航")},
			{Key: constants.SysConfigDefaultCategoryId, Value: "1", Name: chinese("默认节点", "預設節點"), Description: chinese("默认节点", "預設節點")},
			{Key: constants.SysConfigTokenExpireDays, Value: "365", Name: chinese("用户登录有效期(天)", "使用者登入有效期(天)"), Description: chinese("用户登录有效期(天)", "使用者登入有效期(天)")},
			{Key: constants.SysConfigUrlRedirect, Value: "false"},
			{Key: constants.SysConfigEnableHideContent, Value: "false"},
			{Key: constants.SysConfigSiteLogo, Value: ""},
			{Key: constants.SysConfigSiteNotification, Value: ""},
			{Key: constants.SysConfigRecommendTags, Value: ""},
			{Key: constants.SysConfigModules, Value: map[string]bool{"tweet": true, "topic": true, "qa": true, "article": true}},
			{Key: constants.SysConfigSmtpConfig, Value: dto.SmtpConfig{}},
			{Key: constants.SysConfigUploadConfig, Value: dto.UploadConfig{
				EnableUploadMethod: dto.Local,
				AliyunOss: dto.AliyunOssUploadConfig{
					StyleSplitter: "@",
				},
				TencentCos: dto.TencentCosUploadConfig{},
				AwsS3:      dto.AwsS3UploadConfig{},
			}},
		},
	}
}

// chineseSeedText selects localized defaults while preserving the shared seed values.
func chineseSeedText(lang config.Language) func(simplified, traditional string) string {
	if lang == config.LanguageZhTW {
		return func(_, traditional string) string { return traditional }
	}
	return func(simplified, _ string) string { return simplified }
}
