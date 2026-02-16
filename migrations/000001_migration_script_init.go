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

type topicNodeSeed struct {
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

type menuSeed struct {
	ID        int64
	ParentID  int64
	Type      string
	Name      string
	Title     string
	Icon      string
	Path      string
	Component string
	SortNo    int
	Status    int
}

type seedData struct {
	Roles      []roleSeed
	TopicNodes []topicNodeSeed
	SysConfigs []sysConfigSeed
	Menus      []menuSeed
	RoleMenus  []struct {
		RoleCode   string
		MenuSeedID int64
	}
}

const (
	taskMenuRoot = 100
)

func migrate_init() error {
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		now := dates.NowTimestamp()
		seed := seedForLanguage()

		roleIDMap := make(map[string]int64)
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
				roleIDMap[r.Code] = existing.Id
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
			roleIDMap[r.Code] = role.Id
		}

		topicNodeIDMap := make(map[int64]int64)
		for _, n := range seed.TopicNodes {
			existing := repositories.TopicNodeRepository.Take(tx, "name = ?", n.Name)
			if existing != nil {
				existing.Description = n.Description
				existing.Logo = n.Logo
				existing.SortNo = n.SortNo
				existing.Status = n.Status
				if err := repositories.TopicNodeRepository.Update(tx, existing); err != nil {
					return err
				}
				topicNodeIDMap[n.ID] = existing.Id
				continue
			}

			node := &models.TopicNode{
				Model:       models.Model{Id: n.ID},
				Name:        n.Name,
				Description: n.Description,
				Logo:        n.Logo,
				SortNo:      n.SortNo,
				Status:      n.Status,
				CreateTime:  now,
			}
			if err := repositories.TopicNodeRepository.Create(tx, node); err != nil {
				return err
			}
			topicNodeIDMap[n.ID] = node.Id
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

		menuIDMap := make(map[int64]int64)
		for _, m := range seed.Menus {
			existing := repositories.MenuRepository.Get(tx, m.ID)
			if existing != nil {
				existing.ParentId = m.ParentID
				existing.Type = m.Type
				existing.Name = m.Name
				existing.Title = m.Title
				existing.Icon = m.Icon
				existing.Path = m.Path
				existing.Component = m.Component
				existing.SortNo = m.SortNo
				existing.Status = m.Status
				existing.UpdateTime = now
				if err := repositories.MenuRepository.Update(tx, existing); err != nil {
					return err
				}
				menuIDMap[m.ID] = existing.Id
				continue
			}

			menu := &models.Menu{
				Model:      models.Model{Id: m.ID},
				ParentId:   m.ParentID,
				Type:       m.Type,
				Name:       m.Name,
				Title:      m.Title,
				Icon:       m.Icon,
				Path:       m.Path,
				Component:  m.Component,
				SortNo:     m.SortNo,
				Status:     m.Status,
				CreateTime: now,
				UpdateTime: now,
			}
			if err := repositories.MenuRepository.Create(tx, menu); err != nil {
				return err
			}
			menuIDMap[m.ID] = menu.Id
		}

		for _, rm := range seed.RoleMenus {
			roleID := roleIDMap[rm.RoleCode]
			menuID := menuIDMap[rm.MenuSeedID]
			if roleID == 0 || menuID == 0 {
				continue
			}
			if repositories.RoleMenuRepository.Take(tx, "role_id = ? AND menu_id = ?", roleID, menuID) != nil {
				continue
			}
			roleMenu := &models.RoleMenu{
				RoleId:     roleID,
				MenuId:     menuID,
				CreateTime: now,
			}
			if err := repositories.RoleMenuRepository.Create(tx, roleMenu); err != nil {
				return err
			}
		}

		// ensure defaultNodeId sys config points to created node id
		if nodeID := topicNodeIDMap[1]; nodeID > 0 {
			if cfg := repositories.SysConfigRepository.GetByKey(tx, constants.SysConfigDefaultNodeId); cfg != nil {
				cfg.Value = strconv.FormatInt(nodeID, 10)
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
				{ID: 1, Type: 0, Name: "Owner", Code: constants.RoleOwner, SortNo: 0, Remark: "Owner with highest privileges", Status: constants.StatusOk},
				{ID: 2, Type: 0, Name: "Admin", Code: constants.RoleAdmin, SortNo: 1, Remark: "Admin with management privileges", Status: constants.StatusOk},
			},
			TopicNodes: []topicNodeSeed{
				{ID: 1, Name: "Default", Description: "", Logo: "", SortNo: 0, Status: constants.StatusOk},
			},
			SysConfigs: []sysConfigSeed{
				{Key: constants.SysConfigSiteTitle, Value: "BBS-GO Demo Site", Name: "Site Title", Description: "Site Title"},
				{Key: constants.SysConfigSiteDescription, Value: "BBS-GO, an open source community system based on Go language", Name: "Site Description", Description: "Site Description"},
				{Key: constants.SysConfigSiteKeywords, Value: []string{"bbs-go"}, Name: "Site Keywords", Description: "Site Keywords"},
				{Key: constants.SysConfigSiteNavs, Value: []map[string]string{
					{"title": "Topics", "url": "/topics"},
					{"title": "Articles", "url": "/articles"},
					{"title": "Tasks", "url": "/tasks"},
				}, Name: "Site Navigation", Description: "Site Navigation"},
				{Key: constants.SysConfigDefaultNodeId, Value: "1", Name: "Default Node", Description: "Default Node"},
				{Key: constants.SysConfigTokenExpireDays, Value: "365", Name: "User Login Validity Period (Days)", Description: "User Login Validity Period (Days)"},
				{Key: constants.SysConfigUrlRedirect, Value: "false"},
				{Key: constants.SysConfigEnableHideContent, Value: "false"},
				{Key: constants.SysConfigSiteLogo, Value: ""},
				{Key: constants.SysConfigSiteNotification, Value: ""},
				{Key: constants.SysConfigRecommendTags, Value: ""},
				{Key: constants.SysConfigModules, Value: map[string]bool{"tweet": true, "topic": true, "article": true}},
				{Key: constants.SysConfigUploadConfig, Value: dto.UploadConfig{
					EnableUploadMethod: dto.Local,
					AliyunOss:          dto.AliyunOssUploadConfig{},
					TencentCos:         dto.TencentCosUploadConfig{},
					AwsS3:              dto.AwsS3UploadConfig{},
				}},
			},
			Menus:     append(defaultMenusEn(), taskMenusEn()...),
			RoleMenus: defaultRoleMenus(),
		}
	}

	return seedData{
		Roles: []roleSeed{
			{ID: 1, Type: 0, Name: "超级管理员", Code: constants.RoleOwner, SortNo: 0, Remark: "超级管理员，拥有最高权限", Status: constants.StatusOk},
			{ID: 2, Type: 0, Name: "管理员", Code: constants.RoleAdmin, SortNo: 1, Remark: "管理员，拥有管理权限", Status: constants.StatusOk},
		},
		TopicNodes: []topicNodeSeed{
			{ID: 1, Name: "默认节点", Description: "", Logo: "", SortNo: 0, Status: constants.StatusOk},
		},
		SysConfigs: []sysConfigSeed{
			{Key: constants.SysConfigSiteTitle, Value: "bbs-go演示站", Name: "站点标题", Description: "站点标题"},
			{Key: constants.SysConfigSiteDescription, Value: "bbs-go，基于Go语言的开源社区系统", Name: "站点描述", Description: "站点描述"},
			{Key: constants.SysConfigSiteKeywords, Value: []string{"bbs-go"}, Name: "站点关键字", Description: "站点关键字"},
			{Key: constants.SysConfigSiteNavs, Value: []map[string]string{
				{"title": "话题", "url": "/topics"},
				{"title": "文章", "url": "/articles"},
				{"title": "任务", "url": "/tasks"},
			}, Name: "站点导航", Description: "站点导航"},
			{Key: constants.SysConfigDefaultNodeId, Value: "1", Name: "默认节点", Description: "默认节点"},
			{Key: constants.SysConfigTokenExpireDays, Value: "365", Name: "用户登录有效期(天)", Description: "用户登录有效期(天)"},
			{Key: constants.SysConfigUrlRedirect, Value: "false"},
			{Key: constants.SysConfigEnableHideContent, Value: "false"},
			{Key: constants.SysConfigSiteLogo, Value: ""},
			{Key: constants.SysConfigSiteNotification, Value: ""},
			{Key: constants.SysConfigRecommendTags, Value: ""},
			{Key: constants.SysConfigModules, Value: map[string]bool{"tweet": true, "topic": true, "article": true}},
			{Key: constants.SysConfigUploadConfig, Value: dto.UploadConfig{
				EnableUploadMethod: dto.Local,
				AliyunOss: dto.AliyunOssUploadConfig{
					StyleSplitter: "@",
				},
				TencentCos: dto.TencentCosUploadConfig{},
				AwsS3:      dto.AwsS3UploadConfig{},
			}},
		},
		Menus:     append(defaultMenusZh(), taskMenusZh()...),
		RoleMenus: defaultRoleMenus(),
	}
}

