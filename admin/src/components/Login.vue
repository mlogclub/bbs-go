<template>
  <el-dialog
    title="登录"
    :visible.sync="isShowLogin"
    :show-close="false"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
  >
    <el-form status-icon label-width="100px">
      <el-form-item label="用户名" prop="username">
        <el-input
          type="text"
          v-model="username"
          autocomplete="off"
          @keyup.enter="submitForm"
        ></el-input>
      </el-form-item>
      <el-form-item label="密码" prop="password">
        <el-input
          type="password"
          v-model="password"
          autocomplete="off"
          @keyup.enter="submitForm"
        ></el-input>
      </el-form-item>
      <el-form-item label="验证码" prop="captchaCode">
        <el-input
          type="text"
          v-model="captchaCode"
          autocomplete="off"
          @keyup.enter="submitForm"
        ></el-input>
        <a v-if="captchaUrl" @click="showCaptcha">
          <img :src="captchaUrl" style="height: 40px;cursor: pointer;" />
        </a>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="submitForm()">登录</el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>

<script>
import HttpClient from '@/apis/HttpClient'
export default {
  data() {
    return {
      username: '',
      password: '',
      captchaId: '',
      captchaUrl: '',
      captchaCode: ''
    }
  },
  mounted() {
    this.showCaptcha()
  },
  methods: {
    async submitForm() {
      const params = {
        captchaId: this.captchaId,
        captchaCode: this.captchaCode,
        username: this.username,
        password: this.password
      }
      await this.$store.dispatch('Login/doLogin', params)
      await this.showCaptcha()
    },
    async showCaptcha() {
      try {
        const ret = await HttpClient.get('/api/captcha/request')
        this.captchaId = ret.captchaId
        this.captchaUrl = ret.captchaUrl
      } catch (e) {
        this.$message({ message: e.message || e, type: 'success' })
      }
    }
  },
  computed: {
    isShowLogin() {
      return this.$store.state.Login.showLogin
    }
  }
}
</script>

<style scoped></style>
