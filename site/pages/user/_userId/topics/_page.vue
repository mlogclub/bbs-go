<template>
  <section class="main">
    <div class="container-wrapper">
      <div class="columns">
        <div class="column is-9">
          <div class="main-body">
            <nav
              class="breadcrumb"
              aria-label="breadcrumbs"
              style="margin-bottom: 10px;"
            >
              <ul>
                <li><a href="article">首页</a></li>
                <li>
                  <a :href="'/user/' + user.id + '?tab=topics'">{{
                    user.nickname
                  }}</a>
                </li>
                <li class="is-active">
                  <a href="#" aria-current="page">话题列表</a>
                </li>
              </ul>
            </nav>

            <topic-list :topics="topicsPage.results" />
            <pagination
              :page="topicsPage.page"
              :url-prefix="'/user/' + user.id + '/topics/'"
            />
          </div>
        </div>
        <div class="column is-3">
          <div class="main-aside">
            <user-center-sidebar :user="user" :current-user="currentUser" />
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import TopicList from '~/components/TopicList'
import Pagination from '~/components/Pagination'
import UserCenterSidebar from '~/components/UserCenterSidebar'
export default {
  components: {
    TopicList,
    Pagination,
    UserCenterSidebar
  },
  async asyncData({ $axios, params }) {
    const [currentUser, user, topicsPage] = await Promise.all([
      $axios.get('/api/user/current'),
      $axios.get('/api/user/' + params.userId),
      $axios.get('/api/topic/user/topics', {
        params: {
          userId: params.userId,
          page: params.page
        }
      })
    ])
    return {
      currentUser,
      user,
      topicsPage
    }
  },
  computed: {
    // 是否是主人态
    isOwner() {
      return (
        this.user && this.currentUser && this.user.id === this.currentUser.id
      )
    }
  },
  head() {
    return {
      title: this.$siteTitle(this.user.nickname + ' - 话题')
    }
  }
}
</script>

<style lang="scss" scoped></style>
