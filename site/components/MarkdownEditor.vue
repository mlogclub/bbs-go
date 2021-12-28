<template>
  <div class="bbsgoEditor">
    <v-md-editor
      v-model="content"
      :left-toolbar="toolbars"
      :right-toolbar="rightToolbar"
      :height="height"
      :placeholder="placeholder"
      :disabled-menus="[]"
      mode="edit"
      @change="change"
      @upload-image="uploadImage"
      @keydown.ctrl.enter.native="submit"
      @keydown.meta.enter.native="submit"
    ></v-md-editor>
  </div>
</template>

<script>
import Vue from 'vue'
import VMdEditor from '@kangc/v-md-editor'
import '@kangc/v-md-editor/lib/style/base-editor.css'
import githubTheme from '@kangc/v-md-editor/lib/theme/github.js'
import '@kangc/v-md-editor/lib/theme/style/github.css'

// highlightjs
import hljs from 'highlight.js'

VMdEditor.use(githubTheme, {
  Hljs: hljs,
})

Vue.use(VMdEditor)

export default {
  props: {
    value: {
      type: String,
      default: '',
    },
    height: {
      type: String,
      default: '400px', // normal、mini
    },
    placeholder: {
      type: String,
      default: '请输入...',
    },
  },
  data() {
    return {
      width: '100%',
      content: this.value,
    }
  },
  computed: {
    isMobile() {
      return this.$store.state.env.isMobile
    },
    toolbars() {
      if (this.isMobile) {
        return 'h bold italic strikethrough'
      } else {
        return 'undo redo clear | h bold italic strikethrough quote | ul ol table hr | link image code'
      }
    },
    rightToolbar() {
      if (this.$store.state.env.isMobile) {
        return 'fullscreen'
      }
      return 'preview sync-scroll fullscreen'
    },
  },
  mounted() {},
  methods: {
    submit() {
      this.$emit('submit', this.content)
    },
    /**
     * 上传图片
     */
    async uploadImage(event, insertImage, files) {
      if (!files || !files.length) {
        return
      }
      for (let i = 0; i < files.length; i++) {
        const file = files[i]
        const formData = new FormData()
        formData.append('image', file, file.name)
        const ret = await this.$axios.post('/api/upload', formData, {
          headers: { 'Content-Type': 'multipart/form-data' },
        })
        insertImage({
          url: ret.url,
          desc: ' ',
        })
      }
    },
    change(value, render) {
      this.$emit('input', value)
    },
    /**
     * 清空编辑器内容
     */
    clear() {
      this.content = ''
      this.$emit('input', this.content)
    },
    /**
     * 清理缓存
     */
    clearCache() {},
  },
}
</script>

<style lang="scss">
.bbsgoEditor {
  .v-md-editor {
    box-shadow: none !important;
    border: 1px solid var(--border-color2);

    .v-md-editor__toolbar {
      background-color: var(--text-color5);
      padding: 3px;

      .v-md-editor__toolbar-item {
        font-size: 14px !important;
      }
    }

    .v-md-editor__editor-wrapper {
      background-color: var(--text-color5) !important;
    }

    .v-md-editor__preview-wrapper {
      background-color: rgb(251, 251, 251) !important;
    }
  }
}
</style>
