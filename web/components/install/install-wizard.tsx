"use client"

import * as React from "react"
import Link from "@/components/common/link"
import {
  AlertCircle,
  Check,
  CheckCircle,
  ChevronLeft,
  ChevronRight,
  Database,
  Download,
  Globe,
  Loader2,
  Settings,
  Shuffle,
} from "lucide-react"

import { Alert, AlertDescription } from "@/components/ui/alert"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { ApiError, apiFetch } from "@/lib/api/client"
import { createT, type Locale } from "@/lib/i18n"

type Step = "welcome" | "database" | "site" | "admin" | "install" | "complete"
type DbType = "mysql" | "sqlite"

const avatarCount = 128
const avatarBase = "/res/images/avatars"
const avatars = Array.from(
  { length: avatarCount },
  (_, index) => `${avatarBase}/${index}.png`
)

function randomAvatar() {
  return avatars[Math.floor(Math.random() * avatarCount)]
}

function defaultBaseURL() {
  return typeof window === "undefined" ? "/" : window.location.origin
}

export function InstallWizard({
  dockerBuiltinMysql = false,
}: {
  dockerBuiltinMysql?: boolean
}) {
  const [step, setStep] = React.useState<Step>("welcome")
  const [language, setLanguage] = React.useState<Locale>("en-US")
  const t = React.useMemo(() => createT(language), [language])
  const [dbConfig, setDbConfig] = React.useState({
    type: "mysql" as DbType,
    host: "localhost",
    port: "3306",
    database: "",
    username: "",
    password: "",
    path: "",
    inMemory: false,
  })
  const [siteInfo, setSiteInfo] = React.useState({
    title: "",
    description: "",
    baseURL: defaultBaseURL(),
  })
  const [adminInfo, setAdminInfo] = React.useState({
    username: "",
    password: "",
    passwordConfirm: "",
  })
  const [avatar, setAvatar] = React.useState(randomAvatar)
  const [dbError, setDbError] = React.useState("")
  const [dbSuccess, setDbSuccess] = React.useState("")
  const [siteError, setSiteError] = React.useState("")
  const [adminError, setAdminError] = React.useState("")
  const [testingConnection, setTestingConnection] = React.useState(false)
  const [installing, setInstalling] = React.useState(false)
  const [installProgress, setInstallProgress] = React.useState(0)
  const [installMessage, setInstallMessage] = React.useState(
    t("pages.install.install.preparing")
  )
  const [installFailed, setInstallFailed] = React.useState(false)

  const steps: Array<{ key: Step; title: string }> = [
    { key: "welcome", title: t("pages.install.step.welcome") },
    { key: "database", title: t("pages.install.step.database") },
    { key: "site", title: t("pages.install.step.site") },
    { key: "admin", title: t("pages.install.step.admin") },
    { key: "install", title: t("pages.install.step.install") },
    { key: "complete", title: t("pages.install.step.complete") },
  ]
  const currentStepIndex = steps.findIndex((item) => item.key === step)
  const pageBlocking = testingConnection || installing
  const effectiveDbType: DbType = dockerBuiltinMysql ? "mysql" : dbConfig.type
  const isDbFormValid =
    dockerBuiltinMysql ||
    effectiveDbType === "sqlite" ||
    Boolean(
      dbConfig.host && dbConfig.port && dbConfig.database && dbConfig.username
    )
  const isAdminFormValid =
    Boolean(
      adminInfo.username && adminInfo.password && adminInfo.passwordConfirm
    ) && adminInfo.password === adminInfo.passwordConfirm

  function updateDbConfig(values: Partial<typeof dbConfig>) {
    if (dockerBuiltinMysql) return
    setDbConfig((current) => {
      const next = { ...current, ...values }
      if (
        next.type === "sqlite" ||
        (next.host && next.port && next.database && next.username)
      ) {
        setDbError("")
      }
      return next
    })
    setDbSuccess("")
  }

  function validateDbForm() {
    if (dockerBuiltinMysql) {
      return true
    }
    if (effectiveDbType === "sqlite") {
      return true
    }
    if (
      dbConfig.host &&
      dbConfig.port &&
      dbConfig.database &&
      dbConfig.username
    ) {
      return true
    }
    setDbError(t("pages.install.database.validationError"))
    setDbSuccess("")
    return false
  }

  function validateSiteForm() {
    if (!siteInfo.title) {
      setSiteError(t("pages.install.site.validationError"))
      return false
    }
    if (!siteInfo.baseURL) {
      setSiteError(t("pages.install.site.baseURLValidationError"))
      return false
    }
    setSiteError("")
    return true
  }

  function validateAdminForm() {
    if (!adminInfo.username) {
      setAdminError(t("pages.install.admin.usernameError"))
      return false
    }
    if (!adminInfo.password) {
      setAdminError(t("pages.install.admin.passwordError"))
      return false
    }
    if (!adminInfo.passwordConfirm) {
      setAdminError(t("pages.install.admin.confirmPasswordError"))
      return false
    }
    if (adminInfo.password !== adminInfo.passwordConfirm) {
      setAdminError(t("pages.install.admin.passwordMismatchError"))
      return false
    }
    setAdminError("")
    return true
  }

  async function testDbConnection(nextStep?: Step) {
    if (!validateDbForm()) {
      return
    }
    setDbError("")
    setDbSuccess("")
    setTestingConnection(true)
    try {
      await apiFetch<null>("/api/install/test_db_connection", {
        method: "POST",
        body: {
          type: effectiveDbType,
          host: dbConfig.host,
          port: dbConfig.port,
          database: dbConfig.database,
          username: dbConfig.username,
          password: dbConfig.password,
          path: dbConfig.path,
          inMemory: dbConfig.inMemory,
        },
      })
      setDbSuccess(
        dockerBuiltinMysql
          ? t("pages.install.database.dockerBuiltinMysqlConnectSuccess")
          : t("pages.install.database.connectSuccess")
      )
      if (nextStep) {
        setStep(nextStep)
      }
    } catch (error) {
      setDbError(
        error instanceof Error
          ? `${t("pages.install.database.connectFailed")}: ${error.message}`
          : t("pages.install.database.connectFailed")
      )
    } finally {
      setTestingConnection(false)
    }
  }

  async function gotoStep(next: Step) {
    if (pageBlocking) return
    if (next === "site" && step === "database") {
      await testDbConnection("site")
      return
    }
    if (next === "admin" && step === "site" && !validateSiteForm()) {
      return
    }
    setStep(next)
  }

  async function confirmInstall() {
    if (!validateAdminForm()) {
      return
    }

    setInstalling(true)
    setInstallFailed(false)
    setStep("install")
    setInstallProgress(10)
    setInstallMessage(t("pages.install.install.connecting"))

    const timer = window.setInterval(() => {
      setInstallProgress((current) => (current < 90 ? current + 5 : current))
    }, 500)

    try {
      const response = await fetch("/api/install/install", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          siteTitle: siteInfo.title,
          siteDescription: siteInfo.description,
          baseURL: siteInfo.baseURL.trim() || "/",
          dbConfig: {
            type: dbConfig.type,
            host: dbConfig.host,
            port: dbConfig.port,
            database: dbConfig.database,
            username: dbConfig.username,
            password: dbConfig.password,
          },
          username: adminInfo.username,
          password: adminInfo.password,
          avatar,
          language,
        }),
      })
      const data = (await response.json()) as {
        success?: boolean
        message?: string
      }
      if (!data.success) {
        throw new ApiError(data.message || t("pages.install.install.unknown"))
      }
      setInstallProgress(100)
      setInstallMessage(t("pages.install.install.completed"))
      window.setTimeout(() => setStep("complete"), 1000)
    } catch (error) {
      setInstallProgress(100)
      setInstallFailed(true)
      setInstallMessage(
        error instanceof Error
          ? error instanceof ApiError
            ? `${t("pages.install.install.failed")}: ${error.message}`
            : `${t("pages.install.install.requestFailed")}: ${error.message}`
          : t("pages.install.install.requestFailed")
      )
    } finally {
      window.clearInterval(timer)
      setInstalling(false)
    }
  }

  return (
    <div className="page-glow flex min-h-screen items-start justify-center bg-slate-50 p-4 pt-10">
      {pageBlocking ? (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm">
          <div className="flex items-center space-x-3 text-sm font-medium text-white">
            <Loader2 className="h-5 w-5 animate-spin" />
            <span>
              {testingConnection
                ? t("pages.install.database.testing")
                : t("pages.install.install.installing")}
            </span>
          </div>
        </div>
      ) : null}
      <div className="w-full max-w-5xl">
        <div className="mb-5 text-center">
          <h1 className="mb-2 text-3xl font-bold text-gray-900">
            {t("pages.install.title")}
          </h1>
          <div className="mx-auto h-1.5 w-32 rounded-full bg-primary" />
        </div>

        <InstallSteps steps={steps} currentStepIndex={currentStepIndex} />

        <section className="card-animate rounded-2xl border bg-white/80 p-4 shadow-lg backdrop-blur-md md:p-5">
          {step === "welcome" ? (
            <div className="space-y-3">
              <SectionTitle
                title={t("pages.install.step.welcome")}
                description={t("pages.install.welcome.description")}
              />
              <div className="rounded-lg border border-blue-200 bg-blue-50 p-4">
                <h3 className="mb-2.5 font-medium text-blue-900">
                  {t("pages.install.welcome.requirements.title")}
                </h3>
                <ul className="space-y-2 text-blue-800">
                  <RequirementItem>
                    {t("pages.install.welcome.requirements.mysql")}
                  </RequirementItem>
                  <RequirementItem>
                    {t("pages.install.welcome.requirements.site")}
                  </RequirementItem>
                  <RequirementItem>
                    {t("pages.install.welcome.requirements.admin")}
                  </RequirementItem>
                </ul>
              </div>
              <div className="rounded-lg border border-gray-200 bg-gray-50 p-4">
                <h3 className="mb-2.5 font-medium text-gray-900">
                  {t("pages.install.language.title")}
                </h3>
                <p className="mb-3 text-gray-600">
                  {t("pages.install.language.description")}
                </p>
                <div className="space-y-2">
                  {(
                    [
                      ["en-US", t("pages.install.language.english")],
                      ["zh-CN", t("pages.install.language.chinese")],
                    ] as Array<[Locale, string]>
                  ).map(([value, label]) => (
                    <button
                      key={value}
                      type="button"
                      className={`flex w-full cursor-pointer items-center space-x-4 rounded-lg border p-3 text-left transition-colors hover:bg-gray-100/50 ${
                        language === value
                          ? "border-primary/50 bg-primary/10"
                          : "border-border bg-white"
                      }`}
                      onClick={() => setLanguage(value)}
                    >
                      <span className="mr-2 flex h-4 w-4 items-center justify-center rounded-full border border-primary">
                        {language === value ? (
                          <span className="h-2 w-2 rounded-full bg-primary" />
                        ) : null}
                      </span>
                      {label}
                    </button>
                  ))}
                </div>
              </div>
              <FooterActions
                right={
                  <NextButton
                    size="lg"
                    onClick={() => gotoStep("database")}
                    label={t("pages.install.buttons.next")}
                  />
                }
              />
            </div>
          ) : null}

          {step === "database" ? (
            <div className="space-y-3">
              <SectionTitle
                title={t("pages.install.step.database")}
                description={t("pages.install.database.description")}
              />
              <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                <div className="space-y-2 md:col-span-2">
                  <RequiredLabel>
                    {t("pages.install.database.type")}
                  </RequiredLabel>
                  <div className="flex items-center space-x-3">
                    {(["mysql", "sqlite"] as const).map((type) => (
                      <Button
                        key={type}
                        type="button"
                        size="sm"
                        variant="outline"
                        className={
                          effectiveDbType === type
                            ? "bg-primary text-white hover:bg-primary/90 hover:text-white"
                            : ""
                        }
                        disabled={dockerBuiltinMysql}
                        onClick={() => updateDbConfig({ type })}
                      >
                        {type === "mysql" ? "MySQL" : "SQLite"}
                      </Button>
                    ))}
                  </div>
                </div>
                {effectiveDbType === "mysql" ? (
                  <>
                    {dockerBuiltinMysql ? (
                      <Alert className="border-blue-200 bg-blue-50 md:col-span-2">
                        <div className="flex gap-2">
                          <Database className="mt-1 h-4 w-4 text-blue-700" />
                          <AlertDescription className="text-blue-800">
                            {t(
                              "pages.install.database.dockerBuiltinMysqlNotice"
                            )}
                          </AlertDescription>
                        </div>
                      </Alert>
                    ) : null}
                    <FormField
                      required
                      label={t("pages.install.database.host")}
                      id="db-host"
                      value={dockerBuiltinMysql ? "mysql" : dbConfig.host}
                      placeholder={t("pages.install.database.hostPlaceholder")}
                      onChange={(host) => updateDbConfig({ host })}
                      help={t("pages.install.database.hostHelp")}
                      disabled={dockerBuiltinMysql}
                    />
                    <FormField
                      required
                      label={t("pages.install.database.port")}
                      id="db-port"
                      value={dockerBuiltinMysql ? "3306" : dbConfig.port}
                      placeholder={t("pages.install.database.portPlaceholder")}
                      onChange={(port) => updateDbConfig({ port })}
                      disabled={dockerBuiltinMysql}
                    />
                    <FormField
                      required
                      label={t("pages.install.database.name")}
                      id="db-name"
                      value={dockerBuiltinMysql ? "bbsgo" : dbConfig.database}
                      placeholder={t("pages.install.database.namePlaceholder")}
                      onChange={(database) => updateDbConfig({ database })}
                      disabled={dockerBuiltinMysql}
                    />
                    <FormField
                      required
                      label={t("pages.install.database.username")}
                      id="db-user"
                      value={dockerBuiltinMysql ? "bbsgo" : dbConfig.username}
                      placeholder={t(
                        "pages.install.database.usernamePlaceholder"
                      )}
                      onChange={(username) => updateDbConfig({ username })}
                      disabled={dockerBuiltinMysql}
                    />
                    <FormField
                      className="md:col-span-2"
                      type="password"
                      label={t("pages.install.database.password")}
                      id="db-password"
                      value={dockerBuiltinMysql ? "********" : dbConfig.password}
                      placeholder={t(
                        "pages.install.database.passwordPlaceholder"
                      )}
                      onChange={(password) => updateDbConfig({ password })}
                      disabled={dockerBuiltinMysql}
                    />
                  </>
                ) : (
                  <div className="space-y-3 md:col-span-2">
                    <Alert
                      variant="destructive"
                      className="border-amber-200 bg-amber-50 text-amber-900"
                    >
                      <div className="flex gap-2">
                        <AlertCircle className="mt-1 h-4 w-4" />
                        <AlertDescription>
                          {t("pages.install.database.sqliteWarning")}
                        </AlertDescription>
                      </div>
                    </Alert>
                    <p className="text-sm text-gray-600">
                      {t("pages.install.database.sqliteAutoPath")}
                    </p>
                  </div>
                )}
              </div>
              <Feedback error={dbError} success={dbSuccess} />
              <FooterActions
                left={
                  <BackButton
                    size="lg"
                    onClick={() => gotoStep("welcome")}
                    label={t("pages.install.buttons.previous")}
                  />
                }
                right={
                  <div className="flex space-x-3">
                    <Button
                      type="button"
                      variant="secondary"
                      size="lg"
                      disabled={testingConnection || !isDbFormValid}
                      onClick={() => void testDbConnection()}
                    >
                      {testingConnection ? (
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      ) : (
                        <Database className="mr-2 h-4 w-4" />
                      )}
                      {testingConnection
                        ? t("pages.install.database.testing")
                        : t("pages.install.database.testConnection")}
                    </Button>
                    <NextButton
                      size="lg"
                      onClick={() => gotoStep("site")}
                      label={t("pages.install.buttons.next")}
                      disabled={testingConnection || !isDbFormValid}
                    />
                  </div>
                }
              />
            </div>
          ) : null}

          {step === "site" ? (
            <div className="space-y-3">
              <SectionTitle
                title={t("pages.install.step.site")}
                description={t("pages.install.site.description")}
              />
              <FormField
                required
                id="site-name"
                label={t("pages.install.site.title")}
                value={siteInfo.title}
                placeholder={t("pages.install.site.titlePlaceholder")}
                onBlur={validateSiteForm}
                onChange={(title) => {
                  setSiteInfo((current) => ({ ...current, title }))
                  if (title && siteInfo.baseURL) setSiteError("")
                }}
              />
              <div className="space-y-2">
                <Label htmlFor="site-desc" className="text-sm font-medium">
                  {t("pages.install.site.siteDescription")}
                </Label>
                <Textarea
                  id="site-desc"
                  rows={3}
                  value={siteInfo.description}
                  placeholder={t("pages.install.site.descriptionPlaceholder")}
                  onChange={(event) => {
                    const description = event.currentTarget.value
                    setSiteInfo((current) => ({ ...current, description }))
                  }}
                />
              </div>
              <FormField
                required
                id="site-base-url"
                label={t("pages.install.site.baseURL")}
                value={siteInfo.baseURL}
                placeholder={t("pages.install.site.baseURLPlaceholder")}
                onBlur={validateSiteForm}
                onChange={(baseURL) => {
                  setSiteInfo((current) => ({ ...current, baseURL }))
                  if (siteInfo.title && baseURL) setSiteError("")
                }}
              />
              <Feedback error={siteError} />
              <FooterActions
                left={
                  <BackButton
                    size="lg"
                    onClick={() => gotoStep("database")}
                    label={t("pages.install.buttons.previous")}
                  />
                }
                right={
                  <NextButton
                    size="lg"
                    onClick={() => gotoStep("admin")}
                    label={t("pages.install.buttons.next")}
                    disabled={!siteInfo.title}
                  />
                }
              />
            </div>
          ) : null}

          {step === "admin" ? (
            <div className="space-y-3">
              <SectionTitle
                title={t("pages.install.step.admin")}
                description={t("pages.install.admin.description")}
              />
              <FormField
                required
                id="admin-username"
                label={t("pages.install.admin.username")}
                value={adminInfo.username}
                placeholder={t("pages.install.admin.usernamePlaceholder")}
                onBlur={validateAdminForm}
                onChange={(username) => {
                  setAdminInfo((current) => ({ ...current, username }))
                  setAdminError("")
                }}
              />
              <FormField
                required
                id="admin-password"
                type="password"
                label={t("pages.install.admin.password")}
                value={adminInfo.password}
                placeholder={t("pages.install.admin.passwordPlaceholder")}
                onBlur={validateAdminForm}
                onChange={(password) => {
                  setAdminInfo((current) => ({ ...current, password }))
                  setAdminError("")
                }}
              />
              <FormField
                required
                id="admin-password-confirm"
                type="password"
                label={t("pages.install.admin.confirmPassword")}
                value={adminInfo.passwordConfirm}
                placeholder={t(
                  "pages.install.admin.confirmPasswordPlaceholder"
                )}
                onBlur={validateAdminForm}
                onChange={(passwordConfirm) => {
                  setAdminInfo((current) => ({ ...current, passwordConfirm }))
                  setAdminError("")
                }}
              />
              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <Label className="text-sm font-medium">
                    {t("pages.install.admin.avatar")}
                  </Label>
                  <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    onClick={() => setAvatar(randomAvatar())}
                  >
                    <Shuffle className="mr-1 h-4 w-4" />
                    {t("pages.install.admin.randomAvatar")}
                  </Button>
                </div>
                <div className="flex items-center gap-3">
                  <div className="h-16 w-16 overflow-hidden rounded-full border-2 border-primary/60 shadow-sm">                    <img
                      src={avatar}
                      alt={t("pages.install.admin.avatar")}
                      className="h-full w-full object-cover"
                    />
                  </div>
                  <p className="text-sm text-gray-600">
                    {t("pages.install.admin.avatarDescription")}
                  </p>
                </div>
                <div className="grid max-h-56 grid-cols-8 gap-2 overflow-y-auto rounded-xl border bg-gray-50 p-3 sm:grid-cols-10">
                  {avatars.map((item) => (
                    <button
                      key={item}
                      type="button"
                      className={`relative flex aspect-square h-12 w-12 items-center justify-center overflow-hidden rounded-full ring-2 transition hover:ring-primary/60 focus:outline-none ${
                        avatar === item
                          ? "bg-primary/15 shadow-md ring-primary"
                          : "bg-white ring-gray-200"
                      }`}
                      onClick={() => setAvatar(item)}
                    >                      <img
                        src={item}
                        alt={item}
                        className="h-full w-full object-cover"
                      />
                    </button>
                  ))}
                </div>
              </div>
              <Feedback error={adminError} />
              <FooterActions
                left={
                  <BackButton
                    size="lg"
                    onClick={() => gotoStep("site")}
                    label={t("pages.install.buttons.previous")}
                  />
                }
                right={
                  <Button
                    type="button"
                    size="lg"
                    disabled={!isAdminFormValid}
                    onClick={() => void confirmInstall()}
                  >
                    <Download className="mr-2 h-4 w-4" />
                    {t("pages.install.buttons.install")}
                  </Button>
                }
              />
            </div>
          ) : null}

          {step === "install" ? (
            <div className="space-y-3">
              <SectionTitle
                title={t("pages.install.step.install")}
                description={t("pages.install.install.description")}
              />
              <div className="space-y-3">
                <div className="h-4 w-full overflow-hidden rounded-full bg-muted">
                  <div
                    className={`h-full ${installFailed ? "bg-destructive" : "bg-primary"}`}
                    style={{ width: `${installProgress}%` }}
                  />
                </div>
                <div className="text-center">
                  <span className="text-sm font-medium text-gray-700">
                    {installProgress}%
                  </span>
                </div>
              </div>
              <Alert
                variant={installFailed ? "destructive" : "default"}
                className={
                  installFailed ? undefined : "border-blue-200 bg-blue-50"
                }
              >
                <div className="flex gap-2">
                  {installFailed ? (
                    <AlertCircle className="mt-1 h-4 w-4" />
                  ) : (
                    <Loader2 className="mt-1 h-4 w-4 animate-spin" />
                  )}
                  <AlertDescription
                    className={installFailed ? "text-red-800" : "text-blue-800"}
                  >
                    {installMessage}
                  </AlertDescription>
                </div>
              </Alert>
            </div>
          ) : null}

          {step === "complete" ? (
            <div className="space-y-3 text-center">
              <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-green-100">
                <CheckCircle className="h-8 w-8 text-green-600" />
              </div>
              <div>
                <h2 className="mb-4 text-2xl font-semibold text-gray-900">
                  {t("pages.install.step.complete")}
                </h2>
                <p className="mb-6 text-gray-600">
                  {t("pages.install.complete.congratulations")}
                </p>
                <p className="text-gray-600">
                  {t("pages.install.complete.description")}
                </p>
              </div>
              <div className="flex justify-center space-x-4">
                <Button asChild variant="outline" size="lg">
                  <Link href="/">
                    <Globe className="mr-2 h-4 w-4" />
                    {t("pages.install.complete.enterSite")}
                  </Link>
                </Button>
                <Button asChild size="lg">
                  <Link href="/admin">
                    <Settings className="mr-2 h-4 w-4" />
                    {t("pages.install.complete.enterAdmin")}
                  </Link>
                </Button>
              </div>
            </div>
          ) : null}
        </section>
      </div>
    </div>
  )
}

