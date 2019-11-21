<template>
  <section class="main">
    <div class="container-wrapper main-container left-main">
      <div class="left-container">
        <load-more
          v-slot="{ results }"
          :init-data="articlesPage"
          url="/api/article/articles"
        >
          <article-list :articles="results" :show-ad="true" />
        </load-more>
      </div>
      <div class="right-container">
        <weixin-gzh />
        <div style="text-align: center;">
          <!-- 展示广告288x288 -->
          <ins
            class="adsbygoogle"
            style="display:inline-block;width:288px;height:288px"
            data-ad-client="ca-pub-5683711753850351"
            data-ad-slot="4922900917"
          />
          <script>
            ;(adsbygoogle = window.adsbygoogle || []).push({})
          </script>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import ArticleList from '~/components/ArticleList'
import WeixinGzh from '~/components/WeixinGzh'
import LoadMore from '~/components/LoadMore'

export default {
  components: {
    ArticleList,
    LoadMore,
    WeixinGzh
  },
  async asyncData({ $axios }) {
    try {
      const [articlesPage] = await Promise.all([
        $axios.get('/api/article/articles')
      ])
      return {
        articlesPage
      }
    } catch (e) {
      console.error(e)
    }
  },
  head() {
    return {
      title: this.$siteTitle('文章'),
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
