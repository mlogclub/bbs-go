<template>
  <section class="main">
    <div class="container">
      <div class="main-body no-bg">
        <div class="widget signin">
          <div class="widget-content">
            <div class="tabs is-centered">
              <ul>
                <li :class="{ 'is-active': loginMethod === 'password' }">
                  <a @click="switchLoginMethod('password')">账号登录</a>
                </li>
              </ul>
            </div>

            <div v-show="loginMethod === 'weixin'">
              <login-weixin ref="weixinLoginRef" />
            </div>
            <div v-show="loginMethod === 'password'">
              <login-password ref="passwordLoginRef" />
            </div>

            <div class></div>
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

const loginMethod = ref(null);
const weixinLoginRef = ref("password");
const passwordLoginRef = ref(null);

const switchLoginMethod = (method) => {
  loginMethod.value = method;
};

onMounted(() => {
  switchLoginMethod("password");
});
</script>

<style lang="scss" scoped>
.signin {
  max-width: 500px;
  margin: auto !important;
  padding: 20px;
}
</style>
