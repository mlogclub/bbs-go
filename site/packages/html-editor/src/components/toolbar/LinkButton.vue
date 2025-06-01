<template>
  <ToolbarButton
    ref="linkButtonRef"
    @click="toggleLink"
    :isActive="editor?.isActive('link')"
    title="连接"
  >
    <LucideLink ref="linkButtonRef" :size="TOOLBAR_ICON_SIZE" />
  </ToolbarButton>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, nextTick } from "vue";
import type { Editor } from "@tiptap/core";
import { LucideLink } from "lucide-vue-next";
import ToolbarButton from "./ToolbarButton.vue";
import { TOOLBAR_ICON_SIZE } from "../../constants/editor";
import tippy from "tippy.js";
import "tippy.js/dist/tippy.css";

const props = defineProps<{
  editor: Editor | null | undefined;
}>();

const linkButtonRef = ref();
let tippyInstance = null;
const linkUrl = ref("");
const linkText = ref("");

// 创建链接弹窗的 HTML 内容
const createLinkDialogContent = (isEditing: boolean = false) => {
  return `
    <div class="link-dialog-content">
      <div class="link-input-group">
        <label>链接文本:</label>
        <input
          id="link-text-input"
          type="text"
          placeholder="请输入链接显示文本"
          value="${linkText.value}"
        />
      </div>
      <div class="link-input-group">
        <label>链接地址:</label>
        <input
          id="link-url-input"
          type="text"
          placeholder="请输入链接地址"
          value="${linkUrl.value}"
        />
      </div>
      <div class="link-dialog-actions">
        <button id="confirm-link" class="btn-primary">确定</button>
        ${
          isEditing
            ? '<button id="remove-link" class="btn-danger">移除链接</button>'
            : ""
        }
        <button id="cancel-link" class="btn-secondary">取消</button>
      </div>
    </div>
  `;
};

// 初始化 tippy 实例
const initTippy = () => {
  if (!linkButtonRef.value) return;

  // 链接弹窗
  tippyInstance = tippy(linkButtonRef.value.$el, {
    content: "",
    allowHTML: true,
    interactive: true,
    trigger: "manual",
    placement: "bottom-start",
    maxWidth: 320,
    offset: [0, 8],
    hideOnClick: false,
    theme: "link-tippy",
    animation: "shift-away",
    onShow: () => {
      // 聚焦到 URL 输入框
      setTimeout(() => {
        const urlInput = document.getElementById(
          "link-url-input"
        ) as HTMLInputElement;
        if (urlInput) {
          urlInput.focus();
          urlInput.select();
        }
      }, 50);
    },
    onHidden: () => {
      // 清理状态
      linkUrl.value = "";
      linkText.value = "";
    },
  });
};

// 绑定弹窗内的事件
const bindDialogEvents = () => {
  const urlInput = document.getElementById(
    "link-url-input"
  ) as HTMLInputElement;
  const textInput = document.getElementById(
    "link-text-input"
  ) as HTMLInputElement;
  const confirmBtn = document.getElementById("confirm-link");
  const removeBtn = document.getElementById("remove-link");
  const cancelBtn = document.getElementById("cancel-link");

  if (urlInput) {
    urlInput.addEventListener("input", (e) => {
      linkUrl.value = (e.target as HTMLInputElement).value;
    });
    urlInput.addEventListener("keydown", (e) => {
      if (e.key === "Enter") {
        e.preventDefault();
        confirmLink();
      } else if (e.key === "Escape") {
        e.preventDefault();
        closeLinkDialog();
      }
    });
  }

  if (textInput) {
    textInput.addEventListener("input", (e) => {
      linkText.value = (e.target as HTMLInputElement).value;
    });
    textInput.addEventListener("keydown", (e) => {
      if (e.key === "Enter") {
        e.preventDefault();
        confirmLink();
      } else if (e.key === "Escape") {
        e.preventDefault();
        closeLinkDialog();
      }
    });
  }

  if (confirmBtn) {
    confirmBtn.addEventListener("click", confirmLink);
  }

  if (removeBtn) {
    removeBtn.addEventListener("click", removeLink);
  }

  if (cancelBtn) {
    cancelBtn.addEventListener("click", closeLinkDialog);
  }
};

