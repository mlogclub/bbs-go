<template>
  <section class="main">
    <div class="container main-container is-white left-main">
      <div class="left-container">
        <project-list :projects="projectPage.results" />
        <pagination :page="projectPage.page" url-prefix="/projects/" />
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
import ProjectList from '~/components/ProjectList'
import Pagination from '~/components/Pagination'
import WeixinGzh from '~/components/WeixinGzh'
export default {
  components: {
    ProjectList,
    Pagination,
    WeixinGzh
  },
  async asyncData({ $axios, params }) {
    const [projectPage] = await Promise.all([
      $axios.get('/api/project/projects', {
        params: {
          page: params.page || 1
        }
      })
    ])
    return {
      projectPage
    }
  },
  head() {
    return {
      title: this.$siteTitle('开源项目'),
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
