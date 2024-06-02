<template>
  <section class="main">
    <div class="container">
      <div class="main-body no-bg">
        <div class="widget signin">
          <div class="widget-header">登录</div>
          <div class="widget-content">
            <div class="field">
              <label class="label">用户名/邮箱</label>
              <div class="control has-icons-left">
                <input
                  v-model="form.username"
                  class="input is-success"
                  type="text"
                  placeholder="请输入用户名或邮箱"
                  @keyup.enter="signin"
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
                  @keyup.enter="signin"
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
                      @keyup.enter="signin"
                    />
                    <span class="icon is-small is-left"
                      ><i class="iconfont icon-captcha"
                    /></span>
                  </div>
                  <div
                    v-if="form.captchaUrl"
                    class="field login-captcha-img"
                    @click="refreshCaptcha"
                  >
                    <img :src="form.captchaUrl" data-not-lazy />
                  </div>
                </div>
              </div>
            </div>

            <div class="field">
              <button class="button is-link" @click="signin">登录</button>
              <a class="button is-text" @click="toSignup">
                没有账号？点击这里去注册&gt;&gt;
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
useHead({
  title: useSiteTitle("登录"),
});

const route = useRoute();
const form = reactive({
  username: "",
  password: "",
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

async function signin() {
  try {
    if (!form.username) {
      useMsgError("请输入用户名或邮箱");
      return;
    }
    if (!form.password) {
      useMsgError("请输入密码");
      return;
    }
    if (!form.captchaCode) {
      useMsgError("请输入验证码");
      return;
    }

    const userStore = useUserStore();
    const { user, redirect } = await userStore.signin(form);
    if (redirect) {
      useLinkTo(redirect);
    } else {
      useLinkTo(`/user/${user.id}`);
    }
  } catch (e) {
    useCatchError(e);
    await refreshCaptcha();
  }
}

function toSignup() {
  if (form.redirect) {
    useLinkTo(`/user/signup?redirect=${encodeURIComponent(form.redirect)}`);
  } else {
    useLinkTo("/user/signup");
  }
}
</script>

<style lang="scss" scoped></style>
