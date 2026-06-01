import { Flame } from "lucide-react"

interface FlameLevelProps {
  level: number // 0-5
  className?: string
}

/**
 * 火焰等级图标组件
 * level 0: 无火焰
 * level 1: 🔥
 * level 2: 🔥🔥
 * level 3: 🔥🔥🔥
 * level 4: 🔥🔥🔥🔥
 * level 5: 🔥🔥🔥🔥🔥
 */
export function FlameLevel({ level, className = "" }: FlameLevelProps) {
  if (level <= 0) {
    return null
  }

  // 限制最大等级为 5
  const normalizedLevel = Math.min(Math.max(level, 1), 5)

  // 火焰颜色根据等级变化
  const colors = {
    1: "text-orange-400",
    2: "text-orange-500",
    3: "text-orange-600",
    4: "text-red-500",
    5: "text-red-600 animate-pulse",
  }

  return (
    <div className={`flex items-center gap-0.5 ${className}`} title={`热度等级：${normalizedLevel}`}>
      {Array.from({ length: normalizedLevel }).map((_, i) => (
        <Flame
          key={i}
          className={`h-4 w-4 ${colors[normalizedLevel as keyof typeof colors]} fill-current`}
        />
      ))}
    </div>
  )
}
