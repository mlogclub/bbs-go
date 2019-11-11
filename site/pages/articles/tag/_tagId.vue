<template>
  <section class="main">
    <div class="container">
      <div class="left-main-container">
        <div class="m-left">
          <load-more
            v-slot="{results}"
            :init-data="articlesPage"
            :params="{tagId:tag.tagId}"
            url="/api/article/tag/articles"
          >
            <article-list :articles="results" :show-ad="true" />
          </load-more>
        </div>
        <div class="m-right">
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
              (adsbygoogle = window.adsbygoogle || []).push({});
            </script>
          </div>
        </div>
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
  },
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
        tag: tag,
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
