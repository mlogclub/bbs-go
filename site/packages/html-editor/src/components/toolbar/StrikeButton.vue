<template>
  <ToolbarButton
    @click="editor?.chain().focus().toggleStrike().run()"
    :isActive="editor?.isActive('strike')"
    title="删除线"
  >
    <LucideStrikethrough :size="TOOLBAR_ICON_SIZE" />
  </ToolbarButton>
</template>

<script setup lang="ts">
import { Editor } from '@tiptap/core'
import { LucideStrikethrough } from 'lucide-vue-next'
import ToolbarButton from './ToolbarButton.vue'
import { TOOLBAR_ICON_SIZE } from '../../constants/editor'

defineProps<{
  editor: Editor | null | undefined
}>()
</script> 