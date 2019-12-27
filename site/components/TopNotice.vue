<template>
  <div v-if="show" class="container" style="margin-bottom: 10px;">
    <article class="message is-info">
      <div class="message-header">
        <p>温馨提示</p>
        <button @click="close" class="delete" aria-label="delete"></button>
      </div>
      <div class="message-body">
        欢迎访问&nbsp;码农俱乐部，<a href="/user/settings"
          ><strong>点击这里设置您的邮箱</strong></a
        >&nbsp;可以接收站内跟帖、回复邮件提醒，不错过任何一条消息。bbs-go交流群：<strong
          >653248175</strong
        >
      </div>
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
  mounted() {
    this.show = this.isShow()
  },
  methods: {
    close() {
      this.show = false
      this.$cookies.set(closeTimeKey, new Date().getTime())
    },
    isShow() {
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
