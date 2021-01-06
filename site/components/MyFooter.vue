<template>
  <footer class="footer">
    <div class="container content">
      <div>
        <a href="/about">关于</a>
        <a href="/tags">标签</a>
        <a href="/links">友链</a>
      </div>
      <div>
        © 2019 Powered by
        <a href="http://docs.bbs-go.com" target="_blank" class="light"
          >BBS-GO</a
        >
      </div>
    </div>
  </footer>
</template>

<script>
export default {}
</script>

<style lang="scss" scoped>
.footer {
  font-size: 14px;
  color: #666;
  background: none;
  text-align: left;
  margin: 0 10px;
  a {
    color: #666;
    text-decoration: none;
  }

  .light {
    color: #eb5424;
    font-weight: bold;
  }
}
</style>
