<template>
  <section class="main">
    <div class="container">
      <div class="main-body no-bg">
        <div class="widget signin">
          <div class="widget-header">找回密码</div>
          <div class="widget-content">
            <div class="field">
              <label class="label">邮箱</label>
              <div class="control has-icons-left">
                <input
                  v-model="email"
                  class="input is-success"
                  type="text"
                  placeholder="请输入邮箱"
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-email" />
                </span>
              </div>
            </div>

            <div class="field">
              <label class="label">新密码</label>
              <div class="control has-icons-left">
                <input
                  v-model="password"
                  class="input"
                  type="password"
                  placeholder="请输入密码"
                />
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
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-password" />
                </span>
              </div>
            </div>

            <div class="field">
              <label class="label">验证码</label>
              <div class="control has-icons-left">
                <div class="field is-horizontal">
                  <div class="field fpw-email-input">
                    <input
                      v-model="emailCode"
                      class="input"
                      type="text"
                      placeholder="验证码"
                    />
                    <span class="icon is-small is-left"
                      ><i class="iconfont icon-captcha"
                    /></span>
                  </div>
                  <div class="field">
                    <a class="button" @click="sendEmailCode">获取验证码</a>
                  </div>
                </div>
              </div>
            </div>

            <div class="field">
              <div class="control">
                <button class="button is-success" @click="forgotpwd">
                  提交
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  data() {
    return {
      email: '',
      password: '',
      rePassword: '',
      emailCode: '',
    }
  },
  head() {
    return {
      title: this.$siteTitle('找回密码'),
    }
  },
  methods: {
    async forgotpwd() {
      try {
        await this.$axios.post('/api/user/forgotpwd', {
          emailCode: this.emailCode,
          email: this.email,
          password: this.password,
          rePassword: this.rePassword,
        })
        const me = this
        this.$msg({
          message: '找回成功',
          onClose() {
            me.$linkTo('/user/signin')
          },
        })
      } catch (err) {
        this.$message.error(err.message || err)
      }
    },
    async sendEmailCode() {
      try {
        if (this.email === '') {
          this.$message.error('请输入邮箱!')
          return
        }
        await this.$axios.get('/api/user/send_email?email=' + this.email)
        this.$message.success('邮件发送成功')
      } catch (e) {
        this.$message.error(e.message || e)
      }
    },
  },
}
</script>

<style lang="scss" scoped></style>
<style>
.fpw-email-input {
  width: 70%;
  margin-right: 20px;
}
.fpw-email-input .input {
  width: 100% !important;
}
</style>
