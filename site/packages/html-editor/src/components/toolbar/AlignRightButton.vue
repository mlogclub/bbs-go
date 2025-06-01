<template>
  <ToolbarButton
    @click="editor?.chain().focus().setTextAlign('right').run()"
    :isActive="editor?.isActive({ textAlign: 'right' })"
    title="右对齐"
  >
    <LucideAlignRight :size="TOOLBAR_ICON_SIZE" />
  </ToolbarButton>
</template>

<script setup lang="ts">
import { Editor } from '@tiptap/core'
import { LucideAlignRight } from 'lucide-vue-next'
import ToolbarButton from './ToolbarButton.vue'
import { TOOLBAR_ICON_SIZE } from '../../constants/editor'

defineProps<{
  editor: Editor | null | undefined
}>()
</script> 