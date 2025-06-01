<template>
  <section class="main">
    <div class="container">
      <div class="main-body no-bg">
        <div class="widget signin">
          <div class="widget-content">
            <div class="tabs is-centered">
              <ul>
                <li class="is-active">
                  <a @click="switchLoginMethod('password')">账号登录</a>
                </li>
              </ul>
            </div>

            <div>
              <login-password />
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
</script>

<style lang="scss" scoped>
.signin {
  max-width: 500px;
  margin: auto !important;
  padding: 20px;
}
</style>
