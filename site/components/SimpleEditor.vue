<template>
  <div class="simple-editor">
    <div class="simple-editor-toolbar">
      <div class="act-btn">
        <i class="iconfont icon-image" />
      </div>
      <div class="publish-container">
        <span class="tip">还能输入{{ maxWordCount - wordCount }}个字</span>
        <button class="button is-success is-small">发 布</button>
      </div>
    </div>
    <label class="simple-editor-input">
      <textarea
        v-model="post.content"
        placeholder="请输入您要发表的内容 ..."
        @keydown.ctrl.enter="doSubmit"
        @keydown.meta.enter="doSubmit"
        @input="onInput"
      ></textarea>
    </label>
    <div class="simple-editor-image-upload">
      <image-upload
        v-model="post.imageList"
        :file-list="post && post.imageList"
        :on-upload.sync="onUploadImage"
      />
    </div>
  </div>
</template>

<script>
import ImageUpload from '~/components/ImageUpload'
export default {
  name: 'SimpleEditor',
  components: { ImageUpload },
  props: {
    onUpload: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      maxWordCount: 5000,
      onUploadImage: false,
      post: {
        content: '',
        imageList: [],
      },
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
  },
  methods: {
    doSubmit() {
      this.$emit('submit')
    },
    onInput() {
      this.$emit('input', this.post)
      console.log('input', this.post)
    },
  },
}
</script>

<style lang="scss" scoped>
@import './assets/styles/variable.scss';

.simple-editor {
  border: 1px solid hsla(0, 0%, 89.4%, 0.6);
  border-radius: 3px;
  position: relative;
  width: 100%;

  .simple-editor-toolbar {
    width: 100%;
    height: 45px;
    display: flex;
    padding: 0 10px;
    align-items: center;
    background: #ffffff;
    position: sticky;
    top: 65px;
    z-index: 6;
    border-bottom: 1px solid $border-color-base;

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
        color: #d0d4dc;
      }
    }
  }

  .simple-editor-input {
    width: 100%;

    textarea {
      font-family: inherit;
      background: $background-color-editor;
      width: 100%;
      min-height: 200px;
      border: 0;
      outline: 0;
      display: block;
      position: relative;
      resize: none;
      line-height: 16px;
      padding: 15px 15px 20px;
      overflow: hidden;
      overscroll-behavior: contain;
      transition: all 100ms linear;
    }
  }

  .simple-editor-image-upload {
    background: $background-color-editor;
    padding: 20px 20px 30px;
  }
}
</style>
