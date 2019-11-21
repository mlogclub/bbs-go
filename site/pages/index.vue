<template>
  <section class="main">
    <div class="container-wrapper main-container right-main">
      <bbs-left />
      <div class="right-container">
        <topic-list :topics="topicsPage.results" :show-ad="false" />
        <pagination :page="topicsPage.page" url-prefix="/topics/" />
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
    BbsLeft,
    TopicList,
    Pagination
  },
  async asyncData({ $axios, params }) {
    try {
      const [user, topicsPage] = await Promise.all([
        $axios.get('/api/user/current'),
        $axios.get('/api/topic/topics')
      ])
      return { user, topicsPage }
    } catch (e) {
      console.error(e)
    }
  },
  head() {
    return {
      meta: [
        {
          hid: 'description',
          name: 'description',
          content: this.$siteDescription()
        },
        { hid: 'keywords', name: 'keywords', content: this.$siteKeywords() }
      ]
    }
  }
}
</script>

<style lang="scss" scoped></style>
