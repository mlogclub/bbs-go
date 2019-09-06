<template>
  <a class="is-dark" :class="{'button': isButton}" @click="githubLogin">
    <i class="iconfont icon-github" />&nbsp;
    <strong>{{ title }}</strong>
  </a>
</template>

<script>
export default {
  name: 'GithubLogin',
  props: {
    title: {
      type: String,
      default: '登录'
    },
    refUrl: { // 登录来源地址，控制登录成功之后要跳到该地址
      type: String,
      default: ''
    },
    isButton: {
      type: Boolean,
      default: true
    }
  },
  methods: {
    async githubLogin() {
      try {
        if (!this.refUrl && process.client) { // 如果没配置refUrl，那么取当前地址
          this.refUrl = window.location.pathname
        }
        const ret = await this.$axios.get('/api/login/github', {
          params: {
            ref: this.refUrl
          }
        })
        window.location = ret.url
      } catch (e) {
        console.error(e)
        this.$toast.error('登录失败：' + (e.message || e))
      }
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
