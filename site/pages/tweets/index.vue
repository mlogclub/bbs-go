<template>
  <section class="main">
    <div class="container main-container left-main">
      <div class="left-container">
        <div class="main-content no-padding">
          <post-tweets />
        </div>
      </div>
      <topic-side :score-rank="scoreRank" :links="links" />
    </div>
  </section>
</template>

<script>
import TopicSide from '~/components/TopicSide'
import PostTweets from '~/components/PostTweets'

export default {
  components: {
    TopicSide,
    PostTweets
  },
  async asyncData({ $axios, query }) {
    try {
      const [scoreRank, links] = await Promise.all([
        $axios.get('/api/user/score/rank'),
        $axios.get('/api/link/toplinks')
      ])
      return { scoreRank, links }
    } catch (e) {
      console.error(e)
    }
  },
  methods: {},
  head() {
    return {
      title: this.$siteTitle('动态'),
      meta: [
        {
          hid: 'description',
          name: 'description',
          content: this.$siteDescription()
        },
        { hid: 'keywords', name: 'keywords', content: this.$siteKeywords() }
      ]
    }
  }
}
</script>

<style lang="scss" scoped></style>
