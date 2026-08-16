"use client"

import * as React from "react"
import { createPortal } from "react-dom"
import { Extension, Node, mergeAttributes, type CommandProps, type Editor, type Range } from "@tiptap/core"
import { EditorContent, NodeViewWrapper, ReactNodeViewRenderer, ReactRenderer, type ReactNodeViewProps, useEditor } from "@tiptap/react"
import StarterKit from "@tiptap/starter-kit"
import Link from "@tiptap/extension-link"
import Placeholder from "@tiptap/extension-placeholder"
import Underline from "@tiptap/extension-underline"
import TextAlign from "@tiptap/extension-text-align"
import { TextStyle } from "@tiptap/extension-text-style"
import BackgroundColor from "@tiptap/extension-text-style/background-color"
import Color from "@tiptap/extension-color"
import TaskList from "@tiptap/extension-task-list"
import TaskItem from "@tiptap/extension-task-item"
import Typography from "@tiptap/extension-typography"
import HorizontalRule from "@tiptap/extension-horizontal-rule"
import Suggestion, { exitSuggestion, type SuggestionKeyDownProps, type SuggestionProps } from "@tiptap/suggestion"
import { PluginKey } from "prosemirror-state"
import {
  AlignCenter,
  AlignLeft,
  AlignRight,
  Bold,
  Check,
  Code,
  Code2,
  Heading1,
  Heading2,
  Heading3,
  ImageIcon,
  Italic,
  LinkIcon,
  List,
  ListOrdered,
  ListTodo,
  Maximize,
  Minimize,
  MinusSquare,
  Paintbrush,
  Palette,
  Pilcrow,
  Quote,
  Strikethrough,
  Underline as UnderlineIcon,
} from "lucide-react"

import { uploadEditorImage } from "@/components/editor/upload"
import { useI18n } from "@/lib/i18n/provider"
import { searchUsers } from "@/lib/api/users"
import type { SearchUser } from "@/lib/api/types"
import { useToastActions } from "@/lib/toast"
import { cn } from "@/lib/utils"

declare module "@tiptap/core" {
  interface Commands<ReturnType> {
    resizableImage: {
      setResizableImage: (options: { src: string; alt?: string; title?: string; width?: number; height?: number }) => ReturnType
    }
  }
}

const TEXT_COLOR_PALETTE = [
  "#c00000",
  "#ff0000",
  "#ffc000",
  "#ffff00",
  "#a5d610",
  "#00b050",
  "#00b0f0",
  "#0070c0",
  "#002060",
  "#7030a0",
  "#ffffff",
  "#000000",
  "#eeeeee",
  "#525252",
  "#1890ff",
  "#ff7875",
  "#52c41a",
  "#fa8c16",
  "#722ed1",
  "#eb2f96",
]

const BACKGROUND_COLOR_PALETTE = [
  "#FFCCCC",
  "#FFE6CC",
  "#FFFFCC",
  "#CCFFCC",
  "#CCFFFF",
  "#CCE5FF",
  "#E5CCFF",
  "#FFCCFF",
  "#F2F2F2",
  "#E6E6E6",
  "#FFD6CC",
  "#E5FFCC",
  "#CCFFE5",
  "#D6FFFF",
  "#FFE0F2",
  "#FFF0F0",
  "#FFF9E6",
  "#F0FFF0",
  "#F0FFFF",
  "#F5FFFA",
]

type Translate = (key: string) => string

type RichTextEditorLabels = {
  placeholder: string
  toolbar: {
    bold: string
    underline: string
    italic: string
    strike: string
    heading1: string
    heading2: string
    quote: string
    bulletList: string
    orderedList: string
    taskList: string
    alignLeft: string
    alignCenter: string
    alignRight: string
    textColor: string
    backgroundColor: string
    clearColor: string
    inlineCode: string
    codeBlock: string
    link: string
    image: string
    horizontalRule: string
    fullscreen: string
    exitFullscreen: string
    uploading: string
  }
  linkDialog: {
    textLabel: string
    urlLabel: string
    textPlaceholder: string
    urlPlaceholder: string
    confirm: string
    remove: string
    cancel: string
  }
  slash: {
    hintContinueTyping: string
    noMatchingCommand: string
    paragraph: { title: string; description: string }
    heading1: { title: string; description: string }
    heading2: { title: string; description: string }
    heading3: { title: string; description: string }
    bulletList: { title: string; description: string }
    orderedList: { title: string; description: string }
    taskList: { title: string; description: string }
    quote: { title: string; description: string }
    codeBlock: { title: string; description: string }
    horizontalRule: { title: string; description: string }
  }
}

