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
              url="/api/feed/topics"
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
    store.commit('env/setCurrentNodeId', -2) // 设置当前所在node
    let topicsPage, nodes, scoreRank, links
    try {
      // TODO 这里没登陆，或者没有数据的时候页面上要显示相应的引导内容
      if (store.state.user.current) {
        topicsPage = await $axios.get('/api/feed/topics')
      }
    } catch (e) {
      console.log(e.message || e)
    }
    try {
      ;[nodes, scoreRank, links] = await Promise.all([
        $axios.get('/api/topic/nodes'),
        $axios.get('/api/user/score/rank'),
        $axios.get('/api/link/toplinks'),
      ])
    } catch (e) {
      console.error(e)
    }
    return { nodes, topicsPage, scoreRank, links }
  },
  head() {
    return {
      title: this.$siteTitle('关注'),
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
  computed: {
    user() {
      return this.$store.state.user.current
    },
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