func defaultRoleMenus() []struct {
	RoleCode   string
	MenuSeedID int64
} {
	return []struct {
		RoleCode   string
		MenuSeedID int64
	}{
		{constants.RoleOwner, 1},
		{constants.RoleOwner, 2},
		{constants.RoleOwner, 4},
		{constants.RoleOwner, 9},
		{constants.RoleOwner, 5},
		{constants.RoleOwner, 11},
		{constants.RoleOwner, 10},
		{constants.RoleOwner, 7},
		{constants.RoleOwner, 17},
		{constants.RoleOwner, 16},
		{constants.RoleOwner, 6},
		{constants.RoleOwner, 12},
		{constants.RoleOwner, 13},
		{constants.RoleOwner, 14},
		{constants.RoleOwner, 8},
		{constants.RoleOwner, 15},
		{constants.RoleOwner, 18},
		{constants.RoleOwner, 19},
		{constants.RoleOwner, 20},
		{constants.RoleOwner, 21},
		{constants.RoleOwner, 24},
		{constants.RoleOwner, 22},
		{constants.RoleOwner, 23},
		{constants.RoleOwner, 25},
		{constants.RoleOwner, 26},
		{constants.RoleOwner, taskMenuRoot},
		{constants.RoleOwner, 101},
		{constants.RoleOwner, 102},
		{constants.RoleOwner, 103},
		{constants.RoleOwner, 104},
		{constants.RoleOwner, 105},
		{constants.RoleOwner, 106},
		{constants.RoleAdmin, taskMenuRoot},
		{constants.RoleAdmin, 101},
		{constants.RoleAdmin, 102},
		{constants.RoleAdmin, 103},
		{constants.RoleAdmin, 104},
		{constants.RoleAdmin, 105},
		{constants.RoleAdmin, 106},
	}
}

