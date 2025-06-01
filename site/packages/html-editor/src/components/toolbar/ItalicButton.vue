<template>
  <ToolbarButton
    @click="editor?.chain().focus().toggleItalic().run()"
    :isActive="editor?.isActive('italic')"
    title="斜体"
  >
    <LucideItalic :size="TOOLBAR_ICON_SIZE" />
  </ToolbarButton>
</template>

<script setup lang="ts">
import { Editor } from '@tiptap/core'
import { LucideItalic } from 'lucide-vue-next'
import ToolbarButton from './ToolbarButton.vue'
import { TOOLBAR_ICON_SIZE } from '../../constants/editor'

defineProps<{
  editor: Editor | null | undefined
}>()
</script> 