<template>
  <section class="main">
    <div class="container main-container left-main size-320">
      <div class="left-container">
        <div class="main-content">
          <search-topics-nav :nodes="nodes" />
          <div v-if="searchPage && searchPage.results">
            <search-topic-list :search-page="searchPage" />
            <pagination
              :page="searchPage.page"
              :url-prefix="'/search?q=' + keyword + '&p='"
            />
          </div>
          <div v-else class="notification is-info empty-results">
            {{ searchLoading ? '加载中...' : '未搜索到内容' }}
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
  async asyncData({ $axios, query, store }) {
    const keyword = query.q || ''
    const [nodes, scoreRank, links] = await Promise.all([
      $axios.get('/api/topic/nodes'),
      $axios.get('/api/user/score/rank'),
      $axios.get('/api/link/toplinks'),
    ])
    store.dispatch('search/initParams', {
      keyword: query.q || '',
      page: query.p || 1,
    })
    return { keyword, nodes, scoreRank, links }
  },
  computed: {
    searchPage() {
      return this.$store.state.search.searchPage
    },
    searchLoading() {
      return this.$store.state.search.loading
    },
  },
  mounted() {
    this.searchTopic()
  },
  methods: {
    async searchTopic() {
      await this.$store.dispatch('search/searchTopic')
    },
  },
}
</script>

<style lang="scss" scoped>
.search-input {
  background-color: var(--bg-color);
  padding: 10px;
  text-align: center;
}
.empty-results {
  margin-top: 10px;
}
</style>
