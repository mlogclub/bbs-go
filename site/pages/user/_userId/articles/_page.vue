<template>
  <section class="main">
    <div class="container main-container is-white left-main">
      <div class="left-container">
        <nav class="breadcrumb my-breadcrumb">
          <ul>
            <li><a href="article">首页</a></li>
            <li>
              <a :href="'/user/' + user.id + '?tab=articles'">{{
                user.nickname
              }}</a>
            </li>
            <li class="is-active">
              <a href="#" aria-current="page">文章列表</a>
            </li>
          </ul>
        </nav>

        <article-list :articles="articlesPage.results" />
        <pagination
          :page="articlesPage.page"
          :url-prefix="'/user/' + user.id + '/articles/'"
        />
      </div>
      <user-center-sidebar :user="user" />
    </div>
  </section>
</template>

<script>
import ArticleList from '~/components/ArticleList'
import Pagination from '~/components/Pagination'
import UserCenterSidebar from '~/components/UserCenterSidebar'
export default {
  components: {
    ArticleList,
    Pagination,
    UserCenterSidebar
  },
  async asyncData({ $axios, params, error }) {
    let user
    try {
      user = await $axios.get('/api/user/' + params.userId)
    } catch (err) {
      error({
        statusCode: 404,
        message: err.message || '系统错误'
      })
      return
    }

    const [articlesPage] = await Promise.all([
      $axios.get('/api/article/user/articles', {
        params: {
          userId: params.userId,
          page: params.page
        }
      })
    ])

    return {
      user,
      articlesPage
    }
  },
  computed: {
    currentUser() {
      return this.$store.state.user.current
    },
    // 是否是主人态
    isOwner() {
      const current = this.$store.state.user.current
      return this.user && current && this.user.id === current.id
    }
  },
  head() {
    return {
      title: this.$siteTitle(this.user.nickname + ' - 文章')
    }
  }
}
</script>

<style lang="scss" scoped></style>
