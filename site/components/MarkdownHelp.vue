<template>
  <div class="widget markdown-help">
    <div class="widget-header">
      Markdown 语法参考
    </div>
    <div class="widget-content content">
      <ol>
        <li>
          <span style="font-weight: 700; color: red;"
            >将图片复制或拖入编辑器可上传</span
          >
        </li>
        <li>## 单行的标题</li>
        <li>**粗体**</li>
        <li>`console.log('行内代码')`</li>
        <li>```js\n code \n``` 标记代码块</li>
        <li>[内容](链接)</li>
        <li>![文字说明](图片链接)</li>
      </ol>
      <span>
        <a href="https://mlog.club/topic/1598" target="_blank">
          <strong>Markdown语法说明(简体中文版)</strong>
        </a>
      </span>
    </div>
  </div>
</template>

<script>
export default {}
</script>

<style lang="scss" scoped>
.markdown-help {
  .content {
    font-size: 12px !important;
  }
}
</style>
