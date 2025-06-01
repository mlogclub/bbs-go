<template>
  <ToolbarButton
    @click="editor?.chain().focus().toggleTaskList().run()"
    :isActive="editor?.isActive('taskList')"
    title="任务列表"
  >
    <LucideListTodo :size="TOOLBAR_ICON_SIZE" />
  </ToolbarButton>
</template>

<script setup lang="ts">
import { Editor } from "@tiptap/core";
import { LucideListTodo } from "lucide-vue-next";
import ToolbarButton from "./ToolbarButton.vue";
import { TOOLBAR_ICON_SIZE } from "../../constants/editor";

defineProps<{
  editor: Editor | null | undefined;
}>();
</script>
