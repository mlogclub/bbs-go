<template>
  <div class="login-form-wrapper">
    <div class="login-form-title">{{ $t('pages.login.form.title') }}</div>
    <div class="login-form-sub-title">{{
      $t('pages.login.form.subtitle')
    }}</div>
    <div class="login-form-error-msg">{{ errorMessage }}</div>
    <a-form
      ref="loginForm"
      :model="form"
      class="login-form"
      layout="vertical"
      @submit="handleSubmit"
    >
      <a-form-item
        field="username"
        :rules="[
          { required: true, message: $t('pages.login.form.userName.errMsg') },
        ]"
        :validate-trigger="['change', 'blur']"
        hide-label
      >
        <a-input
          v-model="form.username"
          :placeholder="$t('pages.login.form.userName.placeholder')"
        >
          <template #prefix>
            <icon-user />
          </template>
        </a-input>
      </a-form-item>
      <a-form-item
        field="password"
        :rules="[
          { required: true, message: $t('pages.login.form.password.errMsg') },
        ]"
        :validate-trigger="['change', 'blur']"
        hide-label
      >
        <a-input-password
          v-model="form.password"
          :placeholder="$t('pages.login.form.password.placeholder')"
          allow-clear
        >
          <template #prefix>
            <icon-lock />
          </template>
        </a-input-password>
      </a-form-item>

      <a-form-item
        v-if="captchaBase64"
        field="captchaCode"
        :rules="[
          {
            required: true,
            message: $t('pages.login.form.captchaCode.errMsg'),
          },
        ]"
        :validate-trigger="['change', 'blur']"
        hide-label
      >
        <a-input
          v-model="form.captchaCode"
          :placeholder="$t('pages.login.form.captchaCode.placeholder')"
          allow-clear
          class="captcha-code"
        >
          <template #prefix>
            <icon-copy />
          </template>
          <template #append>
            <img
              :src="`data:image/png;base64,${captchaBase64}`"
              @click="refreshCaptcha"
            />
          </template>
        </a-input>
      </a-form-item>

      <a-space :size="16" direction="vertical">
        <a-button type="primary" html-type="submit" long :loading="loading">
          {{ $t('pages.login.form.login.submit') }}
        </a-button>
      </a-space>
    </a-form>
  </div>
</template>

<script lang="ts" setup>
  import { Message } from '@arco-design/web-vue';
  import { ValidatedError } from '@arco-design/web-vue/es/form/interface';
  import { useI18n } from 'vue-i18n';
  import { DEFAULT_ROUTE_NAME } from '@/router/constants';
  import useLoading from '@/hooks/loading';
  import { type LoginData } from '@/api/user';

  const router = useRouter();
  const { t } = useI18n();
  const errorMessage = ref('');
  const { loading, setLoading } = useLoading();
  const userStore = useUserStore();

  const form = reactive({
    username: '',
    password: '',
    captchaId: '',
    captchaCode: '',
  });

  const captchaBase64 = ref('');

  interface User {
    captchaId: string;
    captchaBase64: string;
  }

  const refreshCaptcha = async () => {
    const ret = await axios.get<any, User>('/api/captcha/request');
    form.captchaId = ret.captchaId;
    form.captchaCode = '';

    captchaBase64.value = ret.captchaBase64;
  };

  refreshCaptcha();

  const handleSubmit = async ({
    errors,
  }: {
    errors: Record<string, ValidatedError> | undefined;
  }) => {
    if (loading.value) return;
    if (!errors) {
      setLoading(true);
      try {
        await userStore.login(form as LoginData);
        const { redirect, ...othersQuery } = router.currentRoute.value.query;
        router.push({
          name: (redirect as string) || DEFAULT_ROUTE_NAME,
          query: {
            ...othersQuery,
          },
        });
        Message.success(t('pages.login.form.login.success'));
      } catch (err) {
        refreshCaptcha();
        errorMessage.value = (err as Error).message;
      } finally {
        setLoading(false);
      }
    }
  };
</script>

<style lang="less" scoped>
  .login-form {
    &-wrapper {
      width: 320px;
    }

    &-title {
      color: var(--color-text-1);
      font-weight: 500;
      font-size: 24px;
      line-height: 32px;
    }

    &-sub-title {
      color: var(--color-text-3);
      font-size: 16px;
      line-height: 24px;
    }

    &-error-msg {
      height: 32px;
      color: rgb(var(--red-6));
      line-height: 32px;
    }

    &-password-actions {
      display: flex;
      justify-content: space-between;
    }

    &-register-btn {
      color: var(--color-text-3) !important;
    }
  }

  .captcha-code {
    img {
      height: 30px;
      cursor: pointer;
    }
  }
</style>
