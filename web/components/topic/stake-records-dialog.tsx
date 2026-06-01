"use client"

import { Flame, TrendingUp, Clock, ArrowLeftRight } from "lucide-react"
import { useState, useEffect } from "react"
import { toast } from "sonner"

import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { apiFetch } from "@/lib/api/client"
import type { StakeRecord, HeatPointsQuota } from "@/lib/api/types"
import { FlameLevel } from "@/components/topic/flame-level"

interface StakeRecordsDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function StakeRecordsDialog({ open, onOpenChange }: StakeRecordsDialogProps) {
  const [loading, setLoading] = useState(false)
  const [records, setRecords] = useState<StakeRecord[]>([])
  const [quota, setQuota] = useState<HeatPointsQuota | null>(null)

  useEffect(() => {
    if (open) {
      loadRecords()
      loadQuota()
    }
  }, [open])

  const loadRecords = async () => {
    setLoading(true)
    try {
      const data = await apiFetch<StakeRecord[]>("/api/stake/records", {
        params: { limit: 50 },
      })
      setRecords(data)
    } catch (error) {
      toast.error("加载记录失败")
    } finally {
      setLoading(false)
    }
  }

  const loadQuota = async () => {
    try {
      const data = await apiFetch<HeatPointsQuota>("/api/stake/quota")
      setQuota(data)
    } catch (error) {
      // 静默失败
    }
  }

  const handleRedeem = async (stakeId: number) => {
    try {
      await apiFetch(`/api/stake/redeem/${stakeId}`, { method: "POST" })
      toast.success("赎回成功，收益将在次日结算后到账")
      loadRecords()
    } catch (error: any) {
      toast.error(error.message || "赎回失败")
    }
  }

  const getStatusBadge = (status: number) => {
    const statusMap = {
      0: { label: "质押中", variant: "default" as const },
      1: { label: "已锁定", variant: "secondary" as const },
      2: { label: "已赎回", variant: "outline" as const },
    }
    const config = statusMap[status as keyof typeof statusMap]
    return <Badge variant={config.variant}>{config.label}</Badge>
  }

  const getFlameLevelLabel = (level: number) => {
    const labels = ["无", "🔥", "🔥🔥", "🔥🔥🔥", "🔥🔥🔥🔥", "🔥🔥🔥🔥🔥"]
    return labels[level] || ""
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-4xl max-h-[80vh]">
        <DialogHeader>
          <DialogTitle>我的质押记录</DialogTitle>
          <DialogDescription>
            查看热度点质押记录和收益情况
          </DialogDescription>
        </DialogHeader>

        <div className="mt-4 grid gap-4">
          {/* 统计卡片 */}
          <div className="grid grid-cols-2 gap-4">
            <div className="rounded-lg border bg-card p-4">
              <div className="flex items-center gap-2">
                <Flame className="h-5 w-5 text-orange-500" />
                <div>
                  <p className="text-sm text-muted-foreground">当前持有</p>
                  <p className="text-2xl font-bold">{quota?.heatPoints || 0}</p>
                </div>
              </div>
            </div>
            <div className="rounded-lg border bg-card p-4">
              <div className="flex items-center gap-2">
                <Clock className="h-5 w-5 text-blue-500" />
                <div>
                  <p className="text-sm text-muted-foreground">今日剩余</p>
                  <p className="text-2xl font-bold">
                    {quota?.remainingQuota || 0}/{quota?.totalQuota || 3}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <Tabs defaultValue="active">
            <TabsList>
              <TabsTrigger value="active">进行中</TabsTrigger>
              <TabsTrigger value="history">历史记录</TabsTrigger>
            </TabsList>

            <TabsContent value="active">
              <ScrollArea className="h-[400px]">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>主题</TableHead>
                      <TableHead>质押</TableHead>
                      <TableHead>火焰等级</TableHead>
                      <TableHead>状态</TableHead>
                      <TableHead>操作</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {records
                      .filter((r) => r.status !== 2 && r.settleStatus !== 2)
                      .map((record) => (
                        <TableRow key={record.id}>
                          <TableCell className="max-w-[200px]">
                            <a
                              href={`/topic/${record.topicId}`}
                              className="text-blue-600 hover:underline truncate block"
                              target="_blank"
                              rel="noopener noreferrer"
                            >
                              #{record.topicId} {record.topicTitle}
                            </a>
                          </TableCell>
                          <TableCell className="font-medium">
                            {record.stakedPoints} 点
                          </TableCell>
                          <TableCell>
                            <FlameLevel level={record.flameLevel} />
                          </TableCell>
                          <TableCell>{getStatusBadge(record.status)}</TableCell>
                          <TableCell>
                            {record.settleStatus === 0 && (
                              <Button
                                variant="outline"
                                size="sm"
                                onClick={() => handleRedeem(record.id)}
                              >
                                <ArrowLeftRight className="h-4 w-4 mr-1" />
                                赎回
                              </Button>
                            )}
                            {record.settleStatus === 1 && (
                              <Badge variant="secondary">锁定中</Badge>
                            )}
                          </TableCell>
                        </TableRow>
                      ))}
                    {records.filter((r) => r.status !== 2 && r.settleStatus !== 2)
                      .length === 0 && (
                      <TableRow>
                        <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                          暂无进行中的质押
                        </TableCell>
                      </TableRow>
                    )}
                  </TableBody>
                </Table>
              </ScrollArea>
            </TabsContent>

            <TabsContent value="history">
              <ScrollArea className="h-[400px]">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>主题</TableHead>
                      <TableHead>质押</TableHead>
                      <TableHead>收益</TableHead>
                      <TableHead>赎回日期</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {records
                      .filter((r) => r.status === 2 || r.settleStatus === 2)
                      .map((record) => (
                        <TableRow key={record.id}>
                          <TableCell className="max-w-[200px]">
                            <a
                              href={`/topic/${record.topicId}`}
                              className="text-blue-600 hover:underline truncate block"
                              target="_blank"
                              rel="noopener noreferrer"
                            >
                              #{record.topicId} {record.topicTitle}
                            </a>
                          </TableCell>
                          <TableCell>{record.stakedPoints} 点</TableCell>
                          <TableCell>
                            {record.unsettledInterest && record.unsettledInterest > 0 ? (
                              <span className="text-green-600 font-medium">
                                +{record.unsettledInterest}
                              </span>
                            ) : (
                              <span className="text-muted-foreground">-</span>
                            )}
                          </TableCell>
                          <TableCell className="text-muted-foreground">
                            {record.redeemedAt ? new Date(record.redeemedAt).toLocaleDateString() : "-"}
                          </TableCell>
                        </TableRow>
                      ))}
                    {records.filter((r) => r.status === 2 || r.settleStatus === 2)
                      .length === 0 && (
                      <TableRow>
                        <TableCell colSpan={4} className="text-center py-8 text-muted-foreground">
                          暂无历史记录
                        </TableCell>
                      </TableRow>
                    )}
                  </TableBody>
                </Table>
              </ScrollArea>
            </TabsContent>
          </Tabs>
        </div>
      </DialogContent>
    </Dialog>
  )
}
