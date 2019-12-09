<template>
  <section class="main">
    <div class="container">
      <div class="widget">
        <div class="widget-header">最新话题</div>
        <div class="widget-content">
          <div class="columns">
            <div class="column">
              <topic-list :topics="topics1" :show-ad="false" />
            </div>
            <div class="column">
              <topic-list :topics="topics2" :show-ad="false" />
            </div>
          </div>
        </div>
      </div>

      <div class="widget">
        <div class="widget-header">最新文章</div>
        <div class="widget-content">
          <article-list :articles="articles" :show-ad="false" />
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import TopicList from '~/components/TopicList'
import ArticleList from '~/components/ArticleList'

export default {
  components: { TopicList, ArticleList },
  async asyncData({ $axios, params }) {
    try {
      const [topics, articles] = await Promise.all([
        $axios.get('/api/topic/newest'),
        $axios.get('/api/article/newest')
      ])
      return {
        topics1: topics.slice(0, 5),
        topics2: topics.slice(5, 10),
        articles
      }
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
