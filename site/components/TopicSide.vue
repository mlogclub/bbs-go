<template>
  <div class="right-container">
    <post-btns :current-node-id="currentNodeId" />
    <weixin-gzh />

    <div class="ad">
      <!-- 展示广告 -->
      <adsbygoogle ad-slot="1742173616" />
    </div>
  </div>
</template>

<script>
import WeixinGzh from '~/components/WeixinGzh'
import PostBtns from '~/components/PostBtns'

export default {
  components: { WeixinGzh, PostBtns },
  props: {
    currentNodeId: {
      type: Number,
      default: 0
    }
  }
}
</script>

<style lang="scss" scoped></style>
