"use client"

import * as React from "react"
import {
  CheckIcon,
  ChevronDownIcon,
  GripVerticalIcon,
  PlusIcon,
  RefreshCwIcon,
  SaveIcon,
  Trash2Icon,
} from "lucide-react"

import { adminGet, adminPostJson, type AdminRecord } from "@/lib/api/admin"
import { PERMISSIONS } from "@/lib/auth/permissions.generated"
import { userHasPermission } from "@/lib/auth/roles"
import { useI18n } from "@/lib/i18n/provider"
import { msgError, msgSuccess } from "@/lib/toast"
import { cn } from "@/lib/utils"

import { useCurrentUser } from "@/components/app/app-provider"
import {
  ConfirmDialog,
  type ConfirmDialogState,
} from "@/components/common/confirm-dialog"
import { ErrorPage } from "@/components/common/error-page"
import { DashboardImageUpload } from "@/components/dashboard/image-upload"
import { DashboardSelect } from "@/components/dashboard/dashboard-select"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardAction,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Field as ShadcnField,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { Progress } from "@/components/ui/progress"
import { Separator } from "@/components/ui/separator"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Textarea } from "@/components/ui/textarea"
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip"

type SettingValue =
  | string
  | number
  | boolean
  | null
  | undefined
  | SettingValue[]
  | { [key: string]: SettingValue }

type SettingsState = Record<string, SettingValue>

type LocalizedText = {
  "zh-CN": string
  "en-US": string
}

type NavItem = {
  title: string
  url: string
  openInNewWindow: boolean
  children?: NavItem[]
}

type NotificationTypeRow = {
  key: string
  label: string
  site: boolean
  email: boolean
}

type ScriptInjection = {
  enabled: boolean
  scriptName: string
  type: "external" | "inline"
  src: string
  code: string
  async: boolean
  defer: boolean
  crossorigin: string
}

type FooterLink = {
  text: LocalizedText
  url: string
  visible: boolean
  openInNewWindow: boolean
}

type SearchReindexStatus = {
  running: boolean
  processed: number
  total: number
  topicProcessed?: number
  topicTotal?: number
  articleProcessed?: number
  articleTotal?: number
  userProcessed?: number
  userTotal?: number
  startedAt: number
  finishedAt: number
  error: string
}

type SitemapGenerateStatus = {
  running: boolean
  startedAt: number
  finishedAt: number
  error: string
  sitemapURL: string
}

const SETTINGS_ENDPOINT = "/api/admin/sys-config/configs"
const SAVE_ENDPOINT = "/api/admin/sys-config/save"
const SEARCH_REINDEX_ENDPOINT = "/api/admin/search/reindex"
const SEARCH_REINDEX_STATUS_ENDPOINT = "/api/admin/search/reindex/status"
const SITEMAP_GENERATE_ENDPOINT = "/api/admin/seo/sitemap/generate"
const SITEMAP_STATUS_ENDPOINT = "/api/admin/seo/sitemap/status"
const NOTIFICATION_TYPE_KEYS = [
  "topicComment",
  "commentReply",
  "topicLike",
  "topicFavorite",
  "topicRecommend",
  "topicDelete",
  "articleComment",
  "userLevelUp",
  "userBadgeGrant",
  "qaAnswerAccepted",
] as const
const DEFAULT_ATTACHMENT_TYPES = [
  ".pdf",
  ".doc",
  ".docx",
  ".xls",
  ".xlsx",
  ".ppt",
  ".pptx",
  ".txt",
  ".md",
  ".csv",
  ".zip",
  ".rar",
  ".7z",
  ".tar",
  ".gz",
]

function getPathValue(source: SettingsState, path: string): SettingValue {
  if (!path) return source
  return path.split(".").reduce<SettingValue>((current, part) => {
    if (current && typeof current === "object" && !Array.isArray(current)) {
      return (current as SettingsState)[part]
    }
    return undefined
  }, source)
}

function setPathValue(
  source: SettingsState,
  path: string,
  value: SettingValue
): SettingsState {
  const parts = path.split(".")
  const next: SettingsState = { ...source }
  let cursor = next
  parts.slice(0, -1).forEach((part) => {
    const existing = cursor[part]
    const child =
      existing && typeof existing === "object" && !Array.isArray(existing)
        ? { ...(existing as SettingsState) }
        : {}
    cursor[part] = child
    cursor = child
  })
  cursor[parts[parts.length - 1]] = value
  return next
}

function getObject(value: SettingValue): SettingsState {
  return value && typeof value === "object" && !Array.isArray(value)
    ? (value as SettingsState)
    : {}
}

function getString(value: SettingValue) {
  return typeof value === "string" ? value : ""
}

function getNumber(value: SettingValue) {
  return typeof value === "number" ? value : Number(value || 0)
}

function getStringArray(value: SettingValue) {
  return Array.isArray(value)
    ? value.filter((item): item is string => typeof item === "string")
    : []
}

function isAnyBlank(...values: string[]) {
  return values.some((value) => !value || !value.trim())
}

function isValidOptionalUrl(value: string) {
  const text = value.trim()
  if (!text) return true
  if (text.startsWith("/") || text.startsWith("#")) return true

  try {
    const url = new URL(text)
    return url.protocol === "http:" || url.protocol === "https:"
  } catch {
    return false
  }
}

function createLocalizedText(value?: unknown): LocalizedText {
  const source =
    value && typeof value === "object" && !Array.isArray(value)
      ? (value as Record<string, unknown>)
      : {}
  return {
    "zh-CN": typeof source["zh-CN"] === "string" ? source["zh-CN"] : "",
    "en-US": typeof source["en-US"] === "string" ? source["en-US"] : "",
  }
}

function trimLocalizedText(value: LocalizedText): LocalizedText {
  return {
    "zh-CN": value["zh-CN"].trim(),
    "en-US": value["en-US"].trim(),
  }
}

function hasLocalizedText(value: LocalizedText) {
  return Boolean(value["zh-CN"].trim() || value["en-US"].trim())
}

function normalizeNavs(items: unknown): NavItem[] {
  if (!Array.isArray(items)) return []
  return items.map((item) => {
    const source = getObject(item as SettingValue)
    return {
      title: getString(source.title),
      url: getString(source.url),
      openInNewWindow: Boolean(source.openInNewWindow),
      children: normalizeNavs(source.children),
    }
  })
}

function normalizeScripts(items: unknown): ScriptInjection[] {
  if (!Array.isArray(items)) return []
  return items.map((item) => {
    const source = getObject(item as SettingValue)
    return {
      enabled: Boolean(source.enabled),
      scriptName: getString(source.scriptName || source.remark),
      type: source.type === "inline" ? "inline" : "external",
      src: getString(source.src),
      code: getString(source.code),
      async: Boolean(source.async),
      defer: Boolean(source.defer),
      crossorigin: getString(source.crossorigin),
    }
  })
}

function normalizeFooterLinks(items: unknown): FooterLink[] {
  if (!Array.isArray(items)) return []
  return items.map((item) => {
    const source = getObject(item as SettingValue)
    return {
      text: createLocalizedText(source.text),
      url: getString(source.url),
      visible: source.visible !== false,
      openInNewWindow: Boolean(source.openInNewWindow),
    }
  })
}

function createNav(): NavItem {
  return { title: "", url: "", openInNewWindow: false, children: [] }
}

function createChildNav(): NavItem {
  return { title: "", url: "", openInNewWindow: false }
}

function createScript(): ScriptInjection {
  return {
    enabled: true,
    scriptName: "",
    type: "external",
    src: "",
    code: "",
    async: false,
    defer: false,
    crossorigin: "",
  }
}

function createFooterLink(): FooterLink {
  return {
    text: createLocalizedText(),
    url: "",
    visible: true,
    openInNewWindow: false,
  }
}

function formLabel(children: React.ReactNode) {
  return (
    <FieldLabel className="w-full justify-start text-sm font-normal text-muted-foreground sm:justify-end sm:text-right">
      {children}
    </FieldLabel>
  )
}

