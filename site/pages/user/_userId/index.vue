<template>
  <section class="main">
    <div class="container main-container left-main size-320">
      <div class="left-container">
        <user-profile :user="user" />

        <div class="tabs-warp">
          <div class="tabs">
            <ul>
              <li :class="{ 'is-active': activeTab === 'topics' }">
                <a :href="'/user/' + user.id + '?tab=topics'">
                  <span class="icon is-small">
                    <i class="iconfont icon-topic" aria-hidden="true" />
                  </span>
                  <span>话题</span>
                </a>
              </li>
              <li :class="{ 'is-active': activeTab === 'articles' }">
                <a :href="'/user/' + user.id + '?tab=articles'">
                  <span class="icon is-small">
                    <i class="iconfont icon-article" aria-hidden="true" />
                  </span>
                  <span>文章</span>
                </a>
              </li>
            </ul>
          </div>

          <div v-if="activeTab === 'topics'">
            <div
              v-if="
                topicsPage && topicsPage.results && topicsPage.results.length
              "
            >
              <load-more
                v-if="topicsPage"
                v-slot="{ results }"
                :init-data="topicsPage"
                :url="'/api/topic/user/topics?userId=' + user.id"
              >
                <topic-list :topics="results" :show-avatar="false" />
              </load-more>
            </div>
            <div v-else class="notification is-primary">
              暂无话题
            </div>
          </div>

          <div v-if="activeTab === 'articles'">
            <div
              v-if="
                articlesPage &&
                articlesPage.results &&
                articlesPage.results.length
              "
            >
              <load-more
                v-if="articlesPage"
                v-slot="{ results }"
                :init-data="articlesPage"
                :url="'/api/article/user/articles?userId=' + user.id"
              >
                <article-list :articles="results" />
              </load-more>
            </div>
            <div v-else class="notification is-primary">
              暂无文章
            </div>
          </div>
        </div>
      </div>
      <user-center-sidebar :user="user" />
    </div>
  </section>
</template>

<script>
import TopicList from '~/components/topic/TopicList'
import ArticleList from '~/components/ArticleList'
import UserProfile from '~/components/UserProfile'
import UserCenterSidebar from '~/components/UserCenterSidebar'
import LoadMore from '~/components/LoadMore'

const defaultTab = 'topics'

export default {
  components: {
    TopicList,
    ArticleList,
    UserProfile,
    UserCenterSidebar,
    LoadMore,
  },
  async asyncData({ $axios, params, query, error }) {
    let user
    try {
      user = await $axios.get('/api/user/' + params.userId)
    } catch (err) {
      error({
        statusCode: 404,
        message: err.message || '系统错误',
      })
      return
    }

    const activeTab = query.tab || defaultTab
    let topicsPage = null
    let articlesPage = null
    if (activeTab === 'topics') {
      topicsPage = await $axios.get('/api/topic/user/topics', {
        params: { userId: params.userId },
      })
    } else if (activeTab === 'articles') {
      articlesPage = await $axios.get('/api/article/user/articles', {
        params: { userId: params.userId },
      })
    }
    return {
      activeTab,
      user,
      topicsPage,
      articlesPage,
    }
  },
  data() {
    return {}
  },
  computed: {
    currentUser() {
      return this.$store.state.user.current
    },
    isOwner() {
      const current = this.$store.state.user.current
      return this.user && current && this.user.id === current.id
    },
  },
  head() {
    return {
      title: this.$siteTitle(this.user.nickname),
    }
  },
}
</script>

<style lang="scss" scoped>
.tabs-warp {
  background: #fff;
  padding: 0 10px 10px;

  .tabs {
    margin-bottom: 5px;
  }

  .more {
    text-align: right;
  }
}
</style>
