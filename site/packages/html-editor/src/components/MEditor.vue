<template>
  <div class="m-editor-container" :style="{ height: props.height }">
    <EditorToolbar :editor="editor" :uploadImage="props.uploadImage" />
    <editor-content :editor="editor" class="editor-content" />
  </div>
</template>

<script setup lang="ts">
import "../styles/scrollbar.css";
import "../styles/theme.css";

import { ref, onMounted, onBeforeUnmount, watch } from "vue";
import { useEditor, EditorContent } from "@tiptap/vue-3";

import { Editor } from "@tiptap/core";
import StarterKit from "@tiptap/starter-kit";
import TextAlign from "@tiptap/extension-text-align";
import { TextStyleKit } from "@tiptap/extension-text-style";
import { TaskList, TaskItem } from "@tiptap/extension-list";
import { ResizableImage } from "../extensions/ResizableImage";

import { Typography } from "@tiptap/extension-typography";
import { Placeholder } from "@tiptap/extension-placeholder";

import { SlashSuggestion } from "../extensions/slash";
import { PasteImage } from "../extensions/PasteImage";

import EditorToolbar from "./toolbar/index.vue";
import { uploadImage } from "../utils/imageUtils";
import type { UploadImageFunction } from "../utils/imageUtils";

const props = withDefaults(
  defineProps<{
    modelValue: string;
    height?: string;
    uploadImage?: UploadImageFunction;
  }>(),
  {
    height: "400px",
    uploadImage,
  }
);

const emit = defineEmits<{
  (e: "update:modelValue", value: string): void;
}>();

const editorRef = ref<Editor | null>(null);

const editor = useEditor({
  content: props.modelValue,
  extensions: [
    StarterKit.configure({
      link: {
        openOnClick: false,
        HTMLAttributes: {
          target: "_blank",
          rel: "noopener noreferrer",
        },
      },
    }),
    TextAlign,
    TextStyleKit,
    SlashSuggestion,
    TaskList,
    TaskItem,
    Typography,
    ResizableImage,
    PasteImage.configure({
      uploadImage: props.uploadImage,
    }),
    Placeholder.configure({
      placeholder: "输入 / 插入内容",
    }),
  ],
  onCreate: ({ editor }) => {
    editorRef.value = editor;
  },
  onUpdate: ({ editor }) => {
    emit("update:modelValue", editor.getHTML());
  },
});

watch(
  () => props.modelValue,
  (newValue) => {
    const isSame = newValue === editor.value?.getHTML();
    if (editor.value && !isSame) {
      editor.value.commands.setContent(newValue);
    }
  }
);

onMounted(() => {
  editorRef.value = editor.value;
});

onBeforeUnmount(() => {
  editor.value?.destroy();
});
</script>

