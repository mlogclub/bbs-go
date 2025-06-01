<template>
  <ToolbarButton
    @click="editor?.chain().focus().toggleBold().run()"
    :isActive="editor?.isActive('bold')"
    title="粗体"
  >
    <LucideBold :size="TOOLBAR_ICON_SIZE" />
  </ToolbarButton>
</template>

<script setup lang="ts">
import { Editor } from '@tiptap/core'
import { LucideBold } from 'lucide-vue-next'
import ToolbarButton from './ToolbarButton.vue'
import { TOOLBAR_ICON_SIZE } from '../../constants/editor'

defineProps<{
  editor: Editor | null | undefined
}>()
</script> 