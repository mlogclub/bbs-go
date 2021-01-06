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
        v-model="content"
        placeholder="请输入您要发表的内容 ..."
        @keydown.ctrl.enter="doSubmit"
        @keydown.meta.enter="doSubmit"
      ></textarea>
    </label>
  </div>
</template>

<script>
export default {
  name: 'SimpleEditor',
  data() {
    return {
      content: '',
      maxWordCount: 5000,
    }
  },
  computed: {
    hasContent() {
      return this.content && this.content.length > 0
    },
    wordCount() {
      return this.content ? this.content.length : 0
    },
    user() {
      return this.$store.state.user.current
    },
  },
  methods: {
    async doSubmit() {},
  },
}
</script>

<style lang="scss" scoped>
$border-color-base: rgba(228, 228, 228, 0.6);
$background-color-editor: #f5f6f7;

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
}
</style>
