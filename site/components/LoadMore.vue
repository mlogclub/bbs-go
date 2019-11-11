<template>
  <div
    v-infinite-scroll="loadMore"
    class="load-more"
    infinite-scroll-disabled="disabled"
    infinite-scroll-distance="10"
  >
    <slot :results="results" />
  </div>
</template>

<script>
export default {
  props: {
    // 请求URL
    url: {
      type: String,
      required: true
    },
    // 请求参数
    params: {
      type: Object,
      default: function () {
        return {}
      }
    },
    // 初始化数据
    initData: {
      type: Object,
      default: function () {
        return {
          results: [],
          cursor: ''
        }
      }
    }
  },
  data() {
    return {
      cursor: this.initData.cursor, // 分页标识
      results: this.initData.results, // 列表数据
      laoding: true, // 是否加载中
      hasMore: true, // 是否有更多数据
      busy: false // 是否还有数据
    }
  },
  computed: {
    // 是否禁言自动加载
    disabled() {
      return this.busy || !this.hasMore
    }
  },
  methods: {
    async loadMore() {
      this.busy = true
      try {
        const _params = Object.assign(this.params || {}, {
          cursor: this.cursor
        })
        const ret = await this.$axios.get(this.url, {
          params: _params
        })
        this.cursor = ret.cursor
        if (ret.results && ret.results.length) {
          ret.results.forEach((item) => {
            this.results.push(item)
          })
        } else {
          this.hasMore = false
        }
      } catch (err) {
        console.log(err)
      } finally {
        this.busy = false
      }
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
