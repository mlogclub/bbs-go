<template>
  <nav
    ref="nav"
    class="navbar has-shadow"
    role="navigation"
    aria-label="main navigation"
  >
    <div class="container">
      <div class="navbar-brand">
        <nuxt-link to="/" class="navbar-item">
          <img :alt="config.siteTitle" src="~/assets/images/logo.png" />
        </nuxt-link>
        <a
          :class="{ 'is-active': navbarActive }"
          class="navbar-burger burger"
          data-target="navbarBasic"
          @click="toggleNav"
        >
          <span aria-hidden="true" />
          <span aria-hidden="true" />
          <span aria-hidden="true" />
        </a>
      </div>
      <div :class="{ 'is-active': navbarActive }" class="navbar-menu">
        <div class="navbar-start">
          <nuxt-link
            v-for="(nav, index) in config.siteNavs"
            :key="index"
            :to="nav.url"
            class="navbar-item"
            >{{ nav.title }}</nuxt-link
          >
        </div>

        <div class="navbar-end">
          <div class="navbar-item">
            <search-input />
          </div>

          <div class="navbar-item">
            <create-topic-btn />
          </div>

          <msg-notice v-if="user" />

          <div
            v-if="user"
            class="navbar-item has-dropdown is-hoverable user-menus"
          >
            <nuxt-link :to="'/user/' + user.id" class="navbar-link">
              <no-ssr><avatar :user="user" :size="30" /></no-ssr>
              <span class="user-menus-nickname">{{ user.nickname }}</span>
            </nuxt-link>
            <div class="navbar-dropdown">
              <nuxt-link class="navbar-item" :to="'/user/' + user.id">
                <i class="iconfont icon-username" />&nbsp;个人中心
              </nuxt-link>
              <nuxt-link class="navbar-item" to="/user/favorites">
                <i class="iconfont icon-favorites" />&nbsp;我的收藏
              </nuxt-link>
              <nuxt-link class="navbar-item" to="/user/profile">
                <i class="iconfont icon-username" />&nbsp;编辑资料
              </nuxt-link>
              <a class="navbar-item" @click="signout">
                <i class="iconfont icon-log-out" />&nbsp;退出登录
              </a>
            </div>
          </div>
          <div v-else class="navbar-item">
            <div class="buttons">
              <nuxt-link class="button login-btn" to="/user/signin"
                >登录
              </nuxt-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script>
import UserHelper from '~/common/UserHelper'

export default {
  data() {
    return {
      navbarActive: false,
    }
  },
  computed: {
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
    toggleNav() {
      this.navbarActive = !this.navbarActive
    },
  },
}
</script>

<style lang="scss" scoped>
.navbar {
  /*opacity: 0.99;*/
  /*border-bottom: 1px solid #e7edf3;*/
  background-color: var(--bg-color);

  .navbar-item {
    font-weight: 700;
  }

  .publish {
    color: var(--text-color);
    background-color: #3174dc;
    width: 100px;
    &:hover {
      color: var(--text-color);
      background-color: #4d91fa;
    }
  }

  .login-btn {
    border-color: #000; // TODO
    &:hover {
      color: var(--text-color3);
      border-color: var(--text-color3);
    }
  }
}

.user-menus {
  .user-menus-nickname {
    margin-left: 5px;
  }
}
</style>
