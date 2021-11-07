<template>
  <div>
    <button class="button is-success" @click="follow">
      <i class="iconfont el-icon-plus" />
      <span>关注</span>
    </button>
  </div>
</template>

<script>
export default {
  name: 'FollowBtn',
  props: {
    userId: {
      type: Number,
      required: true,
    },
  },
  computed: {
    user() {
      return this.$store.state.user.current
    },
  },
  methods: {
    async follow() {
      if (!this.user) {
        alert('请先登录')
        return
      }
      await this.$axios.post('/api/fans/follow', {
        userId: this.userId,
      })
      alert('关注成功')
    },
  },
}
</script>

<style lang="scss" scoped></style>
