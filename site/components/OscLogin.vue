<template>
  <a :class="{ button: isButton }" class="is-link" @click="oscLogin">
    <i class="iconfont icon-open-source" />&nbsp;
    <strong>{{ title }}</strong>
  </a>
</template>

<script>
export default {
  props: {
    title: {
      type: String,
      default: 'Oschina 登录',
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
  data() {
    return {
      refUrlValue: this.refUrl,
    }
  },
  methods: {
    async oscLogin() {
      try {
        if (!this.refUrlValue && process.client) {
          // 如果没配置refUrl，那么取当前地址
          this.refUrlValue = window.location.pathname
        }
        const ret = await this.$axios.get('/api/osc/login/authorize', {
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