const toggleLink = () => {
  if (!props.editor || !tippyInstance) return;

  if (props.editor.isActive("link")) {
    // 如果当前有链接，显示编辑对话框
    const { href } = props.editor.getAttributes("link");
    linkUrl.value = href || "";
    linkText.value = props.editor.state.doc.textBetween(
      props.editor.state.selection.from,
      props.editor.state.selection.to
    );

    tippyInstance.setContent(createLinkDialogContent(true));
    tippyInstance.show();

    // 延迟绑定事件，确保 DOM 已更新
    setTimeout(bindDialogEvents, 50);
  } else {
    // 如果没有链接，显示添加对话框
    const selectedText = props.editor.state.doc.textBetween(
      props.editor.state.selection.from,
      props.editor.state.selection.to
    );
    linkUrl.value = "";
    linkText.value = selectedText;

    tippyInstance.setContent(createLinkDialogContent(false));
    tippyInstance.show();

    // 延迟绑定事件，确保 DOM 已更新
    setTimeout(bindDialogEvents, 50);
  }
};

const confirmLink = () => {
  if (!props.editor || !linkUrl.value) return;

  // 如果没有选中文本且没有输入链接文本，使用链接地址作为文本
  if (!linkText.value && props.editor.state.selection.empty) {
    linkText.value = linkUrl.value;
  }

  // 如果有链接文本但没有选中文本，先插入文本
  if (linkText.value && props.editor.state.selection.empty) {
    props.editor.chain().focus().insertContent(linkText.value).run();
    // 选中刚插入的文本
    const { from } = props.editor.state.selection;
    props.editor
      .chain()
      .setTextSelection({ from: from - linkText.value.length, to: from })
      .run();
  }

  // 添加或更新链接
  props.editor
    .chain()
    .focus()
    .extendMarkRange("link")
    .setLink({ href: linkUrl.value, target: "_blank" })
    .run();

  closeLinkDialog();
};

const removeLink = () => {
  if (!props.editor) return;
  props.editor.chain().focus().extendMarkRange("link").unsetLink().run();
  closeLinkDialog();
};

const closeLinkDialog = () => {
  if (tippyInstance) {
    tippyInstance.hide();
  }
};

// 键盘快捷键支持
if (typeof window !== "undefined") {
  const handleKeydown = (e: KeyboardEvent) => {
    if ((e.ctrlKey || e.metaKey) && e.key === "k") {
      e.preventDefault();
      toggleLink();
    }
  };

  document.addEventListener("keydown", handleKeydown);

  onBeforeUnmount(() => {
    document.removeEventListener("keydown", handleKeydown);
  });
}

onMounted(() => {
  nextTick(() => {
    initTippy();
  });
});

onBeforeUnmount(() => {
  if (tippyInstance) {
    tippyInstance.destroy();
  }
});
</script>

<style lang="scss">
/* 全局样式，用于 tippy 弹窗内容 */
.tippy-box[data-theme~="link-tippy"] {
  background-color: var(--editor-bg);
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);

  .tippy-arrow {
    color: var(--editor-bg);
  }
}

/* 暗色主题适配 */
@media (prefers-color-scheme: dark) {
  .tippy-box[data-theme~="link-tippy"] {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  }
}

.link-dialog-content {
  padding: 4px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-width: 280px;

  .link-input-group {
    display: flex;
    flex-direction: column;
    gap: 4px;

    label {
      font-size: 12px;
      font-weight: 500;
      color: var(--editor-text);
    }

    input {
      padding: 8px 12px;
      border: 1px solid var(--editor-border);
      border-radius: 4px;
      background: var(--editor-bg);
      color: var(--editor-text);
      font-size: 14px;
      outline: none;
      transition: border-color 0.2s ease;

      &:focus {
        border-color: var(--editor-focus);
        box-shadow: 0 0 0 2px rgba(96, 165, 250, 0.2);
      }

      &::placeholder {
        color: var(--editor-placeholder);
      }
    }
  }

  .link-dialog-actions {
    display: flex;
    gap: 8px;
    justify-content: flex-end;

    button {
      padding: 6px 12px;
      border: none;
      border-radius: 4px;
      font-size: 12px;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.2s ease;

      &.btn-primary {
        background: var(--editor-focus);
        color: white;

        &:hover {
          opacity: 0.9;
          transform: translateY(-1px);
        }
      }

      &.btn-danger {
        background: #ef4444;
        color: white;

        &:hover {
          background: #dc2626;
          transform: translateY(-1px);
        }
      }

      &.btn-secondary {
        background: var(--editor-hover);
        color: var(--editor-text);
        border: 1px solid var(--editor-border);

        &:hover {
          background: var(--editor-border);
          transform: translateY(-1px);
        }
      }

      &:focus {
        outline: none;
        box-shadow: 0 0 0 2px rgba(96, 165, 250, 0.2);
      }
    }
  }
}
</style>
