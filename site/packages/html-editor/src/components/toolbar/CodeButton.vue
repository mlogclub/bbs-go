<template>
  <ToolbarButton
    @click="editor?.chain().focus().toggleCode().run()"
    :isActive="editor?.isActive('code')"
    title="行内代码"
  >
    <LucideCode :size="TOOLBAR_ICON_SIZE" />
  </ToolbarButton>
</template>

<script setup lang="ts">
import { Editor } from '@tiptap/core'
import { LucideCode } from 'lucide-vue-next'
import ToolbarButton from './ToolbarButton.vue'
import { TOOLBAR_ICON_SIZE } from '../../constants/editor'

defineProps<{
  editor: Editor | null | undefined
}>()
</script> 