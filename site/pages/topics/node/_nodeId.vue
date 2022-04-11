<template>
  <div class="topics-main">
    <sticky-topics :node-id="node.nodeId" />
    <load-more
      v-if="topicsPage"
      v-slot="{ results }"
      :init-data="topicsPage"
      :url="'/api/topic/topics?nodeId=' + node.nodeId"
    >
      <topic-list :topics="results" :show-ad="true" />
    </load-more>
  </div>
</template>

<script>
export default {
  async asyncData({ $axios, params, store }) {
    const nodeId = parseInt(params.nodeId)
    store.commit('env/setCurrentNodeId', nodeId) // 设置当前所在node
    const [node, topicsPage] = await Promise.all([
      $axios.get('/api/topic/node?nodeId=' + nodeId),
      $axios.get('/api/topic/topics?nodeId=' + nodeId),
    ])
    return {
      node,
      topicsPage,
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
}
</script>

<style lang="scss" scoped></style>
