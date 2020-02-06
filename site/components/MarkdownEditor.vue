<template>
  <div
    :id="editorId"
    :style="{ width: width, height: height }"
    class="vditor"
  />
</template>

<script>
export default {
  props: {
    editorId: {
      type: String,
      default: 'vditor'
    },
    value: {
      type: String,
      default: ''
    },
    height: {
      type: String,
      default: '400px' // normal、mini
    },
    placeholder: {
      type: String,
      default: '请输入...'
    }
  },
  data() {
    return {
      isLoading: true,
      vditor: null,
      width: '100%',
      toolbar: [
        'emoji',
        'headings',
        'bold',
        'italic',
        'strike',
        '|',
        'line',
        'quote',
        'list',
        'ordered-list',
        'check',
        'code',
        'inline-code',
        'undo',
        'redo',
        'upload',
        'link',
        'table',
        'both',
        'preview',
        'format',
        'fullscreen'
      ]
    }
  },
  mounted() {
    this.doInit()
    this.$nextTick(async () => {
      if (this.vditor) {
        await this.vditor.getHTML(true)
        this.isLoading = false
      }
    })
  },
  methods: {
    /**
     * 初始化编辑器
     */
    doInit() {
      if (!process.client) {
        return
      }
      const me = this
      const userToken = this.$cookies.get('userToken')
      const options = {
        width: me.width,
        height: me.height,
        toolbar: me.toolbar,
        placeholder: me.placeholder,
        cache: true,
        counter: '999999',
        delay: 500,
        mode: 'markdown-show',
        preview: {
          mode: 'editor',
          hljs: {
            enable: true,
            style: 'GitHub',
            lineNumber: true
          }
        },
        input(val) {
          me.$emit('input', val)
        },
        ctrlEnter(val) {
          me.$emit('input', val)
          me.$emit('submit', val)
        },
        upload: {
          accept: 'image/*',
          url: '/api/upload/editor?userToken=' + userToken,
          linkToImgUrl: '/api/upload/fetch?userToken=' + userToken,
          filename(name) {
            return name.replace(/\?|\\|\/|:|\||<|>|\*|\[|\]|\s+/g, '-')
          }
        }
      }
      this.vditor = new window.Vditor(me.editorId, options)
      if (this.value) {
        this.vditor.setValue(this.value)
      }
    },
    /**
     * 清空编辑器内容
     */
    clear() {
      if (this.vditor) {
        this.value = ''
        this.vditor.setValue('')
        this.clearCache()
      }
    },
    /**
     * 清理缓存
     */
    clearCache() {
      if (this.vditor) {
        this.vditor.clearCache()
      }
    }
  },
  head() {
    return {
      link: [
        {
          rel: 'stylesheet',
          href: '//cdn.jsdelivr.net/npm/vditor/dist/index.classic.css'
        }
      ],
      script: [
        {
          src: '//cdn.jsdelivr.net/npm/vditor/dist/index.min.js'
        }
      ]
    }
  }
}
</script>

<style lang="scss" scoped>
.vditor {
  border: 1px solid #d1d5da;
  width: 100%;
}
</style>
