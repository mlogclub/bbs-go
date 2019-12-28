<template>
  <div id="vditor" :style="{ width: width, height: height }" class="vditor" />
</template>

<script>
export default {
  props: {
    value: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      isLoading: true,
      vditor: null,
      width: '100%',
      height: '400px'
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
    doInit() {
      if (!process.client) {
        return
      }
      const me = this
      const userToken = this.$cookies.get('userToken')
      const options = {
        width: me.width,
        height: me.height,
        cache: false,
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
        ],
        // placeholder: '请输入...',
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
          console.log(val)
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
      this.vditor = new window.Vditor('vditor', options)
      if (this.value) {
        this.vditor.setValue(this.value)
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
          src: '//cdn.jsdelivr.net/npm/vditor/dist/index.min.js',
          defer: true
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
