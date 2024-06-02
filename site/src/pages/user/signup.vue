<template>
  <section class="main">
    <div class="container">
      <div class="main-body no-bg">
        <div class="widget signin">
          <div class="widget-header">注册</div>
          <div class="widget-content">
            <div class="field">
              <label class="label">昵称</label>
              <div class="control has-icons-left">
                <input
                  v-model="form.nickname"
                  class="input is-success"
                  type="text"
                  placeholder="请输入昵称"
                  @keyup.enter="signup"
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-username" />
                </span>
              </div>
            </div>

            <div class="field">
              <label class="label">邮箱</label>
              <div class="control has-icons-left">
                <input
                  v-model="form.email"
                  class="input is-success"
                  type="text"
                  placeholder="请输入邮箱"
                  @keyup.enter="signup"
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-email" />
                </span>
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
                  @keyup.enter="signup"
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
                  v-model="form.rePassword"
                  class="input"
                  type="password"
                  placeholder="请再次输入密码"
                  @keyup.enter="signup"
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
                  <div class="field login-captcha-input">
                    <input
                      v-model="form.captchaCode"
                      class="input"
                      type="text"
                      placeholder="验证码"
                      @keyup.enter="signup"
                    />
                    <span class="icon is-small is-left"
                      ><i class="iconfont icon-captcha"
                    /></span>
                  </div>
                  <div v-if="form.captchaUrl" class="field login-captcha-img">
                    <a @click="refreshCaptcha"
                      ><img :src="form.captchaUrl"
                    /></a>
                  </div>
                </div>
              </div>
            </div>

            <div class="field">
              <div class="control">
                <button class="button is-link" @click="signup">注册</button>
                <a class="button is-text" @click="toSignin">
                  已有账号，前往登录&gt;&gt;
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
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
  captchaId: "",
  captchaUrl: "",
  captchaCode: "",
  redirect: route.query.redirect || "",
});

refreshCaptcha();

async function refreshCaptcha() {
  try {
    const { data: captcha } = await useAsyncData(() => {
      return useMyFetch("/api/captcha/request", {
        params: {
          captchaId: form.captchaId,
        },
      });
    });

    form.captchaId = captcha.value.captchaId;
    form.captchaUrl = captcha.value.captchaUrl;
    form.captchaCode = "";
  } catch (e) {
    useCatchError(e);
  }
}

async function signup() {
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
    await refreshCaptcha();
  }
}

function toSignin() {
  if (form.redirect) {
    useLinkTo(`/user/signin?redirect=${encodeURIComponent(form.redirect)}`);
  } else {
    useLinkTo("/user/signin");
  }
}
</script>

<style lang="scss" scoped></style>
