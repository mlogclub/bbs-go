<template>
  <div>
    <MyHeader />

    <section class="main">
      <div class="container main-container right-main">
        <div class="left-container">
          <div class="profile-edit-tabs-pc">
            <div class="profile-edit-tab-item">
              <nuxt-link to="/user/profile">
                <i class="iconfont icon-username" />
                <span>{{ $t("layout.profile.profile") }}</span>
              </nuxt-link>
            </div>
            <div class="profile-edit-tab-item">
              <nuxt-link to="/user/profile/account">
                <i class="iconfont icon-setting" />
                <span>{{ $t("layout.profile.accountSettings") }}</span>
              </nuxt-link>
            </div>
          </div>
        </div>
        <div class="right-container">
          <div class="profile-edit-tabs-mobile tabs">
            <ul>
              <li :class="{ 'is-active': active === 'profile' }">
                <nuxt-link to="/user/profile">{{
                  $t("layout.profile.profile")
                }}</nuxt-link>
              </li>
              <li :class="{ 'is-active': active === 'account' }">
                <nuxt-link to="/user/profile/account">{{
                  $t("layout.profile.accountSettings")
                }}</nuxt-link>
              </li>
            </ul>
          </div>
          <slot />
        </div>
      </div>
    </section>

    <MyFooter />
  </div>
</template>

<script setup>
const route = useRoute();
const active = computed(() => {
  if (route.path.includes("/user/profile/account")) {
    return "account";
  }
  return "profile";
});
</script>

<style lang="scss" scoped>
.profile-edit-tabs-pc {
  border-radius: var(--border-radius);
  background-color: var(--bg-color);
  padding: 10px;

  .profile-edit-tab-item {
    width: 100%;

    a {
      padding: 10px;

      display: flex;
      align-items: center;
      column-gap: 6px;
      color: var(--text-color);

      &:hover,
      &.active,
      &.router-link-exact-active {
        background: var(--bg-color5);
        color: var(--text-link-color);
      }
    }
  }
}
.profile-edit-tabs-mobile {
  border-radius: var(--border-radius);
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
