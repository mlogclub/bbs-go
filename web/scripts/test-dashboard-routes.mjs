import assert from "node:assert/strict"
import { existsSync, readFileSync } from "node:fs"
import { resolve } from "node:path"

const webRoot = resolve(import.meta.dirname, "..")
const routesDir = resolve(webRoot, "app/routes")
const dashboardComponentsDir = resolve(webRoot, "components/dashboard")
const dashboardDataDir = resolve(dashboardComponentsDir, "data")

const dedicatedRoutes = {
  "dashboard.users.tsx": {
    expectedDefaultExport: "DashboardUsersRoute",
  },
  "dashboard.settings.tsx": {
    removedComponent: "admin-settings-page.tsx",
    expectedDefaultExport: "DashboardSettingsRoute",
    forbiddenImport: "admin-settings-page",
  },
  "dashboard.topics.tsx": {
    removedComponent: "admin-topic-feed-page.tsx",
    expectedDefaultExport: "DashboardTopicsRoute",
    forbiddenImport: "admin-topic-feed-page",
  },
  "dashboard.articles.tsx": {
    removedComponent: "admin-article-feed-page.tsx",
    expectedDefaultExport: "DashboardArticlesRoute",
    forbiddenImport: "admin-article-feed-page",
  },
  "dashboard.levels.tsx": {
    removedComponent: "admin-levels-page.tsx",
    expectedDefaultExport: "DashboardLevelsRoute",
    forbiddenImport: "admin-levels-page",
  },
  "dashboard.user-reports.tsx": {
    expectedDefaultExport: "DashboardUserReportsRoute",
  },
  "dashboard.nodes.tsx": {
    expectedDefaultExport: "DashboardNodesRoute",
  },
  "dashboard.links.tsx": {
    expectedDefaultExport: "DashboardLinksRoute",
  },
  "dashboard.forbidden-words.tsx": {
    expectedDefaultExport: "DashboardForbiddenWordsRoute",
  },
  "dashboard.badges.tsx": {
    expectedDefaultExport: "DashboardBadgesRoute",
  },
  "dashboard.tasks.tsx": {
    expectedDefaultExport: "DashboardTasksRoute",
  },
  "dashboard.roles.tsx": {
    expectedDefaultExport: "DashboardRolesRoute",
  },
  "dashboard.user-badges.tsx": {
    expectedDefaultExport: "DashboardUserBadgesRoute",
  },
  "dashboard.user-exp-logs.tsx": {
    expectedDefaultExport: "DashboardUserExpLogsRoute",
  },
  "dashboard.user-task-logs.tsx": {
    expectedDefaultExport: "DashboardUserTaskLogsRoute",
  },
  "dashboard.email-logs.tsx": {
    expectedDefaultExport: "DashboardEmailLogsRoute",
  },
  "dashboard.content.tsx": {
    expectedDefaultExport: "DashboardContentRoute",
  },
}

for (const [routeFile, routeConfig] of Object.entries(dedicatedRoutes)) {
  const routePath = resolve(routesDir, routeFile)
  assert.equal(
    existsSync(routePath),
    true,
    `${routeFile} should be a dedicated dashboard route`
  )

  const routeSource = readFileSync(routePath, "utf8")
  assert.match(
    routeSource,
    new RegExp(`export default function ${routeConfig.expectedDefaultExport}`),
    `${routeFile} should own its default route component`
  )

  if (routeConfig.forbiddenImport) {
    assert.equal(
      routeSource.includes(routeConfig.forbiddenImport),
      false,
      `${routeFile} should not wrap the old dashboard page component`
    )
  }

  if (routeConfig.removedComponent) {
    assert.equal(
      existsSync(resolve(dashboardComponentsDir, routeConfig.removedComponent)),
      false,
      `${routeConfig.removedComponent} should be folded into its route module`
    )
  }
}

assert.equal(
  existsSync(resolve(routesDir, "dashboard.$.tsx")),
  false,
  "dashboard.$.tsx should be removed after splitting dashboard pages"
)

assert.equal(
  existsSync(resolve(routesDir, "dashboard.comments.tsx")),
  false,
  "dashboard.comments.tsx should be removed because comments are not managed in dashboard"
)

assert.equal(
  existsSync(resolve(dashboardDataDir, "dashboard-data-page-configs.tsx")),
  false,
  "dashboard-data-page-configs.tsx should be removed after moving configs into route modules"
)

for (const sourcePath of [
  resolve(dashboardComponentsDir, "app-sidebar.tsx"),
  resolve(dashboardComponentsDir, "dashboard-overview.tsx"),
]) {
  const source = readFileSync(sourcePath, "utf8")
  assert.equal(
    source.includes("/dashboard/comments"),
    false,
    `${sourcePath} should not link to dashboard comments`
  )
}

const nodesRoute = readFileSync(resolve(routesDir, "dashboard.nodes.tsx"), "utf8")

assert.equal(
  /name:\s*"parentId"[\s\S]*?type:\s*"tree-select"/.test(nodesRoute),
  false,
  "dashboard.nodes.tsx parent node form field should use DashboardSelect via type select"
)

assert.match(
  nodesRoute,
  /name:\s*"nodeId"[\s\S]*?label:\s*dashboardData\.label\(t,\s*"node"\)[\s\S]*?optionsEndpoint:\s*"\/api\/admin\/topic-node\/nodes"/,
  "dashboard.nodes.tsx should filter by selected node id instead of parent id"
)

console.log("dashboard route structure tests passed")
