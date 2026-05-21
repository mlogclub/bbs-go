"use client"

import * as React from "react"
import { ChevronLeft, ChevronRight, X } from "lucide-react"

import { cn } from "@/lib/utils"

type PreviewImage = string | { src?: string; preview?: string }

function imageUrl(image: PreviewImage) {
  if (typeof image === "string") {
    return image
  }
  return image.src || image.preview || ""
}

function normalizeImages(images: PreviewImage[]) {
  return images.map(imageUrl).filter(Boolean)
}

function ImagePreviewDialog({
  images,
  initialIndex,
  onClose,
}: {
  images: string[]
  initialIndex: number
  onClose: () => void
}) {
  const [currentIndex, setCurrentIndex] = React.useState(initialIndex)
  const touchStartRef = React.useRef<{ x: number; y: number } | null>(null)
  const current = images[currentIndex] || images[0] || ""

  const goPrevious = React.useCallback(() => {
    setCurrentIndex((index) => (index > 0 ? index - 1 : images.length - 1))
  }, [images.length])

  const goNext = React.useCallback(() => {
    setCurrentIndex((index) => (index < images.length - 1 ? index + 1 : 0))
  }, [images.length])

  React.useEffect(() => {
    function onKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") {
        onClose()
      } else if (event.key === "ArrowLeft") {
        goPrevious()
      } else if (event.key === "ArrowRight") {
        goNext()
      }
    }

    document.addEventListener("keydown", onKeyDown)
    return () => document.removeEventListener("keydown", onKeyDown)
  }, [goNext, goPrevious, onClose])

  React.useEffect(() => {
    const previous = document.body.style.overflow
    document.body.style.overflow = "hidden"
    return () => {
      document.body.style.overflow = previous
    }
  }, [])

  if (!current) {
    return null
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <button
        type="button"
        aria-label="close"
        className="fixed inset-0 bg-black/80"
        onClick={onClose}
      />
      <button
        type="button"
        aria-label="close"
        className="fixed top-4 right-4 z-[60] flex h-10 w-10 cursor-pointer items-center justify-center rounded-sm bg-white/20 text-white opacity-70 transition-opacity hover:opacity-100 focus:ring-2 focus:ring-white/50 focus:outline-none"
        onClick={onClose}
      >
        <X className="h-5 w-5" />
      </button>
      <div
        className="relative z-10 flex h-full min-h-[50vh] w-full touch-pan-y items-center justify-center"
        style={{ touchAction: "pan-y" }}
        onTouchStart={(event) => {
          const touch = event.touches[0]
          if (touch) {
            touchStartRef.current = { x: touch.clientX, y: touch.clientY }
          }
        }}
        onTouchEnd={(event) => {
          const start = touchStartRef.current
          const touch = event.changedTouches[0]
          touchStartRef.current = null
          if (!start || !touch) {
            return
          }

          const deltaX = start.x - touch.clientX
          const deltaY = start.y - touch.clientY
          if (Math.abs(deltaX) > Math.abs(deltaY) && Math.abs(deltaX) > 50) {
            if (deltaX > 0) {
              goNext()
            } else {
              goPrevious()
            }
          }
        }}
      >
        {images.length > 1 ? (
          <button
            type="button"
            aria-label="previous"
            className="absolute top-1/2 left-4 z-10 hidden -translate-y-1/2 items-center justify-center rounded-full bg-white/20 p-3 text-white transition-colors hover:bg-white/30 focus:ring-2 focus:ring-white/50 focus:outline-none md:flex"
            onClick={goPrevious}
          >
            <ChevronLeft className="h-6 w-6" />
          </button>
        ) : null}        <img
          src={current}
          alt=""
          className="max-h-full max-w-full select-none object-contain"
          style={{ maxHeight: "90vh", maxWidth: "90vw" }}
          draggable={false}
          onClick={(event) => event.stopPropagation()}
        />
        {images.length > 1 ? (
          <button
            type="button"
            aria-label="next"
            className="absolute top-1/2 right-4 z-10 hidden -translate-y-1/2 items-center justify-center rounded-full bg-white/20 p-3 text-white transition-colors hover:bg-white/30 focus:ring-2 focus:ring-white/50 focus:outline-none md:flex"
            onClick={goNext}
          >
            <ChevronRight className="h-6 w-6" />
          </button>
        ) : null}
      </div>
      {images.length > 1 ? (
        <>
          <div className="fixed bottom-4 left-1/2 z-10 -translate-x-1/2 rounded bg-white/20 px-3 py-1 text-sm text-white">
            {currentIndex + 1} / {images.length}
          </div>
          <div className="fixed top-4 left-1/2 z-10 -translate-x-1/2 rounded bg-white/20 px-3 py-1 text-xs text-white opacity-70 md:hidden">
            Swipe to switch
          </div>
        </>
      ) : null}
    </div>
  )
}

export function PreviewableImage({
  src,
  previewSrcList,
  initialIndex = 0,
  className,
  alt = "",
  ...props
}: React.ImgHTMLAttributes<HTMLImageElement> & {
  src: string
  previewSrcList?: PreviewImage[]
  initialIndex?: number
}) {
  const images = normalizeImages(previewSrcList?.length ? previewSrcList : [src])
  const [openIndex, setOpenIndex] = React.useState<number | null>(null)

  return (
    <>      <img
        {...props}
        src={src}
        alt={alt}
        className={cn(
          "cursor-pointer transition-transform hover:scale-105",
          className
        )}
        onClick={(event) => {
          props.onClick?.(event)
          setOpenIndex(initialIndex)
        }}
      />
      {openIndex !== null ? (
        <ImagePreviewDialog
          images={images}
          initialIndex={openIndex}
          onClose={() => setOpenIndex(null)}
        />
      ) : null}
    </>
  )
}

export function HtmlImagePreview({
  html,
  className,
}: {
  html: string
  className?: string
}) {
  const rootRef = React.useRef<HTMLDivElement>(null)
  const [preview, setPreview] = React.useState<{
    images: string[]
    index: number
  } | null>(null)

  return (
    <>
      <div
        ref={rootRef}
        className={className}
        dangerouslySetInnerHTML={{ __html: html }}
        onClick={(event) => {
          const target = event.target
          if (!(target instanceof HTMLImageElement)) {
            return
          }

          const images = Array.from(
            rootRef.current?.querySelectorAll("img") || []
          )
            .map((image) => image.currentSrc || image.src)
            .filter(Boolean)
          const src = target.currentSrc || target.src
          const index = Math.max(0, images.indexOf(src))
          setPreview({ images: images.length ? images : [src], index })
        }}
      />
      {preview ? (
        <ImagePreviewDialog
          images={preview.images}
          initialIndex={preview.index}
          onClose={() => setPreview(null)}
        />
      ) : null}
    </>
  )
}
