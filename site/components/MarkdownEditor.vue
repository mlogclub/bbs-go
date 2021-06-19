<template>
  <div
    :id="editorId"
    :style="{ width: width, height: height }"
    class="vditor"
  />
</template>

<script>
import Vditor from 'vditor'
import 'vditor/src/assets/scss/index.scss'

// import CommonHelper from '~/common/CommonHelper'
// const vditorCss = '//cdn.jsdelivr.net/npm/vditor@3.5.4/dist/index.css'
// const vditorScript = '//cdn.jsdelivr.net/npm/vditor@3.5.4/dist/index.min.js'

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
  // head() {
  //   return {
  //     link: [
  //       {
  //         rel: 'stylesheet',
  //         href: vditorCss,
  //       },
  //     ],
  //     script: [
  //       {
  //         type: 'text/javascript',
  //         src: vditorScript,
  //         callback: () => {
  //           this.createEditor()
  //         },
  //       },
  //     ],
  //   }
  // },
  computed: {
    isMobile() {
      return this.$store.state.env.isMobile
    },
    toolbars() {
      if (this.isMobile) {
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
  },
  mounted() {
    // this.init()
    this.createEditor()
  },
  methods: {
    // init() {
    //   const me = this
    //   if (window.Vditor) {
    //     me.createEditor()
    //   } else {
    //     CommonHelper.addScript(vditorScript, function () {
    //       me.createEditor()
    //     })
    //   }
    // },
    createEditor() {
      console.log('初始化编辑器...')
      if (process.client) {
        const me = this
        me.vditor = new Vditor(
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
        toolbar: me.toolbars,
        mode: 'ir',
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
          mode: 'both',
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
    /**
     * 清空编辑器内容
     */
    clear() {
      if (this.vditor) {
        this.$emit('input', '')
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
}
</script>

<style lang="scss" scoped></style>
