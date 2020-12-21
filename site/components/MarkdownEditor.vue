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
      console.log('do init...')
      const me = this
      if (process.client) {
        me.vditor = new window.Vditor(
          me.editorId,
          me.getOptions(function () {
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
        toolbar: me.getToolbars(),
        mode: 'sv',
        toolbarConfig: {
          pin: false,
          hide: false,
        },
        placeholder: me.placeholder,
        cache: {
          enable: false,
        },
        counter: {
          enable: true,
          type: 'text',
        },
        delay: 200,
        theme: 'classic',
        preview: {
          mode: 'editor',
          markdown: {
            toc: true,
            mark: true,
          },
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
    getToolbars() {
      if (this.$isMobile()) {
        return ['emoji', 'bold', 'italic', 'strike', 'fullscreen']
      } else {
        return [
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
          'upload',
          'table',
          '|',
          'undo',
          'redo',
          '|',
          'outline',
          'edit-mode',
          'preview',
          'both',
          'fullscreen',
          {
            name: 'more',
            toolbar: ['devtools'],
          },
        ]
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
          href: '//cdn.jsdelivr.net/npm/vditor@3.5.5/dist/index.css',
        },
      ],
      script: [
        {
          src: '//cdn.jsdelivr.net/npm/vditor@3.5.5/dist/index.min.js',
        },
      ],
    }
  },
}
</script>

<style lang="scss" scoped></style>
