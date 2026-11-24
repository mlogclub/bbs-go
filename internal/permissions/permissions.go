package permissions

import "strings"

const (
	TypeDashboard = "dashboard"
)

const (
	GroupWorkspace = "workspace"
	GroupContent   = "content"
	GroupCommunity = "community"
	GroupGrowth    = "growth"
	GroupSystem    = "system"
	GroupLogs      = "logs"
)

type PermissionDefinition struct {
	Type        string
	Code        string
	GroupName   string
	SortNo      int
	NameEn      string
	NameZh      string
	NameZhTW    string
	Description string
}

func (p PermissionDefinition) String() string {
	return p.Code
}

func (p PermissionDefinition) IsValid() bool {
	return p.Type != "" &&
		p.Code != "" &&
		p.GroupName != "" &&
		p.SortNo > 0 &&
		p.NameEn != "" &&
		p.NameZh != "" &&
		p.NameZhTW != "" &&
		strings.HasPrefix(p.Code, p.Type+".")
}

var (
	PermissionDashboardView = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.view", GroupName: GroupWorkspace, SortNo: 10, NameEn: "Dashboard Access", NameZh: "进入后台", NameZhTW: "進入後台"}

	PermissionTopicView      = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.topic.view", GroupName: GroupContent, SortNo: 100, NameEn: "View Topics", NameZh: "查看话题", NameZhTW: "檢視話題"}
	PermissionTopicRecommend = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.topic.recommend", GroupName: GroupContent, SortNo: 110, NameEn: "Recommend Topics", NameZh: "推荐话题", NameZhTW: "推薦話題"}
	PermissionTopicSticky    = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.topic.sticky", GroupName: GroupContent, SortNo: 115, NameEn: "Sticky Topics", NameZh: "置顶话题", NameZhTW: "置頂話題"}
	PermissionTopicAudit     = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.topic.audit", GroupName: GroupContent, SortNo: 120, NameEn: "Audit Topics", NameZh: "审核话题", NameZhTW: "稽核話題"}
	PermissionTopicDelete    = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.topic.delete", GroupName: GroupContent, SortNo: 130, NameEn: "Delete Topics", NameZh: "删除话题", NameZhTW: "刪除話題"}
	PermissionTopicSolve     = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.topic.solve", GroupName: GroupContent, SortNo: 140, NameEn: "Solve Topics", NameZh: "标记问答解决", NameZhTW: "標記問答解決"}

	PermissionArticleView   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.article.view", GroupName: GroupContent, SortNo: 200, NameEn: "View Articles", NameZh: "查看文章", NameZhTW: "檢視文章"}
	PermissionArticleUpdate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.article.update", GroupName: GroupContent, SortNo: 210, NameEn: "Update Articles", NameZh: "编辑文章", NameZhTW: "編輯文章"}
	PermissionArticleAudit  = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.article.audit", GroupName: GroupContent, SortNo: 220, NameEn: "Audit Articles", NameZh: "审核文章", NameZhTW: "稽核文章"}
	PermissionArticleDelete = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.article.delete", GroupName: GroupContent, SortNo: 230, NameEn: "Delete Articles", NameZh: "删除文章", NameZhTW: "刪除文章"}
	PermissionArticleTags   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.article.tags", GroupName: GroupContent, SortNo: 240, NameEn: "Update Article Tags", NameZh: "维护文章标签", NameZhTW: "維護文章標籤"}

	PermissionCommentDelete = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.comment.delete", GroupName: GroupContent, SortNo: 310, NameEn: "Delete Comments", NameZh: "删除评论", NameZhTW: "刪除留言"}

	PermissionCategoryView   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.category.view", GroupName: GroupContent, SortNo: 400, NameEn: "View Categories", NameZh: "查看分类", NameZhTW: "檢視分類"}
	PermissionCategoryCreate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.category.create", GroupName: GroupContent, SortNo: 410, NameEn: "Create Categories", NameZh: "创建分类", NameZhTW: "建立分類"}
	PermissionCategoryUpdate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.category.update", GroupName: GroupContent, SortNo: 420, NameEn: "Update Categories", NameZh: "编辑分类", NameZhTW: "編輯分類"}
	PermissionCategoryDelete = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.category.delete", GroupName: GroupContent, SortNo: 430, NameEn: "Delete Categories", NameZh: "删除分类", NameZhTW: "刪除分類"}
	PermissionCategorySort   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.category.sort", GroupName: GroupContent, SortNo: 440, NameEn: "Sort Categories", NameZh: "排序分类", NameZhTW: "排序分類"}

	PermissionLinkView   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.link.view", GroupName: GroupContent, SortNo: 500, NameEn: "View Links", NameZh: "查看链接", NameZhTW: "檢視連結"}
	PermissionLinkCreate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.link.create", GroupName: GroupContent, SortNo: 510, NameEn: "Create Links", NameZh: "创建链接", NameZhTW: "建立連結"}
	PermissionLinkUpdate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.link.update", GroupName: GroupContent, SortNo: 520, NameEn: "Update Links", NameZh: "编辑链接", NameZhTW: "編輯連結"}
	PermissionLinkDelete = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.link.delete", GroupName: GroupContent, SortNo: 530, NameEn: "Delete Links", NameZh: "删除链接", NameZhTW: "刪除連結"}

	PermissionForbiddenWordView   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.forbiddenWord.view", GroupName: GroupContent, SortNo: 600, NameEn: "View Forbidden Words", NameZh: "查看敏感词", NameZhTW: "檢視敏感詞"}
	PermissionForbiddenWordCreate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.forbiddenWord.create", GroupName: GroupContent, SortNo: 610, NameEn: "Create Forbidden Words", NameZh: "创建敏感词", NameZhTW: "建立敏感詞"}
	PermissionForbiddenWordUpdate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.forbiddenWord.update", GroupName: GroupContent, SortNo: 620, NameEn: "Update Forbidden Words", NameZh: "编辑敏感词", NameZhTW: "編輯敏感詞"}
	PermissionForbiddenWordDelete = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.forbiddenWord.delete", GroupName: GroupContent, SortNo: 630, NameEn: "Delete Forbidden Words", NameZh: "删除敏感词", NameZhTW: "刪除敏感詞"}

	PermissionUserView             = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.user.view", GroupName: GroupCommunity, SortNo: 700, NameEn: "View Users", NameZh: "查看用户", NameZhTW: "檢視使用者"}
	PermissionUserCreate           = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.user.create", GroupName: GroupCommunity, SortNo: 710, NameEn: "Create Users", NameZh: "创建用户", NameZhTW: "建立使用者"}
	PermissionUserUpdate           = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.user.update", GroupName: GroupCommunity, SortNo: 720, NameEn: "Update Users", NameZh: "编辑用户", NameZhTW: "編輯使用者"}
	PermissionUserForbidden        = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.user.forbidden", GroupName: GroupCommunity, SortNo: 730, NameEn: "Forbid Users", NameZh: "禁言用户", NameZhTW: "禁言使用者"}
	PermissionUserForbiddenForever = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.user.forbiddenForever", GroupName: GroupCommunity, SortNo: 735, NameEn: "Forbid Users Permanently", NameZh: "永久禁言用户", NameZhTW: "永久禁言使用者"}
	PermissionUserUpdatePassword   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.user.updatePassword", GroupName: GroupCommunity, SortNo: 740, NameEn: "Update Own Password", NameZh: "修改自己的密码", NameZhTW: "修改自己的密碼"}
	PermissionUserResetPassword    = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.user.resetPassword", GroupName: GroupCommunity, SortNo: 750, NameEn: "Reset User Password", NameZh: "重置用户密码", NameZhTW: "重置使用者密碼"}

	PermissionBadgeView   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.badge.view", GroupName: GroupGrowth, SortNo: 800, NameEn: "View Badges", NameZh: "查看徽章", NameZhTW: "檢視徽章"}
	PermissionBadgeCreate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.badge.create", GroupName: GroupGrowth, SortNo: 810, NameEn: "Create Badges", NameZh: "创建徽章", NameZhTW: "建立徽章"}
	PermissionBadgeUpdate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.badge.update", GroupName: GroupGrowth, SortNo: 820, NameEn: "Update Badges", NameZh: "编辑徽章", NameZhTW: "編輯徽章"}
	PermissionBadgeDelete = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.badge.delete", GroupName: GroupGrowth, SortNo: 830, NameEn: "Delete Badges", NameZh: "删除徽章", NameZhTW: "刪除徽章"}

	PermissionLevelView   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.level.view", GroupName: GroupGrowth, SortNo: 900, NameEn: "View Levels", NameZh: "查看等级", NameZhTW: "檢視等級"}
	PermissionLevelUpdate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.level.update", GroupName: GroupGrowth, SortNo: 910, NameEn: "Update Levels", NameZh: "编辑等级", NameZhTW: "編輯等級"}

	PermissionTaskView   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.task.view", GroupName: GroupGrowth, SortNo: 1000, NameEn: "View Tasks", NameZh: "查看任务", NameZhTW: "檢視任務"}
	PermissionTaskCreate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.task.create", GroupName: GroupGrowth, SortNo: 1010, NameEn: "Create Tasks", NameZh: "创建任务", NameZhTW: "建立任務"}
	PermissionTaskUpdate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.task.update", GroupName: GroupGrowth, SortNo: 1020, NameEn: "Update Tasks", NameZh: "编辑任务", NameZhTW: "編輯任務"}
	PermissionTaskDelete = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.task.delete", GroupName: GroupGrowth, SortNo: 1030, NameEn: "Delete Tasks", NameZh: "删除任务", NameZhTW: "刪除任務"}

	PermissionSettingView     = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.setting.view", GroupName: GroupSystem, SortNo: 1100, NameEn: "View Settings", NameZh: "查看设置", NameZhTW: "檢視設定"}
	PermissionSettingUpdate   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.setting.update", GroupName: GroupSystem, SortNo: 1110, NameEn: "Update Settings", NameZh: "编辑设置", NameZhTW: "編輯設定"}
	PermissionSearchReindex   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.search.reindex", GroupName: GroupSystem, SortNo: 1120, NameEn: "Rebuild Search Index", NameZh: "重建搜索索引", NameZhTW: "重建搜尋索引"}
	PermissionSitemapGenerate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.sitemap.generate", GroupName: GroupSystem, SortNo: 1130, NameEn: "Generate Sitemap", NameZh: "生成 Sitemap", NameZhTW: "生成 Sitemap"}

	PermissionRoleView             = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.role.view", GroupName: GroupSystem, SortNo: 1200, NameEn: "View Roles", NameZh: "查看角色", NameZhTW: "檢視角色"}
	PermissionRoleCreate           = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.role.create", GroupName: GroupSystem, SortNo: 1210, NameEn: "Create Roles", NameZh: "创建角色", NameZhTW: "建立角色"}
	PermissionRoleUpdate           = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.role.update", GroupName: GroupSystem, SortNo: 1220, NameEn: "Update Roles", NameZh: "编辑角色", NameZhTW: "編輯角色"}
	PermissionRoleDelete           = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.role.delete", GroupName: GroupSystem, SortNo: 1230, NameEn: "Delete Roles", NameZh: "删除角色", NameZhTW: "刪除角色"}
	PermissionRoleSort             = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.role.sort", GroupName: GroupSystem, SortNo: 1240, NameEn: "Sort Roles", NameZh: "排序角色", NameZhTW: "排序角色"}
	PermissionRolePermissionUpdate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.role.permission.update", GroupName: GroupSystem, SortNo: 1250, NameEn: "Update Role Permissions", NameZh: "编辑角色权限", NameZhTW: "編輯角色權限"}

	PermissionUserBadgeView   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.userBadge.view", GroupName: GroupCommunity, SortNo: 760, NameEn: "View User Badges", NameZh: "查看用户徽章", NameZhTW: "檢視使用者徽章"}
	PermissionUserExpLogView  = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.userExpLog.view", GroupName: GroupCommunity, SortNo: 770, NameEn: "View XP Logs", NameZh: "查看经验日志", NameZhTW: "檢視經驗日誌"}
	PermissionUserTaskLogView = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.userTaskLog.view", GroupName: GroupCommunity, SortNo: 780, NameEn: "View Task Logs", NameZh: "查看任务日志", NameZhTW: "檢視任務日誌"}
	PermissionUserReportView  = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.userReport.view", GroupName: GroupCommunity, SortNo: 790, NameEn: "View User Reports", NameZh: "查看用户举报", NameZhTW: "檢視使用者舉報"}
	PermissionUserReportAudit = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.userReport.audit", GroupName: GroupCommunity, SortNo: 795, NameEn: "Audit User Reports", NameZh: "处理用户举报", NameZhTW: "處理使用者舉報"}

	PermissionEmailLogView   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.emailLog.view", GroupName: GroupSystem, SortNo: 1300, NameEn: "View Email Logs", NameZh: "查看邮件日志", NameZhTW: "檢視郵件日誌"}
	PermissionOperateLogView = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.operateLog.view", GroupName: GroupSystem, SortNo: 1310, NameEn: "View Operation Logs", NameZh: "查看操作日志", NameZhTW: "檢視操作日誌"}
)

