<template>
  <ToolbarButton
    @click="editor?.chain().focus().toggleOrderedList().run()"
    :isActive="editor?.isActive('orderedList')"
    title="有序列表"
  >
    <LucideListOrdered :size="TOOLBAR_ICON_SIZE" />
  </ToolbarButton>
</template>

<script setup lang="ts">
import { Editor } from '@tiptap/core'
import { LucideListOrdered } from 'lucide-vue-next'
import ToolbarButton from './ToolbarButton.vue'
import { TOOLBAR_ICON_SIZE } from '../../constants/editor'

defineProps<{
  editor: Editor | null | undefined
}>()
</script> 