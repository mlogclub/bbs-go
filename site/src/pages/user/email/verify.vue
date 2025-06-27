<template>
  <section class="main">
    <div class="container">
      <div v-if="!data.loading" class="main-body no-bg">
        <article
          class="message"
          :class="{ 'is-success': data.success, 'is-warning': !data.success }"
        >
          <div class="message-header">
            <p>{{ $t("user.email.verify.title") }}</p>
          </div>
          <div class="message-body">
            <div v-if="data.success">
              {{ $t("user.email.verify.success", { email: data.email }) }}
            </div>
            <div v-else>
              {{ $t("user.email.verify.failed")
              }}<span v-if="data.message" class="has-text-danger"
                >&nbsp;{{
                  $t("user.email.verify.reason", { reason: data.message })
                }}</span
              >{{ $t("user.email.verify.retryInstructions") }}&nbsp;<nuxt-link
                to="/user/profile/account"
                style="font-weight: 700"
                >{{ $t("user.email.verify.accountSettings") }}</nuxt-link
              >&nbsp;{{ $t("user.email.verify.retryAction") }}
            </div>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>

<script setup>
const { t } = useI18n();

useHead({
  title: useSiteTitle(t("user.email.verify.title")),
});

const data = reactive({
  loading: true,
  success: false,
  email: "",
  message: "",
});

const verifyEmail = async () => {
  const route = useRoute();
  try {
    const resp = await useHttpPost(
      `/api/user/verify_email?token=${route.query.token}`
    );
    data.success = true;
    data.email = resp.email;
  } catch (e) {
    data.success = false;
    data.message = e.message || "";
  } finally {
    data.loading = false;
  }
};

onMounted(() => {
  verifyEmail();
});
</script>

<style lang="scss" scoped></style>
