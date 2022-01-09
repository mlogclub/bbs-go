<template>
  <section class="main">
    <div class="container main-container left-main size-320">
      <div class="left-container">
        <nuxt-child :key="$route.path" />
      </div>
      <div class="right-container">
        <check-in />
        <div class="ad">
          <!-- 展示广告 -->
          <adsbygoogle ad-slot="1742173616" />
        </div>
        <site-notice />
        <score-rank :score-rank="scoreRank" />
        <friend-links :links="links" />
        <div class="ad">
          <!-- 展示广告 -->
          <adsbygoogle ad-slot="1742173616" />
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  async asyncData({ $axios, store }) {
    try {
      const [nodes, scoreRank, links] = await Promise.all([
        $axios.get('/api/topic/nodes'),
        $axios.get('/api/user/score/rank'),
        $axios.get('/api/link/toplinks'),
      ])
      return { nodes, scoreRank, links }
    } catch (e) {
      console.error(e)
    }
  },
}
</script>

<style lang="scss" scoped></style>
