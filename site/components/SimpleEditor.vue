<template>
  <div class="simple-editor">
    <div class="simple-editor-toolbar">
      <div class="act-btn">
        <i class="iconfont icon-image" @click="switchImageUpload" />
      </div>
      <div class="publish-container">
        <span class="tip">{{ wordCount }} / {{ maxWordCount }}</span>
      </div>
    </div>
    <label class="simple-editor-input">
      <textarea
        v-model="post.content"
        placeholder="请输入您要发表的内容 ..."
        :style="{ 'min-height': height, height: height }"
        @input="onInput"
        @paste="handleParse"
        @drop="handleDrag"
        @keydown.ctrl.enter="doSubmit"
        @keydown.meta.enter="doSubmit"
      ></textarea>
    </label>
    <div v-show="showImageUpload" class="simple-editor-image-upload">
      <image-upload
        ref="imageUploader"
        v-model="post.imageList"
        :on-upload.sync="imageUploading"
        @input="onInput"
      />
    </div>
  </div>
</template>

<script>
export default {
  props: {
    maxWordCount: {
      type: Number,
      default: 5000,
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
    height: {
      type: String,
      default: '200px',
    },
  },
  data() {
    return {
      imageUploading: false, // 图片上传中...
      showImageUpload: false,
      post: this.value,
    }
  },
  computed: {
    hasContent() {
      return this.post.content && this.post.content.length > 0
    },
    wordCount() {
      return this.post.content ? this.post.content.length : 0
    },
    user() {
      return this.$store.state.user.current
    },
    loading() {
      return this.imageUploading
    },
  },
  methods: {
    doSubmit() {
      this.$emit('submit')
    },
    onInput() {
      this.$emit('input', this.post)
    },
    isOnUpload() {
      return this.imageUploading
    },
    switchImageUpload() {
      if (!this.showImageUpload) {
        // 打开文件弹窗
        this.$refs.imageUploader.onClick()
      }
      this.showImageUpload = !this.showImageUpload
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
    clear() {
      this.post.content = ''
      this.post.imageList = []
      this.$refs.imageUploader.clear()
      this.onInput()
    },
  },
}
</script>

<style lang="scss" scoped>
.simple-editor {
  border: 1px solid var(--border-color);
  border-radius: 3px;
  position: relative;
  width: 100%;

  .simple-editor-toolbar {
    width: 100%;
    height: 45px;
    display: flex;
    padding: 0 10px;
    align-items: center;
    background-color: var(--bg-color);
    top: 65px;
    z-index: 6;
    border-bottom: 1px solid var(--border-color);

    .act-btn {
      display: flex;
      padding: 0 10px;

      i {
        cursor: pointer;
        margin-left: 20px;
        font-size: 24px;

        &:first-child {
          margin-left: 0;
        }
      }
    }

    .publish-container {
      margin-left: auto;

      > .button-publish {
        margin-left: auto;

        ::v-deep span {
          font-size: 14px;
        }
      }

      > .tip {
        font-size: 14px;
        margin-right: 10px;
        color: var(--text-color4);
      }
    }
  }

  .simple-editor-input {
    width: 100%;

    textarea {
      font-family: inherit;
      background: var(--bg-color2);
      width: 100%;
      // min-height: 200px;
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
  }

  .simple-editor-image-upload {
    background: var(--bg-color2);
    padding: 20px 20px 20px;
  }
}
</style>
