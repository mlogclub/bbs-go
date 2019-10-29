<template>
  <div>
    <section class="main">
      <div class="container">
        <div class="columns">
          <div class="column is-12">
            <div class="widget">
              <div class="header">
                登录
              </div>
              <div class="content">
                <div class="field">
                  <label class="label">用户名/邮箱</label>
                  <div class="control has-icons-left">
                    <input
                      v-model="username"
                      class="input is-success"
                      type="text"
                      placeholder="请输入用户名或邮箱"
                      @keyup.enter="submitLogin"
                    >
                    <span class="icon is-small is-left"><i class="iconfont icon-username" /></span>
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
                      @keyup.enter="submitLogin"
                    >
                    <span class="icon is-small is-left"><i class="iconfont icon-password" /></span>
                  </div>
                </div>

                <div class="field">
                  <div class="control">
                    <button
                      class="button is-success"
                      @click="submitLogin"
                    >
                      登录
                    </button>
                    <github-login :ref-url="ref" />
                    <qq-login :ref-url="ref" />
                    <!--
                    <a
                      class="button is-text"
                      href="/user/signup"
                    >没有账号？点击这里去注册！</a>
                    -->
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import utils from '~/common/utils'
import GithubLogin from '~/components/GithubLogin'
import QqLogin from '~/components/QqLogin'
export default {
  components: {
    GithubLogin, QqLogin
  },
  data() {
    return {
      username: '',
      password: ''
    }
  },
  head() {
    return {
      title: this.$siteTitle('登录')
    }
  },
  asyncData({ params, query }) {
    return {
      ref: query.ref
    }
  },
  methods: {
    async submitLogin() {
      try {
        const user = await this.$store.dispatch('user/signin', {
          username: this.username,
          password: this.password
        })
        if (this.ref) { // 跳到登录前
          utils.linkTo(this.ref)
        } else { // 跳到个人主页
          utils.linkTo('/user/' + user.id)
        }
      } catch (e) {
        this.$toast.error(e.message || e)
      }
    }
  }
}
</script>
