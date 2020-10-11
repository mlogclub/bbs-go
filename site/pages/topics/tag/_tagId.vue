<template>
  <section class="main">
    <div class="container main-container is-white left-main">
      <div class="left-container">
        <topics-nav :nodes="nodes" />
        <topic-list :topics="topicsPage.results" :show-ad="true" />
        <pagination
          :page="topicsPage.page"
          :url-prefix="'/topics/' + tag.tagId + '?p='"
        />
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
import Pagination from '~/components/Pagination'

export default {
  components: {
    CheckIn,
    SiteNotice,
    ScoreRank,
    FriendLinks,
    TopicsNav,
    TopicList,
    TweetsWidget,
    Pagination,
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
          page: query.p || 1,
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
