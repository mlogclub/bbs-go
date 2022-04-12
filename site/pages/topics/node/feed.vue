<template>
  <div class="topics-main">
    <sticky-topics :node-id="-2" />
    <load-more
      v-if="topicsPage"
      v-slot="{ results }"
      :init-data="topicsPage"
      url="/api/feed/topics"
    >
      <topic-list :topics="results" :show-ad="true" />
    </load-more>
  </div>
</template>

<script>
export default {
  async asyncData({ $axios, store }) {
    store.commit('env/setCurrentNodeId', -2) // 设置当前所在node
    let topicsPage
    try {
      // TODO 这里没登陆，或者没有数据的时候页面上要显示相应的引导内容
      if (store.state.user.current) {
        topicsPage = await $axios.get('/api/feed/topics')
      }
    } catch (e) {
      console.log(e.message || e)
    }
    return { topicsPage }
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
}
</script>

<style lang="scss" scoped></style>
