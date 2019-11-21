<template>
  <section class="main">
    <div class="container-wrapper">
      <div class="columns">
        <div class="column is-9">
          <div class="main-body">
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
                <topic-list :topics="recentTopics" />
                <div class="more">
                  <a :href="'/user/' + user.id + '/topics'">查看更多&gt;&gt;</a>
                </div>
              </div>
              <div
                v-else
                class="notification is-primary"
                style="margin-top: 10px;"
              >
                暂无话题
              </div>
            </div>

            <div v-if="activeTab === 'articles'">
              <div v-if="recentArticles && recentArticles.length">
                <article-list :articles="recentArticles" />
                <div class="more">
                  <a :href="'/user/' + user.id + '/articles'"
                    >查看更多&gt;&gt;</a
                  >
                </div>
              </div>
              <div
                v-else
                class="notification is-primary"
                style="margin-top: 10px;"
              >
                暂无文章
              </div>
            </div>
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
import ArticleList from '~/components/ArticleList'
import UserCenterSidebar from '~/components/UserCenterSidebar'

const defaultTab = 'topics'

export default {
  components: {
    TopicList,
    ArticleList,
    UserCenterSidebar
  },
  async asyncData({ $axios, params, query }) {
    const activeTab = query.tab || defaultTab
    const [currentUser, user] = await Promise.all([
      $axios.get('/api/user/current'),
      $axios.get('/api/user/' + params.userId)
    ])
    let recentTopics = null
    let recentArticles = null
    if (activeTab === 'topics') {
      recentTopics = await $axios.get('/api/topic/user/recent', {
        params: { userId: params.userId }
      })
    } else if (activeTab === 'articles') {
      recentArticles = await $axios.get('/api/article/user/recent', {
        params: { userId: params.userId }
      })
    }
    return {
      activeTab,
      currentUser,
      user,
      recentTopics,
      recentArticles
    }
  },
  data() {
    return {}
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
      title: this.$siteTitle(this.user.nickname)
    }
  }
}
</script>

<style lang="scss" scoped>
.tabs {
  margin-bottom: 5px;
}
.more {
  text-align: right;
}
</style>
