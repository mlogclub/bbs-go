<template>
  <section class="main">
    <div class="container">
      <div class="main-body">
        <div class="widget">
          <div class="widget-header">友情链接</div>
          <div class="widget-content">
            <ul class="links">
              <li
                v-for="link in linksPage.results"
                :key="link.linkId"
                class="link"
              >
                <div class="link-logo">
                  <img v-if="link.logo" :src="link.logo" />
                  <img v-if="!link.logo" src="~/assets/images/net.png" />
                </div>
                <div class="link-content">
                  <a
                    :href="link.url"
                    :title="link.title"
                    class="link-title"
                    target="_blank"
                    >{{ link.title }}</a
                  >
                  <p class="link-summary">
                    {{ link.summary }}
                  </p>
                </div>
              </li>
            </ul>
            <pagination :page="linksPage.page" url-prefix="/links/" />
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import Pagination from '~/components/Pagination'
export default {
  components: {
    Pagination,
  },
  async asyncData({ $axios, params }) {
    const [linksPage] = await Promise.all([
      $axios.get('/api/link/links', {
        params: {
          page: params.page || 1,
        },
      }),
    ])
    return {
      linksPage,
    }
  },
  head() {
    return {
      title: this.$siteTitle('好博客'),
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
