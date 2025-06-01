import { Extension } from '@tiptap/core'
import { Plugin, PluginKey } from '@tiptap/pm/state'
import { UploadImageFunction } from '../utils/imageUtils'

export interface PasteImageOptions {
  /**
   * 是否启用粘贴图片功能
   */
  enabled: boolean
  /**
   * 自定义图片上传函数
   */
  uploadImage: UploadImageFunction
}

/**
 * 粘贴图片扩展
 * 支持粘贴剪贴板中的截图和复制的磁盘图片文件
 */
export const PasteImage = Extension.create<PasteImageOptions>({
  name: 'pasteImage',

  addOptions() {
    return {
      enabled: true,
      uploadImage: undefined,
    }
  },

  onCreate() {
  },

  addProseMirrorPlugins() {
    return [
      new Plugin({
        key: new PluginKey('pasteImage'),
        props: {
          handlePaste: (view, event, slice) => {
            console.log('📋 接收到粘贴事件', { enabled: this.options.enabled, event })

            if (!this.options.enabled) {
              return false
            }

            // 优先使用 files，如果没有则使用 items
            let imageFiles: File[] = []

            // 直接从 files 获取图片文件
            const files = Array.from(event.clipboardData?.files || [])
            imageFiles = files.filter(file => file.type.startsWith('image/'))

            // 如果 files 中没有图片，再尝试从 items 中获取
            if (imageFiles.length === 0) {
              const items = Array.from(event.clipboardData?.items || [])
              imageFiles = items
                .filter(item => item.kind === 'file' && item.type.startsWith('image/'))
                .map(item => item.getAsFile())
                .filter((file): file is File => file !== null)
            }

            if (imageFiles.length === 0) {
              console.log('📋 粘贴事件中没有找到图片文件')
              return false
            }

            console.log('📋 找到图片文件:', imageFiles.map(f => ({ name: f.name, type: f.type, size: f.size })))

            // 阻止默认粘贴行为
            event.preventDefault()

            // 处理图片文件
            imageFiles.forEach((file, index) => {
              setTimeout(async () => {
                try {
                  console.log('开始处理粘贴的图片:', file.name, file.type, file.size)

                  // 使用配置中的uploadImage函数或默认的uploadImage函数
                  const uploadImageFn = this.options.uploadImage
                  const resp = await uploadImageFn(file)

                  // 获取当前光标位置
                  const currentState = view.state
                  const pos = currentState.selection.from + index // 为每个图片偏移位置

                  // 创建图片节点
                  const imageNode = currentState.schema.nodes.resizableImage.create({
                    src: resp.url,
                    alt: resp.name || '',
                    title: resp.name || '',
                  })

                  // 插入图片
                  const tr = currentState.tr.replaceWith(pos, pos, imageNode)
                  view.dispatch(tr)

                  console.log('图片粘贴成功')
                } catch (error) {
                  console.error('粘贴图片失败:', error)
                  alert('图片粘贴失败: ' + (error instanceof Error ? error.message : '未知错误'))
                }
              }, index * 10) // 轻微延迟以确保顺序
            })

            return true
          },

          handleDrop: (view, event, slice, moved) => {
            console.log('🖱️ 接收到拖拽事件', { enabled: this.options.enabled, event })

            if (!this.options.enabled) {
              return false
            }

            const files = Array.from(event.dataTransfer?.files || [])
            const imageFiles = files.filter(file => file.type.startsWith('image/'))

            if (imageFiles.length === 0) {
              console.log('🖱️ 拖拽事件中没有找到图片文件')
              return false
            }

            console.log('🖱️ 找到拖拽图片文件:', imageFiles.map(f => ({ name: f.name, type: f.type, size: f.size })))

            // 阻止默认拖拽行为
            event.preventDefault()

            // 获取拖拽位置
            const coordinates = view.posAtCoords({
              left: event.clientX,
              top: event.clientY,
            })

            if (!coordinates) {
              return false
            }

            // 处理拖拽的图片文件
            imageFiles.forEach((file, index) => {
              setTimeout(async () => {
                try {
                  console.log('开始处理拖拽的图片:', file.name, file.type, file.size)

                  // 使用配置中的uploadImage函数或默认的uploadImage函数
                  const uploadImageFn = this.options.uploadImage
                  const resp = await uploadImageFn(file)

                  // 为每个图片计算插入位置（避免重叠）
                  const insertPos = coordinates.pos + index

                  // 创建图片节点
                  const imageNode = view.state.schema.nodes.resizableImage.create({
                    src: resp.url,
                    alt: resp.name || '',
                    title: resp.name || '',
                  })

                  // 在拖拽位置插入图片
                  const tr = view.state.tr.replaceWith(insertPos, insertPos, imageNode)
                  view.dispatch(tr)

                  console.log('图片拖拽成功')
                } catch (error) {
                  console.error('拖拽图片失败:', error)
                  alert('图片拖拽失败: ' + (error instanceof Error ? error.message : '未知错误'))
                }
              }, index * 10) // 轻微延迟以确保顺序
            })

            return true
          },
        },
      }),
    ]
  },
})
