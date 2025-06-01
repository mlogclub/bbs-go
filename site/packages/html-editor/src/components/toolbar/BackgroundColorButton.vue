<template>
  <div class="background-color-button">
    <ToolbarButton
      @click="toggleColorPicker"
      :isActive="hasBackgroundColor"
      title="背景颜色"
    >
      <div class="button-content">
        <LucidePaintbrush :size="TOOLBAR_ICON_SIZE" />
      </div>
    </ToolbarButton>
    <div class="color-indicator" :style="{ backgroundColor: backgroundColor || 'transparent' }"></div>
    
    <div v-if="isColorPickerOpen" class="color-popup" ref="colorPopup">
      <div class="color-picker-header">
        <span>标准色卡</span>
        <button class="clear-color" @click="clearColor" title="清除颜色">
          <div class="default-color">
            <LucideCheck :size="TOOLBAR_ICON_SIZE" v-if="!backgroundColor" />
          </div>
        </button>
      </div>
      <div class="color-grid">
        <button
          v-for="color in standardColors"
          :key="color"
          class="color-option"
          :style="{ backgroundColor: color }"
          @click="applyColor(color)"
        ></button>
      </div>
      
      <div class="color-section-title">最近使用</div>
      <div class="color-grid recent">
        <button
          v-for="color in recentColors"
          :key="color"
          class="color-option"
          :style="{ backgroundColor: color }"
          @click="applyColor(color)"
        >
          <LucideCheck :size="TOOLBAR_ICON_SIZE" v-if="backgroundColor === color" class="check-icon" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Editor } from '@tiptap/core'
import { LucidePaintbrush, LucideCheck } from 'lucide-vue-next'
import ToolbarButton from './ToolbarButton.vue'
import { TOOLBAR_ICON_SIZE } from '../../constants/editor'
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'

const props = defineProps<{
  editor: Editor | null | undefined
}>()

const backgroundColor = ref('')
const isColorPickerOpen = ref(false)
const recentColors = ref<string[]>([])
const colorPopup = ref<HTMLElement | null>(null)

const standardColors = [
  '#FFCCCC', // 浅红色
  '#FFE6CC', // 浅橙色
  '#FFFFCC', // 浅黄色
  '#CCFFCC', // 浅绿色
  '#CCFFFF', // 浅青色
  '#CCE5FF', // 浅蓝色
  '#E5CCFF', // 浅紫色
  '#FFCCFF', // 浅粉色
  '#F2F2F2', // 更浅的灰色
  '#E6E6E6', // 非常浅的灰色
  '#FFD6CC', // 桃色
  '#E5FFCC', // 嫩绿色
  '#CCFFE5', // 薄荷色
  '#D6FFFF', // 天蓝色
  '#FFE0F2', // 浅粉紫色
  '#FFF0F0', // 极浅红色
  '#FFF9E6', // 极浅黄色
  '#F0FFF0', // 极浅绿色
  '#F0FFFF', // 极浅青色
  '#F5FFFA', // 薄荷奶油色
]

const hasBackgroundColor = computed(() => {
  if (!props.editor) return false
  const attrs = props.editor.getAttributes('textStyle')
  if (attrs.style && attrs.style.includes('background-color')) {
    // 从 style 字符串中提取颜色值
    const match = attrs.style.match(/background-color: (.+)/)
    const color = match ? match[1] : ''
    backgroundColor.value = color
    return true
  }
  return false
})

const closeOnClickOutside = (event: MouseEvent) => {
  const target = event.target as HTMLElement
  const colorButton = document.querySelector('.background-color-button')
  const popup = colorPopup.value
  
  if (colorButton && !colorButton.contains(target) && popup && !popup.contains(target)) {
    isColorPickerOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', closeOnClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', closeOnClickOutside)
})

const toggleColorPicker = async (event: MouseEvent) => {
  event.stopPropagation()
  isColorPickerOpen.value = !isColorPickerOpen.value
  
  if (isColorPickerOpen.value) {
    await nextTick()
    const button = event.currentTarget as HTMLElement
    const popup = colorPopup.value
    if (button && popup) {
      const buttonRect = button.getBoundingClientRect()
      const viewportHeight = window.innerHeight
      const popupHeight = popup.offsetHeight
      
      // 计算弹窗位置
      let top = buttonRect.bottom + 8
      let left = buttonRect.left
      
      // 如果弹窗会超出视口底部，则显示在按钮上方
      if (top + popupHeight > viewportHeight) {
        top = buttonRect.top - popupHeight - 8
      }
      
      // 如果弹窗会超出视口右侧，则向左调整
      if (left + 250 > window.innerWidth) {
        left = window.innerWidth - 250 - 8
      }
      
      popup.style.top = `${top}px`
      popup.style.left = `${left}px`
    }
  }
}

const applyColor = (color: string) => {
  if (!props.editor) return
  backgroundColor.value = color
  props.editor.chain().focus().setBackgroundColor(color).run()
  
  // 添加到最近使用的颜色
  if (!recentColors.value.includes(color)) {
    recentColors.value.unshift(color)
    if (recentColors.value.length > 7) {
      recentColors.value.pop()
    }
  }
}

const clearColor = () => {
  if (!props.editor) return
  backgroundColor.value = ''
  // 移除背景色
  props.editor.chain().focus().setBackgroundColor('').run()
}
</script>

<style scoped>
.background-color-button {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.button-content {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.color-indicator {
  margin-top: -2px;
  width: 24px;
  height: 3px;
  border-radius: 1px;
}

.color-popup {
  position: fixed;
  background: white;
  border: 1px solid #e1e1e1;
  border-radius: 8px;
  box-shadow: 0 3px 10px rgba(0, 0, 0, 0.1);
  width: 280px;
  z-index: 9999;
  padding: 12px;
}

.color-picker-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-size: 14px;
  color: #333;
}

.clear-color {
  background: none;
  border: none;
  padding: 0;
  cursor: pointer;
}

.default-color {
  width: 20px;
  height: 20px;
  border-radius: 3px;
  border: 1px solid #ddd;
  display: flex;
  align-items: center;
  justify-content: center;
  background: white;
}

.check-icon {
  color: white;
}

.color-grid {
  display: grid;
  grid-template-columns: repeat(10, 1fr);
  gap: 6px;
  margin-bottom: 15px;
}

.color-grid.recent {
  grid-template-columns: repeat(7, 1fr);
}

.color-option {
  width: 20px;
  height: 20px;
  border-radius: 3px;
  border: 1px solid rgba(0, 0, 0, 0.1);
  cursor: pointer;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.color-option:hover {
  transform: scale(1.15);
  transition: transform 0.2s;
}

.color-section-title {
  font-size: 14px;
  color: #333;
  margin-bottom: 8px;
}
</style> 