function InstallSteps({
  steps,
  currentStepIndex,
}: {
  steps: Array<{ key: Step; title: string }>
  currentStepIndex: number
}) {
  return (
    <div className="mb-2 flex justify-center overflow-x-auto px-1 py-4">
      <div className="flex items-start">
        {steps.map((step, index) => (
          <div key={step.key} className="flex items-start">
            <div className="flex w-20 flex-col items-center text-center">
              <div
                className={`flex h-10 w-10 items-center justify-center rounded-full text-base font-semibold transition-all duration-300 ease-in-out ${index === currentStepIndex ? "pulse-glow scale-110 bg-primary text-primary-foreground shadow-lg ring-4 ring-primary/20" : index < currentStepIndex ? "bg-green-500 text-white" : "bg-gray-200 text-gray-600"}`}
              >
                {index < currentStepIndex ? (
                  <Check className="h-6 w-6" />
                ) : (
                  index + 1
                )}
              </div>
              <div className="mt-2 text-xs leading-tight font-medium text-gray-600">
                {step.title}
              </div>
            </div>
            {index < steps.length - 1 ? (
              <div
                className={`mx-2 mt-4 h-1 w-12 rounded-full ${index < currentStepIndex ? "bg-green-500" : "bg-gray-200"}`}
              />
            ) : null}
          </div>
        ))}
      </div>
    </div>
  )
}

