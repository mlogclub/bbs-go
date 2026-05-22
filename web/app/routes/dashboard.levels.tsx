"use client"

import * as React from "react"
import { PlusIcon, RefreshCwIcon, SaveIcon, Trash2Icon } from "lucide-react"

import { useCurrentUser } from "@/components/app/app-provider"
import { ErrorPage } from "@/components/common/error-page"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { adminList, adminPostJson, type AdminRecord } from "@/lib/api/admin"
import { userHasPermission } from "@/lib/auth/roles"
import { useI18n } from "@/lib/i18n/provider"
import { PERMISSIONS } from "@/lib/auth/permissions.generated"

type LevelRow = {
  id?: number
  level: number
  needExp: number
  title: string
  status?: number
}

function normalizeRow(record: AdminRecord): LevelRow {
  return {
    id: Number(record.id || 0) || undefined,
    level: Number(record.level || 0),
    needExp: Number(record.needExp || 0),
    title: String(record.title || ""),
    status: Number(record.status || 0),
  }
}

export default function DashboardLevelsRoute() {
  const { t } = useI18n()
  const currentUser = useCurrentUser()
  const [rows, setRows] = React.useState<LevelRow[]>([])
  const [loading, setLoading] = React.useState(true)
  const [saving, setSaving] = React.useState(false)
  const [error, setError] = React.useState<string | null>(null)
  const canView = userHasPermission(currentUser, PERMISSIONS.DASHBOARD_LEVEL_VIEW)
  const canUpdate = userHasPermission(currentUser, PERMISSIONS.DASHBOARD_LEVEL_UPDATE)

  const load = React.useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await adminList("/api/admin/level-config/list", {
        page: 1,
        limit: 200,
      })
      setRows((data.results || []).map(normalizeRow))
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("dashboard.errors.loadFailed")
      )
    } finally {
      setLoading(false)
    }
  }, [t])

  React.useEffect(() => {
    void load()
  }, [load])

  function updateRow(index: number, patch: Partial<LevelRow>) {
    setRows((current) =>
      current.map((row, rowIndex) =>
        rowIndex === index ? { ...row, ...patch } : row
      )
    )
  }

  function addLevel() {
    if (!canUpdate) return
    setRows((current) => [
      ...current,
      {
        level: current.length
          ? Math.max(...current.map((row) => row.level)) + 1
          : 1,
        needExp: 0,
        title: "",
        status: 0,
      },
    ])
  }

  function removeLevel(index: number) {
    if (!canUpdate) return
    setRows((current) => current.filter((_, rowIndex) => rowIndex !== index))
  }

  async function save() {
    if (!canUpdate) return
    const sorted = [...rows].sort((a, b) => a.level - b.level)
    for (let index = 0; index < sorted.length; index += 1) {
      const row = sorted[index]
      if (row.level !== index + 1) {
        setError(t("dashboard.errors.levelMustBeContinuous"))
        return
      }
      if (!row.title.trim()) {
        setError(t("dashboard.errors.required"))
        return
      }
      if (row.needExp < 0) {
        setError(t("dashboard.errors.minValue", { min: 0 }))
        return
      }
      if (index > 0 && row.needExp <= sorted[index - 1].needExp) {
        setError(t("dashboard.errors.levelExpMustIncrease"))
        return
      }
    }

    setSaving(true)
    setError(null)
    try {
      await adminPostJson("/api/admin/level-config/save_all", rows)
      await load()
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("dashboard.errors.saveFailed")
      )
    } finally {
      setSaving(false)
    }
  }

  if (!canView) {
    return <ErrorPage statusCode={403} />
  }

  return (
    <div className="flex flex-1 flex-col gap-4 p-4 pt-4 md:p-6">
      <section className="rounded-lg border bg-[var(--dashboard-panel)] p-3 text-card-foreground shadow-xs">
        <div className="flex flex-wrap justify-end gap-3">
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="icon"
              onClick={() => void load()}
              disabled={loading}
            >
              <RefreshCwIcon />
              <span className="sr-only">{t("dashboard.actions.refresh")}</span>
            </Button>
            {canUpdate ? (
              <>
                <Button variant="outline" onClick={addLevel}>
                  <PlusIcon />
                  {t("dashboard.actions.create")}
                </Button>
                <Button onClick={() => void save()} disabled={loading || saving}>
                  <SaveIcon />
                  {t("dashboard.actions.save")}
                </Button>
              </>
            ) : null}
          </div>
        </div>

        {error ? (
          <div className="mt-3 rounded-md border border-destructive/25 bg-destructive/10 px-3 py-2 text-sm text-destructive">
            {error}
          </div>
        ) : null}
      </section>

      <section className="overflow-hidden rounded-lg border bg-[var(--dashboard-panel)] shadow-xs">
        <div className="overflow-x-auto">
          <table className="w-full min-w-[720px] text-sm">
            <thead className="bg-[var(--dashboard-panel-muted)] text-muted-foreground">
              <tr>
                <th className="h-10 px-3 text-left text-xs font-semibold tracking-wide uppercase">
                  {t("dashboard.fields.level")}
                </th>
                <th className="h-10 px-3 text-left text-xs font-semibold tracking-wide uppercase">
                  {t("dashboard.fields.needExp")}
                </th>
                <th className="h-10 px-3 text-left text-xs font-semibold tracking-wide uppercase">
                  {t("dashboard.fields.title")}
                </th>
                <th className="h-10 w-20 px-3 text-right text-xs font-semibold tracking-wide uppercase">
                  {t("dashboard.actions.title")}
                </th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr>
                  <td
                    colSpan={4}
                    className="h-24 px-3 text-center text-muted-foreground"
                  >
                    {t("dashboard.loading")}
                  </td>
                </tr>
              ) : (
                rows.map((row, index) => (
                  <tr
                    key={row.id ?? index}
                    className="border-t transition-colors hover:bg-muted/30"
                  >
                    <td className="h-11 px-3 py-2 align-middle">
                      <Input
                        type="number"
                        value={row.level}
                        disabled={!canUpdate}
                        onChange={(event) =>
                          updateRow(index, {
                            level: Number(event.target.value),
                          })
                        }
                      />
                    </td>
                    <td className="h-11 px-3 py-2 align-middle">
                      <Input
                        type="number"
                        value={row.needExp}
                        disabled={!canUpdate}
                        onChange={(event) =>
                          updateRow(index, {
                            needExp: Number(event.target.value),
                          })
                        }
                      />
                    </td>
                    <td className="h-11 px-3 py-2 align-middle">
                      <Input
                        value={row.title}
                        disabled={!canUpdate}
                        onChange={(event) =>
                          updateRow(index, { title: event.target.value })
                        }
                      />
                    </td>
                    <td className="h-11 px-3 py-2 text-right align-middle">
                      {canUpdate ? (
                        <Button
                          size="icon-sm"
                          variant="destructive"
                          onClick={() => removeLevel(index)}
                        >
                          <Trash2Icon />
                          <span className="sr-only">
                            {t("dashboard.actions.delete")}
                          </span>
                        </Button>
                      ) : null}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </section>
    </div>
  )
}
