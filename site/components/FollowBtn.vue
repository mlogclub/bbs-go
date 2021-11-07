<template>
  <div>
    <button class="button is-success is-small" @click="follow">
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
        this.$msgSignIn()
        return
      }
      try {
        await this.$axios.post('/api/fans/follow', {
          userId: this.userId,
        })
        this.$message.success('关注成功')
      } catch (e) {
        this.$message.error(e.message || e)
      }
    },
  },
}
</script>

<style lang="scss" scoped></style>