func defaultMenusZh() []menuSeed {
	return []menuSeed{
		{ID: 1, ParentID: 0, Type: constants.MenuTypeMenu, Name: "Dashboard", Title: "仪表盘", Icon: "icon-dashboard", Path: "/dashboard", Component: "dashboard/index", SortNo: 0, Status: constants.StatusOk},
		{ID: 2, ParentID: 0, Type: constants.MenuTypeMenu, Name: "User", Title: "用户管理", Icon: "icon-user", Path: "/user", Component: "user/index", SortNo: 1, Status: constants.StatusOk},
		{ID: 4, ParentID: 0, Type: constants.MenuTypeMenu, Name: "Permission", Title: "权限管理", Icon: "icon-lock", Path: "", Component: "", SortNo: 9, Status: constants.StatusOk},
		{ID: 5, ParentID: 4, Type: constants.MenuTypeMenu, Name: "Role", Title: "角色管理", Icon: "", Path: "/permission/role", Component: "system/role/index", SortNo: 10, Status: constants.StatusOk},
		{ID: 6, ParentID: 4, Type: constants.MenuTypeMenu, Name: "Menu", Title: "菜单管理", Icon: "", Path: "/permission/menu", Component: "system/menu/index", SortNo: 16, Status: constants.StatusOk},
		{ID: 7, ParentID: 4, Type: constants.MenuTypeMenu, Name: "Api", Title: "接口管理", Icon: "", Path: "/permission/api", Component: "system/api/index", SortNo: 13, Status: constants.StatusOk},
		{ID: 8, ParentID: 4, Type: constants.MenuTypeMenu, Name: "Permission", Title: "权限分配", Icon: "", Path: "/permission/index", Component: "system/permission/index", SortNo: 20, Status: constants.StatusOk},
		{ID: 9, ParentID: 2, Type: constants.MenuTypeFunc, Name: "", Title: "编辑", Icon: "", Path: "", Component: "", SortNo: 2, Status: constants.StatusOk},
		{ID: 10, ParentID: 5, Type: constants.MenuTypeFunc, Name: "", Title: "新增", Icon: "", Path: "", Component: "", SortNo: 12, Status: constants.StatusOk},
		{ID: 11, ParentID: 5, Type: constants.MenuTypeFunc, Name: "", Title: "编辑", Icon: "", Path: "", Component: "", SortNo: 11, Status: constants.StatusOk},
		{ID: 12, ParentID: 6, Type: constants.MenuTypeFunc, Name: "", Title: "新增", Icon: "", Path: "", Component: "", SortNo: 17, Status: constants.StatusOk},
		{ID: 13, ParentID: 6, Type: constants.MenuTypeFunc, Name: "", Title: "编辑", Icon: "", Path: "", Component: "", SortNo: 18, Status: constants.StatusOk},
		{ID: 14, ParentID: 6, Type: constants.MenuTypeFunc, Name: "", Title: "排序", Icon: "", Path: "", Component: "", SortNo: 19, Status: constants.StatusOk},
		{ID: 15, ParentID: 8, Type: constants.MenuTypeFunc, Name: "", Title: "保存", Icon: "", Path: "", Component: "", SortNo: 21, Status: constants.StatusOk},
		{ID: 16, ParentID: 7, Type: constants.MenuTypeFunc, Name: "", Title: "新增", Icon: "icon-settings", Path: "", Component: "", SortNo: 15, Status: constants.StatusOk},
		{ID: 17, ParentID: 7, Type: constants.MenuTypeFunc, Name: "", Title: "编辑", Icon: "icon-settings", Path: "", Component: "", SortNo: 14, Status: constants.StatusOk},
		{ID: 18, ParentID: 0, Type: constants.MenuTypeMenu, Name: "System", Title: "系统管理", Icon: "icon-settings", Path: "", Component: "", SortNo: 22, Status: constants.StatusOk},
		{ID: 19, ParentID: 18, Type: constants.MenuTypeMenu, Name: "Settings", Title: "系统设置", Icon: "", Path: "/system/settings", Component: "settings/index", SortNo: 23, Status: constants.StatusOk},
		{ID: 20, ParentID: 18, Type: constants.MenuTypeMenu, Name: "Dict", Title: "字典管理", Icon: "", Path: "/system/dict", Component: "system/dict/index", SortNo: 24, Status: constants.StatusOk},
		{ID: 21, ParentID: 0, Type: constants.MenuTypeMenu, Name: "Article", Title: "文章管理", Icon: "icon-file", Path: "/article", Component: "article/index", SortNo: 6, Status: constants.StatusOk},
		{ID: 22, ParentID: 0, Type: constants.MenuTypeMenu, Name: "Forbidden-word", Title: "违禁词", Icon: "icon-safe", Path: "/forbidden-word", Component: "forbidden-word/index", SortNo: 7, Status: constants.StatusOk},
		{ID: 23, ParentID: 0, Type: constants.MenuTypeMenu, Name: "Link", Title: "友情链接", Icon: "icon-link", Path: "/link", Component: "link/index", SortNo: 8, Status: constants.StatusOk},
		{ID: 24, ParentID: 0, Type: constants.MenuTypeMenu, Name: "", Title: "帖子管理", Icon: "icon-share-alt", Path: "", Component: "", SortNo: 3, Status: constants.StatusOk},
		{ID: 25, ParentID: 24, Type: constants.MenuTypeMenu, Name: "Topic", Title: "帖子管理", Icon: "", Path: "/topic", Component: "topic/index", SortNo: 4, Status: constants.StatusOk},
		{ID: 26, ParentID: 24, Type: constants.MenuTypeMenu, Name: "TopicNode", Title: "节点管理", Icon: "", Path: "/topic-node", Component: "topic-node/index", SortNo: 5, Status: constants.StatusOk},
	}
}

