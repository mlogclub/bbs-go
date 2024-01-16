<template>
  <section class="main">
    <div class="container">
      <div class="main-body no-bg">
        <article
          class="message"
          :class="{ 'is-success': data.success, 'is-warning': !data.success }"
        >
          <div class="message-header">
            <p>邮箱验证</p>
          </div>
          <div class="message-body">
            <div v-if="data.success">
              恭喜，邮箱验证成功。你的邮箱为：{{ data.email }}
            </div>
            <div v-else>
              邮箱验证失败<span v-if="data.message" class="has-text-danger"
                >&nbsp;原因：{{ data.message }}</span
              >，请前往&nbsp;<nuxt-link
                to="/user/profile"
                style="font-weight: 700"
                >个人资料 &gt; 账号设置</nuxt-link
              >&nbsp;页面尝试重新发送验证邮件。
            </div>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>

<script setup>
const data = reactive({
  success: false,
  email: "",
  message: "",
});
const route = useRoute();
try {
  const resp = await useMyFetchPost(
    `/api/user/verify_email?token=${route.query.token}`
  );
  data.success = true;
  data.email = resp.email;
} catch (e) {
  data.success = false;
  data.message = e.message || "";
}
</script>

<style lang="scss" scoped></style>
