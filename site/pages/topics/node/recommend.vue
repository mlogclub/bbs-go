<template>
  <section class="main">
    <div class="container main-container left-main">
      <div class="left-container">
        <div class="main-content">
          <topics-nav :nodes="nodes" :current-node-id="-1" />
          <load-more
            v-if="topicsPage"
            v-slot="{ results }"
            :init-data="topicsPage"
            url="/api/topic/topics?recommend=true"
          >
            <topic-list :topics="results" :show-ad="true" />
          </load-more>
        </div>
      </div>
      <div class="right-container">
        <check-in />
        <site-notice />
        <tweets-widget :tweets="newestTweets" />
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
import TopicsNav from '~/components/TopicsNav'
import TopicList from '~/components/TopicList'
import TweetsWidget from '~/components/TweetsWidget'
import LoadMore from '~/components/LoadMore'

export default {
  components: {
    CheckIn,
    SiteNotice,
    ScoreRank,
    FriendLinks,
    TopicsNav,
    TopicList,
    TweetsWidget,
    LoadMore,
  },
  async asyncData({ $axios, query }) {
    try {
      const [
        nodes,
        topicsPage,
        scoreRank,
        links,
        newestTweets,
      ] = await Promise.all([
        $axios.get('/api/topic/nodes'),
        $axios.get('/api/topic/topics?recommend=true'),
        $axios.get('/api/user/score/rank'),
        $axios.get('/api/link/toplinks'),
        $axios.get('/api/tweet/newest'),
      ])
      return { nodes, topicsPage, scoreRank, links, newestTweets }
    } catch (e) {
      console.error(e)
    }
  },
  methods: {
    twitterCreated(data) {
      if (this.topicsPage) {
        if (this.topicsPage.results) {
          this.topicsPage.results.unshift(data)
        } else {
          this.topicsPage.results = [data]
        }
      }
    },
  },
  head() {
    return {
      title: this.$siteTitle('热门话题'),
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
