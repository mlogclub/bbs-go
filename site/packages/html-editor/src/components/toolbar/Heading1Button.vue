<template>
  <ToolbarButton
    @click="editor?.chain().focus().toggleHeading({ level: 1 }).run()"
    :isActive="editor?.isActive('heading', { level: 1 })"
    title="标题1"
  >
    <LucideHeading1 :size="TOOLBAR_ICON_SIZE" />
  </ToolbarButton>
</template>

<script setup lang="ts">
import { Editor } from '@tiptap/core'
import { LucideHeading1 } from 'lucide-vue-next'
import ToolbarButton from './ToolbarButton.vue'
import { TOOLBAR_ICON_SIZE } from '../../constants/editor'

defineProps<{
  editor: Editor | null | undefined
}>()
</script> 