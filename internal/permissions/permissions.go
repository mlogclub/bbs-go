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
		strings.HasPrefix(p.Code, p.Type+".")
}

var (
	PermissionDashboardView = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.view", GroupName: GroupWorkspace, SortNo: 10, NameEn: "Dashboard Access", NameZh: "进入后台"}

	PermissionTopicView      = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.topic.view", GroupName: GroupContent, SortNo: 100, NameEn: "View Topics", NameZh: "查看话题"}
	PermissionTopicRecommend = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.topic.recommend", GroupName: GroupContent, SortNo: 110, NameEn: "Recommend Topics", NameZh: "推荐话题"}
	PermissionTopicSticky    = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.topic.sticky", GroupName: GroupContent, SortNo: 115, NameEn: "Sticky Topics", NameZh: "置顶话题"}
	PermissionTopicAudit     = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.topic.audit", GroupName: GroupContent, SortNo: 120, NameEn: "Audit Topics", NameZh: "审核话题"}
	PermissionTopicDelete    = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.topic.delete", GroupName: GroupContent, SortNo: 130, NameEn: "Delete Topics", NameZh: "删除话题"}
	PermissionTopicSolve     = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.topic.solve", GroupName: GroupContent, SortNo: 140, NameEn: "Solve Topics", NameZh: "标记问答解决"}

	PermissionArticleView   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.article.view", GroupName: GroupContent, SortNo: 200, NameEn: "View Articles", NameZh: "查看文章"}
	PermissionArticleUpdate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.article.update", GroupName: GroupContent, SortNo: 210, NameEn: "Update Articles", NameZh: "编辑文章"}
	PermissionArticleAudit  = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.article.audit", GroupName: GroupContent, SortNo: 220, NameEn: "Audit Articles", NameZh: "审核文章"}
	PermissionArticleDelete = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.article.delete", GroupName: GroupContent, SortNo: 230, NameEn: "Delete Articles", NameZh: "删除文章"}
	PermissionArticleTags   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.article.tags", GroupName: GroupContent, SortNo: 240, NameEn: "Update Article Tags", NameZh: "维护文章标签"}

	PermissionCommentDelete = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.comment.delete", GroupName: GroupContent, SortNo: 310, NameEn: "Delete Comments", NameZh: "删除评论"}

	PermissionNodeView   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.node.view", GroupName: GroupContent, SortNo: 400, NameEn: "View Nodes", NameZh: "查看节点"}
	PermissionNodeCreate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.node.create", GroupName: GroupContent, SortNo: 410, NameEn: "Create Nodes", NameZh: "创建节点"}
	PermissionNodeUpdate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.node.update", GroupName: GroupContent, SortNo: 420, NameEn: "Update Nodes", NameZh: "编辑节点"}
	PermissionNodeDelete = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.node.delete", GroupName: GroupContent, SortNo: 430, NameEn: "Delete Nodes", NameZh: "删除节点"}
	PermissionNodeSort   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.node.sort", GroupName: GroupContent, SortNo: 440, NameEn: "Sort Nodes", NameZh: "排序节点"}

	PermissionLinkView   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.link.view", GroupName: GroupContent, SortNo: 500, NameEn: "View Links", NameZh: "查看链接"}
	PermissionLinkCreate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.link.create", GroupName: GroupContent, SortNo: 510, NameEn: "Create Links", NameZh: "创建链接"}
	PermissionLinkUpdate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.link.update", GroupName: GroupContent, SortNo: 520, NameEn: "Update Links", NameZh: "编辑链接"}

	PermissionForbiddenWordView   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.forbiddenWord.view", GroupName: GroupContent, SortNo: 600, NameEn: "View Forbidden Words", NameZh: "查看敏感词"}
	PermissionForbiddenWordCreate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.forbiddenWord.create", GroupName: GroupContent, SortNo: 610, NameEn: "Create Forbidden Words", NameZh: "创建敏感词"}
	PermissionForbiddenWordUpdate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.forbiddenWord.update", GroupName: GroupContent, SortNo: 620, NameEn: "Update Forbidden Words", NameZh: "编辑敏感词"}

	PermissionUserView           = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.user.view", GroupName: GroupCommunity, SortNo: 700, NameEn: "View Users", NameZh: "查看用户"}
	PermissionUserCreate         = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.user.create", GroupName: GroupCommunity, SortNo: 710, NameEn: "Create Users", NameZh: "创建用户"}
	PermissionUserUpdate         = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.user.update", GroupName: GroupCommunity, SortNo: 720, NameEn: "Update Users", NameZh: "编辑用户"}
	PermissionUserForbidden      = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.user.forbidden", GroupName: GroupCommunity, SortNo: 730, NameEn: "Forbid Users", NameZh: "禁言用户"}
	PermissionUserUpdatePassword = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.user.updatePassword", GroupName: GroupCommunity, SortNo: 740, NameEn: "Update Own Password", NameZh: "修改自己的密码"}
	PermissionUserResetPassword  = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.user.resetPassword", GroupName: GroupCommunity, SortNo: 750, NameEn: "Reset User Password", NameZh: "重置用户密码"}

	PermissionBadgeView   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.badge.view", GroupName: GroupGrowth, SortNo: 800, NameEn: "View Badges", NameZh: "查看徽章"}
	PermissionBadgeCreate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.badge.create", GroupName: GroupGrowth, SortNo: 810, NameEn: "Create Badges", NameZh: "创建徽章"}
	PermissionBadgeUpdate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.badge.update", GroupName: GroupGrowth, SortNo: 820, NameEn: "Update Badges", NameZh: "编辑徽章"}
	PermissionBadgeDelete = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.badge.delete", GroupName: GroupGrowth, SortNo: 830, NameEn: "Delete Badges", NameZh: "删除徽章"}

	PermissionLevelView   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.level.view", GroupName: GroupGrowth, SortNo: 900, NameEn: "View Levels", NameZh: "查看等级"}
	PermissionLevelUpdate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.level.update", GroupName: GroupGrowth, SortNo: 910, NameEn: "Update Levels", NameZh: "编辑等级"}

	PermissionTaskView   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.task.view", GroupName: GroupGrowth, SortNo: 1000, NameEn: "View Tasks", NameZh: "查看任务"}
	PermissionTaskCreate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.task.create", GroupName: GroupGrowth, SortNo: 1010, NameEn: "Create Tasks", NameZh: "创建任务"}
	PermissionTaskUpdate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.task.update", GroupName: GroupGrowth, SortNo: 1020, NameEn: "Update Tasks", NameZh: "编辑任务"}
	PermissionTaskDelete = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.task.delete", GroupName: GroupGrowth, SortNo: 1030, NameEn: "Delete Tasks", NameZh: "删除任务"}

	PermissionSettingView   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.setting.view", GroupName: GroupSystem, SortNo: 1100, NameEn: "View Settings", NameZh: "查看设置"}
	PermissionSettingUpdate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.setting.update", GroupName: GroupSystem, SortNo: 1110, NameEn: "Update Settings", NameZh: "编辑设置"}

	PermissionRoleView             = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.role.view", GroupName: GroupSystem, SortNo: 1200, NameEn: "View Roles", NameZh: "查看角色"}
	PermissionRoleCreate           = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.role.create", GroupName: GroupSystem, SortNo: 1210, NameEn: "Create Roles", NameZh: "创建角色"}
	PermissionRoleUpdate           = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.role.update", GroupName: GroupSystem, SortNo: 1220, NameEn: "Update Roles", NameZh: "编辑角色"}
	PermissionRoleDelete           = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.role.delete", GroupName: GroupSystem, SortNo: 1230, NameEn: "Delete Roles", NameZh: "删除角色"}
	PermissionRoleSort             = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.role.sort", GroupName: GroupSystem, SortNo: 1240, NameEn: "Sort Roles", NameZh: "排序角色"}
	PermissionRolePermissionUpdate = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.role.permission.update", GroupName: GroupSystem, SortNo: 1250, NameEn: "Update Role Permissions", NameZh: "编辑角色权限"}

	PermissionUserBadgeView   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.userBadge.view", GroupName: GroupCommunity, SortNo: 760, NameEn: "View User Badges", NameZh: "查看用户徽章"}
	PermissionUserExpLogView  = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.userExpLog.view", GroupName: GroupCommunity, SortNo: 770, NameEn: "View XP Logs", NameZh: "查看经验日志"}
	PermissionUserTaskLogView = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.userTaskLog.view", GroupName: GroupCommunity, SortNo: 780, NameEn: "View Task Logs", NameZh: "查看任务日志"}
	PermissionUserReportView  = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.userReport.view", GroupName: GroupCommunity, SortNo: 790, NameEn: "View User Reports", NameZh: "查看用户举报"}

	PermissionEmailLogView   = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.emailLog.view", GroupName: GroupSystem, SortNo: 1300, NameEn: "View Email Logs", NameZh: "查看邮件日志"}
	PermissionOperateLogView = PermissionDefinition{Type: TypeDashboard, Code: "dashboard.operateLog.view", GroupName: GroupSystem, SortNo: 1310, NameEn: "View Operation Logs", NameZh: "查看操作日志"}
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
	PermissionNodeView,
	PermissionNodeCreate,
	PermissionNodeUpdate,
	PermissionNodeDelete,
	PermissionNodeSort,
	PermissionLinkView,
	PermissionLinkCreate,
	PermissionLinkUpdate,
	PermissionForbiddenWordView,
	PermissionForbiddenWordCreate,
	PermissionForbiddenWordUpdate,
	PermissionUserView,
	PermissionUserCreate,
	PermissionUserUpdate,
	PermissionUserForbidden,
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
	PermissionEmailLogView,
	PermissionUserTaskLogView,
	PermissionUserExpLogView,
	PermissionUserBadgeView,
	PermissionUserReportView,
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
