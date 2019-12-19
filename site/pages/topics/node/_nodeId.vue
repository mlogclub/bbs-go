<template>
  <section class="main">
    <div class="container main-container is-white left-main">
      <div class="left-container">
        <topics-nav />
        <topic-list :topics="topicsPage.results" :show-ad="false" />
        <pagination
          :page="topicsPage.page"
          :url-prefix="'/topics/node/' + node.nodeId + '?p='"
        />
      </div>
      <topic-side />
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
  async asyncData({ $axios, params, query }) {
    const [node, user, topicsPage] = await Promise.all([
      $axios.get('/api/topic/node?nodeId=' + params.nodeId),
      $axios.get('/api/user/current'),
      $axios.get('/api/topic/node/topics', {
        params: {
          nodeId: params.nodeId,
          page: query.p || 1
        }
      })
    ])
    return {
      node,
      user,
      topicsPage
    }
  },
  head() {
    return {
      title: this.$siteTitle(this.node.name + ' - 话题'),
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
