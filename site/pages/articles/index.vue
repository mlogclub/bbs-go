<template>
  <section class="main">
    <div class="container main-container is-white left-main">
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
        <!-- 展示广告220*220 -->
        <adsbygoogle
          :ad-style="{
            display: 'inline-block',
            width: '220px',
            height: '220px'
          }"
          ad-slot="1361835285"
        />
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