function SectionTitle({
  title,
  description,
}: {
  title: string
  description?: string
}) {
  return (
    <div className="text-center">
      <h2 className="mb-3 text-2xl font-semibold text-gray-900">{title}</h2>
      {description ? <p className="mb-4 text-gray-600">{description}</p> : null}
    </div>
  )
}

function RequirementItem({ children }: { children: React.ReactNode }) {
  return (
    <li className="flex items-center">
      <div className="mr-3 h-2 w-2 rounded-full bg-blue-500" />
      {children}
    </li>
  )
}

function RequiredLabel({
  children,
  htmlFor,
}: {
  children: React.ReactNode
  htmlFor?: string
}) {
  return (
    <Label htmlFor={htmlFor} className="text-sm font-medium">
      {children}
      <span className="text-red-500">*</span>
    </Label>
  )
}

function FormField({
  className,
  id,
  label,
  value,
  onBlur,
  onChange,
  placeholder,
  required,
  type = "text",
  help,
  disabled,
}: {
  className?: string
  id?: string
  label: string
  value: string
  onBlur?: () => void
  onChange: (value: string) => void
  placeholder?: string
  required?: boolean
  type?: string
  help?: string
  disabled?: boolean
}) {
  return (
    <div className={`space-y-2 ${className ?? ""}`}>
      {required ? (
        <RequiredLabel htmlFor={id}>{label}</RequiredLabel>
      ) : (
        <Label htmlFor={id} className="text-sm font-medium">
          {label}
        </Label>
      )}
      <Input
        id={id}
        required={required}
        type={type}
        value={value}
        placeholder={placeholder}
        disabled={disabled}
        onBlur={onBlur}
        onChange={(event) => onChange(event.currentTarget.value)}
      />
      {help ? <p className="text-xs text-gray-500">{help}</p> : null}
    </div>
  )
}

