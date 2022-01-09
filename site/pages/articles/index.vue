<template>
  <load-more
    v-slot="{ results }"
    :init-data="articlesPage"
    url="/api/article/articles"
  >
    <article-list :articles="results" :show-ad="true" />
  </load-more>
</template>

<script>
export default {
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
