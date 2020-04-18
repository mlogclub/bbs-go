<template>
  <section class="main">
    <div class="container main-container left-main">
      <div class="left-container">
        <div class="main-content no-padding">
          <post-tweets @created="tweetsCreated" />
        </div>

        <load-more
          ref="tweetsLoadMore"
          v-if="tweetsPage"
          v-slot="{ results }"
          :init-data="tweetsPage"
          url="/api/tweets/list"
        >
          <tweets-list :tweets="results" />
        </load-more>
      </div>
      <topic-side :score-rank="scoreRank" :links="links" />
    </div>
  </section>
</template>

<script>
import TopicSide from '~/components/TopicSide'
import PostTweets from '~/components/PostTweets'
import TweetsList from '~/components/TweetsList'
import LoadMore from '~/components/LoadMore'

export default {
  components: {
    TopicSide,
    PostTweets,
    TweetsList,
    LoadMore
  },
  async asyncData({ $axios, query }) {
    try {
      const [tweetsPage, scoreRank, links] = await Promise.all([
        $axios.get('/api/tweets/list'),
        $axios.get('/api/user/score/rank'),
        $axios.get('/api/link/toplinks')
      ])
      return { tweetsPage, scoreRank, links }
    } catch (e) {
      console.error(e)
    }
  },
  methods: {
    tweetsCreated(item) {
      this.$refs.tweetsLoadMore.unshiftResults(item)
    }
  },
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
