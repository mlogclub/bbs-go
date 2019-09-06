<template>
  <div>
    <section class="main">
      <div class="container">
        <div class="columns">
          <div class="column is-9">
            <div class="main-body">
              <div v-if="topics && topics.length" class="widget">
                <div class="header">
                  <span>最新主题</span>
                  <div class="slot">
                    <a href="/topics">查看更多>></a>
                  </div>
                </div>
                <div class="content">
                  <topic-list :topics="topics" />
                </div>
              </div>

              <div v-if="articles && articles.length" class="widget">
                <div class="header">
                  <span>最新文章</span>
                  <div class="slot">
                    <a href="/articles">查看更多>></a>
                  </div>
                </div>
                <div class="content">
                  <article-list :articles="articles" />
                </div>
              </div>
            </div>
          </div>
          <div class="column is-3">
            <div class="main-aside">
              <div style="text-align: center;">
                <!-- 展示广告288x288 -->
                <ins
                  class="adsbygoogle"
                  style="display:inline-block;width:288px;height:288px"
                  data-ad-client="ca-pub-5683711753850351"
                  data-ad-slot="4922900917"
                />
                <script>
                  (adsbygoogle = window.adsbygoogle || []).push({});
                </script>
              </div>
              <active-users :active-users="activeUsers" />
              <active-tags :active-tags="activeTags" />
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import TopicList from '~/components/TopicList'
import ArticleList from '~/components/ArticleList'
import ActiveUsers from '~/components/ActiveUsers'
import ActiveTags from '~/components/ActiveTags'
export default {
  components: {
    TopicList,
    ArticleList,
    ActiveUsers,
    ActiveTags
  },
  head() {
    return {
      meta: [
        { hid: 'description', name: 'description', content: this.$siteDescription() },
        { hid: 'keywords', name: 'keywords', content: this.$siteKeywords() }
      ]
    }
  },
  async asyncData({ $axios, params, store }) {
    const [topics, articles, activeUsers, activeTags] = await Promise.all([
      $axios.get('/api/topic/recent'), // 最新帖子
      $axios.get('/api/article/recent'), // 最新文章
      $axios.get('/api/user/active'), // 活跃用户
      $axios.get('/api/tag/active') // 活跃标签
    ])
    return {
      configs: store.state.config.configs,
      topics: topics,
      articles: articles,
      activeUsers: activeUsers,
      activeTags: activeTags
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
