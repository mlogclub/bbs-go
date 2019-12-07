<template>
  <div class="right-container">
    <post-btns :current-tag-id="currentTagId" />
    <weixin-gzh />

    <!-- 展示广告
    <ins
      class="adsbygoogle"
      style="display:block"
      data-ad-client="ca-pub-5683711753850351"
      data-ad-slot="1742173616"
      data-ad-format="auto"
      data-full-width-responsive="true"
    ></ins>
    <script>
      ;(adsbygoogle = window.adsbygoogle || []).push({})
    </script>
    -->

    <!-- 展示广告190x90 -->
    <ins
      class="adsbygoogle"
      style="display:inline-block;width:190px;height:90px"
      data-ad-client="ca-pub-5683711753850351"
      data-ad-slot="9345305153"
    />
    <script>
      ;(adsbygoogle = window.adsbygoogle || []).push({})
    </script>

    <!-- 展示广告190x480 -->
    <ins
      class="adsbygoogle"
      style="display:inline-block;width:190px;height:480px"
      data-ad-client="ca-pub-5683711753850351"
      data-ad-slot="3438372357"
    />
    <script>
      ;(adsbygoogle = window.adsbygoogle || []).push({})
    </script>
  </div>
</template>

<script>
import WeixinGzh from '~/components/WeixinGzh'
import PostBtns from '~/components/PostBtns'

export default {
  components: { WeixinGzh, PostBtns },
  props: {
    currentTagId: {
      type: Number,
      default: 0
    }
  }
}
</script>

<style lang="scss" scoped></style>