function Feedback({ error, success }: { error?: string; success?: string }) {
  return (
    <div className="min-h-[48px] space-y-2">
      {error ? (
        <Alert variant="destructive">
          <div className="flex gap-2">
            <AlertCircle className="mt-1 h-4 w-4" />
            <AlertDescription>{error}</AlertDescription>
          </div>
        </Alert>
      ) : null}
      {!error && success ? (
        <Alert className="border-green-200 bg-green-50">
          <div className="flex gap-2">
            <CheckCircle className="mt-1 h-4 w-4 text-green-600" />
            <AlertDescription className="text-green-800">
              {success}
            </AlertDescription>
          </div>
        </Alert>
      ) : null}
    </div>
  )
}

function FooterActions({
  left,
  right,
}: {
  left?: React.ReactNode
  right?: React.ReactNode
}) {
  return (
    <div className="flex flex-wrap items-center justify-between gap-3 pt-2">
      <div>{left}</div>
      <div>{right}</div>
    </div>
  )
}

function BackButton({
  label,
  onClick,
  size,
}: {
  label: string
  onClick: () => void
  size?: "default" | "lg"
}) {
  return (
    <Button type="button" variant="outline" size={size} onClick={onClick}>
      <ChevronLeft className="mr-2 h-4 w-4" />
      {label}
    </Button>
  )
}

function NextButton({
  label,
  onClick,
  disabled,
  size,
}: {
  label: string
  onClick: () => void
  disabled?: boolean
  size?: "default" | "lg"
}) {
  return (
    <Button type="button" size={size} disabled={disabled} onClick={onClick}>
      {label}
      <ChevronRight className="ml-2 h-4 w-4" />
    </Button>
  )
}
