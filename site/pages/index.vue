<template>
  <section class="main">
    <div class="container main-container is-white left-main">
      <div class="left-container">
        <topics-nav :nodes="nodes" />
        <topic-list :topics="topicsPage.results" :show-ad="true" />
        <pagination :page="topicsPage.page" url-prefix="/topics?p=" />
      </div>
      <topic-side :score-rank="scoreRank" :links="links" />
    </div>
  </section>
</template>

<script>
import TopicSide from '~/components/TopicSide'
import TopicsNav from '~/components/TopicsNav'
import TopicList from '~/components/TopicList'
import Pagination from '~/components/Pagination'

export default {
  components: {
    TopicSide,
    TopicsNav,
    TopicList,
    Pagination
  },
  async asyncData({ $axios, params }) {
    try {
      const [nodes, topicsPage, scoreRank, links] = await Promise.all([
        $axios.get('/api/topic/nodes'),
        $axios.get('/api/topic/topics'),
        $axios.get('/api/user/score/rank'),
        $axios.get('/api/link/toplinks')
      ])
      return { nodes, topicsPage, scoreRank, links }
    } catch (e) {
      console.error(e)
    }
  },
  methods: {},
  head() {
    return {
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