function createEditorLabels(t: Translate): RichTextEditorLabels {
  const key = (path: string) => `component.richTextEditor.${path}`

  return {
    placeholder: t(key("placeholder")),
    toolbar: {
      bold: t(key("toolbar.bold")),
      underline: t(key("toolbar.underline")),
      italic: t(key("toolbar.italic")),
      strike: t(key("toolbar.strike")),
      heading1: t(key("toolbar.heading1")),
      heading2: t(key("toolbar.heading2")),
      quote: t(key("toolbar.quote")),
      bulletList: t(key("toolbar.bulletList")),
      orderedList: t(key("toolbar.orderedList")),
      taskList: t(key("toolbar.taskList")),
      alignLeft: t(key("toolbar.alignLeft")),
      alignCenter: t(key("toolbar.alignCenter")),
      alignRight: t(key("toolbar.alignRight")),
      textColor: t(key("toolbar.textColor")),
      backgroundColor: t(key("toolbar.backgroundColor")),
      clearColor: t(key("toolbar.clearColor")),
      inlineCode: t(key("toolbar.inlineCode")),
      codeBlock: t(key("toolbar.codeBlock")),
      link: t(key("toolbar.link")),
      image: t(key("toolbar.image")),
      horizontalRule: t(key("toolbar.horizontalRule")),
      fullscreen: t(key("toolbar.fullscreen")),
      exitFullscreen: t(key("toolbar.exitFullscreen")),
      uploading: t(key("toolbar.uploading")),
    },
    linkDialog: {
      textLabel: t(key("linkDialog.textLabel")),
      urlLabel: t(key("linkDialog.urlLabel")),
      textPlaceholder: t(key("linkDialog.textPlaceholder")),
      urlPlaceholder: t(key("linkDialog.urlPlaceholder")),
      confirm: t(key("linkDialog.confirm")),
      remove: t(key("linkDialog.remove")),
      cancel: t(key("linkDialog.cancel")),
    },
    slash: {
      hintContinueTyping: t(key("slash.hintContinueTyping")),
      noMatchingCommand: t(key("slash.noMatchingCommand")),
      paragraph: {
        title: t(key("slash.paragraph.title")),
        description: t(key("slash.paragraph.description")),
      },
      heading1: {
        title: t(key("slash.heading1.title")),
        description: t(key("slash.heading1.description")),
      },
      heading2: {
        title: t(key("slash.heading2.title")),
        description: t(key("slash.heading2.description")),
      },
      heading3: {
        title: t(key("slash.heading3.title")),
        description: t(key("slash.heading3.description")),
      },
      bulletList: {
        title: t(key("slash.bulletList.title")),
        description: t(key("slash.bulletList.description")),
      },
      orderedList: {
        title: t(key("slash.orderedList.title")),
        description: t(key("slash.orderedList.description")),
      },
      taskList: {
        title: t(key("slash.taskList.title")),
        description: t(key("slash.taskList.description")),
      },
      quote: {
        title: t(key("slash.quote.title")),
        description: t(key("slash.quote.description")),
      },
      codeBlock: {
        title: t(key("slash.codeBlock.title")),
        description: t(key("slash.codeBlock.description")),
      },
      horizontalRule: {
        title: t(key("slash.horizontalRule.title")),
        description: t(key("slash.horizontalRule.description")),
      },
    },
  }
}

const RESIZE_HANDLES = ["nw", "n", "ne", "e", "se", "s", "sw", "w"]

function ResizableImageView({
  node,
  selected,
  editor,
  updateAttributes,
  getPos,
}: ReactNodeViewProps) {
  const attrs = node.attrs as { src: string; alt?: string; title?: string; width?: number; height?: number }
  const imageRef = React.useRef<HTMLImageElement>(null)
  const [imageLoaded, setImageLoaded] = React.useState(false)
  const [size, setSize] = React.useState<{ width?: number; height?: number }>({
    width: attrs.width,
    height: attrs.height,
  })
  const aspectRatioRef = React.useRef(1)
  const resizeRef = React.useRef<{
    handle: string
    startX: number
    startWidth: number
    startHeight: number
  } | null>(null)

  function onImageLoad() {
    const image = imageRef.current
    if (!image) return
    setImageLoaded(true)
    aspectRatioRef.current = image.naturalWidth / image.naturalHeight || 1
    if (!size.width && !size.height) {
      const maxWidth = 600
      const width = image.naturalWidth > maxWidth ? maxWidth : image.naturalWidth
      const height = Math.round(width / aspectRatioRef.current)
      setSize({ width, height })
      updateAttributes({ width, height })
    }
  }

  function startResize(event: React.MouseEvent, handle: string) {
    event.preventDefault()
    event.stopPropagation()
    const image = imageRef.current
    if (!image) return
    resizeRef.current = {
      handle,
      startX: event.clientX,
      startWidth: size.width || image.offsetWidth,
      startHeight: size.height || image.offsetHeight,
    }
    editor.view.dom.classList.add("resizing-image")
  }

  React.useEffect(() => {
    function onMouseMove(event: MouseEvent) {
      const current = resizeRef.current
      if (!current) return
      event.preventDefault()
      const deltaX = event.clientX - current.startX
      let width = current.startWidth
      if (["se", "e", "s", "ne"].includes(current.handle)) {
        width = current.startWidth + deltaX
      } else {
        width = current.startWidth - deltaX
      }
      width = Math.max(50, Math.min(800, width))
      const height = Math.round(width / aspectRatioRef.current)
      setSize({ width, height })
    }

    function onMouseUp() {
      const current = resizeRef.current
      if (!current) return
      resizeRef.current = null
      editor.view.dom.classList.remove("resizing-image")
      updateAttributes({ width: size.width, height: size.height })
    }

    document.addEventListener("mousemove", onMouseMove)
    document.addEventListener("mouseup", onMouseUp)
    return () => {
      document.removeEventListener("mousemove", onMouseMove)
      document.removeEventListener("mouseup", onMouseUp)
    }
  }, [editor, size.height, size.width, updateAttributes])

  function selectImage() {
    const pos = getPos()
    if (typeof pos === "number") {
      editor.commands.setNodeSelection(pos)
    }
  }

  return (
    <NodeViewWrapper className={cn("resizable-image-wrapper", selected && "is-selected")}>      <img
        ref={imageRef}
        src={attrs.src}
        alt={attrs.alt || ""}
        title={attrs.title || ""}
        width={size.width}
        height={size.height}
        className="editor-image resizable"
        onLoad={onImageLoad}
        onClick={selectImage}
      />
      {selected && imageLoaded ? (
        <>
          <div className="selection-border">
            <div className="border-line border-top" />
            <div className="border-line border-right" />
            <div className="border-line border-bottom" />
            <div className="border-line border-left" />
          </div>
          {RESIZE_HANDLES.map((handle) => (
            <div key={handle} className={`resize-handle resize-handle-${handle}`} onMouseDown={(event) => startResize(event, handle)} />
          ))}
        </>
      ) : null}
    </NodeViewWrapper>
  )
}

