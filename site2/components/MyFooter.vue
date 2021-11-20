<template>
  <footer class="footer">
    <div class="container content">
      <div>
        <nuxt-link to="/about">关于</nuxt-link>
        <nuxt-link to="/tags">标签</nuxt-link>
        <nuxt-link to="/links">友链</nuxt-link>
      </div>
      <div>
        © 2021 Powered by
        <a href="https://docs.bbs-go.com" target="_blank" class="light"
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
