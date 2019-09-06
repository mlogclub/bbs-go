<template>
  <el-dialog
    title="登录"
    :visible.sync="isShowLogin"
    :show-close="false"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
  >
    <el-form
      status-icon
      label-width="100px"
    >
      <el-form-item label="用户名" prop="username">
        <el-input type="text" v-model="username" autocomplete="off"></el-input>
      </el-form-item>
      <el-form-item label="密码" prop="password">
        <el-input type="password" v-model="password" autocomplete="off"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="submitForm()">登录</el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>

<script>
export default {
  data() {
    return {
      username: '',
      password: ''
    };
  },
  methods: {
    submitForm() {
      const params = {
        username: this.username,
        password: this.password
      }
      this.$store.dispatch('Login/doLogin',  params)
    }
  },
  computed: {
    isShowLogin() {
      return this.$store.state.Login.showLogin;
    }
  }
};
</script>

<style scoped>
</style>
