<template>
  <div v-if="show" class="container" style="margin-bottom: 10px;">
    <article class="message is-info">
      <div class="message-header">
        <p>公告</p>
        <button @click="close" class="delete" aria-label="delete"></button>
      </div>
      <div v-html="config.siteNotification" class="message-body"></div>
    </article>
  </div>
</template>

<script>
const closeTimeKey = 'top.notice.close.time'

export default {
  data() {
    return {
      show: false
    }
  },
  computed: {
    config() {
      return this.$store.state.config.config
    }
  },
  mounted() {
    this.show = this.isShow()
  },
  methods: {
    close() {
      this.show = false
      this.$cookies.set(closeTimeKey, new Date().getTime())
    },
    isShow() {
      if (!this.config.siteNotification) {
        return false
      }
      const closeTime = this.$cookies.get(closeTimeKey) // 上次关闭的时间
      if (!closeTime) {
        // 说明没关闭过
        return true
      }
      if (new Date().getTime() - parseInt(closeTime) >= 86400 * 1000) {
        // 如果关闭时间大于1天
        return true
      }
      return false
    }
  }
}
</script>

<style lang="scss" scoped></style>
