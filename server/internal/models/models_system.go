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
	Type       string `gorm:"size:32" json:"type" form:"type"`                        // 类型（menu/func）
	Name       string `gorm:"size:64" json:"name" form:"name"`                        // 名称
	Title      string `gorm:"size:64" json:"title" form:"title"`                      // 标题
	Icon       string `gorm:"size:1024" json:"icon" form:"icon"`                      // ICON
	Path       string `gorm:"size:1024" json:"path" form:"path"`                      // 路径
	Component  string `gorm:"size:256" json:"component" form:"component"`             // 组件
	SortNo     int    `gorm:"not null;default:0" json:"sortNo" form:"sortNo"`         // 排序
	Status     int    `json:"status" form:"status"`                                   // 状态
	CreateTime int64  `gorm:"not null;default:0" json:"createTime" form:"createTime"` // 创建时间
	UpdateTime int64  `gorm:"not null;default:0" json:"updateTime" form:"updateTime"` // 更新时间
}

// MenuApi 菜单和接口的权限关联
type MenuApi struct {
	Model
	MenuId     int64 `gorm:"not null;default:0;uniqueIndex:idx_menu_api" json:"menuId" form:"menuId"` // 菜单ID
	ApiId      int64 `gorm:"not null;default:0;uniqueIndex:idx_menu_api" json:"apiId" form:"apiId"`   // 接口ID
	CreateTime int64 `gorm:"not null;default:0" json:"createTime" form:"createTime"`                  // 创建时间
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

type Api struct {
	Model
	Name       string `gorm:"size:512;unique" json:"name" form:"name"`                // 名称
	Method     string `gorm:"size:16" json:"method" form:"method"`                    // 方法
	Path       string `gorm:"size:512;unique" json:"path" form:"path"`                // 路径
	CreateTime int64  `gorm:"not null;default:0" json:"createTime" form:"createTime"` // 创建时间
	UpdateTime int64  `gorm:"not null;default:0" json:"updateTime" form:"updateTime"` // 更新时间
}

type DictType struct {
	Model
	Name       string `gorm:"size:32" json:"name" form:"name"`
	Code       string `gorm:"size:64;unique" json:"code" form:"code"`
	Status     int    `gorm:"not null;default:0" json:"status" form:"status"`
	Remark     string `gorm:"size:512" json:"remark" form:"remark"`
	CreateTime int64  `gorm:"not null;default:0" json:"createTime" form:"createTime"` // 创建时间
	UpdateTime int64  `gorm:"not null;default:0" json:"updateTime" form:"updateTime"` // 更新时间
}

type Dict struct {
	Model
	TypeId     int64  `gorm:"uniqueIndex:idx_dict_name" json:"typeId" form:"typeId"`     // 分类
	ParentId   int64  `gorm:"default:0" json:"parentId" form:"parentId"`                 // 上级
	Name       string `gorm:"size:64;uniqueIndex:idx_dict_name" json:"name" form:"name"` // 名称
	Label      string `gorm:"size:64" json:"label" form:"label"`                         // Label
	Value      string `gorm:"type:text" json:"value" form:"value"`                       // Value
	SortNo     int    `gorm:"not null;default:0" json:"sortNo" form:"sortNo"`            // 排序
	Status     int    `gorm:"not null;default:0" json:"status" form:"status"`            // 状态
	CreateTime int64  `gorm:"not null;default:0" json:"createTime" form:"createTime"`    // 创建时间
	UpdateTime int64  `gorm:"not null;default:0" json:"updateTime" form:"updateTime"`    // 更新时间
}
