<template>
  <div class="widget">
    <div class="widget-header">
      <div>
        <span>粉丝</span>
        <span class="count">{{ user.fansCount }}</span>
      </div>
      <div class="slot">
        <nuxt-link to="/">更多</nuxt-link>
      </div>
    </div>
    <div class="widget-content">
      <div v-if="fansList && fansList.length">
        <user-follow-list :users="fansList" @onFollowed="onFollowed" />
      </div>
      <div v-else class="widget-tips">没有更多内容了</div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      fansList: [],
    }
  },
  mounted() {
    this.loadData()
  },
  methods: {
    async loadData() {
      const data = await this.$axios.get(
        '/api/fans/recent/fans?userId=' + this.user.id
      )
      this.fansList = data.results
    },
    async onFollowed(userId, followed) {
      await this.loadData()
    },
  },
}
</script>

<style lang="scss" scoped></style>
