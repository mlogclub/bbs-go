<template>
  <div
    class="text-editor"
    :class="{ 'has-editor-focus': focusHeight > 0, 'is-focus': isFocus }"
    v-click-outside="onBlur"
  >
    <textarea
      ref="textarea"
      v-model="post.content"
      :placeholder="$t('component.textEditor.placeholder')"
      @paste="handleParse"
      @drop="handleDrag"
      @keydown.ctrl.enter="doSubmit"
      @keydown.meta.enter="doSubmit"
      @focus="onFocus"
      @input="onContentChange"
    />
    <div
      ref="imageUploaderContainer"
      v-show="showImageUpload"
      class="text-editor-image-uploader"
    >
      <image-upload
        ref="imageUploader"
        v-model="post.imageList"
        v-model:on-upload="imageUploading"
        size="60px"
      />
    </div>
    <div class="text-editor-bar">
      <div class="text-editor-actions">
        <div
          class="text-editor-action-item"
          :class="{ active: showImageUpload }"
          @click="switchImageUpload"
        >
          <i class="iconfont icon-image" />
        </div>
      </div>
      <div class="text-editor-btn">
        <button class="button is-primary" @click="doSubmit">
          {{ $t("component.textEditor.publish") }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
const { t } = useI18n();
const props = defineProps({
  height: {
    type: Number,
    default: 80,
  },
  focusHeight: {
    type: Number,
    default: 0,
  },
  content: {
    type: String,
    default: "",
  },
  imageList: {
    type: Array,
    default: () => [],
  },
});

const emit = defineEmits(["submit", "update:content", "update:imageList"]);

const post = ref({
  content: props.content,
  imageList: props.imageList,
});
const isFocus = ref(false);
const textarea = ref(null);
const showImageUpload = ref(false); // 是否显示图片上传
const imageUploader = ref(null);
const imageUploaderContainer = ref(null);
const imageUploading = ref(false); // 图片上传中
const dynamicHeight = ref(props.height);
const dynamicFocusHeight = ref(props.focusHeight);

watch(
  () => props.content,
  (newVal, oldVal) => {
    post.value.content = props.content;
  }
);

watch(
  () => props.imageList,
  (newVal, oldVal) => {
    post.value.imageList = props.imageList || [];
  }
);

watch(
  () => showImageUpload.value,
  (newVal, oldVal) => {
    updateHeight();
  }
);

watch(
  () => post.value.content,
  (newVal, oldVal) => {
    emit("update:content", post.value.content);
  }
);

watch(
  () => post.value.imageList,
  (newVal, oldVal) => {
    emit("update:imageList", post.value.imageList);
  }
);

const doSubmit = () => {
  if (imageUploading.value === true) {
    useMsgWarning(t("component.textEditor.pleaseWait"));
    return;
  }
  emit("submit");
};

const handleParse = (e) => {
  const items = e.clipboardData && e.clipboardData.items;
  if (!items || !items.length) {
    return;
  }

  let file = null;
  for (let i = 0; i < items.length; i++) {
    if (items[i].type.includes("image")) {
      file = items[i].getAsFile();
    }
  }

  if (file) {
    e.preventDefault(); // 阻止默认行为即不让剪贴板内容显示出来
    showImageUpload.value = true; // 展开上传面板
    imageUploader.value.addFiles([file]);
    focus();
  }
};
const handleDrag = (e) => {
  e.stopPropagation();
  e.preventDefault();

  const items = e.dataTransfer.items;
  if (!items || !items.length) {
    return;
  }

  const files = [];
  for (let i = 0; i < items.length; i++) {
    if (items[i].type.includes("image")) {
      files.push(items[i].getAsFile());
    }
  }

  if (files && files.length) {
    showImageUpload.value = true; // 展开上传面板
    imageUploader.value.addFiles(files);
    focus();
  }
};
const switchImageUpload = () => {
  if (!showImageUpload.value) {
    // 打开文件弹窗
    imageUploader.value.onClick();
    focus();
  }
  showImageUpload.value = !showImageUpload.value;
};
const focus = () => {
  textarea.value.focus();
};
const onFocus = () => {
  isFocus.value = true;
  if (post.value.imageList && post.value.imageList.length) {
    showImageUpload.value = true;
  }
};
const onBlur = () => {
  isFocus.value = false;
  showImageUpload.value = false;
};
const reset = () => {
  isFocus.value = false;
  showImageUpload.value = false;
};
const updateHeight = () => {
  // 等dom变更完成后在计算高度，否则计算不准
  nextTick(() => {
    const uploadHeight = imageUploaderContainer.value.offsetHeight;
    if (showImageUpload.value) {
      dynamicHeight.value = props.height + uploadHeight;
      dynamicFocusHeight.value = props.focusHeight + uploadHeight;
    } else {
      dynamicHeight.value = props.height;
      dynamicFocusHeight.value = props.focusHeight;
    }
  });
};

const onContentChange = () => {
  if (post.value.content) {
    isFocus.value = true;
  }
};

defineExpose({
  reset,
  focus,
});
</script>

<style lang="scss" scoped>
.text-editor {
  --text-editor-height: v-bind(dynamicHeight + "px");
  --text-editor-focus-height: v-bind(dynamicFocusHeight + "px");

  $bgColor: var(--bg-color3);
  $borderRadius: 8px;

  background: $bgColor;
  border-radius: $borderRadius;
  height: var(--text-editor-height);
  border: 1px solid transparent;

  display: flex;
  flex-direction: column;
  transition: all 200ms;

  &.is-focus {
    &.has-editor-focus {
      height: var(--text-editor-focus-height) !important;
    }

    border: 1px solid var(--border-hover-color);
    background: var(--bg-color) !important;
    textarea,
    .text-editor-bar {
      background: var(--bg-color) !important;
    }
  }

  textarea {
    background: $bgColor;
    border-top-left-radius: $borderRadius;
    border-top-right-radius: $borderRadius;
    height: calc(100% - 36px);
    width: 100%;
    font-family: inherit;
    border: 0;
    outline: 0;
    display: block;
    position: relative;
    resize: none;
    line-height: 1.8;
    padding: 10px;
    overflow: auto;
    overscroll-behavior: contain;
    color: var(--text-color);

    &::-webkit-scrollbar {
      width: 6px;
      height: 6px;
      border-radius: 3px;
      background-color: transparent;
    }
    &::-webkit-scrollbar-thumb {
      background-color: transparent;
      border-radius: 3px;
      border: 2px solid transparent;
    }
    &:hover::-webkit-scrollbar {
      background-color: #e4e4e5;
    }
    &:hover::-webkit-scrollbar-thumb {
      background-color: #c0bebc;
    }
  }

  .text-editor-image-uploader {
    padding: 10px;
  }

  .text-editor-bar {
    background: $bgColor;
    border-bottom-left-radius: $borderRadius;
    border-bottom-right-radius: $borderRadius;
    padding: 3px 10px;
    display: flex;
    align-items: center;
    justify-content: space-between;

    .text-editor-actions {
      .text-editor-action-item {
        cursor: pointer;
        color: var(--text-color3);
        user-select: none;
        display: flex;
        align-items: center;
        column-gap: 5px;

        i {
          font-size: 20px;
        }

        &:hover {
          color: var(--text-link-color);
        }

        &.active {
          color: var(--text-link-color);
          font-weight: 500;
        }
      }
    }
    .text-editor-btn {
      display: flex;
      align-items: center;
      span {
        font-size: 12px;
        color: var(--text-color3);
        margin-right: 5px;
      }
      button {
        font-size: 12px;
        padding: 3px 10px;
      }
    }
  }
}
</style>
