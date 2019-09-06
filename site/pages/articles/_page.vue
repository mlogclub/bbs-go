<template>
  <section class="main">
    <div class="container">
      <div class="columns">
        <div class="column is-9">
          <div class="main-body">
            <article-list :articles="articlesPage.results" :show-ad="true" />
            <pagination :page="articlesPage.page" url-prefix="/articles/" />
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
</template>

<script>
import ArticleList from '~/components/ArticleList'
import Pagination from '~/components/Pagination'
import ActiveUsers from '~/components/ActiveUsers'
import ActiveTags from '~/components/ActiveTags'

export default {
  components: {
    ArticleList,
    Pagination,
    ActiveUsers,
    ActiveTags
  },
  head() {
    return {
      title: this.$siteTitle('文章'),
      meta: [
        { hid: 'description', name: 'description', content: this.$siteDescription() },
        { hid: 'keywords', name: 'keywords', content: this.$siteKeywords() }
      ]
    }
  },
  async asyncData({ $axios, params }) {
    try {
      const [articlesPage, activeUsers, activeTags] = await Promise.all([
        $axios.get('/api/article/articles?page=' + (params.page || 1)),
        $axios.get('/api/user/active'), // 活跃用户
        $axios.get('/api/tag/active') // 活跃标签
      ])
      return {
        articlesPage: articlesPage,
        activeUsers: activeUsers,
        activeTags: activeTags
      }
    } catch (e) {
      console.error(e)
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
