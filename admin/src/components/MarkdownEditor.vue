<template>
  <mavon-editor ref="editor" :value="value" :toolbars="toolbars" :ishljs="true"
                @imgAdd="imgAdd" @fullScreen="fullScreen" @change="change"
                :boxShadow="false" :subfield="true" defaultOpen="edit"
                placeholder="内容"
                codeStyle="atom-one-dark"
                :style="{height: contentHeight, width: '100%'}"
  />
</template>

<script>
export default {
  name: 'MarkdownEditor',
  props: {
    value: {
      type: String,
      default: '',
    },
    initValue: {
      type: String,
      default: '',
    },
    height: {
      type: Number,
      default: 350,
    },
  },
  data() {
    return {
      uploading: false,
      uploadingText: 'loading...', // TODO 2019-02-28 没显示上传文案
      contentHeight: `${this.height}px`,
      fullScreenStatus: false,
      toolbars: {
        bold: true, // 粗体
        italic: true, // 斜体
        header: true, // 标题
        underline: true, // 下划线
        strikethrough: true, // 中划线
        mark: true, // 标记
        superscript: true, // 上角标
        subscript: true, // 下角标
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
        help: true, // 帮助
        /* 1.4.2 */
        navigation: true, // 导航目录
        /* 2.1.8 */
        alignleft: true, // 左对齐
        aligncenter: true, // 居中
        alignright: true, // 右对齐
        /* 2.2.1 */
        subfield: true, // 单双栏模式
        preview: true, // 预览
      },
    };
  },
  mounted() {
  },
  beforeDestroy() {
  },
  methods: {
    /**
       * 添加图片
       */
    async imgAdd(pos, file) {
      // let me = this
      // me.uploading = true
      // me.uploadingText = '正在上传 ' + file.name
    },
    /**
       * 内容变更
       * @param value
       * @param render
       */
    change(value) {
      this.$emit('input', value);
      if (value !== this.initValue) {
        this.$emit('change', value);
      }
    },
    /**
       * 切换全屏
       * @param status
       */
    fullScreen(status) {
      this.fullScreenStatus = status;
      if (status) {
        this.contentHeight = 'auto';
      } else {
        this.contentHeight = `${this.height}px`;
      }
    },
  },
};
</script>

<style scoped>

</style>