func defaultMenusEn() []menuSeed {
	return []menuSeed{
		{ID: 1, ParentID: 0, Type: constants.MenuTypeMenu, Name: "Dashboard", Title: "Dashboard", Icon: "icon-dashboard", Path: "/dashboard", Component: "dashboard/index", SortNo: 0, Status: constants.StatusOk},
		{ID: 2, ParentID: 0, Type: constants.MenuTypeMenu, Name: "User", Title: "User", Icon: "icon-user", Path: "/user", Component: "user/index", SortNo: 1, Status: constants.StatusOk},
		{ID: 4, ParentID: 0, Type: constants.MenuTypeMenu, Name: "Permission", Title: "Permission", Icon: "icon-lock", Path: "", Component: "", SortNo: 9, Status: constants.StatusOk},
		{ID: 5, ParentID: 4, Type: constants.MenuTypeMenu, Name: "Role", Title: "Role", Icon: "", Path: "/permission/role", Component: "system/role/index", SortNo: 10, Status: constants.StatusOk},
		{ID: 6, ParentID: 4, Type: constants.MenuTypeMenu, Name: "Menu", Title: "Menu", Icon: "", Path: "/permission/menu", Component: "system/menu/index", SortNo: 16, Status: constants.StatusOk},
		{ID: 7, ParentID: 4, Type: constants.MenuTypeMenu, Name: "Api", Title: "API", Icon: "", Path: "/permission/api", Component: "system/api/index", SortNo: 13, Status: constants.StatusOk},
		{ID: 8, ParentID: 4, Type: constants.MenuTypeMenu, Name: "Permission", Title: "Permission", Icon: "", Path: "/permission/index", Component: "system/permission/index", SortNo: 20, Status: constants.StatusOk},
		{ID: 9, ParentID: 2, Type: constants.MenuTypeFunc, Name: "", Title: "Edit", Icon: "", Path: "", Component: "", SortNo: 2, Status: constants.StatusOk},
		{ID: 10, ParentID: 5, Type: constants.MenuTypeFunc, Name: "", Title: "Add", Icon: "", Path: "", Component: "", SortNo: 12, Status: constants.StatusOk},
		{ID: 11, ParentID: 5, Type: constants.MenuTypeFunc, Name: "", Title: "Edit", Icon: "", Path: "", Component: "", SortNo: 11, Status: constants.StatusOk},
		{ID: 12, ParentID: 6, Type: constants.MenuTypeFunc, Name: "", Title: "Add", Icon: "", Path: "", Component: "", SortNo: 17, Status: constants.StatusOk},
		{ID: 13, ParentID: 6, Type: constants.MenuTypeFunc, Name: "", Title: "Edit", Icon: "", Path: "", Component: "", SortNo: 18, Status: constants.StatusOk},
		{ID: 14, ParentID: 6, Type: constants.MenuTypeFunc, Name: "", Title: "Sort", Icon: "", Path: "", Component: "", SortNo: 19, Status: constants.StatusOk},
		{ID: 15, ParentID: 8, Type: constants.MenuTypeFunc, Name: "", Title: "Save", Icon: "", Path: "", Component: "", SortNo: 21, Status: constants.StatusOk},
		{ID: 16, ParentID: 7, Type: constants.MenuTypeFunc, Name: "", Title: "Add", Icon: "icon-settings", Path: "", Component: "", SortNo: 15, Status: constants.StatusOk},
		{ID: 17, ParentID: 7, Type: constants.MenuTypeFunc, Name: "", Title: "Edit", Icon: "icon-settings", Path: "", Component: "", SortNo: 14, Status: constants.StatusOk},
		{ID: 18, ParentID: 0, Type: constants.MenuTypeMenu, Name: "System", Title: "System", Icon: "icon-settings", Path: "", Component: "", SortNo: 22, Status: constants.StatusOk},
		{ID: 19, ParentID: 18, Type: constants.MenuTypeMenu, Name: "Settings", Title: "Settings", Icon: "", Path: "/system/settings", Component: "settings/index", SortNo: 23, Status: constants.StatusOk},
		{ID: 20, ParentID: 18, Type: constants.MenuTypeMenu, Name: "Dict", Title: "Dictionary", Icon: "", Path: "/system/dict", Component: "system/dict/index", SortNo: 24, Status: constants.StatusOk},
		{ID: 21, ParentID: 0, Type: constants.MenuTypeMenu, Name: "Article", Title: "Article", Icon: "icon-file", Path: "/article", Component: "article/index", SortNo: 6, Status: constants.StatusOk},
		{ID: 22, ParentID: 0, Type: constants.MenuTypeMenu, Name: "Forbidden-word", Title: "Forbidden Words", Icon: "icon-safe", Path: "/forbidden-word", Component: "forbidden-word/index", SortNo: 7, Status: constants.StatusOk},
		{ID: 23, ParentID: 0, Type: constants.MenuTypeMenu, Name: "Link", Title: "Friend Links", Icon: "icon-link", Path: "/link", Component: "link/index", SortNo: 8, Status: constants.StatusOk},
		{ID: 24, ParentID: 0, Type: constants.MenuTypeMenu, Name: "", Title: "Topic", Icon: "icon-share-alt", Path: "", Component: "", SortNo: 3, Status: constants.StatusOk},
		{ID: 25, ParentID: 24, Type: constants.MenuTypeMenu, Name: "Topic", Title: "Topic", Icon: "", Path: "/topic", Component: "topic/index", SortNo: 4, Status: constants.StatusOk},
		{ID: 26, ParentID: 24, Type: constants.MenuTypeMenu, Name: "TopicNode", Title: "Node", Icon: "", Path: "/topic-node", Component: "topic-node/index", SortNo: 5, Status: constants.StatusOk},
	}
}

