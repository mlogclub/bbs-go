<template>
  <section class="main">
    <div class="container">
      <div class="main-body no-bg">
        <div class="widget forget-password">
          <div class="widget-header">
            <h3 class="title">{{ $t("user.forgetPassword.title") }}</h3>
          </div>
          <div class="widget-content">
            <div class="form-container">
              <div class="form-group">
                <label for="email">{{ $t("user.forgetPassword.emailLabel") }}</label>
                <input
                  id="email"
                  v-model="email"
                  type="email"
                  :placeholder="$t('user.forgetPassword.emailPlaceholder')"
                  class="form-control"
                />
              </div>
              <div class="form-actions">
                <el-button type="primary" @click="submit" :loading="loading">
                  {{ $t("user.forgetPassword.sendResetLinkBtn") }}
                </el-button>
              </div>
              <div class="form-bottom">
                <nuxt-link to="/user/signin">{{ $t("user.forgetPassword.backToLogin") }}</nuxt-link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
const { t } = useI18n();
useHead({
  title: useSiteTitle(t("user.forgetPassword.title")),
});

const email = ref("");
const loading = ref(false);

const submit = async () => {
  if (!email.value) {
    useMsgError(t("user.forgetPassword.emailRequired"));
    return;
  }

  loading.value = true;
  try {
    const response = await useHttpPost('/api/forget-password/send/email', {
      email: email.value
    });

    useMsgSuccess(t("user.forgetPassword.resetLinkSent"));
    // 清空邮箱输入框
    email.value = "";
  } catch (error) {
    useCatchError(error);
  } finally {
    loading.value = false;
  }
};
</script>

<style lang="scss" scoped>
.forget-password {
  max-width: 500px;
  margin: 2rem auto;
  padding: 2rem;
  background-color: var(--bg-color);
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);

  .widget-header {
    text-align: center;
    margin-bottom: 1.5rem;

    .title {
      font-size: 1.5rem;
      color: var(--text-color);
      margin: 0;
    }
  }

  .form-container {
    .form-group {
      margin-bottom: 1.5rem;

      label {
        display: block;
        margin-bottom: 0.5rem;
        font-weight: 500;
        color: var(--text-color);
      }

      input {
        width: 100%;
        padding: 0.75rem;
        border: 1px solid var(--border-color);
        border-radius: 4px;
        background-color: var(--bg-color2);
        color: var(--text-color);

        &:focus {
          outline: none;
          border-color: var(--primary-color);
          background-color: var(--bg-color3);
        }
      }
    }

    .form-actions {
      margin: 1.5rem 0;

      button {
        width: 100%;
        padding: 0.75rem;
        font-size: 1rem;
      }
    }

    .form-bottom {
      text-align: center;
      font-size: 0.9rem;

      a {
        color: var(--primary-color);
        text-decoration: none;

        &:hover {
          text-decoration: underline;
        }
      }
    }
  }
}
</style>