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
        </div>
        <div>
          <span style="line-height: 20px; font-size: 18px;"
            >服务器资源赞助：</span
          >
          <a
            href="https://www.vultr.com/?ref=8358872"
            target="_blank"
            style="line-height: 20px;"
          >
            <img
              style="max-height: 20px;"
              src="https://file.mlog.club/images/misc/logo_onwhite.webp"
            />
          </a>
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
  margin: 0px 10px;
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
