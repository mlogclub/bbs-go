<template>
  <a :class="{ button: isButton }" class="is-info" @click="qqLogin">
    <i class="iconfont icon-qq" />&nbsp;
    <strong>{{ title }}</strong>
  </a>
</template>

<script>
export default {
  name: 'QQLogin',
  props: {
    title: {
      type: String,
      default: 'QQ 登录',
    },
    refUrl: {
      // 登录来源地址，控制登录成功之后要跳到该地址
      type: String,
      default: '',
    },
    isButton: {
      type: Boolean,
      default: true,
    },
  },
  methods: {
    async qqLogin() {
      try {
        if (!this.refUrlValue && process.client) {
          // 如果没配置refUrl，那么取当前地址
          this.refUrlValue = window.location.pathname
        }
        const ret = await this.$axios.get('/api/qq/login/authorize', {
          params: {
            ref: this.refUrlValue,
          },
        })
        window.location = ret.url
      } catch (e) {
        console.error(e)
        this.$message.error('登录失败：' + (e.message || e))
      }
    },
  },
}
</script>

<style lang="scss" scoped></style>
