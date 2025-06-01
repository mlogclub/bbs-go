<template>
  <ToolbarButton
    @click="editor?.chain().focus().setTextAlign('left').run()"
    :isActive="editor?.isActive({ textAlign: 'left' })"
    title="左对齐"
  >
    <LucideAlignLeft :size="TOOLBAR_ICON_SIZE" />
  </ToolbarButton>
</template>

<script setup lang="ts">
import { Editor } from '@tiptap/core'
import { LucideAlignLeft } from 'lucide-vue-next'
import ToolbarButton from './ToolbarButton.vue'
import { TOOLBAR_ICON_SIZE } from '../../constants/editor'

defineProps<{
  editor: Editor | null | undefined
}>()
</script> 