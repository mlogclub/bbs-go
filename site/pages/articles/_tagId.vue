<template>
  <section class="main">
    <div class="container main-container is-white left-main">
      <div class="left-container">
        <load-more
          v-slot="{ results }"
          :init-data="articlesPage"
          :params="{ tagId: tag.tagId }"
          url="/api/article/tag/articles"
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
import LoadMore from '~/components/LoadMore'
import WeixinGzh from '~/components/WeixinGzh'

export default {
  components: { ArticleList, LoadMore, WeixinGzh },
  async asyncData({ $axios, params }) {
    try {
      const [tag, articlesPage] = await Promise.all([
        $axios.get('/api/tag/' + params.tagId),
        $axios.get('/api/article/tag/articles', {
          params: {
            tagId: params.tagId
          }
        })
      ])
      return {
        tag,
        articlesPage
      }
    } catch (e) {
      console.error(e)
    }
  },
  head() {
    return {
      title: this.$siteTitle(this.tag.tagName + ' - 文章'),
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
