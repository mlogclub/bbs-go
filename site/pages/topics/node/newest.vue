<template>
  <div class="topics-main">
    <sticky-topics :node-id="0" />
    <load-more
      v-if="topicsPage"
      v-slot="{ results }"
      :init-data="topicsPage"
      :url="url"
    >
      <topic-list :topics="results" />
    </load-more>
  </div>
</template>

<script>
export default {
  async asyncData({ $axios, store }) {
    const nodeId = 0
    const url = '/api/topic/topics?nodeId=' + nodeId
    store.commit('env/setCurrentNodeId', nodeId) // 设置当前所在node
    try {
      const topicsPage = await $axios.get(url)
      return { topicsPage, url }
    } catch (e) {
      console.error(e)
    }
  },
  head() {
    return {
      title: this.$siteTitle('最新'),
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
