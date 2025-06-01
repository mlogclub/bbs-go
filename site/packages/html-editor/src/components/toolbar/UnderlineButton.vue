<template>
  <ToolbarButton
    @click="editor?.chain().focus().toggleUnderline().run()"
    :isActive="editor?.isActive('underline')"
    title="下划线"
  >
    <LucideUnderline :size="TOOLBAR_ICON_SIZE" />
  </ToolbarButton>
</template>

<script setup lang="ts">
import { Editor } from "@tiptap/core";
import { LucideUnderline } from "lucide-vue-next";
import ToolbarButton from "./ToolbarButton.vue";
import { TOOLBAR_ICON_SIZE } from "../../constants/editor";

defineProps<{
  editor: Editor | null | undefined;
}>();
</script>
