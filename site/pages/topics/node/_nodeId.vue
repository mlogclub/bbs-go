<template>
  <section class="main">
    <top-notice />
    <div class="container main-container is-white left-main">
      <div class="left-container">
        <topics-nav :nodes="nodes" :current-node-id="node.nodeId" />
        <topic-list :topics="topicsPage.results" :show-ad="true" />
        <pagination
          :page="topicsPage.page"
          :url-prefix="'/topics/node/' + node.nodeId + '?p='"
        />
      </div>
      <topic-side :current-node-id="node.nodeId" />
    </div>
  </section>
</template>

<script>
import TopicSide from '~/components/TopicSide'
import TopicsNav from '~/components/TopicsNav'
import TopicList from '~/components/TopicList'
import Pagination from '~/components/Pagination'
import TopNotice from '~/components/TopNotice'

export default {
  components: {
    TopicSide,
    TopicsNav,
    TopicList,
    Pagination,
    TopNotice
  },
  async asyncData({ $axios, params, query }) {
    const [node, user, nodes, topicsPage] = await Promise.all([
      $axios.get('/api/topic/node?nodeId=' + params.nodeId),
      $axios.get('/api/user/current'),
      $axios.get('/api/topic/nodes'),
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
      nodes,
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
