<script setup>
const configStore = useConfigStore();
const userStore = useUserStore();
const route = useRoute();

useHead({
  title: `登录 | ${configStore.config.siteTitle}`,
});

const form = ref({
  username: "",
  password: "",
  captchaId: "",
  captchaUrl: "",
  captchaCode: "",
  redirect: route.query.redirect || "",
});
const loginMethod = computed(() => {
  return configStore.config.loginMethod;
});
// const currentUser = computed(() => {
//   return userStore.user
// })
// const isLogin = computed(() => {
//   return !!currentUser
// })
async function showCaptcha() {
  try {
    const { captchaId, captchaUrl } = await useHttpGet("/api/captcha/request", {
      params: {
        captchaId: form.value.captchaId,
      },
    });
    form.value.captchaId = captchaId;
    form.value.captchaUrl = captchaUrl;
  } catch (e) {
    useMsgError(e.message || e);
  }
}
async function submitLogin() {
  try {
    if (!form.value.username) {
      useMsgError("请输入用户名或邮箱");
      return;
    }
    if (!form.value.password) {
      useMsgError("请输入密码");
      return;
    }
    if (!form.value.captchaCode) {
      useMsgError("请输入验证码");
      return;
    }

    const { user, redirect } = await userStore.signin(form.value);
    if (redirect) {
      useLinkTo(redirect);
    } else {
      useLinkTo("/user/" + user.id);
    }
  } catch (e) {
    useMsgError(e.message || e);
    await showCaptcha();
  }
}
showCaptcha();
</script>

<template>
  <section class="main">
    <div class="container">
      <div class="main-body no-bg">
        <div class="widget signin">
          <div class="widget-header">登录</div>
          <div class="widget-content">
            <template v-if="loginMethod.password">
              <div class="field">
                <label class="label">用户名/邮箱</label>
                <div class="control has-icons-left">
                  <input
                    v-model="form.username"
                    class="input is-success"
                    type="text"
                    placeholder="请输入用户名或邮箱"
                    @keyup.enter="submitLogin"
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
                    v-model="form.password"
                    class="input"
                    type="password"
                    placeholder="请输入密码"
                    @keyup.enter="submitLogin"
                  />
                  <span class="icon is-small is-left"
                    ><i class="iconfont icon-password"
                  /></span>
                </div>
              </div>

              <div class="field">
                <label class="label">验证码</label>
                <div class="control has-icons-left">
                  <div class="field is-horizontal">
                    <div class="field login-captcha-input">
                      <input
                        v-model="form.captchaCode"
                        class="input"
                        type="text"
                        placeholder="验证码"
                        @keyup.enter="submitLogin"
                      />
                      <span class="icon is-small is-left"
                        ><i class="iconfont icon-captcha"
                      /></span>
                    </div>
                    <div
                      v-if="form.captchaUrl"
                      class="field login-captcha-img"
                      @click="showCaptcha"
                    >
                      <img :src="form.captchaUrl" data-not-lazy />
                    </div>
                  </div>
                </div>
              </div>

              <div class="field">
                <button class="button is-success" @click="submitLogin">
                  登录
                </button>
                <nuxt-link class="button is-text" to="/user/signup">
                  没有账号？点击这里去注册&gt;&gt;
                </nuxt-link>
              </div>
            </template>

            <div
              v-if="loginMethod.qq || loginMethod.github || loginMethod.osc"
              class="third-party-line"
            >
              <div class="third-party-title">
                <span>第三方账号登录</span>
              </div>
            </div>

            <div class="third-parties">
              <!-- <github-login v-if="loginMethod.github" :ref-url="ref" />
              <osc-login v-if="loginMethod.osc" :ref-url="ref" />
              <qq-login v-if="loginMethod.qq" :ref-url="ref" /> -->
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<style lang="scss" scoped></style>
