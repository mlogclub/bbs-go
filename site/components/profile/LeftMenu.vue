<template>
  <div class="profile-edit-tabs">
    <div class="profile-edit-tab-item">
      <nuxt-link to="/user/profile" :class="{ active: active === 'edit' }">
        <i class="iconfont icon-username" />
        <span>个人资料</span>
      </nuxt-link>
    </div>
    <div class="profile-edit-tab-item">
      <nuxt-link
        to="/user/profile/account"
        :class="{ active: active === 'account' }"
      >
        <i class="iconfont icon-username" />
        <span>账号设置</span>
      </nuxt-link>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    active: {
      type: String,
      default: 'edit',
    },
  },
}
</script>

<style lang="scss" scoped>
.profile-edit-tabs {
  background-color: #fff;
  padding: 10px;

  .profile-edit-tab-item {
    width: 100%;

    a {
      display: block;
      padding: 10px;

      &:hover,
      &.active {
        background: #e8f3ff;
        color: #1d7dfa;
      }
    }
  }
}
</style>
