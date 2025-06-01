<template>
  <ToolbarButton
    @click="editor?.chain().focus().toggleBlockquote().run()"
    :isActive="editor?.isActive('blockquote')"
    title="引用块"
  >
    <LucideQuote :size="TOOLBAR_ICON_SIZE" />
  </ToolbarButton>
</template>

<script setup lang="ts">
import { Editor } from '@tiptap/core'
import { LucideQuote } from 'lucide-vue-next'
import ToolbarButton from './ToolbarButton.vue'
import { TOOLBAR_ICON_SIZE } from '../../constants/editor'

defineProps<{
  editor: Editor | null | undefined
}>()
</script> 