const ResizableImage = Node.create({
  name: "resizableImage",
  group: "block",
  draggable: true,

  addAttributes() {
    return {
      src: { default: null },
      alt: { default: null },
      title: { default: null },
      width: {
        default: null,
        parseHTML: (element: HTMLElement) => {
          const width = element.getAttribute("width")
          return width ? Number.parseInt(width, 10) : null
        },
      },
      height: {
        default: null,
        parseHTML: (element: HTMLElement) => {
          const height = element.getAttribute("height")
          return height ? Number.parseInt(height, 10) : null
        },
      },
    }
  },

  parseHTML() {
    return [{ tag: "img[src]" }]
  },

  renderHTML({ HTMLAttributes }: { HTMLAttributes: Record<string, unknown> }) {
    return ["img", mergeAttributes(HTMLAttributes)]
  },

  addCommands() {
    return {
      setResizableImage:
        (options: { src: string; alt?: string; title?: string; width?: number; height?: number }) =>
        ({ commands }: CommandProps) =>
          commands.insertContent({
            type: this.name,
            attrs: options,
          }),
    }
  },

  addNodeView() {
    return ReactNodeViewRenderer(ResizableImageView)
  },
})

type SlashCommandItem = {
  title: string
  description: string
  aliases: string[]
  icon: React.ReactNode
  command: ({ editor, range }: { editor: Editor; range: Range }) => void
}

type SlashCommandMenuProps = SuggestionProps<SlashCommandItem, SlashCommandItem> & {
  labels: RichTextEditorLabels
}

type SlashCommandMenuHandle = {
  onKeyDown: (event: KeyboardEvent) => boolean
}

function slashItems(labels: RichTextEditorLabels, locale: string): SlashCommandItem[] {
  const zh = locale === "zh-CN"
  const aliases = zh
    ? {
        paragraph: ["p", "text", "paragraph", "正文", "段落"],
        heading1: ["h1", "一级标题", "heading1"],
        heading2: ["h2", "二级标题", "heading2"],
        heading3: ["h3", "三级标题", "heading3"],
        bulletList: ["ul", "bullet", "list", "无序列表"],
        orderedList: ["ol", "ordered", "numbered", "有序列表"],
        taskList: ["task", "todo", "checklist", "任务列表", "待办"],
        quote: ["quote", "blockquote", "引用", "引用文本"],
        codeBlock: ["code", "codeblock", "代码", "代码块"],
        horizontalRule: ["hr", "line", "divider", "分割线"],
      }
    : {
        paragraph: ["p", "text", "paragraph"],
        heading1: ["h1", "heading1"],
        heading2: ["h2", "heading2"],
        heading3: ["h3", "heading3"],
        bulletList: ["ul", "bullet", "list"],
        orderedList: ["ol", "ordered", "numbered"],
        taskList: ["task", "todo", "checklist"],
        quote: ["quote", "blockquote"],
        codeBlock: ["code", "codeblock"],
        horizontalRule: ["hr", "line", "divider"],
      }

  return [
    {
      title: labels.slash.paragraph.title,
      description: labels.slash.paragraph.description,
      aliases: aliases.paragraph,
      icon: <Pilcrow size={18} />,
      command: ({ editor, range }) => editor.chain().focus().deleteRange(range).setParagraph().run(),
    },
    {
      title: labels.slash.heading1.title,
      description: labels.slash.heading1.description,
      aliases: aliases.heading1,
      icon: <Heading1 size={18} />,
      command: ({ editor, range }) => editor.chain().focus().deleteRange(range).setNode("heading", { level: 1 }).run(),
    },
    {
      title: labels.slash.heading2.title,
      description: labels.slash.heading2.description,
      aliases: aliases.heading2,
      icon: <Heading2 size={18} />,
      command: ({ editor, range }) => editor.chain().focus().deleteRange(range).setNode("heading", { level: 2 }).run(),
    },
    {
      title: labels.slash.heading3.title,
      description: labels.slash.heading3.description,
      aliases: aliases.heading3,
      icon: <Heading3 size={18} />,
      command: ({ editor, range }) => editor.chain().focus().deleteRange(range).setNode("heading", { level: 3 }).run(),
    },
    {
      title: labels.slash.bulletList.title,
      description: labels.slash.bulletList.description,
      aliases: aliases.bulletList,
      icon: <List size={18} />,
      command: ({ editor, range }) => editor.chain().focus().deleteRange(range).toggleBulletList().run(),
    },
    {
      title: labels.slash.orderedList.title,
      description: labels.slash.orderedList.description,
      aliases: aliases.orderedList,
      icon: <ListOrdered size={18} />,
      command: ({ editor, range }) => editor.chain().focus().deleteRange(range).toggleOrderedList().run(),
    },
    {
      title: labels.slash.taskList.title,
      description: labels.slash.taskList.description,
      aliases: aliases.taskList,
      icon: <ListTodo size={18} />,
      command: ({ editor, range }) => editor.chain().focus().deleteRange(range).toggleTaskList().run(),
    },
    {
      title: labels.slash.quote.title,
      description: labels.slash.quote.description,
      aliases: aliases.quote,
      icon: <Quote size={18} />,
      command: ({ editor, range }) => editor.chain().focus().deleteRange(range).toggleBlockquote().run(),
    },
    {
      title: labels.slash.codeBlock.title,
      description: labels.slash.codeBlock.description,
      aliases: aliases.codeBlock,
      icon: <Code2 size={18} />,
      command: ({ editor, range }) => editor.chain().focus().deleteRange(range).toggleCodeBlock().run(),
    },
    {
      title: labels.slash.horizontalRule.title,
      description: labels.slash.horizontalRule.description,
      aliases: aliases.horizontalRule,
      icon: <MinusSquare size={18} />,
      command: ({ editor, range }) => editor.chain().focus().deleteRange(range).setHorizontalRule().run(),
    },
  ]
}

