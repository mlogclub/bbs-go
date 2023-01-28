<template>
  <section class="main">
    <div class="container">
      <div class="main-body no-bg">
        <article
          class="message"
          :class="{ 'is-success': success, 'is-warning': !success }"
        >
          <div class="message-header">
            <p>购买成功</p>
          </div>
          <div class="message-body">
            <div v-if="success">恭喜，购买成功。</div>
            <div v-else>购买失败，请VX联系：YspCoder_</div>
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
      const data = await $axios.get('/api/pay/verify', {
        params: query,
      })
      return { success: data }
    } catch (e) {
      return { success: false, message: e.message || '' }
    }
  },
}
</script>

<style lang="scss" scoped></style>
