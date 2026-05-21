"use client"

import * as React from "react"
import Link from "@/components/common/link"
import {
  ArrowRight,
  CheckSquare,
  ListChecks,
  Medal,
  Target,
  Zap,
} from "lucide-react"

import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { apiFetch } from "@/lib/api/client"
import type { TaskGroupInfo, TaskInfo } from "@/lib/api/types"
import { useI18n } from "@/lib/i18n/provider"
import { useCurrentUser } from "@/components/app/app-provider"
import { useToastActions } from "@/lib/toast"
import { cn } from "@/lib/utils"

function isCompleted(task: TaskInfo) {
  const progress = task.userProgress
  if (!progress) return false
  if (
    (progress.maxFinishCount || 0) > 0 &&
    (progress.finishedCount || 0) >= (progress.maxFinishCount || 0)
  ) {
    return true
  }
  return Boolean(
    progress.eventTarget &&
    (progress.eventProgress || 0) >= progress.eventTarget &&
    (progress.maxFinishCount || 0) <= 1
  )
}

function statusKind(task: TaskInfo) {
  if (isCompleted(task)) return "done"
  if (task.userProgress && (task.userProgress.eventProgress || 0) > 0) {
    return "progress"
  }
  return "idle"
}

function progressPercent(task: TaskInfo) {
  if (isCompleted(task)) return 100
  const progress = task.userProgress
  const target = progress?.eventTarget || task.eventCount || 1
  const done = Math.min(progress?.eventProgress || 0, target)
  return Math.min(100, Math.round((done / target) * 100))
}

function tasksForGroup(tasks: TaskInfo[], group: string) {
  if (!group || group === "all") return tasks
  return tasks.filter((item) => item.groupName === group)
}

export function TasksPageContent({
  groups,
  tasks,
}: {
  groups: TaskGroupInfo[]
  tasks: TaskInfo[]
}) {
  const { t } = useI18n()
  const [currentGroups, setCurrentGroups] = React.useState(groups)
  const [currentTasks, setCurrentTasks] = React.useState(tasks)
  const [activeGroup, setActiveGroup] = React.useState(groups[0]?.key || "all")
  const displayGroups = currentGroups.length
    ? currentGroups
    : [{ key: "all", name: t("user.tasks.groups.all") }]
  const stats = {
    total: currentTasks.length,
    completed: currentTasks.filter((task) => isCompleted(task)).length,
  }

  React.useEffect(() => {
    let mounted = true
    void Promise.all([
      apiFetch<TaskGroupInfo[]>("/api/task/groups").catch(() => []),
      apiFetch<TaskInfo[]>("/api/task/tasks").catch(() => []),
    ]).then(([nextGroups, nextTasks]) => {
      if (!mounted) {
        return
      }

      setCurrentGroups(nextGroups || [])
      setCurrentTasks(nextTasks || [])
      if (nextGroups?.[0]?.key) {
        setActiveGroup(nextGroups[0].key)
      }
    })

    return () => {
      mounted = false
    }
  }, [])

  return (
    <section className="rounded-lg bg-background px-3 py-2">
      <div className="flex flex-col gap-3 border-b border-border pb-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="text-lg font-semibold text-slate-900 dark:text-slate-50">
            {t("user.tasks.title")}
          </h1>
          <p className="text-sm text-slate-500 dark:text-slate-400">
            {t("user.tasks.subtitle")}
          </p>
        </div>
        <div className="flex flex-wrap gap-2 text-xs">
          <span className="inline-flex items-center gap-1 rounded-full bg-slate-100 px-2.5 py-1 font-semibold text-slate-700 dark:bg-slate-800 dark:text-slate-100">
            <ListChecks className="h-3.5 w-3.5" />
            {t("user.tasks.hero.total")} {stats.total}
          </span>
          <span className="inline-flex items-center gap-1 rounded-full bg-emerald-50 px-2.5 py-1 font-semibold text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-100">
            <CheckSquare className="h-3.5 w-3.5" />
            {t("user.tasks.hero.completed")} {stats.completed}
          </span>
        </div>
      </div>

      <div className="mt-4">
        <Tabs value={activeGroup} onValueChange={setActiveGroup}>
          <TabsList aria-label="task groups">
            {displayGroups.map((group) => (
              <TabsTrigger key={group.key} value={group.key}>
                <span>{group.name || group.key}</span>
              </TabsTrigger>
            ))}
          </TabsList>
          {displayGroups.map((group) => (
            <TabsContent
              key={group.key}
              value={group.key}
              className="space-y-3"
            >
              <TaskGrid tasks={tasksForGroup(currentTasks, group.key)} />
            </TabsContent>
          ))}
        </Tabs>
      </div>
    </section>
  )
}

