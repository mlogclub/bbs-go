<template>
  <nav
    ref="nav"
    class="navbar has-shadow"
    role="navigation"
    aria-label="main navigation"
  >
    <div class="container">
      <div class="navbar-brand">
        <a href="/" class="navbar-item">
          <img :alt="config.siteTitle" src="~/assets/images/logo.png" />
        </a>
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
          <a
            v-for="(nav, index) in config.siteNavs"
            :key="index"
            :href="nav.url"
            class="navbar-item"
            >{{ nav.title }}</a
          >
        </div>

        <div class="navbar-end">
          <div class="navbar-item searchFormDiv">
            <form id="searchForm" action="/search">
              <div class="control has-icons-right">
                <input
                  name="q"
                  class="input"
                  type="text"
                  maxlength="30"
                  placeholder="搜索"
                />
                <span class="icon is-medium is-right">
                  <i class="iconfont icon-search" />
                </span>
              </div>
            </form>
          </div>

          <div class="navbar-item">
            <create-topic-btn />
          </div>

          <msg-notice v-if="user" />

          <div v-if="user" class="navbar-item has-dropdown is-hoverable">
            <a :href="'/user/' + user.id" class="navbar-link">
              <strong>{{ user.nickname }}</strong>
            </a>
            <div class="navbar-dropdown">
              <a class="navbar-item" href="/user/favorites">
                <i class="iconfont icon-favorites" />&nbsp;我的收藏
              </a>
              <a class="navbar-item" href="/user/settings">
                <i class="iconfont icon-username" />&nbsp;编辑资料
              </a>
              <a v-if="isOwnerOrAdmin" class="navbar-item" href="/admin">
                <i class="iconfont icon-dashboard" />&nbsp;后台管理
              </a>
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
import MsgNotice from '~/components/MsgNotice'
import CreateTopicBtn from '~/components/topic/CreateTopicBtn'

export default {
  components: {
    MsgNotice,
    CreateTopicBtn,
  },
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

  .navbar-item {
    font-weight: 700;
  }

  .publish {
    color: #fff;
    background-color: #3174dc;
    width: 100px;
    &:hover {
      color: #fff;
      background-color: #4d91fa;
    }
  }

  .login-btn {
    //border-width: 2px;
    border-color: #000;
    &:hover {
      color: #7e7e7e;
      border-color: #7e7e7e;
    }
  }
}

.searchFormDiv {
  @media screen and (max-width: 1024px) {
    display: none;
  }
  #searchForm {
    .input {
      box-shadow: none;
      border-radius: 2px;
      background-color: #fff;
      transition: all 0.4s;
      float: right;
      position: relative;

      &:focus {
        background-color: #fff;
        border-color: #e7672e;
        outline: none;
      }
    }
  }
}
</style>
