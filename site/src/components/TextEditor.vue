<script>
export default {
  props: {
    height: {
      type: Number,
      default: 120,
    },
    value: {
      type: Object,
      default() {
        return {
          content: '',
          imageList: [],
        }
      },
    },
  },
  emits: ['submit', 'update:modelValue'],
  data() {
    return {
      post: this.value,
      showImageUpload: false, // 是否显示图片上传
      imageUploading: false, // 图片上传中
    }
  },
  methods: {
    doSubmit() {
      this.$emit('submit')
    },
    onInput() {
      this.$emit('update:modelValue', this.post)
    },
    isOnUpload() {
      return this.imageUploading
    },
    handleParse(e) {
      const items = e.clipboardData && e.clipboardData.items
      if (!items || !items.length) {
        return
      }

      let file = null
      for (let i = 0; i < items.length; i++) {
        if (items[i].type.includes('image')) {
          file = items[i].getAsFile()
        }
      }

      if (file) {
        e.preventDefault() // 阻止默认行为即不让剪贴板内容显示出来
        this.showImageUpload = true // 展开上传面板
        this.$refs.imageUploader.addFiles([file])
      }
    },
    handleDrag(e) {
      e.stopPropagation()
      e.preventDefault()

      const items = e.dataTransfer.items
      if (!items || !items.length) {
        return
      }

      const files = []
      for (let i = 0; i < items.length; i++) {
        if (items[i].type.includes('image')) {
          files.push(items[i].getAsFile())
        }
      }

      if (files && files.length) {
        this.showImageUpload = true // 展开上传面板
        this.$refs.imageUploader.addFiles(files)
      }
    },
    switchImageUpload() {
      if (!this.showImageUpload) {
        // 打开文件弹窗
        // this.$refs.imageUploader.onClick()
      }
      this.showImageUpload = !this.showImageUpload
    },
    clear() {
      this.post.content = ''
      this.post.imageList = []
      this.showImageUpload = false
      this.$refs.imageUploader.clear()
      this.onInput()
    },
    focus() {
      this.$refs.textarea.focus()
    },
  },
}
</script>

<template>
  <div class="text-editor">
    <textarea
      ref="textarea"
      v-model="post.content"
      placeholder="请输入您要发表的内容 ..."
      :style="{ 'min-height': `${height}px`, 'height': `${height}px` }"
      @input="onInput"
      @paste="handleParse"
      @drop="handleDrag"
      @keydown.ctrl.enter="doSubmit"
      @keydown.meta.enter="doSubmit"
    />
    <div v-show="showImageUpload" class="text-editor-image-uploader">
      <image-upload
        ref="imageUploader"
        v-model="post.imageList"
        v-model:on-upload="imageUploading"
        @input="onInput"
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
          <span>图片</span>
        </div>
      </div>
      <div class="text-editor-btn">
        <span>Ctrl/⌘ + Enter</span>
        <button class="button is-success is-small" @click="doSubmit">
          发布
        </button>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.text-editor {
  border: 1px solid var(--border-color);
  textarea {
    width: 100%;
    font-family: inherit;
    background: var(--bg-color2);
    border: 0;
    outline: 0;
    display: block;
    position: relative;
    resize: none;
    line-height: 1.8;
    padding: 15px 15px 20px;
    overflow: auto;
    overscroll-behavior: contain;
    transition: all 100ms linear;
    color: var(--text-color);
  }

  .text-editor-image-uploader {
    padding: 10px;
  }

  .text-editor-bar {
    background-color: var(--bg-color);
    border-top: 1px solid var(--border-color);
    padding: 5px;
    display: flex;
    align-items: center;
    justify-content: space-between;

    .text-editor-actions {
      .text-editor-action-item {
        cursor: pointer;
        color: var(--text-color3);
        user-select: none;

        i,
        span {
          font-size: 16px;
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
    }
  }
}
</style>