export default function DashboardSettingsRoute() {
  const { t } = useI18n()
  const currentUser = useCurrentUser()
  const [settings, setSettings] = React.useState<SettingsState>({})
  const [categories, setNodes] = React.useState<AdminRecord[]>([])
  const [loading, setLoading] = React.useState(true)
  const [savingSection, setSavingSection] = React.useState<string | null>(null)
  const [reindexing, setReindexing] = React.useState(false)
  const [searchStatus, setSearchStatus] =
    React.useState<SearchReindexStatus | null>(null)
  const [generatingSitemap, setGeneratingSitemap] = React.useState(false)
  const [sitemapStatus, setSitemapStatus] =
    React.useState<SitemapGenerateStatus | null>(null)
  const [confirmState, setConfirmState] =
    React.useState<ConfirmDialogState>(null)
  const [error, setError] = React.useState<string | null>(null)
  const canView = userHasPermission(currentUser, PERMISSIONS.DASHBOARD_SETTING_VIEW)
  const canUpdate = userHasPermission(currentUser, PERMISSIONS.DASHBOARD_SETTING_UPDATE)
  const canReindex = userHasPermission(
    currentUser,
    PERMISSIONS.DASHBOARD_SEARCH_REINDEX
  )
  const canGenerateSitemap = userHasPermission(
    currentUser,
    PERMISSIONS.DASHBOARD_SITEMAP_GENERATE
  )
  const tabKeys = React.useMemo(
    () => [
      "common",
      "content",
      "nav",
      "spam",
      "notification",
      "login",
      "upload",
      "script",
      "page",
      ...(canReindex ? ["search"] : []),
      ...(canGenerateSitemap ? ["sitemap"] : []),
    ],
    [canGenerateSitemap, canReindex]
  )

  const s = React.useCallback(
    (key: string, params?: Record<string, string | number>) =>
      t(`dashboard.settingsForm.${key}`, params),
    [t]
  )

  const loadConfig = React.useCallback(
    async (options?: { silent?: boolean }) => {
      if (!options?.silent) setLoading(true)
      setError(null)
      try {
        const [config, nodeList] = await Promise.all([
          adminGet<AdminRecord>(SETTINGS_ENDPOINT),
          adminGet<AdminRecord[]>("/api/admin/category/options").catch(
            () => []
          ),
        ])
        setSettings(config as SettingsState)
        setNodes(Array.isArray(nodeList) ? nodeList : [])
      } catch (err) {
        setError(
          err instanceof Error ? err.message : t("dashboard.errors.loadFailed")
        )
      } finally {
        if (!options?.silent) setLoading(false)
      }
    },
    [t]
  )

  React.useEffect(() => {
    void loadConfig()
  }, [loadConfig])

  const loadSearchStatus = React.useCallback(async () => {
    if (!canReindex) return
    try {
      const status = await adminGet<SearchReindexStatus>(
        SEARCH_REINDEX_STATUS_ENDPOINT
      )
      setSearchStatus(status)
    } catch {
      // Ignore polling failures; the rebuild action still reports explicit errors.
    }
  }, [canReindex])

  React.useEffect(() => {
    void loadSearchStatus()
  }, [loadSearchStatus])

  React.useEffect(() => {
    if (!canReindex || !searchStatus?.running) return
    const timer = window.setInterval(() => {
      void loadSearchStatus()
    }, 1500)
    return () => window.clearInterval(timer)
  }, [canReindex, loadSearchStatus, searchStatus?.running])

  const loadSitemapStatus = React.useCallback(async () => {
    if (!canGenerateSitemap) return
    try {
      const status = await adminGet<SitemapGenerateStatus>(
        SITEMAP_STATUS_ENDPOINT
      )
      setSitemapStatus(status)
    } catch {
      // Ignore polling failures; the generate action still reports explicit errors.
    }
  }, [canGenerateSitemap])

  React.useEffect(() => {
    void loadSitemapStatus()
  }, [loadSitemapStatus])

  React.useEffect(() => {
    if (!canGenerateSitemap || !sitemapStatus?.running) return
    const timer = window.setInterval(() => {
      void loadSitemapStatus()
    }, 1500)
    return () => window.clearInterval(timer)
  }, [canGenerateSitemap, loadSitemapStatus, sitemapStatus?.running])

  function update(path: string, value: SettingValue) {
    setSettings((current) => setPathValue(current, path, value))
  }

  function showActionError(message: string | null) {
    if (message) msgError(message)
  }

  async function saveSection(
    section: string,
    payload: Record<string, SettingValue>
  ) {
    if (!canUpdate) return
    setSavingSection(section)
    try {
      await adminPostJson(SAVE_ENDPOINT, payload)
      await loadConfig({ silent: true })
      msgSuccess(t("dashboard.messages.saved"))
    } catch (err) {
      msgError(
        err instanceof Error ? err.message : t("dashboard.errors.saveFailed")
      )
    } finally {
      setSavingSection(null)
    }
  }

  async function rebuildSearchIndex() {
    if (!canReindex || searchStatus?.running) return
    setReindexing(true)
    try {
      const status = await adminPostJson<SearchReindexStatus>(
        SEARCH_REINDEX_ENDPOINT,
        {}
      )
      setSearchStatus(status)
      msgSuccess(s("search.started"))
    } catch (err) {
      msgError(
        err instanceof Error ? err.message : t("dashboard.errors.saveFailed")
      )
    } finally {
      setReindexing(false)
    }
  }

  async function generateSitemap() {
    if (!canGenerateSitemap || sitemapStatus?.running) return
    setGeneratingSitemap(true)
    try {
      const status = await adminPostJson<SitemapGenerateStatus>(
        SITEMAP_GENERATE_ENDPOINT,
        {}
      )
      setSitemapStatus(status)
      msgSuccess(s("sitemap.started"))
    } catch (err) {
      msgError(
        err instanceof Error ? err.message : t("dashboard.errors.saveFailed")
      )
    } finally {
      setGeneratingSitemap(false)
    }
  }

  if (!canView) {
    return <ErrorPage statusCode={403} />
  }

  if (loading) {
    return (
      <div className="flex flex-1 flex-col gap-4 p-4 pt-4 md:p-6">
        <div className="rounded-lg border bg-[var(--dashboard-panel)] py-20 text-center text-sm text-muted-foreground shadow-xs">
          {t("dashboard.loading")}
        </div>
      </div>
    )
  }

  return (
    <TooltipProvider>
      <div className="flex flex-1 flex-col gap-4 p-4 pt-4 md:p-6">
        {error ? (
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        ) : null}

        <Tabs defaultValue="common" className="w-full">
          <div className="sticky top-0 z-10 mb-4 overflow-x-auto rounded-lg border bg-[var(--dashboard-panel)]/95 p-2 shadow-xs backdrop-blur supports-[backdrop-filter]:bg-[var(--dashboard-panel)]/85">
            <TabsList className="h-auto min-w-max justify-start bg-transparent p-0">
              {tabKeys.map((key) => (
                <TabsTrigger key={key} value={key}>
                  {s(`${key}.title`)}
                </TabsTrigger>
              ))}
            </TabsList>
          </div>

          <TabsContent value="common" className="mt-0">
            <CommonSettings
              settings={settings}
              saving={savingSection === "common"}
              s={s}
              update={update}
              onSave={() =>
                void saveSection("common", {
                  siteTitle: settings.siteTitle,
                  siteLogo: settings.siteLogo,
                  siteDescription: settings.siteDescription,
                  baseURL: settings.baseURL,
                  siteKeywords: settings.siteKeywords,
                  siteNotification: settings.siteNotification,
                })
              }
            />
          </TabsContent>

          <TabsContent value="content" className="mt-0">
            <ContentSettings
              settings={settings}
              categories={categories}
              saving={savingSection === "content"}
              s={s}
              update={update}
              onError={showActionError}
              onSave={() =>
                void saveSection("content", {
                  recommendTags: settings.recommendTags,
                  defaultCategoryId: settings.defaultCategoryId,
                  urlRedirect: settings.urlRedirect,
                  enableHideContent: settings.enableHideContent,
                  enableQaBounty: settings.enableQaBounty,
                  qaBountyMin: settings.qaBountyMin,
                  qaBountyMax: settings.qaBountyMax,
                  qaBountyRequired: settings.qaBountyRequired,
                  modules: settings.modules,
                  attachmentConfig: settings.attachmentConfig,
                })
              }
            />
          </TabsContent>

          <TabsContent value="nav" className="mt-0">
            <NavSettings
              settings={settings}
              saving={savingSection === "nav"}
              s={s}
              onChange={(items) => update("siteNavs", items)}
              onError={showActionError}
              onSave={() =>
                void saveSection("nav", { siteNavs: settings.siteNavs })
              }
            />
          </TabsContent>

          <TabsContent value="spam" className="mt-0">
            <SpamSettings
              settings={settings}
              saving={savingSection === "spam"}
              s={s}
              update={update}
              onSave={() =>
                void saveSection("spam", {
                  topicCaptcha: settings.topicCaptcha,
                  createTopicEmailVerified: settings.createTopicEmailVerified,
                  createArticleEmailVerified:
                    settings.createArticleEmailVerified,
                  createCommentEmailVerified:
                    settings.createCommentEmailVerified,
                  articlePending: settings.articlePending,
                  userObserveSeconds: settings.userObserveSeconds,
                })
              }
            />
          </TabsContent>

          <TabsContent value="notification" className="mt-0">
            <NotificationSettings
              settings={settings}
              saving={savingSection === "notification"}
              s={s}
              update={update}
              onSave={(payload) => void saveSection("notification", payload)}
            />
          </TabsContent>

          <TabsContent value="login" className="mt-0">
            <LoginSettings
              settings={settings}
              saving={savingSection === "login"}
              s={s}
              update={update}
              onSave={() =>
                void saveSection("login", {
                  loginConfig: settings.loginConfig,
                })
              }
            />
          </TabsContent>

          <TabsContent value="upload" className="mt-0">
            <UploadSettings
              settings={settings}
              saving={savingSection === "upload"}
              s={s}
              update={update}
              onSave={() =>
                void saveSection("upload", {
                  uploadConfig: settings.uploadConfig,
                })
              }
            />
          </TabsContent>

          <TabsContent value="script" className="mt-0">
            <ScriptSettings
              settings={settings}
              saving={savingSection === "script"}
              s={s}
              onChange={(items) => update("scriptInjections", items)}
              onError={showActionError}
              onSave={(payload) => void saveSection("script", payload)}
            />
          </TabsContent>

          <TabsContent value="page" className="mt-0">
            <PageSettings
              settings={settings}
              saving={savingSection === "page"}
              s={s}
              update={update}
              onError={showActionError}
              onSave={(payload) => void saveSection("page", payload)}
            />
          </TabsContent>

          {canReindex ? (
            <TabsContent value="search" className="mt-0">
              <SearchIndexSettings
                saving={reindexing}
                status={searchStatus}
                s={s}
                onReindex={() =>
                  setConfirmState({
                    title: s("search.confirmTitle"),
                    description: s("search.confirmDescription"),
                    confirmText: s("search.confirmAction"),
                    onConfirm: () => void rebuildSearchIndex(),
                  })
                }
              />
            </TabsContent>
          ) : null}

          {canGenerateSitemap ? (
            <TabsContent value="sitemap" className="mt-0">
              <SitemapSettings
                saving={generatingSitemap}
                status={sitemapStatus}
                s={s}
                onGenerate={() =>
                  setConfirmState({
                    title: s("sitemap.confirmTitle"),
                    description: s("sitemap.confirmDescription"),
                    confirmText: s("sitemap.confirmAction"),
                    onConfirm: () => void generateSitemap(),
                  })
                }
              />
            </TabsContent>
          ) : null}
        </Tabs>
        <ConfirmDialog
          state={confirmState}
          onOpenChange={(open) => {
            if (!open) setConfirmState(null)
          }}
        />
      </div>
    </TooltipProvider>
  )
}

