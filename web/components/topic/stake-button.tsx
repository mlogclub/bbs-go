"use client"

import { Flame } from "lucide-react"
import { useState } from "react"
import { toast } from "sonner"

import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { apiFetch } from "@/lib/api/client"
import type { HeatPointsQuota } from "@/lib/api/types"
import { cn } from "@/lib/utils"

interface StakeButtonProps {
  topicId: number
  currentFlameLevel?: number
}

export function StakeButton({ topicId, currentFlameLevel = 0 }: StakeButtonProps) {
  const [open, setOpen] = useState(false)
  const [loading, setLoading] = useState(false)
  const [quota, setQuota] = useState<HeatPointsQuota | null>(null)
  const [amount, setAmount] = useState("")

  const handleOpen = async () => {
    try {
      const data = await apiFetch<HeatPointsQuota>("/api/stake/quota")
      setQuota(data)
      if (data.remainingQuota <= 0) {
        toast.error("今日质押次数已用完")
        return
      }
      setOpen(true)
    } catch (error) {
      toast.error("获取配额失败")
    }
  }

  const handleSubmit = async () => {
    if (!amount || parseInt(amount) < 1) {
      toast.error("请输入有效的质押数量")
      return
    }

    const points = parseInt(amount)
    if (points > (quota?.heatPoints || 0)) {
      toast.error("热度点余额不足")
      return
    }

    setLoading(true)
    try {
      await apiFetch("/api/stake/create", {
        method: "POST",
        body: { topicId, heatPoints: points },
      })
      toast.success("质押成功")
      setOpen(false)
      setAmount("")
      // 刷新页面或触发事件更新 UI
      window.dispatchEvent(new CustomEvent("heat-stake-updated"))
    } catch (error: any) {
      toast.error(error.message || "质押失败")
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      <Button
        variant="outline"
        size="sm"
        onClick={handleOpen}
        className={cn(
          "gap-1.5",
          currentFlameLevel >= 3 && "text-red-600 border-red-200 hover:bg-red-50"
        )}
      >
        <Flame className="h-4 w-4" />
        <span>质押</span>
      </Button>

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>质押热度点</DialogTitle>
            <DialogDescription>
              质押热度点到主题，获得被动收益。当日质押不可赎回，需等到次日结算。
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label>质押数量</Label>
              <Input
                type="number"
                min="1"
                placeholder="输入质押数量"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
              />
              <p className="text-xs text-muted-foreground">
                可用余额：{quota?.heatPoints || 0} 点 | 今日剩余次数：{quota?.remainingQuota || 0}/
                {quota?.totalQuota || 3}
              </p>
            </div>

            {quota && (quota.heatPoints || 0) > 0 && (
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setAmount(Math.min(10, quota.heatPoints).toString())}
                >
                  +10
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setAmount(Math.min(50, quota.heatPoints).toString())}
                >
                  +50
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setAmount(quota.heatPoints.toString())}
                >
                  全部
                </Button>
              </div>
            )}
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setOpen(false)} disabled={loading}>
              取消
            </Button>
            <Button onClick={handleSubmit} disabled={loading}>
              {loading ? "质押中..." : "确认质押"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  )
}
