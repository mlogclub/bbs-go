import { Extension } from '@tiptap/core'
import Suggestion from '@tiptap/suggestion'
import type { Editor } from '@tiptap/core'
import { VueRenderer } from '@tiptap/vue-3'
import SlashCommandsList from './SlashCommands.vue'
import type { CommandItem } from './types'
import { getSuggestionItems } from './types'

export const SlashSuggestion = Extension.create({
  name: 'slash-commands',

  addOptions() {
    return {
      suggestion: {
        char: '/',
        command: ({ editor, range, props }: { editor: Editor; range: { from: number; to: number }; props: any }) => {
          props.command({ editor, range })
        },
        items: ({ query }: { query: string }) => {
          // 过滤出匹配的菜单项
          const matchedItems = getSuggestionItems()
            .filter((item: CommandItem) => {
              // 如果查询为空，显示所有项
              if (!query) return true
              
              // 匹配标题、描述或别名（如果有）
              const searchQuery = query.toLowerCase().trim()
              return (
                item.title.toLowerCase().includes(searchQuery) || 
                item.description.toLowerCase().includes(searchQuery) ||
                (item.aliases && item.aliases.some(alias => alias.toLowerCase().includes(searchQuery)))
              )
            })
            .slice(0, 10)
          
          return matchedItems
        },
      },
    }
  },

  addProseMirrorPlugins() {
    const editor = this.editor

    return [
      Suggestion({
        editor,
        ...this.options.suggestion,
        render: () => {
          let component: VueRenderer
          let popup: HTMLElement | null = null

          // 辅助函数：获取编辑器容器元素
          const getEditorContainer = () => {
            const editorDOM = editor.view.dom
            if (!editorDOM) return null

            // 获取编辑器的容器元素，不依赖于类名，只依赖DOM结构关系
            let editorContent = editorDOM.parentElement
            // 确保找到的容器元素具有相对或绝对定位
            while (editorContent && window.getComputedStyle(editorContent).position === 'static') {
              editorContent = editorContent.parentElement
            }
            
            // 如果没有找到合适的容器，使用编辑器的直接父元素
            return editorContent || editorDOM.parentElement
          }

          // 辅助函数：更新弹出菜单位置
          const updatePopupPosition = (props: { clientRect: (() => DOMRect | null) | null }) => {
            const editorContent = getEditorContainer()
            if (!editorContent || !popup) return

            const editorRect = editorContent.getBoundingClientRect()
            const clientRectObj = props.clientRect?.()
            
            if (!clientRectObj) return

            // 直接使用光标位置来定位弹出窗口
            const left = clientRectObj.left - editorRect.left
            // 使用光标底部位置 + 少量偏移，确保不会遮挡光标
            const top = clientRectObj.bottom - editorRect.top + 5 // 5px的额外间距

            popup.style.left = `${left}px`
            popup.style.top = `${top}px`
          }

          return {
            onStart: (props: { clientRect: (() => DOMRect | null) | null; items: CommandItem[] }) => {
              // 如果没有匹配项，不显示菜单
              if (props.items.length === 0) {
                return
              }
            
              component = new VueRenderer(SlashCommandsList, {
                props,
                editor,
              })

              if (popup && popup.parentElement) {
                popup.parentElement.removeChild(popup)
              }

              const editorContent = getEditorContainer()
              if (!editorContent) return

              popup = document.createElement('div')
              popup.classList.add('slash-commands-popup')
              popup.style.position = 'absolute'
              popup.style.zIndex = '9999'
              
              updatePopupPosition(props)
              
              popup.appendChild(component.element)
              editorContent.appendChild(popup)
            },

            onUpdate(props: { clientRect: (() => DOMRect | null) | null; items: CommandItem[] }) {
              component.updateProps(props)
              
              // 如果没有匹配项，关闭弹出菜单
              if (props.items.length === 0) {
                if (popup && popup.parentElement) {
                  popup.parentElement.removeChild(popup)
                }
                return
              }
              
              updatePopupPosition(props)
            },

            onKeyDown(props: { event: KeyboardEvent }) {
              if (!props.event) {
                return false
              }

              if (props.event.key === 'Escape') {
                if (popup && popup.parentElement) {
                  popup.parentElement.removeChild(popup)
                }
                return true
              }

              if (['ArrowUp', 'ArrowDown', 'Enter'].includes(props.event.key)) {
                return component.ref?.onKeyDown(props.event)
              }

              return false
            },

            onExit() {
              component.destroy()
              if (popup && popup.parentElement) {
                popup.parentElement.removeChild(popup)
              }
            },
          }
        },
      }),
    ]
  },
}) 