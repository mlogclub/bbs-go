<template>
  <section class="main">
    <div class="container-wrapper main-container right-main">
      <bbs-left :current-tag-id="tag.tagId" />
      <div class="right-container">
        <topic-list :topics="topicsPage.results" :show-ad="false" />
        <pagination
          :page="topicsPage.page"
          :url-prefix="'/topics/tag/' + tag.tagId + '/'"
        />
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
    const [tag, user, topicsPage] = await Promise.all([
      $axios.get('/api/tag/' + params.tagId),
      $axios.get('/api/user/current'),
      $axios.get('/api/topic/tag/topics', {
        params: {
          tagId: params.tagId,
          page: params.page || 1
        }
      })
    ])
    return {
      tag,
      user,
      topicsPage
    }
  },
  head() {
    return {
      title: this.$siteTitle(this.tag.tagName + ' - 话题'),
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
