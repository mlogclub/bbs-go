<template>
  <div class="topics-main">
    <sticky-topics :node-id="-1" />
    <load-more
      v-if="topicsPage"
      v-slot="{ results }"
      :init-data="topicsPage"
      url="/api/topic/topics?recommend=true"
    >
      <topic-list :topics="results" :show-ad="true" />
    </load-more>
  </div>
</template>

<script>
export default {
  async asyncData({ $axios, store }) {
    store.commit('env/setCurrentNodeId', -1) // 设置当前所在node
    try {
      const [topicsPage] = await Promise.all([
        $axios.get('/api/topic/topics?recommend=true'),
      ])
      return { topicsPage }
    } catch (e) {
      console.error(e)
    }
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