function updateSlashMenuPosition(element: HTMLElement, clientRect?: (() => DOMRect | null) | null) {
  const rect = clientRect?.()
  if (!rect) {
    return
  }
  const gap = 8
  const offset = 6
  const popupWidth = element.offsetWidth || 320
  const popupHeight = element.offsetHeight || 360
  let left = rect.left
  let top = rect.bottom + offset

  if (left + popupWidth + gap > window.innerWidth) {
    left = window.innerWidth - popupWidth - gap
  }
  if (left < gap) {
    left = gap
  }
  if (top + popupHeight + gap > window.innerHeight) {
    top = rect.top - popupHeight - offset
  }
  if (top < gap) {
    top = gap
  }

  element.style.left = `${left}px`
  element.style.top = `${top}px`
}

const SlashCommandMenu = React.forwardRef<SlashCommandMenuHandle, SlashCommandMenuProps>(function SlashCommandMenu({ items, query, command, labels }, ref) {
  const [selectedIndex, setSelectedIndex] = React.useState(0)
  const itemRefs = React.useRef<Array<HTMLButtonElement | null>>([])

  React.useEffect(() => {
    setSelectedIndex(0)
  }, [items, query])

  React.useEffect(() => {
    setSelectedIndex((index) => {
      if (!items.length) {
        return 0
      }
      return Math.min(index, items.length - 1)
    })
  }, [items])

  React.useEffect(() => {
    itemRefs.current[selectedIndex]?.scrollIntoView({ block: "nearest" })
  }, [selectedIndex])

  React.useImperativeHandle(
    ref,
    () => ({
      onKeyDown(event) {
        if (!items.length) {
          return false
        }
        if (event.key === "ArrowUp") {
          event.preventDefault()
          setSelectedIndex((index) => (index - 1 + items.length) % items.length)
          return true
        }
        if (event.key === "ArrowDown" || event.key === "Tab") {
          event.preventDefault()
          setSelectedIndex((index) => {
            if (event.shiftKey) {
              return (index - 1 + items.length) % items.length
            }
            return (index + 1) % items.length
          })
          return true
        }
        if (event.key === "Enter") {
          event.preventDefault()
          command(items[selectedIndex])
          return true
        }
        const shortcut = Number.parseInt(event.key, 10)
        if (Number.isFinite(shortcut) && shortcut >= 1 && shortcut <= 9 && items[shortcut - 1]) {
          event.preventDefault()
          command(items[shortcut - 1])
          return true
        }
        return false
      },
    }),
    [command, items, selectedIndex]
  )

  return (
    <div className="slash-commands" role="listbox" aria-label="Slash commands">
      <div className="search-hint">{items.length ? labels.slash.hintContinueTyping : labels.slash.noMatchingCommand}</div>
      {items.length ? (
        <div className="slash-items-container">
          {items.map((item, index) => (
            <button
              key={`${item.title}-${item.aliases[0] || index}`}
              type="button"
              role="option"
              aria-selected={index === selectedIndex}
              ref={(element) => {
                itemRefs.current[index] = element
              }}
              className={cn("slash-item", index === selectedIndex && "is-selected")}
              onMouseEnter={() => setSelectedIndex(index)}
              onMouseDown={(event) => {
                event.preventDefault()
                command(item)
              }}
            >
              <span className="item-icon">{item.icon}</span>
              <span className="item-content">
                <span className="item-title">
                  {item.title}
                  {item.aliases[0] ? <span className="item-aliases">/{item.aliases[0]}</span> : null}
                </span>
                <span className="item-description">{item.description}</span>
              </span>
            </button>
          ))}
        </div>
      ) : null}
    </div>
  )
})


type MentionUserItem = {
  username: string
  nickname: string
  avatar: string
}

type MentionMenuProps = SuggestionProps<MentionUserItem, MentionUserItem>

type MentionMenuHandle = {
  onKeyDown: (event: KeyboardEvent) => boolean
}

const MentionMenu = React.forwardRef<MentionMenuHandle, MentionMenuProps>(function MentionMenu({ items, query, command }, ref) {
  const [selectedIndex, setSelectedIndex] = React.useState(0)

  React.useEffect(() => {
    setSelectedIndex(0)
  }, [items, query])

  React.useEffect(() => {
    setSelectedIndex((index) => Math.min(index, Math.max(0, items.length - 1)))
  }, [items])

  React.useImperativeHandle(
    ref,
    () => ({
      onKeyDown(event) {
        if (!items.length) {
          return false
        }
        if (event.key === "ArrowUp") {
          event.preventDefault()
          setSelectedIndex((index) => (index - 1 + items.length) % items.length)
          return true
        }
        if (event.key === "ArrowDown" || event.key === "Tab") {
          event.preventDefault()
          setSelectedIndex((index) => {
            if (event.shiftKey) {
              return (index - 1 + items.length) % items.length
            }
            return (index + 1) % items.length
          })
          return true
        }
        if (event.key === "Enter") {
          event.preventDefault()
          command(items[selectedIndex])
          return true
        }
        return false
      },
    }),
    [command, items, selectedIndex]
  )

  return (
    <div className="mention-menu" role="listbox" aria-label="Mention users">
      {items.length ? (
        <div className="mention-items-container">
          {items.map((item, index) => (
            <button
              key={item.username}
              type="button"
              role="option"
              aria-selected={index === selectedIndex}
              className={cn("mention-item", index === selectedIndex && "is-selected")}
              onMouseEnter={() => setSelectedIndex(index)}
              onMouseDown={(event) => {
                event.preventDefault()
                command(item)
              }}
            >
              <img className="mention-avatar" src={item.avatar || "/default-avatar.png"} alt="" />
              <span className="mention-content">
                <span className="mention-nickname">{item.nickname}</span>
                <span className="mention-username">@{item.username}</span>
              </span>
            </button>
          ))}
        </div>
      ) : (
        <div className="mention-no-results">No users found</div>
      )}
    </div>
  )
})

