<template>

  <section class="main">

    <div class="container">

      <div class="main-body no-bg">

        <div class="widget reset-password">

          <div class="widget-header">

            <h3 class="title">{{ $t("user.resetPassword.title") }}</h3>

          </div>

          <div class="widget-content">

            <form @submit.prevent class="form-container">

              <div class="form-group">

                <label for="password">{{ $t("user.resetPassword.passwordLabel") }}</label>

                <input

                  id="password"

                  v-model="password"

                  type="password"

                  :placeholder="$t('user.resetPassword.passwordPlaceholder')"

                  class="form-control"

                />

              </div>

              <div class="form-group">

                <label for="rePassword">{{ $t("user.resetPassword.rePasswordLabel") }}</label>

                <input

                  id="rePassword"

                  v-model="rePassword"

                  type="password"

                  :placeholder="$t('user.resetPassword.rePasswordPlaceholder')"

                  class="form-control"

                />

              </div>

              <div class="form-actions">

                <el-button type="primary" @click="submit" :loading="loading" native-type="button">

                  {{ $t("user.resetPassword.resetPasswordBtn") }}

                </el-button>

              </div>

              <div class="form-bottom">

                <nuxt-link to="/user/signin">{{ $t("user.resetPassword.backToLogin") }}</nuxt-link>

              </div>

            </form>

          </div>

        </div>

      </div>

    </div>

  </section>

</template>



<script setup>

const route = useRoute();

const { t } = useI18n();

useHead({

  title: useSiteTitle(t("user.resetPassword.title")),

});



const token = route.query.token;

const email = route.query.email;



const password = ref("");

const rePassword = ref("");

const loading = ref(false);



onMounted(() => {

  // 只检查token是否存在，email参数不是必需的

  if (!token) {

    useMsgError(t("user.resetPassword.invalidLink"));

    // 延迟导航以确保错误消息能显示

    setTimeout(() => {

      navigateTo("/user/forget-password");

    }, 2000);

  }

});

const submit = async () => {

  // 只检查token是否存在

  if (!token) {

    useMsgError(t("user.resetPassword.invalidLink"));

    // 延迟导航以确保错误消息能显示

    setTimeout(() => {

      navigateTo("/user/forget-password");

    }, 2000);

    return;

  }



  if (!password.value) {

    useMsgError(t("user.resetPassword.passwordRequired"));

    return;

  }

  

  if (password.value.length < 6) {

    useMsgError(t("user.resetPassword.passwordTooShort"));

    return;

  }

  

  if (password.value !== rePassword.value) {

    useMsgError(t("user.resetPassword.passwordMismatch"));

    return;

  }



  loading.value = true;

  try {

    const response = await useHttpPost('/api/forget-password/reset/password', {

      token: token,

      email: email,

      password: password.value,

      rePassword: rePassword.value

    });



    useMsgSuccess(t("user.resetPassword.resetSuccess"));

    // 重置成功后跳转到登录页

    setTimeout(() => {

      navigateTo("/user/signin");

    }, 2000);

  } catch (error) {

    useCatchError(error);

  } finally {

    loading.value = false;

  }

};
</script>

<style lang="scss" scoped>
.reset-password {
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