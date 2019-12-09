<template>
  <div>
    <footer class="footer">
      <div class="container content">
        <div>
          <a href="/about">关于</a>
          <a href="/tags">标签</a>
        </div>
        <div>
          © 2019 Powered by
          <a href="http://docs.bbs-go.com" target="_blank" class="light"
            >BBS-GO</a
          >
          <a href="http://www.miibeian.gov.cn/" target="_blank" rel="nofollow"
            >鄂ICP备19023811号-3</a
          >
        </div>
      </div>
    </footer>
  </div>
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
