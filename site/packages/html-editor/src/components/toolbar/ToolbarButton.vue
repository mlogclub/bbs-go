<template>
  <button 
    ref="buttonRef"
    @click="onClick" 
    :class="{ 'is-active': isActive }"
  >
    <slot></slot>
  </button>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import tippy from 'tippy.js'
import 'tippy.js/dist/tippy.css'
import type { Instance as TippyInstance } from 'tippy.js'

const props = defineProps<{
  isActive?: boolean
  title?: string
}>()

const emit = defineEmits<{
  (e: 'click', event: MouseEvent): void
}>()

const buttonRef = ref<HTMLElement | null>(null)
let tippyInstance: TippyInstance | null = null

const onClick = (event: MouseEvent) => {
  emit('click', event)
}

onMounted(() => {
  if (buttonRef.value && props.title) {
    tippyInstance = tippy(buttonRef.value, {
      content: props.title,
      placement: 'bottom',
      arrow: true,
      duration: [200, 100]
    })
  }
})

watch(() => props.title, (newTitle) => {
  if (tippyInstance && newTitle) {
    tippyInstance.setContent(newTitle)
  }
})

onBeforeUnmount(() => {
  if (tippyInstance) {
    tippyInstance.destroy()
  }
})
</script>

<style scoped>
button {
  padding: 0.25rem;
  border: none;
  border-radius: 4px;
  background: var(--editor-bg);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  transition: all 0.2s ease;
  color: var(--editor-text);
}

button:hover {
  background: var(--editor-hover);
  border-color: var(--editor-border);
  color: var(--editor-text);
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

button.is-active {
  background: var(--editor-hover);
  border-color: var(--editor-border);
  color: var(--editor-text);
  font-weight: 500;
}

button:active {
  transform: translateY(0);
  box-shadow: none;
}
</style>
