<template>
  <div class="text-editor">
    <textarea
      placeholder="请输入您要发表的内容 ..."
      :style="{ 'min-height': height + 'px', height: height + 'px' }"
      @input="onInput"
      @paste="handleParse"
      @drop="handleDrag"
      @keydown.ctrl.enter="doSubmit"
      @keydown.meta.enter="doSubmit"
    ></textarea>
    <div class="text-editor-bar">
      <div class="text-editor-actions">
        <div class="text-editor-action-item">
          <i class="iconfont icon-image" />
          <span>图片</span>
        </div>
      </div>
      <div>
        <button class="button is-primary is-small">发布</button>
      </div>
    </div>
  </div>
</template>

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
  data() {
    return {
      post: this.value,
    }
  },
  methods: {
    doSubmit() {
      console.log('submit...')
      this.$emit('submit')
    },
    onInput() {
      console.log('input...')
      this.$emit('input', this.post)
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
        this.$refs.imageUploader.addFiles(new Array(file))
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
  },
}
</script>
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
  }

  .text-editor-bar {
    background-color: var(--bg-color);
    border-top: 1px solid var(--border-color);
    padding: 10px;
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
  }
}
</style>