func taskMenusZh() []menuSeed {
	return []menuSeed{
		{ID: taskMenuRoot, ParentID: 0, Type: constants.MenuTypeMenu, Name: "Task", Title: "任务体系", Icon: "icon-user-group", Path: "", Component: "", SortNo: 7, Status: constants.StatusOk},
		{ID: 101, ParentID: taskMenuRoot, Type: constants.MenuTypeMenu, Name: "TaskConfig", Title: "任务管理", Icon: "", Path: "/task-config", Component: "task-config/index", SortNo: 11, Status: constants.StatusOk},
		{ID: 102, ParentID: taskMenuRoot, Type: constants.MenuTypeMenu, Name: "Badge", Title: "勋章管理", Icon: "", Path: "/badge", Component: "badge/index", SortNo: 12, Status: constants.StatusOk},
		{ID: 103, ParentID: taskMenuRoot, Type: constants.MenuTypeMenu, Name: "LevelConfig", Title: "等级配置", Icon: "", Path: "/level-config", Component: "level-config/index", SortNo: 13, Status: constants.StatusOk},
		{ID: 104, ParentID: taskMenuRoot, Type: constants.MenuTypeMenu, Name: "UserTaskLog", Title: "任务流水", Icon: "", Path: "/user-task-log", Component: "user-task-log/index", SortNo: 14, Status: constants.StatusOk},
		{ID: 105, ParentID: taskMenuRoot, Type: constants.MenuTypeMenu, Name: "UserExpLog", Title: "经验流水", Icon: "", Path: "/user-exp-log", Component: "user-exp-log/index", SortNo: 15, Status: constants.StatusOk},
		{ID: 106, ParentID: taskMenuRoot, Type: constants.MenuTypeMenu, Name: "UserBadge", Title: "用户勋章", Icon: "", Path: "/user-badge", Component: "user-badge/index", SortNo: 16, Status: constants.StatusOk},
	}
}

