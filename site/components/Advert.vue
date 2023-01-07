<template>
  <div class="widget">
    <div class="widget-header">广告</div>
    <div class="widget-content content notice-content">
      <el-carousel height="150px" indicator-position="none">
        <el-carousel-item v-for="item in adverts.results" :key="item.index">
          <a :href="item.url"><img :src="item.picUrl" /></a>
        </el-carousel-item>
      </el-carousel>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Advert',
  props: {
    adverts: {
      type: Object,
      default() {
        return null
      },
    },
  },
}
</script>

<style scoped lang="scss">
.el-carousel__indicators {
  margin-left: 0;
}
</style>
