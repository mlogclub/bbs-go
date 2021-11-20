<template>
  <nav
    v-if="page.total > 0"
    class="pagination is-small"
    role="navigation"
    aria-label="pagination"
  >
    <a
      v-if="previousPageUrl"
      :href="previousPageUrl"
      class="pagination-previous"
      >上一页</a
    >
    <a v-else class="pagination-previous" disabled>上一页</a>

    <a v-if="nextPageUrl" :href="nextPageUrl" class="pagination-previous"
      >下一页</a
    >
    <a v-else class="pagination-previous" disabled>下一页</a>

    <ul class="pagination-list">
      <li v-for="p in pageList" :key="p">
        <a
          :class="{ 'is-current': p == page.page }"
          :href="getPageUrl(p)"
          class="pagination-link"
          >{{ p }}</a
        >
      </li>
    </ul>
  </nav>
</template>

<script>
export default {
  props: {
    page: {
      // 页码对象
      type: Object,
      default() {
        return { page: 1, total: 0, limit: 20 }
      },
      required: true,
    },
    urlPrefix: {
      // 分页url前缀
      type: String,
      default: '/',
      required: true,
    },
  },
  computed: {
    pageList() {
      const start = this.page.page - 2
      const end = this.page.page + 2
      const totalPage = this.getTotalPage()
      if (start <= 0) {
        const pages = []
        for (let i = 1; i <= 5 && i <= totalPage; i++) {
          pages.push(i)
        }
        return pages
      } else if (end > totalPage) {
        const pages = []
        let i = totalPage - 5 <= 0 ? 1 : totalPage - 5
        for (; i > 0 && i <= totalPage; i++) {
          pages.push(i)
        }
        return pages
      } else {
        return [
          this.page.page - 2,
          this.page.page - 1,
          this.page.page,
          this.page.page + 1,
          this.page.page + 2,
        ]
      }
    },
    previousPageUrl() {
      return this.getPreviousPageUrl()
    },
    nextPageUrl() {
      return this.getNextPageUrl()
    },
  },
  methods: {
    getNextPageUrl() {
      const nextPage = this.page.page + 1
      if (nextPage > this.getTotalPage()) {
        return ''
      }
      return this.getPageUrl(nextPage)
    },
    getPreviousPageUrl() {
      const previousPage = this.page.page - 1
      if (previousPage <= 0) {
        return ''
      }
      return this.getPageUrl(previousPage)
    },
    getPageUrl(page) {
      if (this.page.page === page) {
        return 'javascript:void(0)'
      }
      return this.urlPrefix + page
    },
    getTotalPage() {
      return this.page.total % this.page.limit > 0
        ? parseInt(this.page.total / this.page.limit) + 1
        : this.page.total / this.page.limit
    },
  },
}
</script>

<style lang="scss" scoped>
.pagination {
  margin: 10px 0;
}
</style>
