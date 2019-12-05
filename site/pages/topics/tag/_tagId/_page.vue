<template>
  <section class="main">
    <div class="container-wrapper main-container left-main">
      <div class="left-container">
        <topics-nav :current-tag-id="tag.tagId" />
        <topic-list :topics="topicsPage.results" :show-ad="true" />
        <pagination
          :page="topicsPage.page"
          :url-prefix="'/topics/tag/' + tag.tagId + '/'"
        />
      </div>
      <topic-side :current-tag-id="tag.tagId" />
    </div>
  </section>
</template>

<script>
import TopicSide from '~/components/TopicSide'
import TopicsNav from '~/components/TopicsNav'
import TopicList from '~/components/TopicList'
import Pagination from '~/components/Pagination'

export default {
  components: {
    TopicSide,
    TopicsNav,
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
