<template>
  <section class="main">
    <div class="container">
      <user-profile :user="currentUser" />
      <div class="container main-container right-main size-320">
        <user-center-sidebar :user="currentUser" />
        <div class="right-container">
          <div class="widget">
            <div class="widget-header">邮箱验证</div>
            <div class="widget-content">
              <div v-if="success">
                恭喜，邮箱验证成功。你的邮箱为：{{
                  currentUser.email
                }}，<nuxt-link to="/user/profile">点击前往资料页</nuxt-link>
              </div>
              <div v-else>
                邮箱验证失败<span v-if="message" class="has-text-danger"
                  >&nbsp;原因：{{ message }}</span
                >，请前往&nbsp;<nuxt-link to="/user/profile">个人资料</nuxt-link
                >&nbsp;页面尝试重新发送验证邮件。
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  middleware: 'authenticated',
  async asyncData({ $axios, query }) {
    try {
      await $axios.get('/api/user/email/verify?token=' + query.token)
      return { success: true }
    } catch (e) {
      return { success: false, message: e.message || '' }
    }
  },
  computed: {
    currentUser() {
      return this.$store.state.user.current
    },
  },
}
</script>

<style lang="scss" scoped></style>
