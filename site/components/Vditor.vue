<template>
  <div>
    <div id="vditor" class="vditor" />
  </div>
</template>

<script>
import Vditor from 'vditor'
import 'vditor/src/assets/scss/dark.scss'

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
      vditor: null
    }
  },
  mounted() {
    this.initVditor()
    this.$nextTick(async () => {
      await this.vditor.getHTML(true)
      this.isLoading = false
    })
  },
  methods: {
    initVditor() {
      const me = this
      const options = {
        cache: false,
        toolbar: ['emoji', 'headings', 'bold', 'italic', 'strike', '|', 'line', 'quote', 'list', 'ordered-list', 'check', 'code',
          'inline-code', 'undo', 'redo', 'upload', 'link', 'table', 'preview', 'fullscreen'],
        // placeholder: '请输入...',
        width: '100%',
        height: 400,
        tab: '\t',
        counter: '999999',
        input: function (val) {
          me.$emit('input', val)
        },
        upload: {
          accept: 'image/*',
          token: 'test',
          url: '{{baseUrl}}/upload/editor',
          linkToImgUrl: '{{baseUrl}}/upload/fetch',
          filename(name) {
            return name.replace(/\?|\\|\/|:|\||<|>|\*|\[|\]|\s+/g, '-')
          }
        }
        // preview: {
        //   delay: 100,
        //   show: true
        // }
      }
      this.vditor = new Vditor('vditor', options)
      // this.vditor.focus()
      if (this.value) {
        this.vditor.setValue(this.value)
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.vditor {
  border: 1px solid #eee;
}
</style>
