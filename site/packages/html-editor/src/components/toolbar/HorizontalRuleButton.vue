<template>
  <ToolbarButton
    @click="editor?.chain().focus().setHorizontalRule().run()"
    title="水平分割线"
  >
    <LucideMinusSquare :size="TOOLBAR_ICON_SIZE" />
  </ToolbarButton>
</template>

<script setup lang="ts">
import { Editor } from '@tiptap/core'
import { LucideMinusSquare } from 'lucide-vue-next'
import ToolbarButton from './ToolbarButton.vue'
import { TOOLBAR_ICON_SIZE } from '../../constants/editor'

defineProps<{
  editor: Editor | null | undefined
}>()
</script> 