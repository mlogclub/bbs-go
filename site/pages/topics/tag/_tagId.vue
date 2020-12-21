<template>
  <section class="main">
    <div class="container main-container left-main">
      <div class="left-container">
        <div class="main-content no-padding">
          <topics-nav :nodes="nodes" />
          <load-more
            v-if="topicsPage"
            v-slot="{ results }"
            :init-data="topicsPage"
            :url="'/api/topic/tag/topics' + tag.tagId"
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
  async asyncData({ $axios, params, query }) {
    const [
      tag,
      nodes,
      topicsPage,
      scoreRank,
      links,
      newestTweets,
    ] = await Promise.all([
      $axios.get('/api/tag/' + params.tagId),
      $axios.get('/api/topic/nodes'),
      $axios.get('/api/topic/tag/topics', {
        params: {
          tagId: params.tagId,
        },
      }),
      $axios.get('/api/user/score/rank'),
      $axios.get('/api/link/toplinks'),
      $axios.get('/api/tweet/newest'),
    ])
    return {
      tag,
      nodes,
      topicsPage,
      scoreRank,
      links,
      newestTweets,
    }
  },
  head() {
    return {
      title: this.$siteTitle(this.tag.tagName + ' - 话题'),
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
