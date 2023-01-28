<template>
  <div class="mobile-sidebar">
    <transition name="fadeLeft">
      <div v-show="show" class="sidebar-container">
        <div v-if="siteNavs && siteNavs.length" class="sidebar-navs">
          <div
            v-for="(nav, index) in siteNavs"
            :key="index"
            class="sidebar-nav-item"
          >
            <i class="iconfont icon-nav" />
            <nuxt-link :to="nav.url">{{ nav.title }}</nuxt-link>
          </div>
        </div>
        <div class="sidebar-message">
          <i class="iconfont icon-message" />
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
              <nuxt-link class="sidebar-menu-item" to="/user/profile"
                >编辑资料</nuxt-link
              >
            </div>
            <div
              v-if="checkIn == null || checkIn?.checkIn == false"
              class="sidebar-menu-item"
            >
              <a @click="doCheckIn">
                <span>立即签到</span>
              </a>
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
  data() {
    return {
      checkIn: null,
    }
  },
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
    siteNavs() {
      const config = this.$store.state.config.config
      return config.siteNavs || []
    },
    isLogin() {
      return this.$store.state.user.current != null
    },
  },
  mounted() {
    this.getCheckIn()
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
    async getCheckIn() {
      try {
        this.checkIn = await this.$axios.get('/api/checkin/checkin')
      } catch (e) {
        console.log(e)
      }
    },
    async doCheckIn() {
      if (!this.isLogin) {
        this.$toSignin()
      }
      try {
        await this.$axios.post('/api/checkin/checkin')
        this.$message.success('签到成功')
        await this.getCheckIn()
      } catch (e) {
        console.error(e)
      }
    },
  },
}
</script>
<style lang="scss" scoped></style>
