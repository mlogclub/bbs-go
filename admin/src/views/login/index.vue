<template>
  <div class="login-container">
    <el-form
      ref="loginForm"
      :model="loginForm"
      label-position="left"
      class="login-form"
      autocomplete="on"
    >
      <div class="title-container">
        <h3 class="title">登录</h3>
      </div>

      <el-form-item prop="username">
        <el-input
          ref="username"
          v-model="loginForm.username"
          placeholder="用户名"
          name="username"
          type="text"
          tabindex="1"
          autocomplete="on"
          @keyup.enter.native="handleLogin"
        />
      </el-form-item>

      <el-tooltip v-model="capsTooltip" content="Caps lock is On" placement="right" manual>
        <el-form-item prop="password">
          <el-input
            ref="password"
            v-model="loginForm.password"
            type="password"
            placeholder="密码"
            name="password"
            tabindex="2"
            autocomplete="on"
            @keyup.native="checkCapslock"
            @blur="capsTooltip = false"
            @keyup.enter.native="handleLogin"
          />
        </el-form-item>
      </el-tooltip>

      <el-form-item prop="captchaCode" class="captcha-code">
        <el-input
          ref="username"
          v-model="loginForm.captchaCode"
          placeholder="验证码"
          name="captchaCode"
          type="text"
          tabindex="3"
          autocomplete="off"
          @keyup.enter.native="handleLogin"
        />
        <div v-if="loginForm.captchaUrl" class="captcha-code-img">
          <a @click="showCaptcha"><img :src="loginForm.captchaUrl" /></a>
        </div>
      </el-form-item>
      <el-button
        :loading="loading"
        type="primary"
        style="width: 100%; margin-bottom: 30px"
        @click.native.prevent="handleLogin"
      >
        Login
      </el-button>
    </el-form>
  </div>
</template>

<script>
export default {
  data() {
    return {
      loginForm: {
        username: "",
        password: "",
        captchaId: "",
        captchaUrl: "",
        captchaCode: "",
      },
      capsTooltip: false,
      loading: false,
      redirect: undefined,
      otherQuery: {},
    };
  },
  watch: {
    $route: {
      handler(route) {
        const { query } = route;
        if (query) {
          this.redirect = query.redirect;
          this.otherQuery = this.getOtherQuery(query);
        }
      },
      immediate: true,
    },
  },
  mounted() {
    if (this.loginForm.username === "") {
      this.$refs.username.focus();
    } else if (this.loginForm.password === "") {
      this.$refs.password.focus();
    }
    this.showCaptcha();
  },
  methods: {
    async showCaptcha() {
      try {
        const ret = await this.axios.get("/api/captcha/request", {
          params: {
            captchaId: this.loginForm.captchaId || "",
          },
        });
        this.loginForm.captchaId = ret.captchaId;
        this.loginForm.captchaUrl =
          process.env.VUE_APP_BASE_API +
          "/api/captcha/show?captchaId=" +
          this.loginForm.captchaId +
          "&timestamp=" +
          new Date().getTime();
      } catch (e) {
        this.$message.error(e.message || e);
      }
    },
    checkCapslock({ shiftKey, key } = {}) {
      if (key && key.length === 1) {
        if ((shiftKey && key >= "a" && key <= "z") || (!shiftKey && key >= "A" && key <= "Z")) {
          this.capsTooltip = true;
        } else {
          this.capsTooltip = false;
        }
      }
      if (key === "CapsLock" && this.capsTooltip === true) {
        this.capsTooltip = false;
      }
    },
    handleLogin() {
      this.$refs.loginForm.validate((valid) => {
        if (valid) {
          this.loading = true;
          this.$store
            .dispatch("user/login", this.loginForm)
            .then(() => {
              this.$router.push({ path: this.redirect || "/", query: this.otherQuery });
              this.loading = false;
            })
            .catch((e) => {
              this.showCaptcha();
              this.loading = false;
            });
          return true;
        }
        return false;
      });
    },
    getOtherQuery(query) {
      return Object.keys(query).reduce((acc, cur) => {
        if (cur !== "redirect") {
          acc[cur] = query[cur];
        }
        return acc;
      }, {});
    },
  },
};
</script>

<style lang="scss">
.login-container {
  min-height: 100%;
  width: 100%;
  background-color: #fff;
  overflow: hidden;

  .login-form {
    position: relative;
    width: 520px;
    max-width: 100%;
    padding: 160px 35px 0;
    margin: 0 auto;
    overflow: hidden;

    .captcha-code {
      & > div {
        display: flex;
        .captcha-code-img {
          // margin-left: 10px;
          img {
            height: 36px;
          }
        }
      }
    }
  }

  .tips {
    font-size: 14px;
    color: #fff;
    margin-bottom: 10px;

    span {
      &:first-of-type {
        margin-right: 16px;
      }
    }
  }

  .title-container {
    position: relative;

    .title {
      font-size: 26px;
      color: #000;
      margin: 0 auto 40px auto;
      text-align: center;
      font-weight: bold;
    }
  }
}
</style>