function createMentionSuggestion() {
  return Extension.create({
    name: "at-mention",

    addProseMirrorPlugins() {
      const editor = this.editor
      return [
        Suggestion({
          pluginKey: new PluginKey("mention-suggestion"),
          editor,
          char: "@",
          allowSpaces: false,
          command: ({ editor, range, props }: { editor: Editor; range: Range; props: MentionUserItem }) => {
            // Insert @username text
            editor
              .chain()
              .focus()
              .deleteRange(range)
              .insertContent(`@${props.username} `)
              .run()
          },
          items: async ({ query }: { query: string }) => {
            try {
              const result = await searchUsers({ keyword: query })
              const users = result?.results || []
              return users
                .map((u: SearchUser) => {
                  const user = u.user || u
                  return {
                    username: (user as any).username || u.username || "",
                    nickname: (user as any).nickname || u.nickname || "",
                    avatar: (user as any).smallAvatar || (user as any).avatar || "",
                  }
                })
                .filter((item: MentionUserItem) => item.username)
                .slice(0, 8)
            } catch {
              return []
            }
          },
          render: () => {
            let renderer: ReactRenderer<MentionMenuHandle, MentionMenuProps> | null = null
            let currentProps: SuggestionProps<MentionUserItem, MentionUserItem> | null = null

            function syncPosition() {
              if (renderer && currentProps) {
                const rect = currentProps.clientRect?.()
                if (!rect || !renderer.element) return
                const gap = 8
                const offset = 6
                const popupWidth = renderer.element.offsetWidth || 280
                const popupHeight = renderer.element.offsetHeight || 300
                let left = rect.left
                let top = rect.bottom + offset
                if (left + popupWidth + gap > window.innerWidth) {
                  left = window.innerWidth - popupWidth - gap
                }
                if (left < gap) left = gap
                if (top + popupHeight + gap > window.innerHeight) {
                  top = rect.top - popupHeight - offset
                }
                if (top < gap) top = gap
                renderer.element.style.left = `${left}px`
                renderer.element.style.top = `${top}px`
              }
            }

            return {
              onStart(props) {
                currentProps = props
                renderer = new ReactRenderer(MentionMenu, {
                  editor: props.editor,
                  props,
                })
                renderer.element.classList.add("mention-popup")
                const portalTarget = document.fullscreenElement || document.body
                portalTarget.appendChild(renderer.element)
                syncPosition()
                window.addEventListener("resize", syncPosition)
                window.addEventListener("scroll", syncPosition, true)
              },
              onUpdate(props) {
                currentProps = props
                renderer?.updateProps(props)
                syncPosition()
              },
              onKeyDown({ event, view }: SuggestionKeyDownProps) {
                if (event.key === "Escape") {
                  exitSuggestion(view)
                  return true
                }
                return renderer?.ref?.onKeyDown(event) || false
              },
              onExit() {
                window.removeEventListener("resize", syncPosition)
                window.removeEventListener("scroll", syncPosition, true)
                renderer?.destroy()
                renderer?.element.remove()
                renderer = null
                currentProps = null
              },
            }
          },
        }),
      ]
    },
  })
}

function createSlashSuggestion(labels: RichTextEditorLabels, locale: string) {
  return Extension.create({
    name: "slash-commands",

    addProseMirrorPlugins() {
      const editor = this.editor
      return [
        Suggestion({
          editor,
          char: "/",
          command: ({ editor, range, props }: { editor: Editor; range: Range; props: SlashCommandItem }) => {
            props.command({ editor, range })
          },
          items: ({ query }: { query: string }) => {
            const searchQuery = query.toLowerCase().trim()
            return slashItems(labels, locale)
              .filter((item) => {
                if (!searchQuery) return true
                return item.title.toLowerCase().includes(searchQuery) || item.description.toLowerCase().includes(searchQuery) || item.aliases.some((alias) => alias.toLowerCase().includes(searchQuery))
              })
              .slice(0, 10)
          },
          render: () => {
            let renderer: ReactRenderer<SlashCommandMenuHandle, SlashCommandMenuProps> | null = null
            let currentProps: SuggestionProps<SlashCommandItem, SlashCommandItem> | null = null

            function syncPosition() {
              if (renderer && currentProps) {
                updateSlashMenuPosition(renderer.element, currentProps.clientRect)
              }
            }

            function mount(props: SuggestionProps<SlashCommandItem, SlashCommandItem>) {
              currentProps = props
              renderer = new ReactRenderer(SlashCommandMenu, {
                editor: props.editor,
                props: { ...props, labels },
              })
              renderer.element.classList.add("slash-commands-popup")
              const portalTarget = document.fullscreenElement || document.body
              portalTarget.appendChild(renderer.element)
              syncPosition()
              window.addEventListener("resize", syncPosition)
              window.addEventListener("scroll", syncPosition, true)
            }

            function cleanup() {
              window.removeEventListener("resize", syncPosition)
              window.removeEventListener("scroll", syncPosition, true)
              renderer?.destroy()
              renderer?.element.remove()
              renderer = null
              currentProps = null
            }

            function update(props: SuggestionProps<SlashCommandItem, SlashCommandItem>) {
              currentProps = props
              renderer?.updateProps({ ...props, labels })
              syncPosition()
            }

            return {
              onStart(props) {
                mount(props)
              },
              onUpdate(props) {
                update(props)
              },
              onKeyDown({ event, view }: SuggestionKeyDownProps) {
                if (event.key === "Escape") {
                  exitSuggestion(view)
                  cleanup()
                  return true
                }
                return renderer?.ref?.onKeyDown(event) || false
              },
              onExit() {
                cleanup()
              },
            }
          },
        }),
      ]
    },
  })
}

