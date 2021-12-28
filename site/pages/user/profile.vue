<template>
  <section class="main">
    <div class="container main-container right-main">
      <div class="left-container">
        <div class="profile-edit-tabs-pc">
          <div class="profile-edit-tab-item">
            <nuxt-link to="/user/profile">
              <i class="iconfont icon-username" />
              <span>个人资料</span>
            </nuxt-link>
          </div>
          <div class="profile-edit-tab-item">
            <nuxt-link to="/user/profile/account">
              <i class="iconfont icon-username" />
              <span>账号设置</span>
            </nuxt-link>
          </div>
        </div>
      </div>
      <div class="right-container">
        <div class="profile-edit-tabs-mobile tabs">
          <ul>
            <li :class="{ 'is-active': active === 'profile' }">
              <nuxt-link to="/user/profile">个人资料</nuxt-link>
            </li>
            <li :class="{ 'is-active': active === 'account' }">
              <nuxt-link to="/user/profile/account">账号设置</nuxt-link>
            </li>
          </ul>
        </div>
        <nuxt-child />
      </div>
    </div>
  </section>
</template>

<script>
export default {
  middleware: 'authenticated',
  computed: {
    active() {
      if (this.$route.path.includes('/user/profile/account')) {
        return 'account'
      }
      return 'profile'
    },
  },
}
</script>

<style lang="scss" scoped>
.profile-edit-tabs-pc {
  background-color: var(--bg-color);
  padding: 10px;

  .profile-edit-tab-item {
    width: 100%;

    a {
      display: block;
      padding: 10px;

      &:hover,
      &.active,
      &.nuxt-link-exact-active {
        background: #e8f3ff;
        color: var(--text-link-color);
      }
    }
  }
}
.profile-edit-tabs-mobile {
  background-color: var(--bg-color);
  margin-bottom: 10px !important;
  display: none;
}

@media screen and (max-width: 1024px) {
  .profile-edit-tabs-mobile {
    display: block;
  }
}
</style>
