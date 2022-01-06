<template>
  <div>
    <button
      class="button follow-btn"
      :class="{ 'is-followed': followed }"
      @click="follow"
    >
      <i class="iconfont icon-add" />
      <span>{{ followed ? '已关注' : '关注' }}</span>
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
    followed: {
      type: Boolean,
      default: false,
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
        if (this.followed) {
          await this.$axios.post('/api/fans/unfollow', {
            userId: this.userId,
          })
          this.$emit('onFollowed', this.userId, false)
          this.$message.success('取消关注成功')
        } else {
          await this.$axios.post('/api/fans/follow', {
            userId: this.userId,
          })
          this.$emit('onFollowed', this.userId, true)
          this.$message.success('关注成功')
        }
      } catch (e) {
        this.$message.error(e.message || e)
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.follow-btn {
  font-size: 12px;
  height: 25px;
  background-color: #2469f6; // TODO
  border-color: #2469f6;
  color: var(--text-color5);

  &:hover,
  &.is-followed {
    background-color: #7ba5f9; // TODO
    border-color: #7ba5f9;
  }
  i {
    font-size: 12px;
    margin-right: 5px;
  }
}
</style>
