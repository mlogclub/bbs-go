<template>
  <div class="mobile-sidebar">
    <transition name="fadeLeft">
      <div v-show="show" class="sidebar-container">
        <div class="sidebar-message">
          <i class="iconfont icon-bell" />
          <nuxt-link to="/user/messages">消息</nuxt-link>
        </div>
        <template v-if="user">
          <div class="sidebar-userinfo">
            <i class="iconfont icon-username" />
            <span>{{ user.nickname }}</span>
          </div>
          <div class="sidebar-menus">
            <div class="sidebar-menu-item">
              <nuxt-link :to="'/user/' + user.id">个人中心</nuxt-link>
            </div>
            <div class="sidebar-menu-item">
              <nuxt-link class="sidebar-menu-item" to="/user/favorites"
                >我的收藏</nuxt-link
              >
            </div>
            <div class="sidebar-menu-item">
              <nuxt-link class="sidebar-menu-item" to="/user/settings"
                >编辑资料</nuxt-link
              >
            </div>
            <div class="sidebar-menu-item">
              <nuxt-link
                v-if="isOwnerOrAdmin"
                class="sidebar-menu-item"
                to="/admin"
                >后台管理</nuxt-link
              >
            </div>
          </div>
        </template>
        <template v-else>
          <nuxt-link
            class="sidebar-login-btn button is-primary"
            to="/user/signin"
            >登录
          </nuxt-link>
        </template>
      </div>
    </transition>
  </div>
</template>

<script>
import UserHelper from '~/common/UserHelper'
export default {
  computed: {
    show() {
      return this.$store.state.env.showMobileSidebar
    },
    user() {
      return this.$store.state.user.current
    },
    isOwnerOrAdmin() {
      return UserHelper.isOwner(this.user) || UserHelper.isAdmin(this.user)
    },
    config() {
      return this.$store.state.config.config
    },
  },
  methods: {
    async signout() {
      try {
        await this.$store.dispatch('user/signout')
        this.$linkTo('/')
      } catch (e) {
        console.error(e)
      }
    },
  },
}
</script>
<style lang="scss" scoped></style>
