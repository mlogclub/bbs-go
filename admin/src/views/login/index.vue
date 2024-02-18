<template>
  <div class="login-container">
    <div class="logo">
      <img
        alt="logo"
        src="//p3-armor.byteimg.com/tos-cn-i-49unhts6dw/dfdba5317c0c20ce20e64fac803d52bc.svg~tplv-49unhts6dw-image.image"
      />
      <div class="logo-text">BBS-GO</div>
    </div>
    <LoginBanner />
    <div class="content">
      <div class="content-inner">
        <a-spin :loading="loading" tip="Loading">
          <LoginForm />
        </a-spin>
      </div>
      <div class="footer">
        <Footer />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import Footer from '@/components/footer/index.vue';
  import { DEFAULT_ROUTE_NAME } from '@/router/constants';
  import LoginBanner from './components/banner.vue';
  import LoginForm from './components/login-form.vue';

  const userStore = useUserStore();
  const router = useRouter();
  const loading = ref(true);

  onMounted(async () => {
    try {
      const user = await userStore.info();
      if (user) {
        const { redirect, ...othersQuery } = router.currentRoute.value.query;
        router.push({
          name: (redirect as string) || DEFAULT_ROUTE_NAME,
          query: {
            ...othersQuery,
          },
        });
      }
    } finally {
      loading.value = false;
    }
  });
</script>

<style lang="less" scoped>
  .login-container {
    display: flex;
    height: 100vh;

    .banner {
      width: 550px;
      background: linear-gradient(163.85deg, #1d2129 0%, #00308f 100%);
    }

    .content {
      position: relative;
      display: flex;
      flex: 1;
      align-items: center;
      justify-content: center;
      padding-bottom: 40px;
    }

    .footer {
      position: absolute;
      right: 0;
      bottom: 0;
      width: 100%;
    }
  }

  .logo {
    position: fixed;
    top: 24px;
    left: 22px;
    z-index: 1;
    display: inline-flex;
    align-items: center;

    &-text {
      margin-right: 4px;
      margin-left: 4px;
      color: var(--color-fill-1);
      font-size: 20px;
    }
  }
</style>

<style lang="less" scoped>
  // responsive
  @media (max-width: @screen-lg) {
    .container {
      .banner {
        width: 25%;
      }
    }
  }
</style>
