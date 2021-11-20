<template>
  <section class="main">
    <div class="container main-container left-main size-360">
      <div class="left-container">
        <div class="main-content no-padding no-bg topics-wrapper">
          <div class="topics-nav">
            <topics-nav :nodes="nodes" />
          </div>
          <div class="topics-main">
            <load-more
              v-if="topicsPage"
              v-slot="{ results }"
              :init-data="topicsPage"
              url="/api/topic/topics?recommend=true"
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
  async asyncData({ $axios, store }) {
    store.commit('env/setCurrentNodeId', -1) // 设置当前所在node
    try {
      const [nodes, topicsPage, scoreRank, links] = await Promise.all([
        $axios.get('/api/topic/nodes'),
        $axios.get('/api/topic/topics?recommend=true'),
        $axios.get('/api/user/score/rank'),
        $axios.get('/api/link/toplinks'),
      ])
      return { nodes, topicsPage, scoreRank, links }
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
