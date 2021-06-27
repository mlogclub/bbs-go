<template>
  <section class="main">
    <div class="container main-container left-main size-360">
      <div class="left-container">
        <div class="main-content no-padding no-bg topics-wrapper">
          <div class="topics-nav"><topics-nav :nodes="nodes" /></div>
          <div class="topics-main">
            <load-more
              v-if="topicsPage"
              v-slot="{ results }"
              :init-data="topicsPage"
              :url="'/api/topic/tag/topics' + tag.tagId"
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
    store.commit('env/setCurrentNodeId', 0) // 设置当前所在node
    const [tag, nodes, topicsPage, scoreRank, links] = await Promise.all([
      $axios.get('/api/tag/' + params.tagId),
      $axios.get('/api/topic/nodes'),
      $axios.get('/api/topic/tag/topics', {
        params: {
          tagId: params.tagId,
        },
      }),
      $axios.get('/api/user/score/rank'),
      $axios.get('/api/link/toplinks'),
    ])
    return {
      tag,
      nodes,
      topicsPage,
      scoreRank,
      links,
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
