<template>
  <section class="main">
    <div class="container">
      <div class="error">
        <div>
          <img src="~/assets/images/logo.png" style="max-width: 100px;">
        </div>
        <div class="description">
          <span v-if="error.statusCode === 404">
            {{ error.message ? error.message : '页面没找到' }}
          </span>
          <span v-else>
            页面异常
          </span>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  props: {
    error: {
      type: Object,
      default: null
    }
  },
  layout: 'default' // 你可以为错误页面指定自定义的布局
}
</script>

<style lang="scss" scoped>
.error {
  text-align: center;
  vertical-align: center;
  padding: 100px 0;

  .description span {
    font-size: 18px;
    font-weight: bold;
    line-height: 22px;
    color: rgb(230, 76, 76);
  }
}
</style>
