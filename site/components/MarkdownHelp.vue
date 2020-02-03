<template>
  <div class="widget markdown-help">
    <div class="widget-header">
      Markdown 语法参考
    </div>
    <div class="widget-content content">
      <ol>
        <li><tt>### 单行的标题</tt></li>
        <li><tt>**粗体**</tt></li>
        <li><tt>`console.log('行内代码')`</tt></li>
        <li><tt>```js\n code \n```</tt> 标记代码块</li>
        <li><tt>[内容](链接)</tt></li>
        <li><tt>![文字说明](图片链接)</tt></li>
      </ol>
      <span>
        <a href="https://mlog.club/article/5522" target="_blank"
          >Markdown 文档</a
        >
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
