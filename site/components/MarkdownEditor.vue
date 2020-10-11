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
      default: 'vditor',
    },
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
      isLoading: true,
      vditor: null,
      width: '100%',
    }
  },
  mounted() {
    this.doInit()
  },
  methods: {
    /**
     * 初始化编辑器
     */
    doInit() {
      const me = this
      if (process.client) {
        this.vditor = new window.Vditor(
          this.editorId,
          this.getOptions(function () {
            me.vditor.setValue(me.value)
          })
        )
      }
    },
    getOptions(afterFunc) {
      const me = this
      const userToken = me.$cookies.get('userToken')
      return {
        width: me.width,
        height: me.height,
        toolbarConfig: {
          pin: true,
        },
        toolbar: [
          'emoji',
          'headings',
          'bold',
          'italic',
          'strike',
          'link',
          '|',
          'list',
          'ordered-list',
          'check',
          'outdent',
          'indent',
          '|',
          'quote',
          'line',
          'code',
          'inline-code',
          'insert-before',
          'insert-after',
          '|',
          'upload',
          'record',
          'table',
          '|',
          'undo',
          'redo',
          '|',
          'edit-mode',
          'fullscreen',
          {
            name: 'more',
            toolbar: [
              'both',
              'code-theme',
              'content-theme',
              'export',
              'outline',
              'preview',
              'devtools',
              'info',
              'help',
            ],
          },
        ],
        placeholder: me.placeholder,
        cache: {
          enable: false,
        },
        counter: '999999',
        delay: 500,
        theme: 'classic',
        preview: {
          mode: 'editor',
          hljs: {
            enable: true,
            style: 'github',
            lineNumber: true,
          },
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
          },
        },
        after: afterFunc || function () {},
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
    },
  },
  head() {
    return {
      link: [
        {
          rel: 'stylesheet',
          href: '//cdn.jsdelivr.net/npm/vditor@3.5.4/dist/index.css',
        },
      ],
      script: [
        {
          src: '//cdn.jsdelivr.net/npm/vditor@3.5.4/dist/index.min.js',
        },
      ],
    }
  },
}
</script>

<style lang="scss" scoped></style>
