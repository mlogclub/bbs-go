<template>
  <section class="main">
    <div class="container main-container is-white left-main">
      <div class="left-container">
        <project-list :projects="projectPage.results" />
        <pagination :page="projectPage.page" url-prefix="/projects/" />
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
import ProjectList from '~/components/ProjectList'
import Pagination from '~/components/Pagination'
import SiteNotice from '~/components/SiteNotice'
export default {
  components: {
    ProjectList,
    Pagination,
    SiteNotice,
  },
  async asyncData({ $axios, params }) {
    const [projectPage] = await Promise.all([
      $axios.get('/api/project/projects', {
        params: {
          page: params.page || 1,
        },
      }),
    ])
    return {
      projectPage,
    }
  },
  head() {
    return {
      title: this.$siteTitle('开源项目'),
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
