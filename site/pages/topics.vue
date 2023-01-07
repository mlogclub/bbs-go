<template>
  <section class="main">
    <div class="container main-container left-main size-320">
      <div class="left-container">
        <div class="main-content no-padding no-bg topics-wrapper">
          <div class="topics-nav"><topics-nav :nodes="nodes" /></div>
          <nuxt-child :key="$route.path" />
        </div>
      </div>
      <div class="right-container">
        <site-notice />
        <advert :adverts="adverts" />
        <check-in />
        <score-rank :score-rank="scoreRank" />
        <friend-links :links="links" />
      </div>
    </div>
  </section>
</template>

<script>
export default {
  async asyncData({ $axios, store }) {
    try {
      store.commit('env/setCurrentNodeId', 0) // 设置当前所在node
      const [nodes, scoreRank, links, adverts] = await Promise.all([
        $axios.get('/api/topic/node_navs'),
        $axios.get('/api/user/score/rank'),
        $axios.get('/api/link/toplinks'),
        $axios.get('/api/advert/list'),
      ])
      return { nodes, scoreRank, links, adverts }
    } catch (e) {
      console.error(e)
    }
  },
}
</script>

<style lang="scss" scoped></style>
