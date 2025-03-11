<script setup>
const configStore = useConfigStore();

defineProps({
  error: {
    type: Object,
    default: null,
  },
});

definePageMeta({
  layout: "default",
});

const handleError = () => {
  clearError({ redirect: "/" });
};
</script>

<template>
  <div>
    <MyHeader />
    <section class="main">
      <div class="container">
        <div class="error">
          <div>
            <img
              v-if="configStore.config.siteLogo"
              :src="configStore.config.siteLogo"
              style="max-width: 100px"
            />
            <img
              v-else
              src="~/assets/images/logo.png"
              style="max-width: 100px"
            />
          </div>
          <div class="description">
            <div v-if="error.message">
              {{ error.message }}
            </div>

            <template v-else>
              <div v-if="error.statusCode === 404">页面没找到</div>
              <div v-if="error.statusCode === 403">Forbidden</div>
              <div v-else>{{ error.statusCode }} 异常</div>
            </template>
          </div>
          <div class="report">
            <a @click="handleError">返回首页</a>
          </div>
        </div>
      </div>
    </section>
    <MyFooter />
  </div>
</template>

<style lang="scss" scoped>
.error {
  text-align: center;
  vertical-align: center;
  padding: 100px 0;

  .description {
    margin-top: 30px;
    div {
      font-size: 18px;
      font-weight: bold;
      line-height: 22px;
      color: rgb(230, 76, 76);
    }
  }

  .report {
    margin-top: 20px;
    a {
      font-size: 15px;
      font-weight: 500;
    }
  }
}
</style>
