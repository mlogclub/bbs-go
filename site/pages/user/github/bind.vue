<template>
  <section class="main">
    <div class="container">
      <div class="card">
        <header class="card-header">
          <p class="card-header-title">
            绑定账号
          </p>
        </header>
        <div class="card-content">
          <div class="tabs is-centered">
            <ul>
              <li :class="{'is-active': bindType === 'login'}">
                <a @click="switchTo('login')">
                  <span>绑定已有账号</span>
                </a>
              </li>
              <li :class="{'is-active': bindType === 'signup'}">
                <a @click="switchTo('signup')">
                  <span>注册并绑定</span>
                </a>
              </li>
            </ul>
          </div>

          <div v-if="bindType === 'login'">
            <div class="field">
              <label class="label">
                <span style="color:red;">*&nbsp;</span>用户名/邮箱
              </label>
              <div class="control has-icons-left">
                <input
                  v-model="usernameOrEmail"
                  class="input is-success"
                  type="text"
                  placeholder="请输入用户名或邮箱"
                  @keydown.enter="submitForm"
                >
                <span class="icon is-small is-left">
                  <i class="iconfont icon-username" />
                </span>
              </div>
            </div>

            <div class="field">
              <label class="label">
                <span style="color:red;">*&nbsp;</span>密码
              </label>
              <div class="control has-icons-left">
                <input v-model="password" class="input" type="password" placeholder="请输入密码" @keydown.enter="submitForm">
                <span class="icon is-small is-left">
                  <i class="iconfont icon-password" />
                </span>
              </div>
            </div>
          </div>

          <div v-if="bindType === 'signup'">
            <div class="field">
              <label class="label">
                <span style="color:red;">*&nbsp;</span>用户名
              </label>
              <div class="control has-icons-left">
                <input v-model="username" class="input is-success" type="text" placeholder="请输入用户名" @keydown.enter="submitForm">
                <span class="icon is-small is-left">
                  <i class="iconfont icon-username" />
                </span>
              </div>
            </div>

            <div class="field">
              <label class="label">
                <span style="color:red;">*&nbsp;</span>邮箱
              </label>
              <div class="control has-icons-left">
                <input v-model="email" class="input is-success" type="text" placeholder="请输入邮箱" @keydown.enter="submitForm">
                <span class="icon is-small is-left">
                  <i class="iconfont icon-email" />
                </span>
              </div>
            </div>

            <div class="field">
              <label class="label">
                <span style="color:red;">*&nbsp;</span>昵称
              </label>
              <div class="control has-icons-left">
                <input v-model="nickname" class="input is-success" type="text" placeholder="请输入昵称" @keydown.enter="submitForm">
                <span class="icon is-small is-left">
                  <i class="iconfont icon-username" />
                </span>
              </div>
            </div>

            <div class="field">
              <label class="label">
                <span style="color:red;">*&nbsp;</span>密码
              </label>
              <div class="control has-icons-left">
                <input v-model="password" class="input" type="password" placeholder="请输入密码" @keydown.enter="submitForm">
                <span class="icon is-small is-left">
                  <i class="iconfont icon-password" />
                </span>
              </div>
            </div>

            <div class="field">
              <label class="label">
                <span style="color:red;">*&nbsp;</span>确认密码
              </label>
              <div class="control has-icons-left">
                <input v-model="rePassword" class="input" type="password" placeholder="请再次输入密码" @keydown.enter="submitForm">
                <span class="icon is-small is-left">
                  <i class="iconfont icon-password" />
                </span>
              </div>
            </div>
          </div>

          <div class="field">
            <label class="label">&nbsp;</label>
            <div class="control">
              <a class="button is-success" @click="submitForm">绑定</a>
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
  data() {
    return {
      usernameOrEmail: '', // 用户名或邮箱
      username: '', // 用户名
      email: '', // 邮箱
      nickname: '', // 昵称
      password: '', // 密码
      rePassword: '', // 确认密码
      bindType: 'login' // login/signup
    }
  },
  head() {
    return {
      title: this.$siteTitle('绑定账号')
    }
  },
  async asyncData({ $axios, params, query }) {
    const githubUser = await $axios.get(
      '/api/login/github/user/' + query.githubId
    )
    return {
      githubUser: githubUser,
      ref: query.ref
    }
  },
  methods: {
    async submitForm() {
      try {
        const ret = await this.$axios.post('/api/login/github/bind', {
          bindType: this.bindType,
          githubId: this.githubUser.id,
          username: this.bindType === 'login' ? this.usernameOrEmail : this.username,
          email: this.email,
          password: this.password,
          rePassword: this.rePassword,
          nickname: this.nickname,
          ref: this.ref // 成功之后的跳转地址
        })

        this.$store.dispatch('user/loginSuccess', ret)

        this.$toast.success('绑定成功', {
          duration: 1000,
          onComplete: function () {
            if (ret.ref) {
              utils.linkTo(ret.ref)
            } else {
              utils.linkTo('/user/' + ret.user.id)
            }
          }
        })
      } catch (e) {
        console.error(e)
        this.$toast.error('绑定失败：' + (e.message || e))
      }
    },
    switchTo(bindType) {
      this.bindType = bindType
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
