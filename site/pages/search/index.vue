<template>
  <section class="main">
    <div class="container main-container left-main">
      <div class="left-container">
        <div class="main-content no-padding">
          <div class="search-input">
            <form method="get" action="/search">
              <div class="field has-addons">
                <div class="control is-expanded">
                  <input
                    v-model="keyword"
                    name="q"
                    class="input"
                    type="text"
                    maxlength="30"
                    placeholder="请输入搜索关键字"
                  />
                </div>
                <div class="control">
                  <button type="submit" class="button is-info">
                    <span class="icon">
                      <i class="iconfont icon-search" />
                    </span>
                    <span>搜索</span>
                  </button>
                </div>
              </div>
            </form>
          </div>
        </div>

        <div v-if="docsPage && docsPage.results" class="topic-search-items">
          <div
            v-for="doc in docsPage.results"
            :key="doc.id"
            class="topic-search-item"
          >
            <a target="_blank" :href="'/topic/' + doc.topicId">
              <h1 class="topic-search-item-title" v-html="doc.title"></h1>
            </a>
            <p class="topic-search-item-summary" v-html="doc.summary"></p>
            <span class="topic-search-item-time">{{
              doc.createTime | formatDate
            }}</span>
          </div>
        </div>
        <div v-else-if="keyword" class="notification is-primary">
          没搜索到内容，请换一个搜索关键字
        </div>

        <pagination
          :page="docsPage.page"
          :url-prefix="'/search?q=' + keyword + '&p='"
        />
      </div>
      <div class="right-container">
        <check-in />
        <site-notice />
      </div>
    </div>
  </section>
</template>

<script>
import CheckIn from '@/components/CheckIn'
import SiteNotice from '@/components/SiteNotice'
import Pagination from '@/components/Pagination'

export default {
  components: {
    CheckIn,
    SiteNotice,
    Pagination,
  },
  async asyncData({ $axios, query }) {
    try {
      const keyword = query.q || ''
      const page = query.p || 1
      const [docsPage] = await Promise.all([
        $axios.get('/api/search/topic', {
          params: {
            keyword,
            page,
          },
        }),
      ])
      return { keyword, docsPage }
    } catch (e) {
      console.error(e)
    }
  },
}
</script>

<style lang="scss" scoped>
.search-input {
  background: #fff;
  padding: 10px;
  text-align: center;
}
.topic-search-items {
  background: #fff;
  .topic-search-item {
    padding: 10px;
    margin-bottom: 8px;
    background: #fff;
    box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
    //border-radius: 0.2rem;
    .topic-search-item-title {
      font-weight: bold;
      margin-bottom: 6px;
    }
    .topic-search-item-summary {
      font-size: 80%;
    }
    .topic-search-item-time {
      font-size: 80%;
      color: #8a9aa9;
    }
  }
}
</style>
