<template>
  <section class="main">
    <div class="container-wrapper main-container left-main">
      <div class="left-container">
        <project-list :projects="projectPage.results" />
        <pagination :page="projectPage.page" url-prefix="/projects/" />
      </div>
      <div class="right-container">
        <!-- 展示广告190x90 -->
        <ins
          class="adsbygoogle"
          style="display:inline-block;width:190px;height:90px"
          data-ad-client="ca-pub-5683711753850351"
          data-ad-slot="9345305153"
        />
        <script>
          ;(adsbygoogle = window.adsbygoogle || []).push({})
        </script>

        <!-- 展示广告190x190 -->
        <ins
          class="adsbygoogle"
          style="display:inline-block;width:190px;height:190px"
          data-ad-client="ca-pub-5683711753850351"
          data-ad-slot="5685455263"
        />
        <script>
          ;(adsbygoogle = window.adsbygoogle || []).push({})
        </script>

        <!-- 展示广告190x480 -->
        <ins
          class="adsbygoogle"
          style="display:inline-block;width:190px;height:480px"
          data-ad-client="ca-pub-5683711753850351"
          data-ad-slot="3438372357"
        />
        <script>
          ;(adsbygoogle = window.adsbygoogle || []).push({})
        </script>
      </div>
    </div>
  </section>
</template>

<script>
import ProjectList from '~/components/ProjectList'
import Pagination from '~/components/Pagination'
export default {
  components: {
    ProjectList,
    Pagination
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
