<template>
  <section class="main">
    <div class="container main-container is-white left-main">
      <div class="left-container">
        <div class="widget">
          <div class="widget-header">邮箱验证</div>
          <div class="widget-content">
            <div v-if="success">
              恭喜，邮箱验证成功。你的邮箱为：{{ currentUser.email }}，<a
                href="/user/settings"
                >点击前往资料页</a
              >
            </div>
            <div v-else>
              邮箱验证失败<span v-if="message" class="has-text-danger"
                >&nbsp;原因：{{ message }}</span
              >，请前往&nbsp;<a href="/user/settings">编辑资料</a
              >&nbsp;页面尝试重新发送验证邮件。
            </div>
          </div>
        </div>
      </div>
      <user-center-sidebar :user="currentUser" />
    </div>
  </section>
</template>

<script>
import UserCenterSidebar from '~/components/UserCenterSidebar'
export default {
  middleware: 'authenticated',
  components: { UserCenterSidebar },
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
