<template>
  <section class="main">
    <div class="container">
      <div class="main-body no-bg">
        <div class="widget signup">
          <div class="widget-header" style="text-align: center">注册账号</div>
          <div class="widget-content">
            <form class="signup-form" @submit="clickSignup">
              <div class="field">
                <label class="label">
                  <span>昵称</span>
                  <span class="is-danger">*</span>
                </label>
                <div class="control">
                  <input
                    v-model="form.nickname"
                    class="input"
                    type="text"
                    placeholder="请输入昵称"
                  />
                </div>
              </div>

              <div class="field">
                <label class="label">
                  <span>邮箱</span>
                  <span class="is-danger">*</span>
                </label>
                <div class="control">
                  <input
                    v-model="form.email"
                    class="input"
                    type="email"
                    placeholder="请输入邮箱"
                  />
                </div>
              </div>

              <div class="field">
                <label class="label">
                  <span>密码</span>
                  <span class="is-danger">*</span>
                </label>
                <div class="control">
                  <input
                    v-model="form.password"
                    class="input"
                    type="password"
                    placeholder="请输入密码"
                  />
                </div>
                <p class="help">密码长度必须不少于 6 个字。</p>
              </div>

              <div class="field">
                <label class="label">
                  <span>确认密码</span>
                  <span class="is-danger">*</span>
                </label>
                <div class="control">
                  <input
                    v-model="form.rePassword"
                    class="input"
                    type="password"
                    placeholder="请再次输入密码"
                  />
                </div>
              </div>

              <div class="signup-btn">
                <el-button type="primary" @click="clickSignup">注册</el-button>
              </div>

              <div class="signup-bottom">
                <a @click="toSignin">已有账号，前往登录&gt;&gt;</a>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <CaptchaDialog ref="captchaDialog" />
  </section>
</template>

<script setup>
useHead({
  title: useSiteTitle("注册"),
});

const route = useRoute();
const form = reactive({
  nickname: "",
  email: "",
  password: "",
  rePassword: "",
  redirect: route.query.redirect || "",
  captchaId: "",
  captchaCode: "",
  captchaProtocol: 2,
});

const captchaDialog = ref(null);

const clickSignup = async () => {
  if (!form.nickname) {
    useMsgError("请输入昵称");
    return;
  }
  if (!form.email) {
    useMsgError("请输入邮箱");
    return;
  }
  if (!form.password) {
    useMsgError("请输入密码");
    return;
  }
  if (form.password !== form.rePassword) {
    useMsgError("两次输入密码不一致");
    return;
  }
  captchaDialog.value.show().then(async (captcha) => {
    form.captchaId = captcha.captchaId;
    form.captchaCode = captcha.captchaCode;

    try {
      const userStore = useUserStore();
      const { user, redirect } = await userStore.signup(form);

      if (redirect) {
        useLinkTo(redirect);
      } else {
        useLinkTo(`/user/${user.id}`);
      }
    } catch (err) {
      useCatchError(err);
    }
  });
};

function toSignin() {
  if (form.redirect) {
    useLinkTo(`/user/signin?redirect=${encodeURIComponent(form.redirect)}`);
  } else {
    useLinkTo("/user/signin");
  }
}
</script>

<style lang="scss" scoped>
.signup {
  max-width: 600px;
  margin: auto !important;

  .widget-header {
    justify-content: center;
  }
}
.signup-form {
  @media screen and (min-width: 768px) {
    padding: 20px;
  }

  .field {
    margin-bottom: 20px;
    .label {
      display: flex;
      align-items: center;
      column-gap: 6px;

      span {
        font-size: 15px;
        font-weight: 500;

        &.is-danger {
          line-height: 24px;
          color: red;
        }
      }
    }

    .help {
      color: var(--text-color3);
    }
  }

  .signup-btn {
    margin-top: 25px;
    width: 100%;
    button {
      width: 100%;
      height: 40px;
    }
  }

  .signup-bottom {
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
