<template>
  <ToolbarButton
    @click="editor?.chain().focus().toggleBulletList().run()"
    :isActive="editor?.isActive('bulletList')"
    title="无序列表"
  >
    <LucideList :size="TOOLBAR_ICON_SIZE" />
  </ToolbarButton>
</template>

<script setup lang="ts">
import { Editor } from '@tiptap/core'
import { LucideList } from 'lucide-vue-next'
import ToolbarButton from './ToolbarButton.vue'
import { TOOLBAR_ICON_SIZE } from '../../constants/editor'

defineProps<{
  editor: Editor | null | undefined
}>()
</script> 