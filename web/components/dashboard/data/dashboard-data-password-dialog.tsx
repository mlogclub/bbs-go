"use client"

import { CopyIcon } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { msgSuccess } from "@/lib/toast"

export function DashboardDataPasswordDialog({
  password,
  title,
  passwordLabel,
  copyLabel,
  copiedMessage,
  cancelLabel,
  onClose,
}: {
  password: string | null
  title: string
  passwordLabel: string
  copyLabel: string
  copiedMessage: string
  cancelLabel: string
  onClose: () => void
}) {
  if (!password) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-background/80 p-4 backdrop-blur-sm">
      <div className="w-full max-w-md rounded-xl border bg-card p-4 shadow-lg">
        <div className="mb-4 flex items-center justify-between gap-3">
          <h2 className="text-lg font-semibold">{title}</h2>
          <Button type="button" variant="ghost" onClick={onClose}>
            {cancelLabel}
          </Button>
        </div>
        <div className="grid gap-3">
          <Label>{passwordLabel}</Label>
          <div className="flex gap-2">
            <Input readOnly value={password} />
            <Button
              type="button"
              variant="outline"
              onClick={() => {
                void navigator.clipboard?.writeText(password)
                msgSuccess(copiedMessage)
              }}
            >
              <CopyIcon />
              {copyLabel}
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}

