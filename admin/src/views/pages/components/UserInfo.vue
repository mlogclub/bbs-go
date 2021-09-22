<template>
  <div class="user-info-component">
    <avatar :user="user" />
    <div class="user-info-component-right">
      <div>{{ user.id }}</div>
      <span>{{ user.nickname }}</span>
    </div>
  </div>
</template>

<script>
import Avatar from "@/components/Avatar";
export default {
  components: { Avatar },
  props: {
    user: {
      type: Object,
      required: true,
      default() {
        return null;
      },
    },
  },
};
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
