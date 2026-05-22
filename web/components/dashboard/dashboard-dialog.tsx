"use client"

import type * as React from "react"
import { useState } from "react"

import { cn } from "@/lib/utils"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Maximize2Icon, Minimize2Icon, XIcon } from "lucide-react"

const dialogSizeClassName = {
  sm: "max-w-md sm:max-w-md",
  md: "max-w-xl sm:max-w-xl",
  lg: "max-w-2xl sm:max-w-2xl",
  xl: "max-w-4xl sm:max-w-4xl",
  xxl: "max-w-5xl sm:max-w-5xl",
} as const

type DashboardDialogSize = keyof typeof dialogSizeClassName

type DashboardDialogProps = React.ComponentProps<typeof Dialog> & {
  title: React.ReactNode
  description?: React.ReactNode
  size?: DashboardDialogSize
  children: React.ReactNode
  footer?: React.ReactNode
  contentClassName?: string
  headerClassName?: string
  bodyClassName?: string
  footerClassName?: string
  showCloseButton?: boolean
  closeOnEsc?: boolean
  allowFullscreen?: boolean
  defaultFullscreen?: boolean
  bodyScrollable?: boolean
}

function DashboardDialog({
  open,
  onOpenChange,
  title,
  description,
  size = "md",
  children,
  footer,
  contentClassName,
  headerClassName,
  bodyClassName,
  footerClassName,
  showCloseButton = true,
  closeOnEsc = false,
  allowFullscreen = true,
  defaultFullscreen = false,
  bodyScrollable = true,
}: DashboardDialogProps) {
  const [fullscreen, setFullscreen] = useState(defaultFullscreen)

  function handleOpenChange(nextOpen: boolean) {
    if (!nextOpen) {
      setFullscreen(defaultFullscreen)
    }
    onOpenChange?.(nextOpen)
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <style>{`
        .dashboard-dialog-native-scrollbar {
          scrollbar-width: thin;
          scrollbar-color: var(--border) transparent;
        }

        .dashboard-dialog-native-scrollbar::-webkit-scrollbar {
          width: 10px;
        }

        .dashboard-dialog-native-scrollbar::-webkit-scrollbar-track {
          background: transparent;
        }

        .dashboard-dialog-native-scrollbar::-webkit-scrollbar-thumb {
          background: var(--border);
          border: 2px solid transparent;
          border-radius: 9999px;
          background-clip: content-box;
        }

        .dashboard-dialog-native-scrollbar::-webkit-scrollbar-thumb:hover {
          background: color-mix(
            in srgb,
            var(--border) 80%,
            var(--foreground)
          );
          border: 2px solid transparent;
          background-clip: content-box;
        }
      `}</style>
      <DialogContent
        className={cn(
          "flex max-h-[calc(100vh-2rem)] flex-col gap-0 overflow-hidden p-0",
          fullscreen
            ? "top-5 left-5 h-[calc(100vh-40px)] max-h-[calc(100vh-40px)] w-[calc(100vw-40px)] max-w-[calc(100vw-40px)] translate-x-0 translate-y-0 rounded-xl sm:max-w-[calc(100vw-40px)]"
            : dialogSizeClassName[size],
          contentClassName,
        )}
        showCloseButton={false}
        onEscapeKeyDown={(event) => {
          if (!closeOnEsc) {
            event.preventDefault()
          }
        }}
      >
        {(allowFullscreen || showCloseButton) && (
          <div className="absolute top-2 right-2 z-10 flex items-center gap-1">
            {allowFullscreen ? (
              <Button
                type="button"
                variant="ghost"
                size="icon-sm"
                onClick={() => setFullscreen((value) => !value)}
              >
                {fullscreen ? <Minimize2Icon /> : <Maximize2Icon />}
                <span className="sr-only">
                  {fullscreen ? "退出全屏" : "全屏显示"}
                </span>
              </Button>
            ) : null}
            {showCloseButton ? (
              <DialogClose asChild>
                <Button type="button" variant="ghost" size="icon-sm">
                  <XIcon />
                  <span className="sr-only">关闭</span>
                </Button>
              </DialogClose>
            ) : null}
          </div>
        )}
        <DialogHeader className={cn("shrink-0 px-6 py-3", headerClassName)}>
          <DialogTitle>{title}</DialogTitle>
          {description ? (
            <DialogDescription>{description}</DialogDescription>
          ) : null}
        </DialogHeader>
        {bodyScrollable ? (
          <div
            className={cn(
              "dashboard-dialog-native-scrollbar min-h-0 flex-1 overflow-y-auto",
              bodyClassName
            )}
          >
            <div className="w-full space-y-4 p-6">{children}</div>
          </div>
        ) : (
          <div className={cn("min-h-0 flex-1", bodyClassName)}>{children}</div>
        )}
        {footer ? (
          <DialogFooter
            className={cn("mx-0 mb-0 shrink-0 px-6 py-4", footerClassName)}
          >
            {footer}
          </DialogFooter>
        ) : null}
      </DialogContent>
    </Dialog>
  )
}

export { DashboardDialog }
