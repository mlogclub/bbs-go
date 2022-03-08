<template>
  <section class="main">
    <div class="container">
      <div class="main-body no-bg">
        <article
          class="message"
          :class="{ 'is-success': success, 'is-warning': !success }"
        >
          <div class="message-header">
            <p>邮箱验证</p>
          </div>
          <div class="message-body">
            <div v-if="success">
              恭喜，邮箱验证成功。你的邮箱为：{{ email }}
            </div>
            <div v-else>
              邮箱验证失败<span v-if="message" class="has-text-danger"
                >&nbsp;原因：{{ message }}</span
              >，请前往&nbsp;<nuxt-link
                to="/user/profile"
                style="font-weight: 700"
                >个人资料 &gt; 账号设置</nuxt-link
              >&nbsp;页面尝试重新发送验证邮件。
            </div>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  async asyncData({ $axios, query }) {
    try {
      const data = await $axios.post(
        '/api/user/verify_email?token=' + query.token
      )
      return { success: true, email: data.email }
    } catch (e) {
      return { success: false, message: e.message || '' }
    }
  },
}
</script>

<style lang="scss" scoped></style>