function ToolbarButton({
  title,
  active,
  disabled,
  children,
  onClick,
}: {
  title: string
  active?: boolean
  disabled?: boolean
  children: React.ReactNode
  onClick: () => void
}) {
  return (
    <button type="button" className={cn("m-editor-toolbar-button", active && "is-active")} title={title} disabled={disabled} onClick={onClick}>
      {children}
    </button>
  )
}

function ToolbarDivider() {
  return <span className="m-editor-toolbar-divider" />
}

function ColorButton({
  title,
  type,
  palette,
  activeColor,
  onApply,
  onClear,
  children,
}: {
  title: string
  type: "text" | "background"
  palette: string[]
  activeColor: string
  onApply: (color: string) => void
  onClear: () => void
  children: React.ReactNode
}) {
  const [open, setOpen] = React.useState(false)
  const buttonRef = React.useRef<HTMLDivElement>(null)
  const popupRef = React.useRef<HTMLDivElement>(null)
  const [popupPosition, setPopupPosition] = React.useState({ top: 0, left: 0 })
  const [mounted, setMounted] = React.useState(false)

  React.useEffect(() => {
    setMounted(true)
  }, [])

  const updatePopupPosition = React.useCallback(() => {
    const button = buttonRef.current
    if (!button) {
      return
    }
    const rect = button.getBoundingClientRect()
    const popupWidth = popupRef.current?.offsetWidth || 220
    const left = Math.min(Math.max(rect.left, 8), window.innerWidth - popupWidth - 8)
    setPopupPosition({ top: rect.bottom + 6, left })
  }, [])

  React.useEffect(() => {
    if (!open) {
      return
    }
    updatePopupPosition()
    window.addEventListener("resize", updatePopupPosition)
    window.addEventListener("scroll", updatePopupPosition, true)
    return () => {
      window.removeEventListener("resize", updatePopupPosition)
      window.removeEventListener("scroll", updatePopupPosition, true)
    }
  }, [open, updatePopupPosition])

  React.useEffect(() => {
    if (!open) {
      return
    }
    const closeOnOutsideEvent = (event: PointerEvent | FocusEvent) => {
      const target = event.target
      if (!(target instanceof globalThis.Node)) {
        return
      }
      if (buttonRef.current?.contains(target) || popupRef.current?.contains(target)) {
        return
      }
      setOpen(false)
    }
    const closeOnEscape = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setOpen(false)
      }
    }
    document.addEventListener("pointerdown", closeOnOutsideEvent, true)
    document.addEventListener("focusin", closeOnOutsideEvent, true)
    document.addEventListener("keydown", closeOnEscape, true)
    return () => {
      document.removeEventListener("pointerdown", closeOnOutsideEvent, true)
      document.removeEventListener("focusin", closeOnOutsideEvent, true)
      document.removeEventListener("keydown", closeOnEscape, true)
    }
  }, [open])

  const portalTarget = mounted ? document.fullscreenElement || document.body : null
  const popup =
    open && portalTarget
      ? createPortal(
          <div ref={popupRef} className="color-popup" style={{ top: popupPosition.top, left: popupPosition.left }}>
            <div className="color-picker-header">
              <span>{title}</span>
              <button
                type="button"
                className="clear-color"
                title="clear"
                onMouseDown={(event) => event.preventDefault()}
                onClick={() => {
                  setOpen(false)
                  onClear()
                }}
              >
                <div className="default-color">{!activeColor ? <Check size={16} /> : null}</div>
              </button>
            </div>
            <div className="color-grid">
              {palette.map((color) => (
                <button
                  key={color}
                  type="button"
                  className="color-option"
                  style={{ backgroundColor: color }}
                  onMouseDown={(event) => event.preventDefault()}
                  onClick={() => {
                    setOpen(false)
                    onApply(color)
                  }}
                >
                  {activeColor === color ? <Check size={16} className="check-icon" /> : null}
                </button>
              ))}
            </div>
          </div>,
          portalTarget
        )
      : null

  return (
    <div ref={buttonRef} className={`${type}-color-button m-editor-color-button`}>
      <ToolbarButton title={title} active={Boolean(activeColor)} onClick={() => setOpen((value) => !value)}>
        <span className="button-content">{children}</span>
      </ToolbarButton>
      <div className="color-indicator" style={{ backgroundColor: activeColor || "transparent" }} />
      {popup}
    </div>
  )
}

