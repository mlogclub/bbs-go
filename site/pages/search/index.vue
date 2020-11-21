<template>
  <section class="main">
    <div class="container main-container left-main">
      <div class="left-container">
        <div class="main-content no-padding">
          <form method="get" action="/search">
            <div class="control has-icons-right">
              <input
                v-model="keyword"
                name="q"
                class="input"
                type="text"
                maxlength="30"
                placeholder="搜索"
              />
              <span class="icon is-medium is-right">
                <i class="iconfont icon-search" />
              </span>
            </div>
          </form>
        </div>

        <div v-if="docsPage && docsPage.results">
          <div v-for="doc in docsPage.results" :key="doc.id">
            <p v-html="doc.title"></p>
          </div>
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
      console.log(docsPage)
      return { keyword, docsPage }
    } catch (e) {
      console.error(e)
    }
  },
}
</script>

<style scoped></style>
