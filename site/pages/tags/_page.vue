<template>
  <section class="main">
    <div class="container">
      <div class="columns">
        <div class="column is-9">
          <div class="main-body">
            <div class="widget">
              <div class="header">
                <span>标签</span>
              </div>
              <div class="content">
                <div class="tags are-medium">
                  <span v-for="tag in tagsPage.results" :key="tag.tagId" class="tag is-normal">
                    <a
                      :href="'/articles/tag/' + tag.tagId "
                      :title="tag.tagName"
                      target="_blank"
                    >{{ tag.tagName }}</a>
                  </span>
                </div>
              </div>
              <pagination :page="tagsPage.page" url-prefix="/tags/" />
            </div>
          </div>
        </div>
        <div class="column is-3">
          <div class="main-aside">
            <!-- todo -->
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
    Pagination
  },
  head() {
    return {
      title: this.$siteTitle('标签')
    }
  },
  async asyncData({ $axios, params }) {
    const [tagsPage] = await Promise.all([
      $axios.get('/api/tag/tags', {
        params: {
          page: params.page
        }
      })
    ])
    return {
      tagsPage: tagsPage
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
