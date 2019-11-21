<template>
  <section class="main">
    <div class="container-wrapper">
      <div class="widget">
        <div class="widget-header">
          关于
        </div>
        <div class="widget-content">
          <h3>简介</h3>
          <p>
            本站是基于
            <a href="http://docs.bbs-go.com" target="_blank">bbs-go</a>
            搭建，旨在学习交流，欢迎各路大神提供意见和建议。
          </p>
          <p>
            交流群：<br /><img
              style="max-width: 200px;"
              src="https://file.mlog.club/mlog.club/bbs-go-qq.jpg"
            />
          </p>

          <h3>公众号</h3>
          <p>
            <img
              style="height: 150px; width: 150px;"
              src="https://file.mlog.club/mlog.club/qrcode.jpg"
            />
          </p>

          <h3>版权申明</h3>
          <p>
            <a href="https://mlog.club">https://mlog.club</a>
            部分文章来源于公众号，仅提供微信公众号所发布的文章链接索引，文章版权归属原作者或所属媒体所有。
          </p>

          <h3>豁免申明</h3>
          <p>
            <a href="https://mlog.club">https://mlog.club</a>
            是一个公益、共享网络平台，目的是为公众提供有价值的微信文章索引，不声明也不保证其内容的有效性、正确性与可靠性，
            所收录的文章索引及图片资源与言论信息仅代表原微信公众号意见，与<a
              href="https://mlog.club"
              >https://mlog.club</a
            >
            立场无关。
          </p>

          <h3>文章收录</h3>
          <p>
            如果您希望我们收录某公众号文章，请与我们联系并注明请求收录的公众号ID。如果您不想您的公众号文章被收录，请与我们联系，我们将停止收录您的帐号文章。
          </p>

          <h3>与我联系</h3>
          <p>
            如有任何想法、意见、问题可邮件联系 mlog1#qq.com ( @替换# )
          </p>

          <ins
            class="adsbygoogle"
            style="display:block"
            data-ad-client="ca-pub-5683711753850351"
            data-ad-slot="1742173616"
            data-ad-format="auto"
            data-full-width-responsive="true"
          />
          <script>
            ;(adsbygoogle = window.adsbygoogle || []).push({})
          </script>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  head() {
    return {
      title: this.$siteTitle('关于')
    }
  }
}
</script>

<style lang="scss" scoped>
.widget {
  font-size: 13px;
}
</style>
