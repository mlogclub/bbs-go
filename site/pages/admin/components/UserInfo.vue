<template>
  <div class="user-info-component">
    <img :src="user.smallAvatar" class="avatar" />
    <div class="user-info-component-right">
      <div>{{ user.id }}</div>
      <span>{{ user.nickname }}</span>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    user: {
      type: Object,
      required: true,
      default() {
        return null
      },
    },
  },
}
</script>

<style lang="scss" scoped>
.user-info-component {
  display: flex;

  .user-info-component-right {
    margin-left: 5px;
    padding: 3px 0;

    font-weight: 700;
    font-size: 15px;
  }
}
</style>
