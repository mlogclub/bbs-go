<template>
  <section class="main">
    <div class="container main-container left-main">
      <div class="left-container">
        <div class="main-content no-padding">
          <post-tweets @created="tweetsCreated" />
        </div>

        <load-more
          v-if="tweetsPage"
          ref="tweetsLoadMore"
          v-slot="{ results }"
          :init-data="tweetsPage"
          url="/api/tweet/list"
        >
          <tweets-list :tweets="results" />
        </load-more>
      </div>
      <div class="right-container">
        <check-in />
        <site-notice />
        <score-rank :score-rank="scoreRank" />
        <friend-links :links="links" />
      </div>
    </div>
  </section>
</template>

<script>
import CheckIn from '~/components/CheckIn'
import SiteNotice from '~/components/SiteNotice'
import ScoreRank from '~/components/ScoreRank'
import FriendLinks from '~/components/FriendLinks'
import PostTweets from '~/components/PostTweets'
import TweetsList from '~/components/TweetsList'
import LoadMore from '~/components/LoadMore'

export default {
  components: {
    CheckIn,
    SiteNotice,
    ScoreRank,
    FriendLinks,
    PostTweets,
    TweetsList,
    LoadMore,
  },
  async asyncData({ $axios, query }) {
    try {
      const [tweetsPage, scoreRank, links] = await Promise.all([
        $axios.get('/api/tweet/list'),
        $axios.get('/api/user/score/rank'),
        $axios.get('/api/link/toplinks'),
      ])
      return { tweetsPage, scoreRank, links }
    } catch (e) {
      console.error(e)
    }
  },
  methods: {
    tweetsCreated(item) {
      this.$refs.tweetsLoadMore.unshiftResults(item)
    },
  },
  head() {
    return {
      title: this.$siteTitle('动态'),
      meta: [
        {
          hid: 'description',
          name: 'description',
          content: this.$siteDescription(),
        },
        { hid: 'keywords', name: 'keywords', content: this.$siteKeywords() },
      ],
    }
  },
}
</script>

<style lang="scss" scoped></style>
