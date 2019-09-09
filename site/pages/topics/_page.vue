<template>
  <section class="main">
    <div class="container">
      <div class="right-main-container">
        <bbs-left />
        <div class="m-right">
          <topic-list :topics="topicsPage.results" :show-ad="false" />
          <pagination :page="topicsPage.page" url-prefix="/topics/" />
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import BbsLeft from '~/components/BbsLeft'
import TopicList from '~/components/TopicList'
import Pagination from '~/components/Pagination'

export default {
  components: {
    BbsLeft, TopicList, Pagination
  },
  head() {
    return {
      title: this.$siteTitle('话题'),
      meta: [
        { hid: 'description', name: 'description', content: this.$siteDescription() },
        { hid: 'keywords', name: 'keywords', content: this.$siteKeywords() }
      ]
    }
  },
  async asyncData({ $axios, params }) {
    try {
      const [user, topicsPage] = await Promise.all([
        $axios.get('/api/user/current'),
        $axios.get('/api/topic/topics?page=' + (params.page || 1))
      ])
      return { user: user, topicsPage: topicsPage }
    } catch (e) {
      console.error(e)
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
