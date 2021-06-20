<template>
  <section class="main">
    <div class="container">
      <div class="error">
        <div>
          <img src="~/assets/images/logo.png" style="max-width: 100px" />
        </div>
        <div class="description">
          <span v-if="error.statusCode === 404">{{
            error.message ? error.message : '页面没找到'
          }}</span>
          <span v-if="error.statusCode === 403">{{
            error.message ? error.message : 'forbidden'
          }}</span>
          <span v-else>{{ error.statusCode }} 页面异常</span>
        </div>
        <div class="report">
          <nuxt-link to="/topic/create" target="_blank"
            >点击这里反馈该问题>></nuxt-link
          >
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  layout: 'no-footer', // 你可以为错误页面指定自定义的布局
  props: {
    error: {
      type: Object,
      default: null,
    },
  },
}
</script>

<style lang="scss" scoped>
.error {
  text-align: center;
  vertical-align: center;
  padding: 100px 0;

  .description {
    margin-top: 30px;
    span {
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
      font-weight: 700;
    }
  }
}
</style>
