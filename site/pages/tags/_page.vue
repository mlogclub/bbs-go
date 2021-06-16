<template>
  <section class="main">
    <div class="container">
      <div class="columns">
        <div class="column is-9">
          <div class="main-body">
            <div class="widget">
              <div class="widget-header">
                <span>标签</span>
              </div>
              <div class="widget-content">
                <div class="tags are-medium">
                  <span
                    v-for="tag in tagsPage.results"
                    :key="tag.tagId"
                    class="tag is-normal"
                  >
                    <nuxt-link
                      :to="'/articles/' + tag.tagId"
                      :title="tag.tagName"
                      target="_blank"
                      >{{ tag.tagName }}</nuxt-link
                    >
                  </span>
                </div>
              </div>
              <pagination :page="tagsPage.page" url-prefix="/tags/" />
            </div>
          </div>
        </div>
        <div class="column is-3">
          <div class="main-aside">
            <site-notice />

            <div class="ad">
              <!-- 展示广告 -->
              <adsbygoogle ad-slot="1742173616" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  async asyncData({ $axios, params }) {
    const [tagsPage] = await Promise.all([
      $axios.get('/api/tag/tags', {
        params: {
          page: params.page,
        },
      }),
    ])
    return {
      tagsPage,
    }
  },
  head() {
    return {
      title: this.$siteTitle('标签'),
    }
  },
}
</script>

<style lang="scss" scoped></style>
