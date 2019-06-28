<template>
  <div>
    <el-dialog title="登录" :visible.sync="isShowLogin" :show-close="false" :close-on-click-modal="false"
               :close-on-press-escape="false" @open="open">
      <iframe id="loginFrame" :src="frameUrl" scrolling="no" frameborder="0" style="width: 100%; height: 300px;"/>
    </el-dialog>
  </div>
</template>

<script>
  import cookieManager from '../apis/CookieManager'
  import config from '../apis/Config'

  export default {
    name: 'Login',
    data() {
      return {
        frameUrl: ''
      }
    },
    methods: {
      // 登录弹窗打开之后
      open() {
        let me = this
        me.frameUrl = config.host + '/oauth/client'

        // 方法绑定到window
        window['loginSuccess'] = function (accessToken, refreshToken, expiry) {
          cookieManager.setCookie('accessToken', accessToken, 7)
          cookieManager.setCookie('refreshToken', refreshToken, 7)
          cookieManager.setCookie('expiry', expiry, 7)
          me.$store.dispatch('Login/hideLogin')
          me.frameUrl = ''
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

<style scoped>

</style>
