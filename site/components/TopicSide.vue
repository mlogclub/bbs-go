<template>
  <div class="right-container">
    <post-btns :current-tag-id="currentTagId" />
    <weixin-gzh />

    <!-- 展示广告220*220 -->
    <adsbygoogle
      :ad-style="{
        display: 'inline-block',
        width: '220px',
        height: '220px'
      }"
      ad-slot="1361835285"
    />

    <!-- 展示广告220*220
    <adsbygoogle
      :ad-style="{
        display: 'inline-block',
        width: '220px',
        height: '220px'
      }"
      ad-slot="1361835285"
    />
    -->
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
