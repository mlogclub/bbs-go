<template>
  <ToolbarButton
    :title="isFullscreen ? '退出全屏' : '全屏'"
    :is-active="isFullscreen"
    @click="toggleFullscreen"
  >
    <template v-if="isFullscreen">
      <LucideMinimize :size="TOOLBAR_ICON_SIZE" />
    </template>
    <template v-else>
      <LucideMaximize :size="TOOLBAR_ICON_SIZE" />
    </template>
  </ToolbarButton>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { Editor } from "@tiptap/core";
import ToolbarButton from "./ToolbarButton.vue";
import { LucideMaximize, LucideMinimize } from "lucide-vue-next";
import { TOOLBAR_ICON_SIZE } from "../../constants/editor";

const props = defineProps<{
  editor: Editor | null | undefined;
}>();

const isFullscreen = ref(false);

onMounted(() => {
  document.addEventListener("fullscreenchange", handleFullscreenChange);
});

onUnmounted(() => {
  document.removeEventListener("fullscreenchange", handleFullscreenChange);
});

function handleFullscreenChange() {
  isFullscreen.value = !!document.fullscreenElement;
}

function toggleFullscreen() {
  if (!props.editor) return;

  if (!document.fullscreenElement) {
    // 从当前编辑器元素向上查找编辑器容器
    const editorDOM = props.editor.view.dom;
    if (!editorDOM || !(editorDOM instanceof HTMLElement)) return;

    // 首先尝试找到 .editor-container 类的容器
    let container = editorDOM.closest(".m-editor-container");

    if (!container) {
      // 如果找不到特定类的容器，向上寻找父元素直到找到合适的容器
      let parent = editorDOM.parentElement;
      while (parent && !parent.classList.contains("m-editor-container")) {
        if (parent.tagName === "BODY") break;
        parent = parent.parentElement;
      }
      container = parent;
    }

    // 如果还是找不到合适的容器，就直接使用编辑器元素
    const targetElement = (container || editorDOM) as HTMLElement;

    targetElement.requestFullscreen().catch((err) => {
      console.error(`无法切换到全屏模式: ${err.message}`);
    });
  } else {
    document.exitFullscreen();
  }
}
</script>
