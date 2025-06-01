/**
 * 图片上传相关工具函数
 */

export type UploadImageResponse = {
  url: string;
  name?: string;
  width?: number;
  height?: number;
};

export type UploadImageFunction = (file: File) => Promise<UploadImageResponse>;

// 支持的图片格式
export const SUPPORTED_IMAGE_TYPES = [
  'image/jpeg',
  'image/jpg',
  'image/png',
  'image/gif',
  'image/webp',
  'image/svg+xml'
]

// 最大文件大小 (5MB)
export const MAX_FILE_SIZE = 5 * 1024 * 1024

/**
 * 检查文件大小是否符合要求
 */
export function isValidFileSize(file: File): boolean {
  return file.size <= MAX_FILE_SIZE
}

/**
 * 将文件转换为Base64格式
 */
export function fileToBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => {
      if (typeof reader.result === 'string') {
        resolve(reader.result)
      } else {
        reject(new Error('Failed to convert file to base64'))
      }
    }
    reader.onerror = () => reject(new Error('Failed to read file'))
    reader.readAsDataURL(file)
  })
}

/**
 * 创建文件选择器
 */
export function createFileInput(): HTMLInputElement {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = SUPPORTED_IMAGE_TYPES.join(',')
  input.style.display = 'none'
  return input
}

/**
 * 格式化文件大小显示
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes'

  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * 图片上传处理函数
 * 这里使用 base64 作为示例，实际项目中应该上传到服务器
 */
export async function uploadImage(file: File): Promise<UploadImageResponse> {
  if (!isValidFileSize(file)) {
    throw new Error(`文件大小超过限制。最大支持：${formatFileSize(MAX_FILE_SIZE)}`)
  }
  try {
    // 这里使用 base64 作为演示
    const base64 = await fileToBase64(file)
    return {
      url: base64,
      name: file.name,
    }
  } catch (error) {
    throw new Error('图片上传失败：' + (error as Error).message)
  }
}