func taskMenusEn() []menuSeed {
	return []menuSeed{
		{ID: taskMenuRoot, ParentID: 0, Type: constants.MenuTypeMenu, Name: "Growth", Title: "Growth", Icon: "icon-user-group", Path: "", Component: "", SortNo: 7, Status: constants.StatusOk},
		{ID: 101, ParentID: taskMenuRoot, Type: constants.MenuTypeMenu, Name: "TaskConfig", Title: "Tasks", Icon: "", Path: "/task-config", Component: "task-config/index", SortNo: 11, Status: constants.StatusOk},
		{ID: 102, ParentID: taskMenuRoot, Type: constants.MenuTypeMenu, Name: "Badge", Title: "Badges", Icon: "", Path: "/badge", Component: "badge/index", SortNo: 12, Status: constants.StatusOk},
		{ID: 103, ParentID: taskMenuRoot, Type: constants.MenuTypeMenu, Name: "LevelConfig", Title: "Levels", Icon: "", Path: "/level-config", Component: "level-config/index", SortNo: 13, Status: constants.StatusOk},
		{ID: 104, ParentID: taskMenuRoot, Type: constants.MenuTypeMenu, Name: "UserTaskLog", Title: "Task Logs", Icon: "", Path: "/user-task-log", Component: "user-task-log/index", SortNo: 14, Status: constants.StatusOk},
		{ID: 105, ParentID: taskMenuRoot, Type: constants.MenuTypeMenu, Name: "UserExpLog", Title: "Exp Logs", Icon: "", Path: "/user-exp-log", Component: "user-exp-log/index", SortNo: 15, Status: constants.StatusOk},
		{ID: 106, ParentID: taskMenuRoot, Type: constants.MenuTypeMenu, Name: "UserBadge", Title: "User Badges", Icon: "", Path: "/user-badge", Component: "user-badge/index", SortNo: 16, Status: constants.StatusOk},
	}
}
