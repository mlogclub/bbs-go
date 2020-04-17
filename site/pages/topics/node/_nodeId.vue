<template>
  <section class="main">
    <div class="container main-container left-main">
      <div class="left-container">
        <div class="main-content">
          <topics-nav :nodes="nodes" :current-node-id="node.nodeId" />
          <topic-list :topics="topicsPage.results" :show-ad="true" />
          <pagination
            :page="topicsPage.page"
            :url-prefix="'/topics/node/' + node.nodeId + '?p='"
          />
        </div>
      </div>
      <topic-side
        :current-node-id="node.nodeId"
        :score-rank="scoreRank"
        :links="links"
      />
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
    const [node, nodes, topicsPage, scoreRank, links] = await Promise.all([
      $axios.get('/api/topic/node?nodeId=' + params.nodeId),
      $axios.get('/api/topic/nodes'),
      $axios.get('/api/topic/node/topics', {
        params: {
          nodeId: params.nodeId,
          page: query.p || 1
        }
      }),
      $axios.get('/api/user/score/rank'),
      $axios.get('/api/link/toplinks')
    ])
    return {
      node,
      nodes,
      topicsPage,
      scoreRank,
      links
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
