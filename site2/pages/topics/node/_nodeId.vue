<template>
  <section class="main">
    <div class="container main-container left-main size-360">
      <div class="left-container">
        <div class="main-content no-padding no-bg topics-wrapper">
          <div class="topics-nav">
            <topics-nav :nodes="nodes" :current-node-id="node.nodeId" />
          </div>
          <div class="topics-main">
            <load-more
              v-if="topicsPage"
              v-slot="{ results }"
              :init-data="topicsPage"
              :url="'/api/topic/topics?nodeId=' + node.nodeId"
            >
              <topic-list :topics="results" :show-ad="true" />
            </load-more>
          </div>
        </div>
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
export default {
  async asyncData({ $axios, params, store }) {
    const nodeId = parseInt(params.nodeId)
    store.commit('env/setCurrentNodeId', nodeId) // 设置当前所在node
    const [node, nodes, topicsPage, scoreRank, links] = await Promise.all([
      $axios.get('/api/topic/node?nodeId=' + nodeId),
      $axios.get('/api/topic/nodes'),
      $axios.get('/api/topic/topics?nodeId=' + nodeId),
      $axios.get('/api/user/score/rank'),
      $axios.get('/api/link/toplinks'),
    ])
    return {
      node,
      nodes,
      topicsPage,
      scoreRank,
      links,
    }
  },
  head() {
    return {
      title: this.$siteTitle(this.node.name + ' - 话题'),
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
}
</script>

<style lang="scss" scoped></style>
