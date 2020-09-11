<template>
  <section class="main">
    <div class="container main-container left-main">
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
            <div v-if="recentTopics && recentTopics.length">
              <topic-list :topics="recentTopics" :show-avatar="false" />
              <div class="more">
                <a :href="'/user/' + user.id + '/topics'">查看更多&gt;&gt;</a>
              </div>
            </div>
            <div v-else class="notification is-primary">
              暂无话题
            </div>
          </div>

          <div v-if="activeTab === 'articles'">
            <div v-if="recentArticles && recentArticles.length">
              <article-list :articles="recentArticles" />
              <div class="more">
                <a :href="'/user/' + user.id + '/articles'">查看更多&gt;&gt;</a>
              </div>
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
import TopicList from '~/components/TopicList'
import ArticleList from '~/components/ArticleList'
import UserProfile from '~/components/UserProfile'
import UserCenterSidebar from '~/components/UserCenterSidebar'

const defaultTab = 'topics'

export default {
  components: {
    TopicList,
    ArticleList,
    UserProfile,
    UserCenterSidebar,
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
    let recentTopics = null
    let recentArticles = null
    if (activeTab === 'topics') {
      recentTopics = await $axios.get('/api/topic/user/recent', {
        params: { userId: params.userId },
      })
    } else if (activeTab === 'articles') {
      recentArticles = await $axios.get('/api/article/user/recent', {
        params: { userId: params.userId },
      })
    }
    return {
      activeTab,
      user,
      recentTopics,
      recentArticles,
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
