"use client"

import { Flame, TrendingUp, Clock, Zap } from "lucide-react"
import { useState, useEffect } from "react"
import { toast } from "sonner"

import { apiClient } from "@/lib/api"
import type { HeatPointsQuota } from "@/lib/api/types"

export function HeatPointsPanel() {
  const [quota, setQuota] = useState<HeatPointsQuota | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadQuota()
    
    // 监听质押更新事件
    const handler = () => loadQuota()
    window.addEventListener("heat-stake-updated" as any, handler)
    return () => window.removeEventListener("heat-stake-updated" as any, handler)
  }, [])

  const loadQuota = async () => {
    try {
      const data = await apiClient.get<HeatPointsQuota>("/api/stake/quota")
      setQuota(data)
    } catch (error) {
      toast.error("获取热度点数据失败")
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="animate-pulse text-muted-foreground">加载中...</div>
  }

  const utilization = quota && quota.heatPoints > 0
    ? Math.round(((quota.heatPoints - (quota.stakedPoints || 0)) / quota.heatPoints) * 100)
    : 0

  return (
    <div className="space-y-4">
      {/* 余额卡片 */}
      <div className="rounded-lg border bg-gradient-to-br from-orange-50 to-red-50 p-6">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm text-muted-foreground">热度点余额</p>
            <p className="text-4xl font-bold text-orange-600">{quota?.heatPoints || 0}</p>
          </div>
          <Flame className="h-16 w-16 text-orange-400 opacity-50" />
        </div>
      </div>

      {/* 统计网格 */}
      <div className="grid grid-cols-2 gap-4">
        <div className="rounded-lg border p-4">
          <div className="flex items-center gap-2">
            <Zap className="h-5 w-5 text-yellow-500" />
            <div>
              <p className="text-xs text-muted-foreground">已质押</p>
              <p className="text-lg font-semibold">{quota?.stakedPoints || 0}</p>
            </div>
          </div>
        </div>

        <div className="rounded-lg border p-4">
          <div className="flex items-center gap-2">
            <TrendingUp className="h-5 w-5 text-green-500" />
            <div>
              <p className="text-xs text-muted-foreground">利用率</p>
              <p className="text-lg font-semibold">{utilization}%</p>
            </div>
          </div>
        </div>

        <div className="rounded-lg border p-4">
          <div className="flex items-center gap-2">
            <Clock className="h-5 w-5 text-blue-500" />
            <div>
              <p className="text-xs text-muted-foreground">今日剩余</p>
              <p className="text-lg font-semibold">
                {quota?.remainingQuota || 0}/{quota?.totalQuota || 3}
              </p>
            </div>
          </div>
        </div>

        <div className="rounded-lg border p-4">
          <div className="flex items-center gap-2">
            <Flame className="h-5 w-5 text-orange-500" />
            <div>
              <p className="text-xs text-muted-foreground">待结算</p>
              <p className="text-lg font-semibold">{quota?.pendingInterest || 0}</p>
            </div>
          </div>
        </div>
      </div>

      {/* 提示信息 */}
      <div className="rounded-lg bg-blue-50 p-4 text-sm text-blue-700">
        <p className="font-medium mb-2">💡 热度点小贴士</p>
        <ul className="list-disc list-inside space-y-1 text-xs">
          <li>质押可获得每日收益，利率随热度变化</li>
          <li>当日质押不可赎回，需等到次日结算</li>
          <li>长期不质押会以每日 2% 速率衰减</li>
          <li>火焰等级越高，收益波动越大</li>
        </ul>
      </div>
    </div>
  )
}
