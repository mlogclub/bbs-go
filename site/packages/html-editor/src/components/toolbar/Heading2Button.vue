<template>
  <ToolbarButton
    @click="editor?.chain().focus().toggleHeading({ level: 2 }).run()"
    :isActive="editor?.isActive('heading', { level: 2 })"
    title="标题2"
  >
    <LucideHeading2 :size="TOOLBAR_ICON_SIZE" />
  </ToolbarButton>
</template>

<script setup lang="ts">
import { Editor } from '@tiptap/core'
import { LucideHeading2 } from 'lucide-vue-next'
import ToolbarButton from './ToolbarButton.vue'
import { TOOLBAR_ICON_SIZE } from '../../constants/editor'

defineProps<{
  editor: Editor | null | undefined
}>()
</script> 