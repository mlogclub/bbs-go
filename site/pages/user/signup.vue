<template>
  <section class="main">
    <div class="container">
      <div class="main-body">
        <div class="widget no-margin">
          <div class="header">
            注册
          </div>
          <div class="content">
            <div class="field">
              <label class="label">昵称</label>
              <div class="control has-icons-left">
                <input
                  v-model="nickname"
                  class="input is-success"
                  type="text"
                  placeholder="请输入昵称"
                  @keyup.enter="signup"
                >
                <span class="icon is-small is-left">
                  <i class="iconfont icon-username" />
                </span>
              </div>
            </div>

            <div class="field">
              <label class="label">用户名</label>
              <div class="control has-icons-left">
                <input
                  v-model="username"
                  class="input is-success"
                  type="text"
                  placeholder="请输入用户名"
                  @keyup.enter="signup"
                >
                <span class="icon is-small is-left">
                  <i class="iconfont icon-username" />
                </span>
              </div>
            </div>

            <div class="field">
              <label class="label">密码</label>
              <div class="control has-icons-left">
                <input
                  v-model="password"
                  class="input"
                  type="password"
                  placeholder="请输入密码"
                  @keyup.enter="signup"
                >
                <span class="icon is-small is-left">
                  <i class="iconfont icon-password" />
                </span>
              </div>
            </div>

            <div class="field">
              <label class="label">确认密码</label>
              <div class="control has-icons-left">
                <input
                  v-model="rePassword"
                  class="input"
                  type="password"
                  placeholder="请再次输入密码"
                  @keyup.enter="signup"
                >
                <span class="icon is-small is-left">
                  <i class="iconfont icon-password" />
                </span>
              </div>
            </div>

            <div class="field">
              <div class="control">
                <button class="button is-success" @click="signup">
                  注册
                </button>
                <github-login :ref-url="ref" />
                <qq-login :ref-url="ref" />
                <nuxt-link class="button is-text" to="/user/signin">
                  已有账号，前往登录&gt;&gt;
                </nuxt-link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import utils from '~/common/utils'
export default {
  head() {
    return {
      title: this.$siteTitle('注册')
    }
  },
  data() {
    return {
      nickname: '',
      username: '',
      password: '',
      rePassword: ''
    }
  },
  asyncData({ params, query }) {
    return {
      ref: query.ref
    }
  },
  methods: {
    async signup() {
      try {
        await this.$store.dispatch('user/signup', {
          nickname: this.nickname,
          username: this.username,
          password: this.password,
          rePassword: this.rePassword,
          ref: this.ref
        })
        if (this.ref) { // 跳到登录前
          utils.linkTo(this.ref)
        } else { // 跳到个人主页
          utils.linkTo('/user/settings')
        }
      } catch (err) {
        this.$toast.error(err.message || err)
      }
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
