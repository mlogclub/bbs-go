<template>
  <ToolbarButton
    @click="editor?.chain().focus().toggleCodeBlock().run()"
    :isActive="editor?.isActive('codeBlock')"
    title="代码块"
  >
    <LucideCode2 :size="TOOLBAR_ICON_SIZE" />
  </ToolbarButton>
</template>

<script setup lang="ts">
import { Editor } from '@tiptap/core'
import { LucideCode2 } from 'lucide-vue-next'
import ToolbarButton from './ToolbarButton.vue'
import { TOOLBAR_ICON_SIZE } from '../../constants/editor'

defineProps<{
  editor: Editor | null | undefined
}>()
</script> 