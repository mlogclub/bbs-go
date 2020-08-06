<template>
  <section class="main">
    <div class="container main-container is-white left-main size-320">
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
        <site-notice />

        <div class="ad">
          <!-- 展示广告 -->
          <adsbygoogle ad-slot="1742173616" />
        </div>

        <div class="ad">
          <!-- 展示广告 -->
          <adsbygoogle ad-slot="1742173616" />
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import ArticleList from '~/components/ArticleList'
import SiteNotice from '~/components/SiteNotice'
import LoadMore from '~/components/LoadMore'

export default {
  components: {
    ArticleList,
    SiteNotice,
    LoadMore,
  },
  async asyncData({ $axios }) {
    try {
      const [articlesPage] = await Promise.all([
        $axios.get('/api/article/articles'),
      ])
      return {
        articlesPage,
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
          content: this.$siteDescription(),
        },
        { hid: 'keywords', name: 'keywords', content: this.$siteKeywords() },
      ],
    }
  },
}
</script>

<style lang="scss" scoped></style>