function TaskGrid({ tasks }: { tasks: TaskInfo[] }) {
  const { t } = useI18n()

  if (!tasks.length) {
    return (
      <div className="rounded-xl border border-dashed border-slate-200 bg-slate-50 p-8 text-center text-slate-600 dark:border-slate-700 dark:bg-slate-900/40 dark:text-slate-300">
        <div className="mb-2 text-3xl">🪁</div>
        <div className="mb-2 font-semibold">{t("user.tasks.emptyTitle")}</div>
        <Link
          className="inline-flex items-center rounded-lg bg-indigo-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-indigo-500"
          href="/"
        >
          {t("user.tasks.emptyAction")}
        </Link>
      </div>
    )
  }

  return (
    <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
      {tasks.map((task) => (
        <TaskCard key={task.id} task={task} />
      ))}
    </div>
  )
}

function TaskCard({ task }: { task: TaskInfo }) {
  const { t } = useI18n()
  const user = useCurrentUser()
  const { msgSignIn } = useToastActions()
  const kind = statusKind(task)

  function statusLabel() {
    if (kind === "done") return t("user.tasks.status.completed")
    if (kind === "progress") return t("user.tasks.status.inProgress")
    return t("user.tasks.status.ready")
  }

  function progressLabel() {
    const progress = task.userProgress
    if (!progress) return t("user.tasks.progress.notStarted")
    const target = progress.eventTarget || task.eventCount || 1
    const done = isCompleted(task)
      ? target
      : Math.min(progress.eventProgress || 0, target)
    return t("user.tasks.progress.counter", { progress: done, target })
  }

  function finishLabel() {
    const progress = task.userProgress
    if (!progress) return ""
    if ((progress.maxFinishCount || 0) > 0) {
      return t("user.tasks.progress.finishedCount", {
        count: Math.min(
          progress.finishedCount || 0,
          progress.maxFinishCount || 0
        ),
        max: progress.maxFinishCount || 0,
      })
    }
    return t("user.tasks.progress.unlimited")
  }

  const action = task.actionUrl ? (
    user ? (
      <Link
        className="inline-flex items-center gap-1 rounded-lg bg-slate-900 px-2.5 py-1.5 text-[11px] font-semibold text-white shadow-sm transition hover:-translate-y-0.5 hover:bg-indigo-600 hover:shadow-md dark:bg-indigo-600 dark:hover:bg-indigo-500"
        href={task.actionUrl}
      >
        <ArrowRight className="h-3.5 w-3.5" />
        {task.btnName || t("user.tasks.actions.go")}
      </Link>
    ) : (
      <button
        type="button"
        className="inline-flex items-center gap-1 rounded-lg bg-slate-900 px-2.5 py-1.5 text-[11px] font-semibold text-white shadow-sm transition hover:-translate-y-0.5 hover:bg-indigo-600 hover:shadow-md dark:bg-indigo-600 dark:hover:bg-indigo-500"
        onClick={() => msgSignIn()}
      >
        <ArrowRight className="h-3.5 w-3.5" />
        {task.btnName || t("user.tasks.actions.go")}
      </button>
    )
  ) : (
    <span
      className="inline-flex items-center gap-1 rounded-lg px-2.5 py-1.5 text-[11px] font-semibold opacity-0"
      aria-hidden="true"
    >
      <ArrowRight className="h-3.5 w-3.5" />
      {t("user.tasks.actions.go")}
    </span>
  )

  return (
    <div className="relative flex flex-col gap-3 overflow-hidden rounded-2xl border border-slate-100/70 bg-white/95 p-3.5 shadow-sm ring-1 ring-transparent transition hover:-translate-y-0.5 hover:shadow-md hover:ring-indigo-100 dark:border-slate-800/80 dark:bg-slate-900/80 dark:ring-slate-800/60">
      <div className="pointer-events-none absolute inset-0 bg-gradient-to-br from-slate-50/80 via-white/70 to-indigo-50/60 dark:from-slate-900/60 dark:via-slate-900/50 dark:to-indigo-900/40" />
      <div className="relative flex items-start justify-between gap-2.5">
        <div className="flex items-center gap-3">
          <div className="space-y-1">
            <h3 className="text-[15px] font-semibold text-slate-900 dark:text-slate-50">
              {task.title}
            </h3>
          </div>
        </div>
        <span
          className={cn(
            "inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-[11px] font-semibold backdrop-blur",
            kind === "done" &&
              "border-emerald-300/70 bg-emerald-50/90 text-emerald-800 dark:border-emerald-700/80 dark:bg-emerald-900/40 dark:text-emerald-100",
            kind === "progress" &&
              "border-blue-300/70 bg-blue-50/90 text-blue-800 dark:border-blue-700/80 dark:bg-blue-900/40 dark:text-blue-100",
            kind === "idle" &&
              "border-slate-300/70 bg-slate-100/90 text-slate-700 dark:border-slate-700/80 dark:bg-slate-800/70 dark:text-slate-100"
          )}
        >
          {statusLabel()}
        </span>
      </div>
      <p className="relative line-clamp-2 min-h-[42px] text-[13px] leading-relaxed text-slate-600 dark:text-slate-300">
        {task.description}
      </p>
      <div className="relative flex flex-wrap gap-1.5 text-[11px] text-slate-600 dark:text-slate-300">
        {task.score ? (
          <RewardPill icon={<Target className="h-3.5 w-3.5" />}>
            {t("user.tasks.reward.score", { score: task.score })}
          </RewardPill>
        ) : null}
        {task.exp ? (
          <RewardPill icon={<Zap className="h-3.5 w-3.5" />}>
            {t("user.tasks.reward.exp", { exp: task.exp })}
          </RewardPill>
        ) : null}
        {task.badgeId ? (
          <RewardPill icon={<Medal className="h-3.5 w-3.5" />}>
            {t("user.tasks.reward.badge")}
          </RewardPill>
        ) : null}
      </div>
      <div className="relative mt-auto space-y-2.5">
        <div className="flex items-center justify-between text-[11px] font-medium text-slate-600 dark:text-slate-300">
          <span>{progressLabel()}</span>
          <span className="flex items-center gap-1 font-semibold text-slate-800 dark:text-slate-100">
            {progressPercent(task)}%
          </span>
        </div>
        <div className="h-2 w-full overflow-hidden rounded-full bg-slate-200/80 dark:bg-slate-800">
          <div
            className="h-full rounded-full bg-gradient-to-r from-indigo-500 via-indigo-500/80 to-emerald-400"
            style={{ width: `${progressPercent(task)}%` }}
          />
        </div>
        <div className="flex items-center justify-between gap-3 text-[11px] text-slate-500 dark:text-slate-400">
          <span className="truncate">{finishLabel()}</span>
          <div className="flex min-h-[26px] items-center justify-end">
            {action}
          </div>
        </div>
      </div>
    </div>
  )
}

function RewardPill({
  icon,
  children,
}: {
  icon: React.ReactNode
  children: React.ReactNode
}) {
  return (
    <span className="inline-flex items-center gap-1 rounded-full border border-slate-200/80 bg-white/70 px-2 py-0.5 font-semibold text-slate-700 shadow-sm backdrop-blur dark:border-slate-700 dark:bg-slate-800/70 dark:text-slate-100">
      {icon}
      {children}
    </span>
  )
}