var Permissions = []PermissionDefinition{
	PermissionDashboardView,
	PermissionTopicView,
	PermissionTopicRecommend,
	PermissionTopicSticky,
	PermissionTopicAudit,
	PermissionTopicDelete,
	PermissionTopicSolve,
	PermissionArticleView,
	PermissionArticleUpdate,
	PermissionArticleAudit,
	PermissionArticleDelete,
	PermissionArticleTags,
	PermissionCommentDelete,
	PermissionCategoryView,
	PermissionCategoryCreate,
	PermissionCategoryUpdate,
	PermissionCategoryDelete,
	PermissionCategorySort,
	PermissionLinkView,
	PermissionLinkCreate,
	PermissionLinkUpdate,
	PermissionLinkDelete,
	PermissionForbiddenWordView,
	PermissionForbiddenWordCreate,
	PermissionForbiddenWordUpdate,
	PermissionForbiddenWordDelete,
	PermissionUserView,
	PermissionUserCreate,
	PermissionUserUpdate,
	PermissionUserForbidden,
	PermissionUserForbiddenForever,
	PermissionUserUpdatePassword,
	PermissionUserResetPassword,
	PermissionBadgeView,
	PermissionBadgeCreate,
	PermissionBadgeUpdate,
	PermissionBadgeDelete,
	PermissionLevelView,
	PermissionLevelUpdate,
	PermissionTaskView,
	PermissionTaskCreate,
	PermissionTaskUpdate,
	PermissionTaskDelete,
	PermissionRoleView,
	PermissionRoleCreate,
	PermissionRoleUpdate,
	PermissionRoleDelete,
	PermissionRoleSort,
	PermissionRolePermissionUpdate,
	PermissionSettingView,
	PermissionSettingUpdate,
	PermissionSearchReindex,
	PermissionSitemapGenerate,
	PermissionEmailLogView,
	PermissionUserTaskLogView,
	PermissionUserExpLogView,
	PermissionUserBadgeView,
	PermissionUserReportView,
	PermissionUserReportAudit,
	PermissionOperateLogView,
}

func FindByCode(code string) (PermissionDefinition, bool) {
	for _, permission := range Permissions {
		if permission.Code == code {
			return permission, true
		}
	}
	return PermissionDefinition{}, false
}