function CommonSettings({
  settings,
  saving,
  s,
  update,
  onSave,
}: SettingsProps) {
  return (
    <SettingsForm
      onSave={onSave}
      saving={saving}
      submitLabel={s("common.submit")}
    >
      <Field label={s("common.siteTitle")}>
        <Input
          value={getString(settings.siteTitle)}
          placeholder={s("common.placeholder.siteTitle")}
          onChange={(event) => update("siteTitle", event.target.value)}
        />
      </Field>
      <Field label={s("common.siteLogo")}>
        <DashboardImageUpload
          value={getString(settings.siteLogo)}
          onChange={(value) => update("siteLogo", value)}
        />
      </Field>
      <Field label={s("common.siteDescription")}>
        <Textarea
          value={getString(settings.siteDescription)}
          placeholder={s("common.placeholder.siteDescription")}
          onChange={(event) => update("siteDescription", event.target.value)}
        />
      </Field>
      <Field label={s("common.baseURL")}>
        <Input
          value={getString(settings.baseURL)}
          placeholder={s("common.placeholder.baseURL")}
          onChange={(event) => update("baseURL", event.target.value)}
        />
      </Field>
      <Field label={s("common.siteKeywords")}>
        <TagsInput
          value={getStringArray(settings.siteKeywords)}
          placeholder={s("common.placeholder.siteKeywords")}
          onChange={(value) => update("siteKeywords", value)}
        />
      </Field>
      <Field label={s("common.siteNotification")}>
        <Textarea
          value={getString(settings.siteNotification)}
          placeholder={s("common.placeholder.siteNotification")}
          onChange={(event) => update("siteNotification", event.target.value)}
        />
      </Field>
    </SettingsForm>
  )
}

type SettingsProps = {
  settings: SettingsState
  saving: boolean
  s: (key: string, params?: Record<string, string | number>) => string
  update: (path: string, value: SettingValue) => void
  onSave: () => void
}