export function RichTextEditor({
  value,
  height = "400px",
  onChange,
}: {
  value: string
  placeholder?: string
  height?: string
  onChange: (value: string) => void
}) {
  const { locale, t } = useI18n()
  const labels = React.useMemo(() => createEditorLabels(t), [t])
  const { catchError } = useToastActions()
  const containerRef = React.useRef<HTMLDivElement>(null)
  const fileInputRef = React.useRef<HTMLInputElement>(null)
  const lastExternalValueRef = React.useRef(value)
  const [uploading, setUploading] = React.useState(false)
  const [isFullscreen, setIsFullscreen] = React.useState(false)
  const [linkOpen, setLinkOpen] = React.useState(false)
  const [linkText, setLinkText] = React.useState("")
  const [linkUrl, setLinkUrl] = React.useState("")

  const editor = useEditor({
    immediatelyRender: false,
    extensions: [
      StarterKit.configure({
        link: false,
        horizontalRule: false,
      }),
      Link.configure({
        openOnClick: false,
        autolink: true,
        linkOnPaste: true,
        HTMLAttributes: {
          target: "_blank",
          rel: "noopener noreferrer",
        },
      }),
      ResizableImage,
      Underline,
      TextAlign.configure({
        types: ["heading", "paragraph"],
      }),
      TextStyle,
      Color,
      BackgroundColor,
      TaskList,
      TaskItem.configure({
        nested: true,
      }),
      Typography,
      HorizontalRule,
      createSlashSuggestion(labels, locale),
      createMentionSuggestion(),
      Placeholder.configure({
        placeholder: labels.placeholder,
      }),
    ],
    content: value || "",
    editorProps: {
      attributes: {
        class: "tiptap",
      },
      handlePaste(view, event) {
        const items = event.clipboardData?.items
        if (!items?.length) {
          return false
        }
        const files = Array.from(items)
          .filter((item) => item.type.includes("image"))
          .map((item) => item.getAsFile())
          .filter(Boolean) as File[]
        if (!files.length) {
          return false
        }
        event.preventDefault()
        void uploadImages(files)
        return true
      },
      handleDrop(view, event) {
        const files = Array.from(event.dataTransfer?.files || []).filter((file) => file.type.includes("image"))
        if (!files.length) {
          return false
        }
        event.preventDefault()
        void uploadImages(files)
        return true
      },
    },
    onUpdate({ editor: currentEditor }) {
      const html = currentEditor.getHTML()
      lastExternalValueRef.current = html
      onChange(html === "<p></p>" ? "" : html)
    },
  })

  async function uploadImages(files: File[]) {
    setUploading(true)
    try {
      const urls = await Promise.all(files.map((file) => uploadEditorImage(file)))
      urls.forEach((url, index) => {
        editor?.chain().focus().setResizableImage({ src: url, alt: files[index]?.name || "", title: files[index]?.name || "" }).run()
      })
    } catch (error) {
      catchError(error)
    } finally {
      setUploading(false)
    }
  }

  React.useEffect(() => {
    if (!editor || value === lastExternalValueRef.current) {
      return
    }
    lastExternalValueRef.current = value
    editor.commands.setContent(value || "", { emitUpdate: false })
  }, [editor, value])

  React.useEffect(() => {
    const onFullscreenChange = () => setIsFullscreen(Boolean(document.fullscreenElement))
    document.addEventListener("fullscreenchange", onFullscreenChange)
    return () => document.removeEventListener("fullscreenchange", onFullscreenChange)
  }, [])

  function toggleFullscreen() {
    const el = containerRef.current
    if (!el) return
    if (!document.fullscreenElement) {
      void el.requestFullscreen()
    } else {
      void document.exitFullscreen()
    }
  }

  function openLinkDialog() {
    if (!editor) return
    const previousUrl = editor.getAttributes("link").href as string | undefined
    const selectedText = editor.state.doc.textBetween(editor.state.selection.from, editor.state.selection.to, " ")
    setLinkText(selectedText)
    setLinkUrl(previousUrl || "")
    setLinkOpen(true)
  }

  function applyLink() {
    if (!editor) return
    if (!linkUrl) {
      editor.chain().focus().extendMarkRange("link").unsetLink().run()
      setLinkOpen(false)
      return
    }
    if (linkText && editor.state.selection.empty) {
      editor.chain().focus().insertContent(`<a href="${linkUrl}" target="_blank" rel="noopener noreferrer">${linkText}</a>`).run()
    } else {
      editor.chain().focus().extendMarkRange("link").setLink({ href: linkUrl }).run()
    }
    setLinkOpen(false)
  }

  const activeTextColor = (editor?.getAttributes("textStyle").color as string | undefined) || ""
  const activeBackgroundColor = (editor?.getAttributes("textStyle").backgroundColor as string | undefined) || ""
  const toolbar = labels.toolbar

  return (
    <div ref={containerRef} className="m-editor-container" style={{ height }}>
      <div className="editor-toolbar">
        <div className="editor-toolbar-btns editor-toolbar-left">
          <ToolbarButton title={toolbar.bold} active={editor?.isActive("bold")} onClick={() => editor?.chain().focus().toggleBold().run()}>
            <Bold size={16} />
          </ToolbarButton>
          <ToolbarButton title={toolbar.underline} active={editor?.isActive("underline")} onClick={() => editor?.chain().focus().toggleUnderline().run()}>
            <UnderlineIcon size={16} />
          </ToolbarButton>
          <ToolbarButton title={toolbar.italic} active={editor?.isActive("italic")} onClick={() => editor?.chain().focus().toggleItalic().run()}>
            <Italic size={16} />
          </ToolbarButton>
          <ToolbarButton title={toolbar.strike} active={editor?.isActive("strike")} onClick={() => editor?.chain().focus().toggleStrike().run()}>
            <Strikethrough size={16} />
          </ToolbarButton>
          <ToolbarDivider />
          <ToolbarButton title={toolbar.heading1} active={editor?.isActive("heading", { level: 1 })} onClick={() => editor?.chain().focus().toggleHeading({ level: 1 }).run()}>
            <Heading1 size={16} />
          </ToolbarButton>
          <ToolbarButton title={toolbar.heading2} active={editor?.isActive("heading", { level: 2 })} onClick={() => editor?.chain().focus().toggleHeading({ level: 2 }).run()}>
            <Heading2 size={16} />
          </ToolbarButton>
          <ToolbarButton title={toolbar.quote} active={editor?.isActive("blockquote")} onClick={() => editor?.chain().focus().toggleBlockquote().run()}>
            <Quote size={16} />
          </ToolbarButton>
          <ToolbarButton title={toolbar.bulletList} active={editor?.isActive("bulletList")} onClick={() => editor?.chain().focus().toggleBulletList().run()}>
            <List size={16} />
          </ToolbarButton>
          <ToolbarButton title={toolbar.orderedList} active={editor?.isActive("orderedList")} onClick={() => editor?.chain().focus().toggleOrderedList().run()}>
            <ListOrdered size={16} />
          </ToolbarButton>
          <ToolbarButton title={toolbar.taskList} active={editor?.isActive("taskList")} onClick={() => editor?.chain().focus().toggleTaskList().run()}>
            <ListTodo size={16} />
          </ToolbarButton>
          <ToolbarDivider />
          <ToolbarButton title={toolbar.alignLeft} active={editor?.isActive({ textAlign: "left" })} onClick={() => editor?.chain().focus().setTextAlign("left").run()}>
            <AlignLeft size={16} />
          </ToolbarButton>
          <ToolbarButton title={toolbar.alignCenter} active={editor?.isActive({ textAlign: "center" })} onClick={() => editor?.chain().focus().setTextAlign("center").run()}>
            <AlignCenter size={16} />
          </ToolbarButton>
          <ToolbarButton title={toolbar.alignRight} active={editor?.isActive({ textAlign: "right" })} onClick={() => editor?.chain().focus().setTextAlign("right").run()}>
            <AlignRight size={16} />
          </ToolbarButton>
          <ToolbarDivider />
          <ColorButton title={toolbar.textColor} type="text" palette={TEXT_COLOR_PALETTE} activeColor={activeTextColor} onApply={(color) => editor?.chain().focus().setColor(color).run()} onClear={() => editor?.chain().focus().unsetColor().run()}>
            <Palette size={16} />
          </ColorButton>
          <ColorButton title={toolbar.backgroundColor} type="background" palette={BACKGROUND_COLOR_PALETTE} activeColor={activeBackgroundColor} onApply={(color) => editor?.chain().focus().setBackgroundColor(color).run()} onClear={() => editor?.chain().focus().unsetBackgroundColor().run()}>
            <Paintbrush size={16} />
          </ColorButton>
          <ToolbarDivider />
          <ToolbarButton title={toolbar.inlineCode} active={editor?.isActive("code")} onClick={() => editor?.chain().focus().toggleCode().run()}>
            <Code size={16} />
          </ToolbarButton>
          <ToolbarButton title={toolbar.codeBlock} active={editor?.isActive("codeBlock")} onClick={() => editor?.chain().focus().toggleCodeBlock().run()}>
            <Code2 size={16} />
          </ToolbarButton>
          <ToolbarDivider />
          <div className="link-button">
            <ToolbarButton title={toolbar.link} active={editor?.isActive("link")} onClick={openLinkDialog}>
              <LinkIcon size={16} />
            </ToolbarButton>
            {linkOpen ? (
              <div className="link-dialog-content">
                <div className="link-input-group">
                  <label>{labels.linkDialog.textLabel}</label>
                  <input value={linkText} placeholder={labels.linkDialog.textPlaceholder} onChange={(event) => setLinkText(event.currentTarget.value)} />
                </div>
                <div className="link-input-group">
                  <label>{labels.linkDialog.urlLabel}</label>
                  <input value={linkUrl} placeholder={labels.linkDialog.urlPlaceholder} onChange={(event) => setLinkUrl(event.currentTarget.value)} />
                </div>
                <div className="link-dialog-actions">
                  <button type="button" className="btn-primary" onClick={applyLink}>{labels.linkDialog.confirm}</button>
                  <button type="button" className="btn-danger" onClick={() => { editor?.chain().focus().extendMarkRange("link").unsetLink().run(); setLinkOpen(false) }}>{labels.linkDialog.remove}</button>
                  <button type="button" className="btn-secondary" onClick={() => setLinkOpen(false)}>{labels.linkDialog.cancel}</button>
                </div>
              </div>
            ) : null}
          </div>
          <div className="image-upload-button">
            <ToolbarButton title={toolbar.image} disabled={!editor || uploading} onClick={() => fileInputRef.current?.click()}>
              <ImageIcon size={16} />
            </ToolbarButton>
            {uploading ? (
              <div className="upload-progress">
                <div className="upload-spinner" />
                <span>{toolbar.uploading}</span>
              </div>
            ) : null}
          </div>
          <ToolbarButton title={toolbar.horizontalRule} onClick={() => editor?.chain().focus().setHorizontalRule().run()}>
            <MinusSquare size={16} />
          </ToolbarButton>
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            className="hidden"
            onChange={(event) => {
              const files = Array.from(event.currentTarget.files || [])
              if (files.length) void uploadImages(files)
              event.currentTarget.value = ""
            }}
          />
        </div>
        <div className="editor-toolbar-btns editor-toolbar-right">
          <ToolbarButton title={isFullscreen ? toolbar.exitFullscreen : toolbar.fullscreen} active={isFullscreen} onClick={toggleFullscreen}>
            {isFullscreen ? <Minimize size={16} /> : <Maximize size={16} />}
          </ToolbarButton>
        </div>
      </div>
      <EditorContent editor={editor} className="editor-content" />
    </div>
  )
}
