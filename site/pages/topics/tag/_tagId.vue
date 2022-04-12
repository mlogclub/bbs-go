<template>
  <div class="topics-main">
    <sticky-topics :node-id="0" />
    <load-more
      v-if="topicsPage"
      v-slot="{ results }"
      :init-data="topicsPage"
      :url="'/api/topic/tag/topics' + tag.tagId"
    >
      <topic-list :topics="results" :show-ad="true" />
    </load-more>
  </div>
</template>

<script>
export default {
  async asyncData({ $axios, params, store }) {
    store.commit('env/setCurrentNodeId', 0) // 设置当前所在node
    const [tag, topicsPage] = await Promise.all([
      $axios.get('/api/tag/' + params.tagId),
      $axios.get('/api/topic/tag/topics', {
        params: {
          tagId: params.tagId,
        },
      }),
    ])
    return {
      tag,
      topicsPage,
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
