<template>
  <section class="main">
    <div class="container main-container left-main">
      <div class="left-container">
        <div
          v-if="(topics1 && topics1.length) || (topics2 && topics2.length)"
          class="widget"
        >
          <div class="widget-header">最新话题</div>
          <div class="widget-content">
            <div class="columns">
              <div class="column">
                <topic-list :topics="topics1" :show-ad="false" />
              </div>
              <div class="column">
                <topic-list :topics="topics2" :show-ad="false" />
              </div>
            </div>
          </div>
          <div class="widget-footer is-right">
            <a href="/topics" class="more-text">查看更多话题...</a>
          </div>
        </div>

        <div v-if="hotArticles && hotArticles.length" class="widget">
          <div class="widget-header">热门文章</div>
          <div class="widget-content">
            <article-list :articles="hotArticles" :show-ad="false" />
          </div>
          <div class="widget-footer is-right">
            <a href="/articles" class="more-text">查看更多文章...</a>
          </div>
        </div>

        <div v-if="articles && articles.length" class="widget">
          <div class="widget-header">最新文章</div>
          <div class="widget-content">
            <article-list :articles="articles" :show-ad="false" />
          </div>
          <div class="widget-footer is-right">
            <a href="/articles" class="more-text">查看更多文章...</a>
          </div>
        </div>
      </div>
      <div class="right-container">
        <post-btns />
        <weixin-gzh />

        <div class="widget">
          <div class="widget-header">新入驻</div>
          <div class="widget-content">
            <user-list :users="users" />
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import TopicList from '~/components/TopicList'
import ArticleList from '~/components/ArticleList'
import UserList from '~/components/UserList'
import WeixinGzh from '~/components/WeixinGzh'
import PostBtns from '~/components/PostBtns'
export default {
  components: { TopicList, ArticleList, UserList, WeixinGzh, PostBtns },
  async asyncData({ $axios, params }) {
    try {
      const [topics, articles, hotArticles, users] = await Promise.all([
        $axios.get('/api/topic/newest'),
        $axios.get('/api/article/newest'),
        $axios.get('/api/article/hot'),
        $axios.get('/api/user/newest')
      ])
      return {
        topics1: topics.slice(0, 5),
        topics2: topics.slice(5, 10),
        articles,
        hotArticles,
        users
      }
    } catch (e) {
      console.error(e)
    }
  },
  head() {
    return {
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

<style lang="scss" scoped>
.more-text {
  font-size: 14px;
  font-weight: bold;
  &:hover {
    color: #eb5424;
    text-decoration: underline;
  }
}
</style>
