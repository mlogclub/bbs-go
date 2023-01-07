<template>
  <section class="main">
    <div class="container main-container left-main size-320">
      <div class="left-container">
        <nuxt-child :key="$route.path" />
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
      const [scoreRank, links, adverts] = await Promise.all([
        $axios.get('/api/user/score/rank'),
        $axios.get('/api/link/toplinks'),
        $axios.get('/api/advert/list'),
      ])
      return { scoreRank, links, adverts }
    } catch (e) {
      console.error(e)
    }
  },
}
</script>

<style lang="scss" scoped></style>
