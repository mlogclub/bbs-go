"use client"

import * as React from "react"
import Link from "@/components/common/link"
import { ArrowRight, Medal, Sparkles, Target, Zap } from "lucide-react"

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
    <section className="rounded-md bg-background">
      <Tabs
        value={activeGroup}
        onValueChange={setActiveGroup}
        className="gap-0"
      >
        <div className="flex flex-col gap-3 border-b border-border px-4 py-3 lg:flex-row lg:items-center lg:justify-between">
          <div className="min-w-0">
            <div className="flex flex-wrap items-center gap-x-3 gap-y-1">
              <h1 className="text-base font-semibold text-foreground">
                {t("user.tasks.title")}
              </h1>
              <span className="inline-flex items-center rounded-md border border-border bg-muted/40 px-2 py-0.5 text-xs font-medium text-muted-foreground">
                {stats.completed} / {stats.total}
              </span>
            </div>
            <p className="mt-0.5 text-xs text-muted-foreground">
              {t("user.tasks.subtitle")}
            </p>
          </div>
          <div className="flex justify-start lg:justify-end">
            <TabsList aria-label="task groups">
              {displayGroups.map((group) => (
                <TabsTrigger key={group.key} value={group.key}>
                  <span>{group.name || group.key}</span>
                </TabsTrigger>
              ))}
            </TabsList>
          </div>
        </div>

        <div className="px-4 py-3">
          {displayGroups.map((group) => (
            <TabsContent key={group.key} value={group.key}>
              <TaskList tasks={tasksForGroup(currentTasks, group.key)} />
            </TabsContent>
          ))}
        </div>
      </Tabs>
    </section>
  )
}

function TaskList({ tasks }: { tasks: TaskInfo[] }) {
  const { t } = useI18n()

  if (!tasks.length) {
    return (
      <div className="rounded-md border border-dashed border-border bg-muted/30 p-8 text-center text-muted-foreground">
        <Sparkles className="mx-auto mb-2 h-5 w-5" />
        <div className="mb-2 text-sm font-semibold text-foreground">
          {t("user.tasks.emptyTitle")}
        </div>
        <Link
          className="inline-flex items-center rounded-md bg-primary px-3 py-1.5 text-xs font-medium text-primary-foreground transition hover:opacity-90"
          href="/"
        >
          {t("user.tasks.emptyAction")}
        </Link>
      </div>
    )
  }

  return (
    <div className="divide-y divide-border/70">
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
        className="inline-flex h-8 items-center gap-1.5 rounded-md bg-primary px-3 text-xs font-medium text-primary-foreground transition hover:opacity-90"
        href={task.actionUrl}
      >
        <ArrowRight className="h-3.5 w-3.5" />
        {task.btnName || t("user.tasks.actions.go")}
      </Link>
    ) : (
      <button
        type="button"
        className="inline-flex h-8 items-center gap-1.5 rounded-md bg-primary px-3 text-xs font-medium text-primary-foreground transition hover:opacity-90"
        onClick={() => msgSignIn()}
      >
        <ArrowRight className="h-3.5 w-3.5" />
        {task.btnName || t("user.tasks.actions.go")}
      </button>
    )
  ) : (
    <span
      className="inline-flex h-8 items-center gap-1.5 rounded-md px-3 text-xs font-medium opacity-0"
      aria-hidden="true"
    >
      <ArrowRight className="h-3.5 w-3.5" />
      {t("user.tasks.actions.go")}
    </span>
  )

  return (
    <div className="px-3 py-2.5 transition hover:bg-muted/30">
      <div className="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <div className="min-w-0 flex-1">
          <div className="flex flex-wrap items-center gap-2">
            <h3 className="text-sm font-semibold text-foreground">
              {task.title}
            </h3>
            <span
              className={cn(
                "inline-flex items-center rounded-md px-1.5 py-0.5 text-[11px] font-medium",
                kind === "done" &&
                  "bg-emerald-50 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-300",
                kind === "progress" &&
                  "bg-sky-50 text-sky-700 dark:bg-sky-950/40 dark:text-sky-300",
                kind === "idle" && "bg-muted/40 text-muted-foreground"
              )}
            >
              {statusLabel()}
            </span>
          </div>
          <p className="mt-1 line-clamp-1 text-xs leading-5 text-muted-foreground">
            {task.description}
          </p>
          <div className="mt-2 flex flex-wrap gap-x-3 gap-y-1 text-xs text-muted-foreground">
            {task.score ? (
              <RewardItem icon={<Target className="h-3.5 w-3.5" />}>
                {t("user.tasks.reward.score", { score: task.score })}
              </RewardItem>
            ) : null}
            {task.exp ? (
              <RewardItem icon={<Zap className="h-3.5 w-3.5" />}>
                {t("user.tasks.reward.exp", { exp: task.exp })}
              </RewardItem>
            ) : null}
            {task.badgeId ? (
              <RewardItem icon={<Medal className="h-3.5 w-3.5" />}>
                {t("user.tasks.reward.badge")}
              </RewardItem>
            ) : null}
            {finishLabel() ? <span>{finishLabel()}</span> : null}
          </div>
        </div>
        <div className="w-full shrink-0 space-y-2 lg:w-56">
          <div className="flex items-center justify-between text-xs text-muted-foreground">
            <span>{progressLabel()}</span>
            <span className="font-medium text-foreground">
              {progressPercent(task)}%
            </span>
          </div>
          <div className="h-1.5 w-full overflow-hidden rounded-full bg-muted">
            <div
              className={cn(
                "h-full rounded-full",
                kind === "done" ? "bg-emerald-500" : "bg-primary"
              )}
              style={{ width: `${progressPercent(task)}%` }}
            />
          </div>
        </div>
        <div className="flex shrink-0 justify-end lg:w-24">{action}</div>
      </div>
    </div>
  )
}

function RewardItem({
  icon,
  children,
}: {
  icon: React.ReactNode
  children: React.ReactNode
}) {
  return (
    <span className="inline-flex items-center gap-1">
      {icon}
      {children}
    </span>
  )
}
