<template>
  <section class="main">
    <div class="container">
      <div class="main-body">
        <div class="widget">
          <div class="widget-header">
            登录
          </div>
          <div class="widget-content">
            <div class="field">
              <label class="label">用户名/邮箱</label>
              <div class="control has-icons-left">
                <input
                  v-model="username"
                  @keyup.enter="submitLogin"
                  class="input is-success"
                  type="text"
                  placeholder="请输入用户名或邮箱"
                />
                <span class="icon is-small is-left"
                  ><i class="iconfont icon-username"
                /></span>
              </div>
            </div>

            <div class="field">
              <label class="label">密码</label>
              <div class="control has-icons-left">
                <input
                  v-model="password"
                  @keyup.enter="submitLogin"
                  class="input"
                  type="password"
                  placeholder="请输入密码"
                />
                <span class="icon is-small is-left"
                  ><i class="iconfont icon-password"
                /></span>
              </div>
            </div>

            <div class="field">
              <div class="control">
                <button @click="submitLogin" class="button is-success">
                  登录
                </button>
                <github-login :ref-url="ref" />
                <qq-login :ref-url="ref" />
                <nuxt-link class="button is-text" to="/user/signup">
                  没有账号？点击这里去注册&gt;&gt;
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
import GithubLogin from '~/components/GithubLogin'
import QqLogin from '~/components/QqLogin'
export default {
  components: {
    GithubLogin,
    QqLogin
  },
  asyncData({ params, query }) {
    return {
      ref: query.ref
    }
  },
  data() {
    return {
      username: '',
      password: ''
    }
  },
  methods: {
    async submitLogin() {
      try {
        const user = await this.$store.dispatch('user/signin', {
          username: this.username,
          password: this.password,
          ref: this.ref
        })
        if (this.ref) {
          // 跳到登录前
          utils.linkTo(this.ref)
        } else {
          // 跳到个人主页
          utils.linkTo('/user/' + user.id)
        }
      } catch (e) {
        this.$toast.error(e.message || e)
      }
    }
  },
  head() {
    return {
      title: this.$siteTitle('登录')
    }
  }
}
</script>
