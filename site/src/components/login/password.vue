<template>
  <div class="sms-login">
    <div class="login-field">
      <input
        v-model="form.username"
        type="text"
        placeholder="请输入用户名或邮箱"
        @keyup.enter="signin"
      />
    </div>

    <div class="login-field">
      <input
        v-model="form.password"
        type="password"
        placeholder="请输入密码"
        @keyup.enter="signin"
      />
    </div>

    <div class="login-btn">
      <el-button type="primary" @click="clickLogin">登录</el-button>
    </div>

    <!-- <div class="field">
      <button class="button is-link" @click="signin">登录</button>
      <a class="button is-text" @click="toSignup">
        没有账号？点击这里去注册&gt;&gt;
      </a>
    </div> -->

    <div class="login-bottom">
      <a @click="toSignup">没有账号？点击这里去注册&gt;&gt;</a>
    </div>

    <CaptchaDialog ref="captchaDialog" />
  </div>
</template>

<script setup>
const route = useRoute();
const form = reactive({
  username: "",
  password: "",
  redirect: route.query.redirect || "",
  captchaId: null,
  captchaCode: null,
  captchaProtocol: 2,
});

const captchaDialog = ref(null);

const clickLogin = async () => {
  if (!form.username) {
    useMsgError("请输入用户名或邮箱");
    return;
  }
  if (!form.password) {
    useMsgError("请输入密码");
    return;
  }

  captchaDialog.value.show().then(async (captcha) => {
    form.captchaId = captcha.captchaId;
    form.captchaCode = captcha.captchaCode;

    try {
      const userStore = useUserStore();
      const { user, redirect } = await userStore.signin(form);

      if (redirect) {
        useLinkTo(redirect);
      } else {
        useLinkTo(`/user/${user.id}`);
      }
    } catch (e) {
      useCatchError(e);
    }
  });
};

const toSignup = async () => {
  if (form.redirect) {
    useLinkTo(`/user/signup?redirect=${encodeURIComponent(form.redirect)}`);
  } else {
    useLinkTo("/user/signup");
  }
};
</script>

<style lang="scss" scoped>
.sms-login {
  max-width: 400px;
  margin: auto;
  .login-field {
    width: 100%;
    height: 39px;
    margin: 40px 0;
    display: flex;
    align-items: center;
    background-color: var(--bg-color2);
    border: 1px solid var(--border-color);
    border-radius: 3px;

    &:has(input:focus) {
      background-color: var(--bg-color3);
      border: 1px solid var(--border-hover-color);
      input {
        background-color: var(--bg-color3);
      }

      .phone-area {
        background-color: var(--bg-color3);
      }
    }

    input {
      padding: 0 15px;
      width: 100%;
      height: 37px;
      min-width: max-content;
      border: none;
      outline: none;
      background-color: var(--bg-color2);
      border-radius: 3px;

      &:-webkit-autofill,
      &:-webkit-autofill:hover,
      &:-webkit-autofill:focus,
      &:-webkit-autofill:active {
        -webkit-box-shadow: 0 0 0 30px white inset !important;
        -webkit-text-fill-color: black !important;
      }
      &:-webkit-autofill:selected {
        background-color: var(--bg-color2);
      }
    }
    span {
      color: var(--text-color);
      font-size: 14px;
      min-width: 36px;
    }
    img {
      height: 38px;
      cursor: pointer;
    }
    a {
      min-width: max-content;
      font-size: 13px;
      font-weight: 500;
    }
  }

  .login-btn {
    width: 100%;
    button {
      width: 100%;
      height: 40px;
    }
  }

  .login-bottom {
    margin: 20px 0;
    font-size: 13px;
    display: flex;
    justify-content: center;
    a {
      color: var(--text-color3);
    }
  }
}
</style>
