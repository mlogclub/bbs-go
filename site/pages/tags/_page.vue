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
            <active-users :active-users="activeUsers" />
            <active-tags :active-tags="activeTags" />
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import ActiveUsers from '~/components/ActiveUsers'
import ActiveTags from '~/components/ActiveTags'
import Pagination from '~/components/Pagination'
export default {
  components: {
    ActiveUsers,
    ActiveTags,
    Pagination
  },
  head() {
    return {
      title: this.$siteTitle('标签')
    }
  },
  async asyncData({ $axios, params }) {
    const [tagsPage, activeUsers, activeTags] = await Promise.all([
      $axios.get('/api/tag/tags', {
        params: {
          page: params.page
        }
      }),
      $axios.get('/api/user/active'), // 活跃用户
      $axios.get('/api/tag/active') // 活跃标签
    ])
    return {
      tagsPage: tagsPage,
      activeUsers: activeUsers,
      activeTags: activeTags
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
