<template>
  <div class="widget">
    <div class="widget-header">
      <span>关注</span>
      <span class="count">{{ user.followCount }}</span>
    </div>
    <div class="widget-content">
      <div v-if="followList && followList.length">
        <user-follow-list :users="followList" @onFollowed="onFollowed" />
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
      followList: [],
    }
  },
  mounted() {
    this.loadData()
  },
  methods: {
    async loadData() {
      const data = await this.$axios.get(
        '/api/fans/recent/follow?userId=' + this.user.id
      )
      this.followList = data.results
    },
    async onFollowed(userId, followed) {
      await this.loadData()
    },
  },
}
</script>

<style lang="scss" scoped></style>
