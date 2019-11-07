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
            <!-- todo -->
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import ArticleList from '~/components/ArticleList'
import Pagination from '~/components/Pagination'

export default {
  components: {
    ArticleList,
    Pagination
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
      const [articlesPage] = await Promise.all([
        $axios.get('/api/article/articles?page=' + (params.page || 1))
      ])
      return {
        articlesPage: articlesPage
      }
    } catch (e) {
      console.error(e)
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
