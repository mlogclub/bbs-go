<template>
  <div>
    <div class="comment-header">
      评论<span v-if="commentCount > 0">({{ commentCount }})</span>
    </div>

    <template v-if="isLogin">
      <div v-if="isNeedEmailVerify" class="comment-not-login">
        <div class="comment-login-div">
          请先前往
          <nuxt-link style="font-weight: 700" to="/user/profile"
            >个人中心 &gt; 个人资料</nuxt-link
          >页面设置邮箱，并完成邮箱认证。
        </div>
      </div>
      <template v-else>
        <comment-input
          ref="input"
          :mode="mode"
          :entity-id="entityId"
          :entity-type="entityType"
          @created="commentCreated"
        />
      </template>
    </template>
    <div v-else class="comment-not-login">
      <div class="comment-login-div">
        请
        <a style="font-weight: 700" @click="toLogin">登录</a>后发表观点
      </div>
    </div>

    <comment-list
      ref="list"
      :entity-id="entityId"
      :entity-type="entityType"
      :comments-page="commentsPage"
      @reply="reply"
    />
  </div>
</template>

<script>
export default {
  props: {
    mode: {
      type: String,
      default: 'markdown',
    },
    entityType: {
      type: String,
      default: '',
      required: true,
    },
    entityId: {
      type: Number,
      default: 0,
      required: true,
    },
    commentsPage: {
      type: Object,
      default() {
        return {}
      },
    },
    commentCount: {
      type: Number,
      default: 0,
    },
    showAd: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isLogin() {
      return this.$store.state.user.current != null
    },
    user() {
      return this.$store.state.user.current
    },
    config() {
      return this.$store.state.config.config
    },
    // 是否需要先邮箱认证
    isNeedEmailVerify() {
      return (
        this.config.createCommentEmailVerified &&
        this.user &&
        !this.user.emailVerified
      )
    },
  },
  methods: {
    commentCreated(data) {
      this.$refs.list.append(data)
    },
    reply(quote) {
      this.$refs.input.reply(quote)
    },
    toLogin() {
      this.$toSignin()
    },
  },
}
</script>
<style lang="scss" scoped>
.comment-header {
  display: flex;
  padding-top: 20px;
  margin: 0 10px;
  // border-top: 1px solid rgba(228, 228, 228, 0.6);
  color: #6d6d6d;
  font-size: 16px;
}

.comment-not-login {
  margin: 10px;
  border: 1px solid #f0f0f0;
  border-radius: 0;
  overflow: hidden;
  position: relative;
  padding: 10px;
  box-sizing: border-box;

  .comment-login-div {
    color: #d5d5d5;
    cursor: pointer;
    border-radius: 3px;
    padding: 0 10px;

    a {
      margin-left: 10px;
      margin-right: 10px;
    }
  }
}
</style>
