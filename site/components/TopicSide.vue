<template>
  <div class="right-container">
    <post-btns :current-node-id="currentNodeId" />
    <site-notice />

    <div class="ad">
      <!-- 展示广告 -->
      <adsbygoogle ad-slot="1742173616" />
    </div>
  </div>
</template>

<script>
import PostBtns from '~/components/PostBtns'
import SiteNotice from '~/components/SiteNotice'

export default {
  components: { PostBtns, SiteNotice },
  props: {
    currentNodeId: {
      type: Number,
      default: 0
    }
  }
}
</script>

<style lang="scss" scoped></style>