<style lang="scss">
.m-editor-container {
  border-radius: 2px;
  background: var(--editor-bg);
  border: 1px solid var(--editor-border);
  display: flex;
  flex-direction: column;
  // box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.editor-content {
  background: var(--editor-bg);
  border-radius: 0 0 8px 8px;
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  scrollbar-width: thin;
  scrollbar-color: var(--editor-scrollbar-thumb) var(--editor-scrollbar-track);
  height: 100%;

  .tiptap {
    flex: 1;
    outline: none;
    padding: 12px 12px;
    color: var(--editor-text);
    line-height: 1.7;
    font-size: 14px;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC",
      "Hiragino Sans GB", "Microsoft YaHei", "Helvetica Neue", Helvetica, Arial,
      sans-serif;

    :first-child {
      margin-top: 0;
    }

    /* 段落样式 */
    p {
      margin: 0 0 16px 0;
      line-height: 1.7;

      &:last-child {
        margin-bottom: 0;
      }

      &.is-editor-empty:first-child::before {
        color: var(--editor-placeholder);
        content: attr(data-placeholder);
        float: left;
        height: 0;
        pointer-events: none;
      }
    }

    /* 标题样式 */
    h1,
    h2,
    h3,
    h4,
    h5,
    h6 {
      font-weight: 600;
      margin: 32px 0 16px 0;
      line-height: 1.4;

      &:first-child {
        margin-top: 0;
      }
    }

    h1 {
      font-size: 1.75em;
      border-bottom: 2px solid var(--editor-border);
      padding-bottom: 8px;
    }

    h2 {
      font-size: 1.375em;
      border-bottom: 1px solid var(--editor-border);
      padding-bottom: 6px;
    }

    h3 {
      font-size: 1.125em;
    }

    h4 {
      font-size: 1em;
    }

    h5 {
      font-size: 0.875em;
    }

    h6 {
      font-size: 0.75em;
      color: #6b7280;
    }

    /* 列表样式 */
    ul,
    ol {
      margin: 16px 0;
      padding-left: 24px;

      li {
        margin: 8px 0;
        line-height: 1.7;

        p {
          margin: 4px 0;
        }

        /* 嵌套列表 */
        ul,
        ol {
          margin: 8px 0;
        }
      }
    }

    ul {
      list-style-type: disc;

      ul {
        list-style-type: circle;

        ul {
          list-style-type: square;
        }
      }
    }

    ol {
      list-style-type: decimal;

      ol {
        list-style-type: lower-alpha;

        ol {
          list-style-type: lower-roman;
        }
      }
    }

    /* 任务列表样式 */
    ul[data-type="taskList"] {
      list-style: none;
      margin: 16px 0;
      padding-left: 0;

      li {
        display: flex;
        align-items: flex-start;
        margin: 4px 0;
        padding: 4px 0;

        &:hover {
          background-color: var(--editor-hover);
          border-radius: 4px;
          margin-left: -8px;
          margin-right: -8px;
          padding-left: 8px;
          padding-right: 8px;
        }

        > label {
          flex: 0 0 auto;
          margin-right: 8px;
          margin-top: 2px;
          user-select: none;

          input[type="checkbox"] {
            cursor: pointer;
            width: 14px;
            height: 14px;
            accent-color: var(--editor-focus);
          }
        }

        > div {
          flex: 1 1 auto;
          min-width: 0;
        }

        /* 完成的任务样式 */
        &[data-checked="true"] {
          opacity: 0.6;

          > div {
            text-decoration: line-through;
          }
        }
      }

      /* 嵌套任务列表 */
      ul[data-type="taskList"] {
        margin: 8px 0 8px 24px;
      }
    }

    /* 引用样式 */
    blockquote {
      border-left: 4px solid var(--editor-blockquote-border);
      background-color: var(--editor-blockquote-bg);
      margin: 16px 0;
      padding: 12px 16px;
      font-style: italic;
      border-radius: 0 4px 4px 0;

      p {
        margin: 0;

        &:not(:last-child) {
          margin-bottom: 8px;
        }
      }
    }

    /* 代码样式 */
    code {
      background-color: var(--editor-code-bg);
      border-radius: 4px;
      padding: 2px 6px;
      font-family: "SF Mono", Monaco, Inconsolata, "Roboto Mono",
        "Source Code Pro", Menlo, Consolas, "DejaVu Sans Mono", monospace;
      font-size: 0.8em;
      color: var(--editor-code-text);
    }

    pre {
      background-color: var(--editor-code-bg);
      border-radius: 8px;
      margin: 16px 0;
      padding: 16px;
      overflow-x: auto;
      border: 1px solid var(--editor-border);

      code {
        background-color: transparent;
        padding: 0;
        border-radius: 0;
        color: inherit;
        font-size: 13px;
        line-height: 1.5;
      }
    }

    /* 水平分割线 */
    hr {
      border: none;
      border-top: 2px solid var(--editor-border);
      margin: 32px 0;
      opacity: 0.6;
    }

    /* 强调文本 */
    strong {
      font-weight: 600;
    }

    em {
      font-style: italic;
    }

    u {
      text-decoration: underline;
    }

    s {
      text-decoration: line-through;
    }

    /* 链接样式 */
    a {
      color: var(--editor-link-color, #3b82f6);
      text-decoration: none;
      border-bottom: 1px solid transparent;
      transition: all 0.2s ease;
      cursor: pointer;

      &:hover {
        border-bottom-color: var(--editor-link-color, #3b82f6);
        text-decoration: underline;
      }

      &:active {
        color: var(--editor-link-color, #2563eb);
      }
    }

    /* 表格样式 */
    table {
      border-collapse: collapse;
      width: 100%;
      margin: 16px 0;
      border-radius: 8px;
      overflow: hidden;
      border: 1px solid var(--editor-border);

      th,
      td {
        border: 1px solid var(--editor-border);
        padding: 12px 16px;
        text-align: left;
        line-height: 1.5;
      }

      th {
        background-color: var(--editor-table-header);
        font-weight: 600;
      }

      tr:nth-child(even) {
        background-color: var(--editor-hover);
      }

      tr:hover {
        background-color: rgba(var(--editor-focus), 0.1);
      }
    }

    /* 文本对齐 */
    .has-text-align-left {
      text-align: left;
    }

    .has-text-align-center {
      text-align: center;
    }

    .has-text-align-right {
      text-align: right;
    }

    .has-text-align-justify {
      text-align: justify;
    }

    /* 图片样式 */
    img {
      max-width: 100%;
      height: auto;
      border-radius: 4px;
      margin: 0;
      cursor: pointer;
    }

    /* 选择样式 */
    ::selection {
      background-color: var(--editor-selection);
    }

    /* 焦点样式 */
    &:focus {
      outline: none;
    }

    /* 段落间距优化 */
    > * + * {
      margin-top: 16px;
    }

    > h1 + *,
    > h2 + *,
    > h3 + *,
    > h4 + *,
    > h5 + *,
    > h6 + * {
      margin-top: 16px;
    }

    /* 响应式设计 */
    @media (max-width: 768px) {
      padding: 16px 16px;
      font-size: 13px;

      h1 {
        font-size: 1.5em;
      }

      h2 {
        font-size: 1.25em;
      }

      h3 {
        font-size: 1em;
      }

      ul,
      ol {
        padding-left: 20px;
      }

      blockquote {
        margin-left: 0;
        margin-right: 0;
        padding: 12px 12px;
      }

      pre {
        padding: 12px;
      }

      table {
        font-size: 12px;

        th,
        td {
          padding: 8px 12px;
        }
      }
    }
  }
}
</style>
