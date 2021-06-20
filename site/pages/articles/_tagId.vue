<template>
  <section class="main">
    <div class="container main-container is-white left-main size-320">
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
export default {
  async asyncData({ $axios, params }) {
    try {
      const [tag, articlesPage] = await Promise.all([
        $axios.get('/api/tag/' + params.tagId),
        $axios.get('/api/article/tag/articles', {
          params: {
            tagId: params.tagId,
          },
        }),
      ])
      return {
        tag,
        articlesPage,
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
          content: this.$siteDescription(),
        },
        { hid: 'keywords', name: 'keywords', content: this.$siteKeywords() },
      ],
    }
  },
}
</script>

<style lang="scss" scoped></style>