function SearchIndexSettings({
  saving,
  status,
  s,
  onReindex,
}: {
  saving: boolean
  status: SearchReindexStatus | null
  s: (key: string, params?: Record<string, string | number>) => string
  onReindex: () => void
}) {
  const running = Boolean(status?.running)
  const processed = status?.processed ?? 0
  const total = status?.total ?? 0
  const topicProcessed = status?.topicProcessed ?? 0
  const topicTotal = status?.topicTotal ?? 0
  const articleProcessed = status?.articleProcessed ?? 0
  const articleTotal = status?.articleTotal ?? 0
  const userProcessed = status?.userProcessed ?? 0
  const userTotal = status?.userTotal ?? 0
  const progress =
    total > 0 ? Math.min(100, Math.round((processed / total) * 100)) : 0
  const badgeVariant: React.ComponentProps<typeof Badge>["variant"] =
    status?.error ? "destructive" : running ? "secondary" : "outline"
  const statusText = status?.error
    ? s("search.statusFailed")
    : running
      ? s("search.statusRunning")
      : status?.finishedAt
        ? s("search.statusFinished")
        : s("search.statusIdle")
  const showStatus = running || Boolean(status?.finishedAt || status?.error)

  return (
    <Card size="sm" className="bg-[var(--dashboard-panel)] shadow-xs">
      <CardHeader className="border-b">
        <CardTitle>{s("search.title")}</CardTitle>
      </CardHeader>
      <CardContent className="grid gap-5 p-5">
        <div className="grid gap-2 sm:grid-cols-[184px_minmax(0,520px)] sm:gap-4">
          <div className="text-sm font-medium text-muted-foreground sm:text-right">
            {s("search.indexedContent")}
          </div>
          <div className="text-sm">
            {s("search.indexedContentAll")}
          </div>
        </div>
        {showStatus ? (
          <div className="grid gap-2 sm:grid-cols-[184px_minmax(0,520px)] sm:gap-4">
            <div className="text-sm font-medium text-muted-foreground sm:text-right">
              {s("search.status")}
            </div>
            <div className="grid gap-2">
              <div className="flex flex-wrap items-center gap-2">
                <Badge variant={badgeVariant}>{statusText}</Badge>
                <span className="text-sm text-muted-foreground">
                  {s("search.progressText", {
                    processed,
                    total,
                    percent: progress,
                  })}
                </span>
              </div>
              <Progress value={progress} />
              <div className="grid gap-1 text-xs text-muted-foreground sm:grid-cols-3">
                <span>
                  {s("search.topicProgress", {
                    processed: topicProcessed,
                    total: topicTotal,
                  })}
                </span>
                <span>
                  {s("search.articleProgress", {
                    processed: articleProcessed,
                    total: articleTotal,
                  })}
                </span>
                <span>
                  {s("search.userProgress", {
                    processed: userProcessed,
                    total: userTotal,
                  })}
                </span>
              </div>
              {status?.error ? (
                <p className="text-sm text-destructive">{status.error}</p>
              ) : null}
            </div>
          </div>
        ) : null}
        <div className="grid gap-2 sm:grid-cols-[184px_minmax(0,520px)] sm:gap-4">
          <div className="text-sm font-medium text-muted-foreground sm:text-right">
            {s("search.operation")}
          </div>
          <div className="grid gap-3">
            <p className="text-sm text-muted-foreground">
              {s("search.description")}
            </p>
            <Button
              className="w-fit"
              disabled={saving || running}
              onClick={onReindex}
            >
              <RefreshCwIcon
                className={running ? "animate-spin" : undefined}
              />
              {running ? s("search.rebuilding") : s("search.rebuild")}
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

function SitemapSettings({
  saving,
  status,
  s,
  onGenerate,
}: {
  saving: boolean
  status: SitemapGenerateStatus | null
  s: (key: string, params?: Record<string, string | number>) => string
  onGenerate: () => void
}) {
  const running = Boolean(status?.running)
  const badgeVariant: React.ComponentProps<typeof Badge>["variant"] =
    status?.error ? "destructive" : running ? "secondary" : "outline"
  const statusText = status?.error
    ? s("sitemap.statusFailed")
    : running
      ? s("sitemap.statusRunning")
      : status?.finishedAt
        ? s("sitemap.statusFinished")
        : s("sitemap.statusIdle")
  const showStatus = running || Boolean(status?.finishedAt || status?.error)
  const sitemapURL = status?.sitemapURL || "/sitemap.xml"

  return (
    <Card size="sm" className="bg-[var(--dashboard-panel)] shadow-xs">
      <CardHeader className="border-b">
        <CardTitle>{s("sitemap.title")}</CardTitle>
      </CardHeader>
      <CardContent className="grid gap-5 p-5">
        <div className="grid gap-2 sm:grid-cols-[184px_minmax(0,520px)] sm:gap-4">
          <div className="text-sm font-medium text-muted-foreground sm:text-right">
            {s("sitemap.content")}
          </div>
          <div className="text-sm">
            {s("sitemap.contentAll")}
          </div>
        </div>
        {showStatus ? (
          <div className="grid gap-2 sm:grid-cols-[184px_minmax(0,520px)] sm:gap-4">
            <div className="text-sm font-medium text-muted-foreground sm:text-right">
              {s("sitemap.status")}
            </div>
            <div className="grid gap-2">
              <div className="flex flex-wrap items-center gap-2">
                <Badge variant={badgeVariant}>{statusText}</Badge>
                {status?.sitemapURL ? (
                  <a
                    className="text-sm text-primary underline-offset-4 hover:underline"
                    href="/sitemap.xml"
                    target="_blank"
                    rel="noreferrer"
                  >
                    {s("sitemap.open")}
                  </a>
                ) : null}
              </div>
              {status?.error ? (
                <p className="text-sm text-destructive">{status.error}</p>
              ) : null}
            </div>
          </div>
        ) : null}
        <div className="grid gap-2 sm:grid-cols-[184px_minmax(0,520px)] sm:gap-4">
          <div className="text-sm font-medium text-muted-foreground sm:text-right">
            {s("sitemap.operation")}
          </div>
          <div className="grid gap-3">
            <p className="text-sm text-muted-foreground">
              {s("sitemap.description")}
            </p>
            <div className="flex flex-wrap items-center gap-3">
              <Button
                className="w-fit"
                disabled={saving || running}
                onClick={onGenerate}
              >
                <RefreshCwIcon
                  className={running ? "animate-spin" : undefined}
                />
                {running ? s("sitemap.generating") : s("sitemap.generate")}
              </Button>
              <span className="text-sm text-muted-foreground">
                {s("sitemap.publicURL")}: {sitemapURL}
              </span>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

function ContentSettings({
  settings,
  categories,
  saving,
  s,
  update,
  onError,
  onSave,
}: SettingsProps & {
  categories: AdminRecord[]
  onError: (message: string | null) => void
}) {
  const modules = getObject(settings.modules)
  const attachment = getObject(settings.attachmentConfig)
  const attachmentAllowedTypes = Array.isArray(attachment.allowedTypes)
    ? getStringArray(attachment.allowedTypes)
    : DEFAULT_ATTACHMENT_TYPES
  const qaBountyMin = getNumber(settings.qaBountyMin)
  const qaBountyMax = getNumber(settings.qaBountyMax)
  const [validationError, setValidationError] = React.useState<string | null>(
    null
  )

  function submit() {
    if (qaBountyMin > 0 && qaBountyMax > 0 && qaBountyMin > qaBountyMax) {
      const message = s("content.message.qaBountyRangeInvalid")
      setValidationError(message)
      onError(message)
      return
    }
    setValidationError(null)
    onSave()
  }

  return (
    <SettingsForm
      onSave={submit}
      saving={saving}
      submitLabel={s("common.submit")}
    >
      {validationError ? <ValidationAlert message={validationError} /> : null}
      <SectionTitle>{s("content.sectionGeneral")}</SectionTitle>
      <Field label={s("content.modules")}>
        <div className="flex flex-wrap gap-2">
          {(["tweet", "topic", "qa", "article"] as const).map((key) => (
            <ToggleBox
              key={key}
              checked={Boolean(modules[key])}
              label={s(`content.${key}`)}
              onChange={(checked) => update(`modules.${key}`, checked)}
            />
          ))}
        </div>
      </Field>
      <Field label={s("content.defaultCategoryId")}>
        <DashboardSelect
          value={settings.defaultCategoryId as string | number | null | undefined}
          options={categories.map((category) => ({
            label: String(category.name ?? category.title ?? category.id),
            value: Number(category.id),
          }))}
          placeholder={s("content.placeholder.defaultCategoryId")}
          onValueChange={(value) =>
            update("defaultCategoryId", value === undefined ? undefined : Number(value))
          }
        />
      </Field>
      <Field label={s("content.recommendTags")}>
        <TagsInput
          value={getStringArray(settings.recommendTags)}
          placeholder={s("content.placeholder.recommendTags")}
          onChange={(value) => update("recommendTags", value)}
        />
      </Field>
      <Field label={s("content.urlRedirect")}>
        <SwitchWithTooltip
          checked={Boolean(settings.urlRedirect)}
          tooltip={s("content.urlRedirectTooltip")}
          onChange={(checked) => update("urlRedirect", checked)}
        />
      </Field>
      <Field label={s("content.enableHideContent")}>
        <SwitchWithTooltip
          checked={Boolean(settings.enableHideContent)}
          tooltip={s("content.enableHideContentTooltip")}
          onChange={(checked) => update("enableHideContent", checked)}
        />
      </Field>

      <SectionTitle>{s("content.sectionQaBounty")}</SectionTitle>
      <Field label={s("content.enableQaBounty")}>
        <SwitchWithTooltip
          checked={settings.enableQaBounty !== false}
          tooltip={s("content.enableQaBountyTooltip")}
          onChange={(checked) => update("enableQaBounty", checked)}
        />
      </Field>
      <Field label={s("content.qaBountyMin")}>
        <TooltipNumberInput
          value={qaBountyMin}
          min={0}
          tooltip={s("content.qaBountyMinTooltip")}
          onChange={(value) => update("qaBountyMin", value)}
        />
      </Field>
      <Field label={s("content.qaBountyMax")}>
        <TooltipNumberInput
          value={qaBountyMax}
          min={0}
          tooltip={s("content.qaBountyMaxTooltip")}
          onChange={(value) => update("qaBountyMax", value)}
        />
      </Field>
      <Field label={s("content.qaBountyRequired")}>
        <SwitchWithTooltip
          checked={Boolean(settings.qaBountyRequired)}
          tooltip={s("content.qaBountyRequiredTooltip")}
          onChange={(checked) => update("qaBountyRequired", checked)}
        />
      </Field>

      <SectionTitle>{s("content.sectionAttachment")}</SectionTitle>
      <Field label={s("content.attachmentEnabled")}>
        <SwitchWithTooltip
          checked={attachment.enabled !== false}
          tooltip={s("content.attachmentEnabledTooltip")}
          onChange={(checked) => update("attachmentConfig.enabled", checked)}
        />
      </Field>
      <Field label={s("content.attachmentAllowedTypes")}>
        <TagsInput
          value={attachmentAllowedTypes}
          placeholder={s("content.placeholder.attachmentAllowedTypes")}
          onChange={(value) => update("attachmentConfig.allowedTypes", value)}
        />
      </Field>
      <Field label={s("content.attachmentMaxSizeMB")}>
        <TooltipNumberInput
          value={getNumber(attachment.maxSizeMB || 10)}
          min={1}
          max={100}
          tooltip={s("content.attachmentMaxSizeMBTooltip")}
          onChange={(value) => update("attachmentConfig.maxSizeMB", value)}
        />
      </Field>
      <Field label={s("content.attachmentMaxCount")}>
        <TooltipNumberInput
          value={getNumber(attachment.maxCount || 5)}
          min={1}
          max={20}
          tooltip={s("content.attachmentMaxCountTooltip")}
          onChange={(value) => update("attachmentConfig.maxCount", value)}
        />
      </Field>
    </SettingsForm>
  )
}

function NavSettings({
  settings,
  saving,
  s,
  onChange,
  onError,
  onSave,
}: {
  settings: SettingsState
  saving: boolean
  s: SettingsProps["s"]
  onChange: (items: NavItem[]) => void
  onError: (message: string | null) => void
  onSave: () => void
}) {
  const navs = normalizeNavs(settings.siteNavs)
  const [selectedIndex, setSelectedIndex] = React.useState(0)
  const [draggingIndex, setDraggingIndex] = React.useState<number | null>(null)
  const [validationError, setValidationError] = React.useState<string | null>(
    null
  )
  const selected = navs[selectedIndex]

  React.useEffect(() => {
    if (!navs.length) {
      setSelectedIndex(-1)
      return
    }
    if (selectedIndex < 0 || selectedIndex >= navs.length) {
      setSelectedIndex(0)
    }
  }, [navs.length, selectedIndex])

  function commit(items: NavItem[]) {
    onChange(items)
  }

  function updateParent(index: number, value: Partial<NavItem>) {
    commit(
      navs.map((item, itemIndex) =>
        itemIndex === index ? { ...item, ...value } : item
      )
    )
  }

  function addParent() {
    const next = [...navs, createNav()]
    commit(next)
    setSelectedIndex(next.length - 1)
  }

  function removeParent(index: number) {
    const next = navs.filter((_, itemIndex) => itemIndex !== index)
    commit(next)
    setSelectedIndex(Math.min(index, next.length - 1))
  }

  function moveParent(from: number, to: number) {
    if (
      from < 0 ||
      to < 0 ||
      from >= navs.length ||
      to >= navs.length ||
      from === to
    )
      return
    const next = [...navs]
    const [item] = next.splice(from, 1)
    next.splice(to, 0, item)
    commit(next)
    setSelectedIndex(to)
  }

  function updateChildren(children: NavItem[]) {
    if (!selected) return
    updateParent(selectedIndex, { children })
  }

  function submit() {
    const hasEmptyTitle = navs.some((nav) => isAnyBlank(nav.title))
    if (hasEmptyTitle) {
      const message = s("nav.message.validationTitleRequired")
      setValidationError(message)
      onError(message)
      return
    }
    const hasInvalidPrimaryUrl = navs.some(
      (nav) => !(nav.children || []).length && isAnyBlank(nav.url)
    )
    if (hasInvalidPrimaryUrl) {
      const message = s("nav.message.validationPrimaryUrlRequired")
      setValidationError(message)
      onError(message)
      return
    }
    const hasInvalidChild = navs.some((nav) =>
      (nav.children || []).some((child) => isAnyBlank(child.title, child.url))
    )
    if (hasInvalidChild) {
      const message = s("nav.message.validationChildRequired")
      setValidationError(message)
      onError(message)
      return
    }
    setValidationError(null)
    onSave()
  }

  return (
    <div className="flex flex-col gap-4">
      {validationError ? <ValidationAlert message={validationError} /> : null}
      <div className="grid gap-4 lg:grid-cols-[minmax(260px,36%)_1fr]">
        <div className="rounded-lg border bg-[var(--dashboard-panel)] p-3 shadow-xs">
          <div className="mb-3 flex items-center justify-between">
            <h3 className="text-sm font-medium">{s("nav.primaryList")}</h3>
            <Button size="icon-sm" onClick={addParent}>
              <PlusIcon />
              <span className="sr-only">{s("script.add")}</span>
            </Button>
          </div>
          <div className="flex max-h-[560px] flex-col gap-2 overflow-auto">
            {navs.length ? (
              navs.map((item, index) => (
                <button
                  key={`parent-${index}`}
                  type="button"
                  draggable
                  onClick={() => setSelectedIndex(index)}
                  onDragStart={(event) => {
                    setDraggingIndex(index)
                    event.dataTransfer.effectAllowed = "move"
                    event.dataTransfer.setData("text/plain", String(index))
                  }}
                  onDragOver={(event) => event.preventDefault()}
                  onDrop={() => {
                    if (draggingIndex !== null) moveParent(draggingIndex, index)
                    setDraggingIndex(null)
                  }}
                  onDragEnd={() => setDraggingIndex(null)}
                  className={cn(
                    "flex items-center gap-2 rounded-md border p-3 text-left text-sm transition-colors",
                    selectedIndex === index && "border-primary bg-muted",
                    draggingIndex === index && "opacity-60"
                  )}
                >
                  <GripVerticalIcon className="size-4 text-primary" />
                  <span className="min-w-0 flex-1">
                    <span className="block truncate font-medium">
                      {item.title || s("nav.untitled")}
                    </span>
                    <span className="block truncate text-xs text-muted-foreground">
                      {item.url || s("nav.groupMenu")}
                    </span>
                  </span>
                  <Button
                    type="button"
                    size="icon-xs"
                    variant="destructive"
                    onClick={(event) => {
                      event.stopPropagation()
                      removeParent(index)
                    }}
                  >
                    <Trash2Icon />
                  </Button>
                </button>
              ))
            ) : (
              <div className="rounded-md border border-dashed py-10 text-center text-sm text-muted-foreground">
                {s("nav.empty")}
              </div>
            )}
          </div>
        </div>

        <div className="rounded-lg border bg-[var(--dashboard-panel)] p-3 shadow-xs">
          {selected ? (
            <div className="grid gap-4">
              <h3 className="text-sm font-medium">{s("nav.detailTitle")}</h3>
              <div className="grid gap-3 md:grid-cols-[1fr_2fr_auto]">
                <ShadcnField>
                  <FieldLabel>{s("nav.tableTitle")}</FieldLabel>
                  <Input
                    value={selected.title}
                    onChange={(event) =>
                      updateParent(selectedIndex, { title: event.target.value })
                    }
                  />
                </ShadcnField>
                <ShadcnField>
                  <FieldLabel>{s("nav.tableUrl")}</FieldLabel>
                  <Input
                    value={selected.url}
                    placeholder={s("nav.groupUrlPlaceholder")}
                    onChange={(event) =>
                      updateParent(selectedIndex, { url: event.target.value })
                    }
                  />
                </ShadcnField>
                <ShadcnField>
                  <FieldLabel>{s("nav.tableOpenInNewWindow")}</FieldLabel>
                  <div className="flex min-h-9 items-center">
                    <SwitchControl
                      checked={selected.openInNewWindow}
                      onChange={(checked) =>
                        updateParent(selectedIndex, {
                          openInNewWindow: checked,
                        })
                      }
                    />
                  </div>
                </ShadcnField>
              </div>

              <div className="flex items-center justify-between">
                <h3 className="text-sm font-medium">
                  {s("nav.childrenTitle")}
                </h3>
                <Button
                  size="icon-sm"
                  onClick={() =>
                    updateChildren([
                      ...(selected.children || []),
                      createChildNav(),
                    ])
                  }
                >
                  <PlusIcon />
                </Button>
              </div>
              <ChildrenTable
                items={selected.children || []}
                s={s}
                onChange={updateChildren}
              />
            </div>
          ) : (
            <div className="rounded-md border border-dashed py-20 text-center text-sm text-muted-foreground">
              {s("nav.noSelection")}
            </div>
          )}
        </div>
      </div>
      <Button className="w-fit" disabled={saving} onClick={submit}>
        <SaveIcon />
        {s("nav.submit")}
      </Button>
    </div>
  )
}

function ChildrenTable({
  items,
  s,
  onChange,
}: {
  items: NavItem[]
  s: SettingsProps["s"]
  onChange: (items: NavItem[]) => void
}) {
  const [draggingIndex, setDraggingIndex] = React.useState<number | null>(null)

  function updateChild(index: number, value: Partial<NavItem>) {
    onChange(
      items.map((item, itemIndex) =>
        itemIndex === index ? { ...item, ...value } : item
      )
    )
  }

  function moveChild(from: number, to: number) {
    if (
      from < 0 ||
      to < 0 ||
      from >= items.length ||
      to >= items.length ||
      from === to
    )
      return
    const next = [...items]
    const [item] = next.splice(from, 1)
    next.splice(to, 0, item)
    onChange(next)
  }

  if (!items.length) {
    return (
      <div className="rounded-md border border-dashed py-8 text-center text-sm text-muted-foreground">
        {s("nav.childrenEmpty")}
      </div>
    )
  }

  return (
    <div className="overflow-x-auto rounded-lg border bg-[var(--dashboard-panel)] shadow-xs">
      <table className="w-full min-w-[760px] text-sm">
        <thead className="bg-[var(--dashboard-panel-muted)] text-muted-foreground">
          <tr>
            <th className="h-10 w-12 px-3 text-left text-xs font-semibold tracking-wide uppercase">
              {s("nav.sort")}
            </th>
            <th className="h-10 w-56 px-3 text-left text-xs font-semibold tracking-wide uppercase">
              {s("nav.tableTitle")}
            </th>
            <th className="h-10 px-3 text-left text-xs font-semibold tracking-wide uppercase">
              {s("nav.tableUrl")}
            </th>
            <th className="h-10 w-40 px-3 text-left text-xs font-semibold tracking-wide uppercase">
              {s("nav.tableOpenInNewWindow")}
            </th>
            <th className="h-10 w-24 px-3 text-right text-xs font-semibold tracking-wide uppercase">
              {s("nav.operation")}
            </th>
          </tr>
        </thead>
        <tbody>
          {items.map((item, index) => (
            <tr
              key={`child-${index}`}
              draggable
              onDragStart={(event) => {
                setDraggingIndex(index)
                event.dataTransfer.effectAllowed = "move"
              }}
              onDragOver={(event) => event.preventDefault()}
              onDrop={() => {
                if (draggingIndex !== null) moveChild(draggingIndex, index)
                setDraggingIndex(null)
              }}
              onDragEnd={() => setDraggingIndex(null)}
              className={cn(
                "border-t transition-colors hover:bg-muted/30",
                draggingIndex === index && "opacity-60"
              )}
            >
              <td className="h-11 px-3 py-2 align-middle">
                <GripVerticalIcon className="size-4 text-primary" />
              </td>
              <td className="h-11 px-3 py-2 align-middle">
                <Input
                  value={item.title}
                  onChange={(event) =>
                    updateChild(index, { title: event.target.value })
                  }
                />
              </td>
              <td className="h-11 px-3 py-2 align-middle">
                <Input
                  value={item.url}
                  onChange={(event) =>
                    updateChild(index, { url: event.target.value })
                  }
                />
              </td>
              <td className="h-11 px-3 py-2 align-middle">
                <SwitchControl
                  checked={item.openInNewWindow}
                  onChange={(checked) =>
                    updateChild(index, { openInNewWindow: checked })
                  }
                />
              </td>
              <td className="h-11 px-3 py-2 text-right align-middle">
                <Button
                  size="icon-sm"
                  variant="destructive"
                  onClick={() =>
                    onChange(
                      items.filter((_, itemIndex) => itemIndex !== index)
                    )
                  }
                >
                  <Trash2Icon />
                </Button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function SpamSettings({ settings, saving, s, update, onSave }: SettingsProps) {
  return (
    <SettingsForm
      onSave={onSave}
      saving={saving}
      submitLabel={s("spam.submit")}
    >
      {[
        ["topicCaptcha", "topicCaptcha"],
        ["createTopicEmailVerified", "createTopicEmailVerified"],
        ["createArticleEmailVerified", "createArticleEmailVerified"],
        ["createCommentEmailVerified", "createCommentEmailVerified"],
        ["articlePending", "articlePending"],
      ].map(([path, key]) => (
        <Field key={path} label={s(`spam.${key}`)}>
          <SwitchWithTooltip
            checked={Boolean(settings[path])}
            tooltip={s(`spam.${key}Tooltip`)}
            onChange={(checked) => update(path, checked)}
          />
        </Field>
      ))}
      <Field label={s("spam.userObserveSeconds")}>
        <TooltipNumberInput
          value={getNumber(settings.userObserveSeconds)}
          min={0}
          max={720}
          tooltip={s("spam.userObserveSecondsTooltip")}
          onChange={(value) => update("userObserveSeconds", value)}
        />
      </Field>
    </SettingsForm>
  )
}

function NotificationSettings({
  settings,
  saving,
  s,
  update,
  onSave,
}: Omit<SettingsProps, "onSave"> & {
  onSave: (payload: Record<string, SettingValue>) => void
}) {
  const smtp = getObject(settings.smtpConfig)
  const rawTypes = getObject(settings.notificationTypes)
  const rows = NOTIFICATION_TYPE_KEYS.map((key) => {
    const config = getObject(rawTypes[key])
    return {
      key,
      label: s(`notification.types.${key}`),
      site: config.site !== false,
      email:
        key === "topicDelete" ? config.email === true : config.email !== false,
    }
  })

  function updateType(key: string, field: "site" | "email", value: boolean) {
    update(`notificationTypes.${key}.${field}`, value)
  }

  function submit() {
    onSave({
      emailNoticeIntervalSeconds: settings.emailNoticeIntervalSeconds,
      emailWhitelist: settings.emailWhitelist,
      smtpConfig: settings.smtpConfig,
      notificationTypes: Object.fromEntries(
        rows.map((row) => [row.key, { site: row.site, email: row.email }])
      ) as SettingValue,
    })
  }

  return (
    <SettingsForm
      onSave={submit}
      saving={saving}
      submitLabel={s("notification.submit")}
    >
      <Field label={s("notification.emailNoticeIntervalSeconds")}>
        <TooltipNumberInput
          value={getNumber(settings.emailNoticeIntervalSeconds)}
          min={0}
          tooltip={s("notification.emailNoticeIntervalSecondsTooltip")}
          onChange={(value) => update("emailNoticeIntervalSeconds", value)}
        />
      </Field>
      <Field label={s("notification.emailWhitelist")}>
        <TagsInput
          value={getStringArray(settings.emailWhitelist)}
          placeholder={s("notification.placeholder.emailWhitelist")}
          onChange={(value) => update("emailWhitelist", value)}
        />
      </Field>

      <SectionTitle>{s("notification.typesTitle")}</SectionTitle>
      <div className="overflow-x-auto rounded-lg border bg-[var(--dashboard-panel)] shadow-xs">
        <table className="w-full min-w-[480px] text-sm">
          <thead className="bg-[var(--dashboard-panel-muted)] text-muted-foreground">
            <tr>
              <th className="h-10 px-3 text-left text-xs font-semibold tracking-wide uppercase">
                {s("notification.typeName")}
              </th>
              <th className="h-10 w-32 px-3 text-left text-xs font-semibold tracking-wide uppercase">
                {s("notification.siteColumn")}
              </th>
              <th className="h-10 w-32 px-3 text-left text-xs font-semibold tracking-wide uppercase">
                {s("notification.emailColumn")}
              </th>
            </tr>
          </thead>
          <tbody>
            {rows.map((row) => (
              <tr
                key={row.key}
                className="border-t transition-colors hover:bg-muted/30"
              >
                <td className="h-11 px-3 py-2 align-middle">{row.label}</td>
                <td className="h-11 px-3 py-2 align-middle">
                  <SwitchControl
                    checked={row.site}
                    onChange={(checked) => updateType(row.key, "site", checked)}
                  />
                </td>
                <td className="h-11 px-3 py-2 align-middle">
                  <SwitchControl
                    checked={row.email}
                    onChange={(checked) =>
                      updateType(row.key, "email", checked)
                    }
                  />
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <SectionTitle>{s("notification.smtpSectionTitle")}</SectionTitle>
      {(["host", "port", "username", "password"] as const).map((key) => (
        <Field key={key} label={s(`notification.smtp.${key}`)}>
          <Input
            className={key === "port" ? "max-w-44" : undefined}
            type={
              key === "password"
                ? "password"
                : key === "port"
                  ? "number"
                  : "text"
            }
            min={key === "port" ? 0 : undefined}
            value={getString(smtp[key])}
            onChange={(event) =>
              update(`smtpConfig.${key}`, event.target.value)
            }
          />
        </Field>
      ))}
      <Field label={s("notification.smtp.ssl")}>
        <SwitchControl
          checked={Boolean(smtp.ssl)}
          onChange={(checked) => update("smtpConfig.ssl", checked)}
        />
      </Field>
    </SettingsForm>
  )
}

function LoginSettings({ settings, saving, s, update, onSave }: SettingsProps) {
  const login = getObject(settings.loginConfig)

  return (
    <SettingsForm
      onSave={onSave}
      saving={saving}
      submitLabel={s("login.submit")}
    >
      <LoginCard title={s("login.weixinLogin")}>
        <Field label={s("login.enabled")}>
          <SwitchControl
            checked={Boolean(getPathValue(login, "weixinLogin.enabled"))}
            onChange={(checked) =>
              update("loginConfig.weixinLogin.enabled", checked)
            }
          />
        </Field>
        <Field label={s("login.appId")}>
          <Input
            value={getString(getPathValue(login, "weixinLogin.appId"))}
            onChange={(event) =>
              update("loginConfig.weixinLogin.appId", event.target.value)
            }
          />
        </Field>
        <Field label={s("login.appSecret")}>
          <Input
            type="password"
            value={getString(getPathValue(login, "weixinLogin.appSecret"))}
            onChange={(event) =>
              update("loginConfig.weixinLogin.appSecret", event.target.value)
            }
          />
        </Field>
      </LoginCard>

      <LoginCard title={s("login.smsLogin")}>
        <Field label={s("login.enabled")}>
          <SwitchControl
            checked={Boolean(getPathValue(login, "smsLogin.enabled"))}
            onChange={(checked) =>
              update("loginConfig.smsLogin.enabled", checked)
            }
          />
        </Field>
        {[
          ["accessKeyId", "aliyun.accessKeyId"],
          ["accessKeySecret", "aliyun.accessKeySecret"],
          ["signName", "aliyun.signName"],
          ["templateCode", "aliyun.templateCode"],
        ].map(([labelKey, path]) => (
          <Field key={path} label={s(`login.${labelKey}`)}>
            <Input
              type={labelKey === "accessKeySecret" ? "password" : "text"}
              value={getString(getPathValue(login, `smsLogin.${path}`))}
              onChange={(event) =>
                update(`loginConfig.smsLogin.${path}`, event.target.value)
              }
            />
          </Field>
        ))}
      </LoginCard>

      {(["googleLogin", "githubLogin"] as const).map((provider) => (
        <LoginCard key={provider} title={s(`login.${provider}`)}>
          <Field label={s("login.enabled")}>
            <SwitchControl
              checked={Boolean(getPathValue(login, `${provider}.enabled`))}
              onChange={(checked) =>
                update(`loginConfig.${provider}.enabled`, checked)
              }
            />
          </Field>
          <Field label={s("login.clientId")}>
            <Input
              value={getString(getPathValue(login, `${provider}.clientId`))}
              onChange={(event) =>
                update(`loginConfig.${provider}.clientId`, event.target.value)
              }
            />
          </Field>
          <Field label={s("login.clientSecret")}>
            <Input
              type="password"
              value={getString(getPathValue(login, `${provider}.clientSecret`))}
              onChange={(event) =>
                update(
                  `loginConfig.${provider}.clientSecret`,
                  event.target.value
                )
              }
            />
          </Field>
        </LoginCard>
      ))}

      <Field label={s("login.passwordLogin")}>
        <SwitchControl
          checked={Boolean(getPathValue(login, "passwordLogin.enabled"))}
          onChange={(checked) =>
            update("loginConfig.passwordLogin.enabled", checked)
          }
        />
      </Field>
    </SettingsForm>
  )
}

function UploadSettings({
  settings,
  saving,
  s,
  update,
  onSave,
}: SettingsProps) {
  const config = getObject(settings.uploadConfig)
  const method = getString(config.enableUploadMethod) || "Local"

  return (
    <SettingsForm
      onSave={onSave}
      saving={saving}
      submitLabel={s("upload.submit")}
    >
      <Card size="sm">
        <CardHeader>
          <CardTitle>{s("upload.uploadConfig")}</CardTitle>
        </CardHeader>
        <CardContent>
          <Field label={s("upload.enableUploadMethod")}>
            <RadioGroup
              value={method}
              options={[
                ["Local", s("upload.local")],
                ["AliyunOss", s("upload.aliyunOss")],
                ["TencentCos", s("upload.tencentCos")],
                ["AwsS3", s("upload.awsS3")],
              ]}
              showCheckbox
              onChange={(value) =>
                update("uploadConfig.enableUploadMethod", value)
              }
            />
          </Field>
        </CardContent>
      </Card>

      {method === "Local" ? (
        <Alert>
          <AlertDescription>{s("upload.localTip")}</AlertDescription>
        </Alert>
      ) : null}

      {method === "AliyunOss" ? (
        <ProviderCard title={s("upload.aliyunOss")}>
          {[
            ["host", "host"],
            ["bucket", "bucket"],
            ["endpoint", "endpoint"],
            ["accessKeyId", "accessKeyId"],
            ["accessKeySecret", "accessKeySecret"],
          ].map(([labelKey, field]) => (
            <Field key={field} label={s(`upload.${labelKey}`)}>
              <Input
                type={field === "accessKeySecret" ? "password" : "text"}
                value={getString(getPathValue(config, `aliyunOss.${field}`))}
                placeholder={s(`upload.placeholder.${labelKey}`)}
                onChange={(event) =>
                  update(`uploadConfig.aliyunOss.${field}`, event.target.value)
                }
              />
            </Field>
          ))}
          <SectionTitle>{s("upload.imageStyleConfig")}</SectionTitle>
          {[
            "styleSplitter",
            "styleAvatar",
            "stylePreview",
            "styleSmall",
            "styleDetail",
          ].map((field) => (
            <Field key={field} label={s(`upload.${field}`)}>
              <Input
                value={getString(getPathValue(config, `aliyunOss.${field}`))}
                placeholder={s(`upload.placeholder.${field}`)}
                onChange={(event) =>
                  update(`uploadConfig.aliyunOss.${field}`, event.target.value)
                }
              />
            </Field>
          ))}
        </ProviderCard>
      ) : null}

      {method === "TencentCos" ? (
        <ProviderCard title={s("upload.tencentCos")}>
          {[
            ["bucket", "bucket", "tencentBucket"],
            ["region", "region", "region"],
            ["secretId", "secretId", "secretId"],
            ["secretKey", "secretKey", "secretKey"],
          ].map(([field, labelKey, placeholderKey]) => (
            <Field key={field} label={s(`upload.${labelKey}`)}>
              <Input
                type={field === "secretKey" ? "password" : "text"}
                value={getString(getPathValue(config, `tencentCos.${field}`))}
                placeholder={s(`upload.placeholder.${placeholderKey}`)}
                onChange={(event) =>
                  update(`uploadConfig.tencentCos.${field}`, event.target.value)
                }
              />
            </Field>
          ))}
        </ProviderCard>
      ) : null}

      {method === "AwsS3" ? (
        <ProviderCard title={s("upload.awsS3")}>
          {[
            ["region", "region", "awsRegion"],
            ["bucket", "bucket", "awsBucket"],
            ["accessKeyId", "accessKeyId", "awsAccessKeyId"],
            ["accessKeySecret", "accessKeySecret", "awsSecretAccessKey"],
          ].map(([field, labelKey, placeholderKey]) => (
            <Field key={field} label={s(`upload.${labelKey}`)}>
              <Input
                type={field === "accessKeySecret" ? "password" : "text"}
                value={getString(getPathValue(config, `awsS3.${field}`))}
                placeholder={s(`upload.placeholder.${placeholderKey}`)}
                onChange={(event) =>
                  update(`uploadConfig.awsS3.${field}`, event.target.value)
                }
              />
            </Field>
          ))}
        </ProviderCard>
      ) : null}
    </SettingsForm>
  )
}

function ScriptSettings({
  settings,
  saving,
  s,
  onChange,
  onError,
  onSave,
}: {
  settings: SettingsState
  saving: boolean
  s: SettingsProps["s"]
  onChange: (items: ScriptInjection[]) => void
  onError: (message: string | null) => void
  onSave: (payload: Record<string, SettingValue>) => void
}) {
  const scripts = normalizeScripts(settings.scriptInjections)
  const [collapsedScripts, setCollapsedScripts] = React.useState<boolean[]>([])
  const [validationError, setValidationError] = React.useState<string | null>(
    null
  )

  function updateScript(index: number, value: Partial<ScriptInjection>) {
    onChange(
      scripts.map((item, itemIndex) =>
        itemIndex === index ? { ...item, ...value } : item
      )
    )
  }

  function toggleScriptCollapsed(index: number) {
    setCollapsedScripts((current) => {
      const next = [...current]
      next[index] = !next[index]
      return next
    })
  }

  function removeScript(index: number) {
    setCollapsedScripts((current) =>
      current.filter((_, itemIndex) => itemIndex !== index)
    )
    onChange(scripts.filter((_, itemIndex) => itemIndex !== index))
  }

  function submit() {
    if (scripts.some((item) => !item.scriptName.trim())) {
      const message = s("script.message.scriptNameRequired")
      setValidationError(message)
      onError(message)
      return
    }
    const scriptInjections = scripts.map((item) => ({
      ...item,
      scriptName: item.scriptName.trim(),
    }))
    setValidationError(null)
    onSave({ scriptInjections })
  }

  return (
    <div className="grid w-full gap-4">
      <Alert>
        <AlertDescription>{s("script.tip")}</AlertDescription>
      </Alert>
      {validationError ? <ValidationAlert message={validationError} /> : null}
      {scripts.map((item, index) => {
        const collapsed = Boolean(collapsedScripts[index])
        return (
          <Card
            key={`script-${index}`}
            size="sm"
            className={cn(
              "gap-4 bg-[var(--dashboard-panel)] shadow-xs",
              collapsed && "gap-0"
            )}
          >
            <CardHeader
              className={cn("border-b pb-4", collapsed && "border-b-0 pb-0")}
            >
              <CardTitle>
                <Input
                  value={item.scriptName}
                  placeholder={s("script.placeholder.scriptName")}
                  onChange={(event) =>
                    updateScript(index, { scriptName: event.target.value })
                  }
                />
              </CardTitle>
              <CardAction className="flex items-center gap-2">
                <Button
                  variant="destructive"
                  size="sm"
                  onClick={() => removeScript(index)}
                >
                  <Trash2Icon />
                  {s("script.remove")}
                </Button>
                <Button
                  variant="ghost"
                  size="icon-sm"
                  onClick={() => toggleScriptCollapsed(index)}
                >
                  <ChevronDownIcon
                    className={cn(
                      "transition-transform",
                      collapsed && "-rotate-90"
                    )}
                  />
                </Button>
              </CardAction>
            </CardHeader>
            {collapsed ? null : (
              <CardContent className="grid gap-5">
                <Field label={s("script.enabled")}>
                  <SwitchControl
                    checked={item.enabled}
                    onChange={(checked) =>
                      updateScript(index, { enabled: checked })
                    }
                  />
                </Field>
                <Field label={s("script.type")}>
                  <RadioGroup
                    value={item.type}
                    options={[
                      ["external", s("script.external")],
                      ["inline", s("script.inline")],
                    ]}
                    onChange={(value) =>
                      updateScript(index, {
                        type: value === "inline" ? "inline" : "external",
                      })
                    }
                  />
                </Field>
                {item.type === "external" ? (
                  <>
                    <Field label={s("script.src")}>
                      <Input
                        value={item.src}
                        placeholder={s("script.placeholder.src")}
                        onChange={(event) =>
                          updateScript(index, { src: event.target.value })
                        }
                      />
                    </Field>
                    <Field label={s("script.async")}>
                      <SwitchControl
                        checked={item.async}
                        onChange={(checked) =>
                          updateScript(index, { async: checked })
                        }
                      />
                    </Field>
                    <Field label={s("script.defer")}>
                      <SwitchControl
                        checked={item.defer}
                        onChange={(checked) =>
                          updateScript(index, { defer: checked })
                        }
                      />
                    </Field>
                    <Field label={s("script.crossorigin")}>
                      <Input
                        value={item.crossorigin}
                        placeholder={s("script.placeholder.crossorigin")}
                        onChange={(event) =>
                          updateScript(index, {
                            crossorigin: event.target.value,
                          })
                        }
                      />
                    </Field>
                  </>
                ) : (
                  <Field label={s("script.code")}>
                    <Textarea
                      className="min-h-40 font-mono"
                      value={item.code}
                      placeholder={s("script.placeholder.code")}
                      onChange={(event) =>
                        updateScript(index, { code: event.target.value })
                      }
                    />
                  </Field>
                )}
              </CardContent>
            )}
          </Card>
        )
      })}
      <div className="flex flex-wrap gap-2 rounded-lg border bg-[var(--dashboard-panel)] px-5 py-4 shadow-xs">
        <Button
          variant="outline"
          onClick={() => onChange([...scripts, createScript()])}
        >
          <PlusIcon />
          {s("script.add")}
        </Button>
        <Button disabled={saving} onClick={submit}>
          <SaveIcon />
          {s("script.submit")}
        </Button>
      </div>
    </div>
  )
}

function PageSettings({
  settings,
  saving,
  s,
  update,
  onError,
  onSave,
}: Omit<SettingsProps, "onSave"> & {
  onError: (message: string | null) => void
  onSave: (payload: Record<string, SettingValue>) => void
}) {
  const about = createLocalizedText(getObject(settings.aboutPageConfig).content)
  const footerLinks = normalizeFooterLinks(settings.footerLinks)
  const [validationError, setValidationError] = React.useState<string | null>(
    null
  )

  function updateAbout(locale: keyof LocalizedText, value: string) {
    update("aboutPageConfig.content", { ...about, [locale]: value })
  }

  function updateFooter(index: number, value: Partial<FooterLink>) {
    update(
      "footerLinks",
      footerLinks.map((item, itemIndex) =>
        itemIndex === index ? { ...item, ...value } : item
      )
    )
  }

  function submit() {
    const nextFooterLinks = footerLinks.map((item) => ({
      text: trimLocalizedText(item.text),
      url: item.url.trim(),
      visible: item.visible,
      openInNewWindow: item.openInNewWindow,
    }))
    if (nextFooterLinks.some((item) => !hasLocalizedText(item.text))) {
      const message = s("page.message.footerLinkInvalid")
      setValidationError(message)
      onError(message)
      return
    }
    if (nextFooterLinks.some((item) => !isValidOptionalUrl(item.url))) {
      const message = s("page.message.footerLinkUrlInvalid")
      setValidationError(message)
      onError(message)
      return
    }
    setValidationError(null)
    onSave({
      aboutPageConfig: {
        content: trimLocalizedText(about),
      },
      footerLinks: nextFooterLinks,
    })
  }

  return (
    <div className="grid w-full gap-4">
      <Alert>
        <AlertDescription>{s("page.aboutHelp")}</AlertDescription>
      </Alert>
      {validationError ? <ValidationAlert message={validationError} /> : null}

      <Card size="sm" className="gap-4 bg-[var(--dashboard-panel)] shadow-xs">
        <CardHeader className="border-b pb-4">
          <CardTitle>{s("page.aboutTitle")}</CardTitle>
        </CardHeader>
        <CardContent className="grid gap-5">
          <Field label={s("page.zhCNContent")} wide>
            <Textarea
              className="min-h-52 w-full"
              value={about["zh-CN"]}
              placeholder={s("page.placeholder.aboutContent")}
              onChange={(event) => updateAbout("zh-CN", event.target.value)}
            />
          </Field>
          <Field label={s("page.enUSContent")} wide>
            <Textarea
              className="min-h-52 w-full"
              value={about["en-US"]}
              placeholder={s("page.placeholder.aboutContent")}
              onChange={(event) => updateAbout("en-US", event.target.value)}
            />
          </Field>
        </CardContent>
      </Card>

      <Card size="sm" className="gap-4 bg-[var(--dashboard-panel)] shadow-xs">
        <CardHeader className="border-b pb-4">
          <CardTitle>{s("page.footerLinksTitle")}</CardTitle>
          <CardAction>
            <Button
              size="sm"
              onClick={() =>
                update("footerLinks", [...footerLinks, createFooterLink()])
              }
            >
              <PlusIcon />
              {s("page.addFooterLink")}
            </Button>
          </CardAction>
        </CardHeader>
        <CardContent className="grid gap-3">
          {footerLinks.length ? (
            footerLinks.map((item, index) => (
              <div
                key={`footer-${index}`}
                className="overflow-hidden rounded-md border bg-[var(--dashboard-panel)]"
              >
                <div className="flex items-center justify-between gap-3 border-b bg-muted/30 px-4 py-3">
                  <h4 className="text-sm font-medium">
                    {s("page.footerLinkItem", { index: index + 1 })}
                  </h4>
                  <Button
                    size="sm"
                    variant="destructive"
                    onClick={() =>
                      update(
                        "footerLinks",
                        footerLinks.filter(
                          (_, itemIndex) => itemIndex !== index
                        )
                      )
                    }
                  >
                    <Trash2Icon />
                    {s("page.remove")}
                  </Button>
                </div>
                <div className="grid gap-5 p-4">
                  <Field label={s("page.zhCNText")}>
                    <Input
                      value={item.text["zh-CN"]}
                      placeholder={s("page.placeholder.linkText")}
                      onChange={(event) =>
                        updateFooter(index, {
                          text: { ...item.text, "zh-CN": event.target.value },
                        })
                      }
                    />
                  </Field>
                  <Field label={s("page.enUSText")}>
                    <Input
                      value={item.text["en-US"]}
                      placeholder={s("page.placeholder.linkText")}
                      onChange={(event) =>
                        updateFooter(index, {
                          text: { ...item.text, "en-US": event.target.value },
                        })
                      }
                    />
                  </Field>
                  <Field label={s("page.urlOptional")}>
                    <Input
                      value={item.url}
                      placeholder={s("page.placeholder.linkUrl")}
                      onChange={(event) =>
                        updateFooter(index, { url: event.target.value })
                      }
                    />
                  </Field>
                  <Field label={s("page.visible")}>
                    <SwitchControl
                      checked={item.visible}
                      onChange={(checked) =>
                        updateFooter(index, { visible: checked })
                      }
                    />
                  </Field>
                  <Field label={s("page.openInNewWindow")}>
                    <SwitchControl
                      checked={item.openInNewWindow}
                      onChange={(checked) =>
                        updateFooter(index, { openInNewWindow: checked })
                      }
                    />
                  </Field>
                </div>
              </div>
            ))
          ) : (
            <div className="rounded-md border border-dashed py-10 text-center text-sm text-muted-foreground">
              {s("page.footerLinksEmpty")}
            </div>
          )}
        </CardContent>
      </Card>

      <div className="rounded-lg border bg-[var(--dashboard-panel)] px-5 py-4 shadow-xs">
        <Button disabled={saving} onClick={submit}>
          <SaveIcon />
          {s("common.submit")}
        </Button>
      </div>
    </div>
  )
}

function SettingsForm({
  children,
  saving,
  submitLabel,
  onSave,
}: {
  children: React.ReactNode
  saving: boolean
  submitLabel: string
  onSave: () => void
}) {
  return (
    <div className="w-full overflow-hidden rounded-lg border bg-[var(--dashboard-panel)] shadow-xs">
      <FieldGroup className="gap-5 p-5">{children}</FieldGroup>
      <div className="grid gap-2 border-t bg-[var(--dashboard-panel-muted)]/60 px-5 py-4 sm:grid-cols-[184px_minmax(0,520px)] sm:gap-4">
        <div className="hidden sm:block" />
        <Button className="w-fit" disabled={saving} onClick={onSave}>
          <SaveIcon />
          {submitLabel}
        </Button>
      </div>
    </div>
  )
}

function Field({
  label,
  children,
  wide = false,
}: {
  label: React.ReactNode
  children: React.ReactNode
  wide?: boolean
}) {
  return (
    <ShadcnField
      className={cn(
        "grid gap-2 sm:gap-4",
        wide
          ? "sm:grid-cols-[184px_minmax(0,1fr)]"
          : "sm:grid-cols-[184px_minmax(0,520px)]"
      )}
    >
      <div className="flex min-h-9 items-center self-start sm:justify-end">
        {formLabel(label)}
      </div>
      <div className="flex min-h-9 min-w-0 items-center self-start">
        {children}
      </div>
    </ShadcnField>
  )
}

function SectionTitle({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex items-center gap-3 pt-2">
      <h3 className="shrink-0 text-sm font-semibold">{children}</h3>
      <Separator className="flex-1" />
    </div>
  )
}

function ValidationAlert({ message }: { message: string }) {
  return (
    <Alert variant="destructive">
      <AlertDescription>{message}</AlertDescription>
    </Alert>
  )
}

function LoginCard({
  title,
  children,
}: {
  title: string
  children: React.ReactNode
}) {
  return (
    <Card size="sm" className="gap-4 bg-[var(--dashboard-panel)] shadow-xs">
      <CardHeader className="border-b pb-4">
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent className="grid gap-5">{children}</CardContent>
    </Card>
  )
}

function ProviderCard({
  title,
  children,
}: {
  title: string
  children: React.ReactNode
}) {
  return (
    <Card size="sm" className="gap-4 bg-[var(--dashboard-panel)] shadow-xs">
      <CardHeader className="border-b pb-4">
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent className="grid gap-5">{children}</CardContent>
    </Card>
  )
}

function SwitchControl({
  checked,
  onChange,
}: {
  checked: boolean
  onChange: (checked: boolean) => void
}) {
  return (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      onClick={() => onChange(!checked)}
      className={cn(
        "relative block h-6 w-11 shrink-0 rounded-full border transition-colors",
        checked ? "border-primary bg-primary" : "border-input bg-muted"
      )}
    >
      <span
        className={cn(
          "absolute top-0.5 left-0.5 size-5 rounded-full bg-background shadow transition-transform",
          checked && "translate-x-5"
        )}
      />
    </button>
  )
}

function ToggleBox({
  checked,
  label,
  onChange,
}: {
  checked: boolean
  label: string
  onChange: (checked: boolean) => void
}) {
  return (
    <button
      type="button"
      onClick={() => onChange(!checked)}
      className={cn(
        "inline-flex items-center gap-2 rounded-md border px-3 py-2 text-sm transition-colors",
        checked ? "border-primary bg-primary/10 text-primary" : "bg-background"
      )}
    >
      <span
        className={cn(
          "flex size-4 items-center justify-center rounded-sm border",
          checked
            ? "border-primary bg-primary text-primary-foreground"
            : "border-input bg-background"
        )}
      >
        {checked ? <CheckIcon className="size-3" /> : null}
      </span>
      {label}
    </button>
  )
}

function SwitchWithTooltip({
  checked,
  tooltip,
  onChange,
}: {
  checked: boolean
  tooltip: string
  onChange: (checked: boolean) => void
}) {
  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <span className="w-fit">
          <SwitchControl checked={checked} onChange={onChange} />
        </span>
      </TooltipTrigger>
      <TooltipContent>{tooltip}</TooltipContent>
    </Tooltip>
  )
}

function TooltipNumberInput({
  value,
  min,
  max,
  tooltip,
  onChange,
}: {
  value: number
  min?: number
  max?: number
  tooltip: string
  onChange: (value: number) => void
}) {
  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <Input
          className="max-w-44"
          type="number"
          value={Number.isFinite(value) ? value : 0}
          min={min}
          max={max}
          onChange={(event) => onChange(Number(event.target.value || 0))}
        />
      </TooltipTrigger>
      <TooltipContent>{tooltip}</TooltipContent>
    </Tooltip>
  )
}

function RadioGroup({
  value,
  options,
  showCheckbox = false,
  onChange,
}: {
  value: string
  options: Array<[string, string]>
  showCheckbox?: boolean
  onChange: (value: string) => void
}) {
  return (
    <div className="flex flex-wrap gap-2">
      {options.map(([optionValue, label]) => (
        <button
          key={optionValue}
          type="button"
          onClick={() => onChange(optionValue)}
          className={cn(
            "inline-flex items-center gap-2 rounded-md border px-3 py-2 text-sm",
            value === optionValue
              ? "border-primary bg-primary/10 text-primary"
              : "bg-background"
          )}
        >
          {showCheckbox ? (
            <span
              className={cn(
                "flex size-4 items-center justify-center rounded-sm border",
                value === optionValue
                  ? "border-primary bg-primary text-primary-foreground"
                  : "border-input bg-background"
              )}
            >
              {value === optionValue ? <CheckIcon className="size-3" /> : null}
            </span>
          ) : null}
          {label}
        </button>
      ))}
    </div>
  )
}

function TagsInput({
  value,
  placeholder,
  onChange,
}: {
  value: string[]
  placeholder: string
  onChange: (value: string[]) => void
}) {
  const [draft, setDraft] = React.useState("")

  function addDraft() {
    const items = draft
      .split(",")
      .map((item) => item.trim())
      .filter(Boolean)
    if (!items.length) return
    onChange(Array.from(new Set([...value, ...items])))
    setDraft("")
  }

  return (
    <div className="grid w-full gap-2">
      <div className="flex min-h-9 w-full flex-wrap gap-2 rounded-md border bg-background p-2">
        {value.map((item) => (
          <span
            key={item}
            className="inline-flex items-center gap-1 rounded-md bg-muted px-2 py-1 text-sm"
          >
            {item}
            <button
              type="button"
              className="text-muted-foreground hover:text-foreground"
              onClick={() =>
                onChange(value.filter((current) => current !== item))
              }
            >
              ×
            </button>
          </span>
        ))}
        <input
          className="min-w-44 flex-1 bg-transparent px-1 py-1 text-sm outline-none"
          value={draft}
          placeholder={placeholder}
          onChange={(event) => setDraft(event.target.value)}
          onKeyDown={(event) => {
            if (event.key === "Enter") {
              event.preventDefault()
              addDraft()
            } else if (event.key === "Backspace" && !draft && value.length) {
              event.preventDefault()
              onChange(value.slice(0, -1))
            }
          }}
          onBlur={addDraft}
        />
      </div>
    </div>
  )
}
