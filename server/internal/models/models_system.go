package models

type Role struct {
	Model
	Type       int    `gorm:"not null;default:1" json:"type" form:"type"`             // 角色类型（0：系统角色、1：自定义角色）
	Name       string `gorm:"size:64" json:"name" form:"name"`                        // 角色名称
	Code       string `gorm:"unique;size:64" json:"code" form:"code"`                 // 角色编码
	SortNo     int    `json:"sortNo" form:"sortNo"`                                   // 排序
	Remark     string `gorm:"size:256" json:"remark" form:"remark"`                   // 备注
	Status     int    `json:"status" form:"status"`                                   // 状态
	CreateTime int64  `gorm:"not null;default:0" json:"createTime" form:"createTime"` // 创建时间
	UpdateTime int64  `gorm:"not null;default:0" json:"updateTime" form:"updateTime"` // 更新时间
}

type Menu struct {
	Model
	ParentId   int64  `json:"parentId" form:"parentId"`                               // 上级菜单
	Name       string `gorm:"size:256" json:"name" form:"name"`                       // 名称
	Title      string `gorm:"size:64" json:"title" form:"title"`                      // 标题
	Icon       string `gorm:"size:1024" json:"icon" form:"icon"`                      // ICON
	Path       string `gorm:"size:1024" json:"path" form:"path"`                      // 路径
	SortNo     int    `gorm:"not null;default:0" json:"sortNo" form:"sortNo"`         // 排序
	Status     int    `json:"status" form:"status"`                                   // 状态
	CreateTime int64  `gorm:"not null;default:0" json:"createTime" form:"createTime"` // 创建时间
	UpdateTime int64  `gorm:"not null;default:0" json:"updateTime" form:"updateTime"` // 更新时间
}

type UserRole struct {
	Model
	UserId     int64 `gorm:"uniqueIndex:idx_user_role" json:"userId" form:"userId"`
	RoleId     int64 `gorm:"uniqueIndex:idx_user_role" json:"roleId" form:"roleId"`
	CreateTime int64 `gorm:"not null;default:0" json:"createTime" form:"createTime"` // 创建时间
}

type RoleMenu struct {
	Model
	RoleId     int64 `gorm:"uniqueIndex:idx_role_menu" json:"roleId" form:"roleId"`
	MenuId     int64 `gorm:"uniqueIndex:idx_role_menu" json:"menuId" form:"menuId"`
	CreateTime int64 `gorm:"not null;default:0" json:"createTime" form:"createTime"` // 创建时间
}
