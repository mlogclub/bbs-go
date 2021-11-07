<template>
  <div class="widget">
    <div class="widget-header">
      <span>关注</span>
      <span class="count">{{ user.followCount }}</span>
    </div>
    <div class="widget-content">
      <div class="widget-tips">没有更多内容了</div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
}
</script>

<style lang="scss" scoped></style>
