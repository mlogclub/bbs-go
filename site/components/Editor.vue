<template>
  <div ref="editor" :style="{width: width}">
    <no-ssr>
      <mavon-editor
        ref="editor"
        class="md-editor"
        :value="value"
        :toolbars="toolbars"
        :ishljs="true"
        :box-shadow="false"
        :subfield="false"
        :scroll-style="true"
        code-style="atom-one-dark"
        @change="change"
        @imgAdd="imgAdd"
      />
    </no-ssr>
  </div>
</template>

<script>
import Vue from 'vue'
import MavonEditor from 'mavon-editor'
import 'mavon-editor/dist/css/index.css'
Vue.use(MavonEditor)

export default {
  props: {
    value: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      width: 'auto',
      toolbars: {
        bold: true, // 粗体
        italic: true, // 斜体
        header: true, // 标题
        underline: true, // 下划线
        strikethrough: true, // 中划线
        mark: true, // 标记
        // superscript: true, // 上角标
        // subscript: true, // 下角标
        quote: true, // 引用
        ol: true, // 有序列表
        ul: true, // 无序列表
        link: true, // 链接
        imagelink: true, // 图片链接
        code: true, // code
        table: true, // 表格
        fullscreen: true, // 全屏编辑
        // readmodel: true, // 沉浸式阅读
        // htmlcode: true, // 展示html源码
        // help: true, // 帮助
        // /* 1.3.5 */
        // undo: true, // 上一步
        // redo: true, // 下一步
        // trash: true, // 清空
        // save: true, // 保存（触发events中的save事件）
        /* 1.4.2 */
        navigation: true, // 导航目录
        /* 2.1.8 */
        alignleft: true, // 左对齐
        aligncenter: true, // 居中
        alignright: true, // 右对齐
        /* 2.2.1 */
        subfield: true, // 单双栏模式
        preview: true // 预览
      }
    }
  },
  mounted() {
    this.initEditorSize()
  },
  methods: {
    /**
     * 内容变更
     * @param value
     * @param render
     */
    change(value) {
      this.$emit('input', value)
    },

    /**
     * 图片上传
     */
    async imgAdd(pos, file) {
      try {
        const formData = new FormData()
        formData.append('image', file, file.name)
        const ret = await this.$axios.post('/api/upload', formData, {
          headers: { 'Content-Type': 'multipart/form-data' }
        })
        this.$refs.editor.$img2Url(pos, ret.url)
      } catch (e) {
        console.error(e)
      }
    },

    initEditorSize() {
      if (!process.client) {
        return
      }
      const wrapper = this.$refs.editor
      const parentElement = wrapper.parentElement
      if (!parentElement) {
        return
      }
      const me = this
      me.width = parentElement.clientWidth + 'px'
      window.addEventListener('resize', function () {
        me.width = parentElement.clientWidth + 'px'
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.md-editor {
  width: 100%;
  height: 100%;
  height: 450px;
}
</style>
