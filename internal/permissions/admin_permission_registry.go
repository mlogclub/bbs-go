package permissions

import (
	"bbs-go/internal/pkg/urls"
	"strings"
)

type adminPermissionRule struct {
	Method      string
	Pattern     string
	Permissions []PermissionDefinition
}

var (
	adminPathMatcher = urls.NewAntPathMatcher()

	adminPermissionRules = []adminPermissionRule{
		{Method: "GET", Pattern: "/api/admin/common/**", Permissions: []PermissionDefinition{PermissionDashboardView}},

		{Method: "GET", Pattern: "/api/admin/topic/*", Permissions: []PermissionDefinition{PermissionTopicView}},
		{Method: "POST", Pattern: "/api/admin/topic/list", Permissions: []PermissionDefinition{PermissionTopicView}},
		{Method: "POST", Pattern: "/api/admin/topic/recommend", Permissions: []PermissionDefinition{PermissionTopicRecommend}},
		{Method: "POST", Pattern: "/api/admin/topic/audit", Permissions: []PermissionDefinition{PermissionTopicAudit}},
		{Method: "POST", Pattern: "/api/admin/topic/delete", Permissions: []PermissionDefinition{PermissionTopicDelete}},
		{Method: "POST", Pattern: "/api/admin/topic/undelete", Permissions: []PermissionDefinition{PermissionTopicDelete}},
		{Method: "POST", Pattern: "/api/admin/topic/mark_solved", Permissions: []PermissionDefinition{PermissionTopicSolve}},
		{Method: "POST", Pattern: "/api/admin/topic/mark_unsolved", Permissions: []PermissionDefinition{PermissionTopicSolve}},

		{Method: "GET", Pattern: "/api/admin/article/*", Permissions: []PermissionDefinition{PermissionArticleView}},
		{Method: "POST", Pattern: "/api/admin/article/list", Permissions: []PermissionDefinition{PermissionArticleView}},
		{Method: "GET", Pattern: "/api/admin/article/tags", Permissions: []PermissionDefinition{PermissionArticleView}},
		{Method: "POST", Pattern: "/api/admin/article/update", Permissions: []PermissionDefinition{PermissionArticleUpdate}},
		{Method: "POST", Pattern: "/api/admin/article/audit", Permissions: []PermissionDefinition{PermissionArticleAudit}},
		{Method: "POST", Pattern: "/api/admin/article/delete", Permissions: []PermissionDefinition{PermissionArticleDelete}},
		{Method: "POST", Pattern: "/api/admin/article/tags", Permissions: []PermissionDefinition{PermissionArticleTags}},

		{Method: "GET", Pattern: "/api/admin/category/*", Permissions: []PermissionDefinition{PermissionCategoryView}},
		{Method: "GET", Pattern: "/api/admin/category/options", Permissions: []PermissionDefinition{PermissionCategoryView}},
		{Method: "POST", Pattern: "/api/admin/category/list", Permissions: []PermissionDefinition{PermissionCategoryView}},
		{Method: "POST", Pattern: "/api/admin/category/create", Permissions: []PermissionDefinition{PermissionCategoryCreate}},
		{Method: "POST", Pattern: "/api/admin/category/update", Permissions: []PermissionDefinition{PermissionCategoryUpdate}},
		{Method: "POST", Pattern: "/api/admin/category/delete", Permissions: []PermissionDefinition{PermissionCategoryDelete}},
		{Method: "POST", Pattern: "/api/admin/category/update_sort", Permissions: []PermissionDefinition{PermissionCategorySort}},

		{Method: "GET", Pattern: "/api/admin/link/*", Permissions: []PermissionDefinition{PermissionLinkView}},
		{Method: "POST", Pattern: "/api/admin/link/list", Permissions: []PermissionDefinition{PermissionLinkView}},
		{Method: "POST", Pattern: "/api/admin/link/create", Permissions: []PermissionDefinition{PermissionLinkCreate}},
		{Method: "POST", Pattern: "/api/admin/link/update", Permissions: []PermissionDefinition{PermissionLinkUpdate}},
		{Method: "POST", Pattern: "/api/admin/link/delete", Permissions: []PermissionDefinition{PermissionLinkDelete}},

		{Method: "GET", Pattern: "/api/admin/forbidden-word/*", Permissions: []PermissionDefinition{PermissionForbiddenWordView}},
		{Method: "POST", Pattern: "/api/admin/forbidden-word/list", Permissions: []PermissionDefinition{PermissionForbiddenWordView}},
		{Method: "POST", Pattern: "/api/admin/forbidden-word/create", Permissions: []PermissionDefinition{PermissionForbiddenWordCreate}},
		{Method: "POST", Pattern: "/api/admin/forbidden-word/update", Permissions: []PermissionDefinition{PermissionForbiddenWordUpdate}},
		{Method: "POST", Pattern: "/api/admin/forbidden-word/delete", Permissions: []PermissionDefinition{PermissionForbiddenWordDelete}},

		{Method: "GET", Pattern: "/api/admin/user/*", Permissions: []PermissionDefinition{PermissionUserView}},
		{Method: "GET", Pattern: "/api/admin/user/synccount", Permissions: []PermissionDefinition{PermissionUserUpdate}},
		{Method: "POST", Pattern: "/api/admin/user/list", Permissions: []PermissionDefinition{PermissionUserView}},
		{Method: "POST", Pattern: "/api/admin/user/create", Permissions: []PermissionDefinition{PermissionUserCreate}},
		{Method: "POST", Pattern: "/api/admin/user/update", Permissions: []PermissionDefinition{PermissionUserUpdate}},
		{Method: "POST", Pattern: "/api/admin/user/forbidden", Permissions: []PermissionDefinition{PermissionUserForbidden, PermissionUserForbiddenForever}},
		{Method: "POST", Pattern: "/api/admin/user/update_password", Permissions: []PermissionDefinition{PermissionUserUpdatePassword}},
		{Method: "POST", Pattern: "/api/admin/user/reset_password", Permissions: []PermissionDefinition{PermissionUserResetPassword}},

		{Method: "GET", Pattern: "/api/admin/role/roles", Permissions: []PermissionDefinition{PermissionRoleView, PermissionUserUpdate}},
		{Method: "GET", Pattern: "/api/admin/role/*", Permissions: []PermissionDefinition{PermissionRoleView}},
		{Method: "POST", Pattern: "/api/admin/role/list", Permissions: []PermissionDefinition{PermissionRoleView}},
		{Method: "POST", Pattern: "/api/admin/role/create", Permissions: []PermissionDefinition{PermissionRoleCreate}},
		{Method: "POST", Pattern: "/api/admin/role/update", Permissions: []PermissionDefinition{PermissionRoleUpdate}},
		{Method: "POST", Pattern: "/api/admin/role/delete", Permissions: []PermissionDefinition{PermissionRoleDelete}},
		{Method: "POST", Pattern: "/api/admin/role/update_sort", Permissions: []PermissionDefinition{PermissionRoleSort}},
		{Method: "POST", Pattern: "/api/admin/role/update_permissions", Permissions: []PermissionDefinition{PermissionRolePermissionUpdate}},

		{Method: "GET", Pattern: "/api/admin/sys-config/**", Permissions: []PermissionDefinition{PermissionSettingView}},
		{Method: "POST", Pattern: "/api/admin/sys-config/list", Permissions: []PermissionDefinition{PermissionSettingView}},
		{Method: "POST", Pattern: "/api/admin/sys-config/save", Permissions: []PermissionDefinition{PermissionSettingUpdate}},
		{Method: "GET", Pattern: "/api/admin/search/reindex/status", Permissions: []PermissionDefinition{PermissionSearchReindex}},
		{Method: "POST", Pattern: "/api/admin/search/reindex", Permissions: []PermissionDefinition{PermissionSearchReindex}},
		{Method: "GET", Pattern: "/api/admin/seo/sitemap/status", Permissions: []PermissionDefinition{PermissionSitemapGenerate}},
		{Method: "POST", Pattern: "/api/admin/seo/sitemap/generate", Permissions: []PermissionDefinition{PermissionSitemapGenerate}},

		{Method: "GET", Pattern: "/api/admin/badge/*", Permissions: []PermissionDefinition{PermissionBadgeView}},
		{Method: "POST", Pattern: "/api/admin/badge/list", Permissions: []PermissionDefinition{PermissionBadgeView}},
		{Method: "POST", Pattern: "/api/admin/badge/create", Permissions: []PermissionDefinition{PermissionBadgeCreate}},
		{Method: "POST", Pattern: "/api/admin/badge/update", Permissions: []PermissionDefinition{PermissionBadgeUpdate}},
		{Method: "POST", Pattern: "/api/admin/badge/delete", Permissions: []PermissionDefinition{PermissionBadgeDelete}},
		{Method: "POST", Pattern: "/api/admin/badge/update_sort", Permissions: []PermissionDefinition{PermissionBadgeUpdate}},

		{Method: "GET", Pattern: "/api/admin/level-config/*", Permissions: []PermissionDefinition{PermissionLevelView}},
		{Method: "POST", Pattern: "/api/admin/level-config/list", Permissions: []PermissionDefinition{PermissionLevelView}},
		{Method: "POST", Pattern: "/api/admin/level-config/save_all", Permissions: []PermissionDefinition{PermissionLevelUpdate}},

		{Method: "GET", Pattern: "/api/admin/task-config/*", Permissions: []PermissionDefinition{PermissionTaskView}},
		{Method: "GET", Pattern: "/api/admin/task-config/groups", Permissions: []PermissionDefinition{PermissionTaskView}},
		{Method: "POST", Pattern: "/api/admin/task-config/list", Permissions: []PermissionDefinition{PermissionTaskView}},
		{Method: "POST", Pattern: "/api/admin/task-config/create", Permissions: []PermissionDefinition{PermissionTaskCreate}},
		{Method: "POST", Pattern: "/api/admin/task-config/update", Permissions: []PermissionDefinition{PermissionTaskUpdate}},
		{Method: "POST", Pattern: "/api/admin/task-config/delete", Permissions: []PermissionDefinition{PermissionTaskDelete}},

		{Method: "GET", Pattern: "/api/admin/email-log/*", Permissions: []PermissionDefinition{PermissionEmailLogView}},
		{Method: "POST", Pattern: "/api/admin/email-log/list", Permissions: []PermissionDefinition{PermissionEmailLogView}},
		{Method: "GET", Pattern: "/api/admin/user-task-log/*", Permissions: []PermissionDefinition{PermissionUserTaskLogView}},
		{Method: "POST", Pattern: "/api/admin/user-task-log/list", Permissions: []PermissionDefinition{PermissionUserTaskLogView}},
		{Method: "GET", Pattern: "/api/admin/user-exp-log/*", Permissions: []PermissionDefinition{PermissionUserExpLogView}},
		{Method: "POST", Pattern: "/api/admin/user-exp-log/list", Permissions: []PermissionDefinition{PermissionUserExpLogView}},
		{Method: "GET", Pattern: "/api/admin/user-badge/*", Permissions: []PermissionDefinition{PermissionUserBadgeView}},
		{Method: "POST", Pattern: "/api/admin/user-badge/list", Permissions: []PermissionDefinition{PermissionUserBadgeView}},
		{Method: "GET", Pattern: "/api/admin/user-report/*", Permissions: []PermissionDefinition{PermissionUserReportView}},
		{Method: "POST", Pattern: "/api/admin/user-report/list", Permissions: []PermissionDefinition{PermissionUserReportView}},
		{Method: "GET", Pattern: "/api/admin/operate-log/*", Permissions: []PermissionDefinition{PermissionOperateLogView}},
		{Method: "POST", Pattern: "/api/admin/operate-log/list", Permissions: []PermissionDefinition{PermissionOperateLogView}},
	}
)

func GetAdminPermissionCode(method, path string) (string, bool) {
	codes, ok := GetAdminPermissionCodes(method, path)
	if !ok || len(codes) == 0 {
		return "", false
	}
	return codes[0], true
}

func GetAdminPermissionCodes(method, path string) ([]string, bool) {
	method = strings.ToUpper(method)
	for _, rule := range adminPermissionRules {
		if rule.Method == method && adminPathMatcher.Match(rule.Pattern, path) {
			codes := make([]string, 0, len(rule.Permissions))
			for _, permission := range rule.Permissions {
				codes = append(codes, permission.Code)
			}
			return codes, true
		}
	}
	return nil, false
}
