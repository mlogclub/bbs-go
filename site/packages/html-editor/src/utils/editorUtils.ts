import { Editor } from '@tiptap/core'
import { computed } from 'vue'

// 检查编辑器是否准备就绪
export const useEditorReady = (editor: Editor | null | undefined) => {
  return computed(() => !!editor)
}

// 通用的链接设置函数
export const setLink = (editor: Editor | null | undefined) => {
  if (!editor) return
  const previousUrl = editor.getAttributes('link').href
  const url = window.prompt('URL', previousUrl)

  if (url === null) {
    return
  }

  if (url === '') {
    editor.chain().focus().extendMarkRange('link').unsetLink().run()
    return
  }

  editor.chain().focus().extendMarkRange('link').setLink({ href: url }).run()
}
