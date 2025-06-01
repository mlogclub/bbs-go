<template>
  <ToolbarButton
    @click="editor?.chain().focus().setTextAlign('center').run()"
    :isActive="editor?.isActive({ textAlign: 'center' })"
    title="居中对齐"
  >
    <LucideAlignCenter :size="TOOLBAR_ICON_SIZE" />
  </ToolbarButton>
</template>

<script setup lang="ts">
import { Editor } from '@tiptap/core'
import { LucideAlignCenter } from 'lucide-vue-next'
import ToolbarButton from './ToolbarButton.vue'
import { TOOLBAR_ICON_SIZE } from '../../constants/editor'

defineProps<{
  editor: Editor | null | undefined
}>()
</script> 