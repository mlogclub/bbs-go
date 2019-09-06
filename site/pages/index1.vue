<template>
  <section class="main">
    <div class="container">
      <div class="columns">
        <div class="column is-9">
          <div class="main-body">
            <topic-list :topics="topicsPage.results" :show-ad="true" />
            <pagination :page="topicsPage.page" url-prefix="/topics/" />
          </div>
        </div>
        <div class="column is-3">
          <div class="main-aside">
            <div class="widget">
              <div class="content">
                <a class="button is-success" href="/topic/create">
                  <i class="iconfont icon-topic" />&nbsp;
                  <strong>发表主题</strong>
                </a>
              </div>
            </div>
            <WidgetUser :user="user" />
            <div style="text-align: center;">
              <!-- 展示广告288x288 -->
              <ins
                class="adsbygoogle"
                style="display:inline-block;width:288px;height:288px"
                data-ad-client="ca-pub-5683711753850351"
                data-ad-slot="4922900917"
              />
              <script>
                (adsbygoogle = window.adsbygoogle || []).push({});
              </script>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import TopicList from '~/components/TopicList'
import Pagination from '~/components/Pagination'
import WidgetUser from '~/components/WidgetUser'

export default {
  components: {
    TopicList, Pagination, WidgetUser
  },
  head() {
    return {
      meta: [
        { hid: 'description', name: 'description', content: this.$siteDescription() },
        { hid: 'keywords', name: 'keywords', content: this.$siteKeywords() }
      ]
    }
  },
  async asyncData({ $axios, params }) {
    try {
      const [user, topicsPage] = await Promise.all([
        $axios.get('/api/user/current'),
        $axios.get('/api/topic/topics')
      ])
      return { user: user, topicsPage: topicsPage }
    } catch (e) {
      console.error(e)
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
