<template>
  <nav
    ref="nav"
    class="navbar is-white is-fixed-top"
    role="navigation"
    aria-label="main navigation"
  >
    <div class="container">
      <div class="navbar-brand">
        <a href="/" class="navbar-item">
          <img src="~/assets/images/logo.png" />
        </a>
        <a
          :class="{ 'is-active': navbarActive }"
          @click="toggleNav"
          class="navbar-burger burger"
          data-target="navbarBasic"
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
          <div class="navbar-item">
            <form
              id="searchForm"
              action="https://www.google.com/search"
              target="_blank"
            >
              <div class="control has-icons-right">
                <input name="q" type="hidden" value="site:mlog.club" />
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

          <!--
          <div class="navbar-item">
            <div class="buttons">
              <a class="button is-success" href="/topic/create">
                <i class="iconfont icon-topic">
                  <strong>发帖/提问</strong>
                </i>
              </a>
            </div>
          </div>
          -->

          <div v-if="user" class="navbar-item has-dropdown is-hoverable">
            <a :href="'/user/' + user.id" class="navbar-link">
              <strong>{{ user.nickname }}</strong>
            </a>
            <div class="navbar-dropdown">
              <a class="navbar-item" href="/topic/create">
                <i class="iconfont icon-topic" />&nbsp;发帖/提问
              </a>
              <a class="navbar-item" href="/article/create">
                <i class="iconfont icon-publish" />&nbsp;发表文章
              </a>
              <a class="navbar-item" href="/user/messages">
                <i class="iconfont icon-message" />&nbsp;消息
              </a>
              <a class="navbar-item" href="/user/favorites">
                <i class="iconfont icon-favorites" />&nbsp;收藏
              </a>
              <a class="navbar-item" href="/user/settings">
                <i class="iconfont icon-username" />&nbsp;编辑资料
              </a>
              <a @click="signout" class="navbar-item">
                <i class="iconfont icon-log-out" />&nbsp;退出登录
              </a>
            </div>
          </div>
          <div v-if="!user" class="navbar-item">
            <div class="buttons">
              <github-login />
              <qq-login />
            </div>
          </div>

          <div v-if="user && msgcount > 0" class="navbar-item">
            <a :href="'/user/messages'" class="msgcount">{{
              msgcount > 9 ? '9+' : msgcount
            }}</a>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script>
import utils from '~/common/utils'
import GithubLogin from '~/components/GithubLogin'
import QqLogin from '~/components/QqLogin'
export default {
  components: { GithubLogin, QqLogin },
  data() {
    return {
      msgcount: 0,
      navbarActive: false
    }
  },
  computed: {
    user() {
      return this.$store.state.user.current
    },
    config() {
      return this.$store.state.config.config
    }
  },
  mounted() {
    window.addEventListener('scroll', this.handleScroll)
    this.getMsgcount()
  },
  methods: {
    async getMsgcount() {
      this.msgcount = await this.$axios.get('/api/user/msgcount')
    },
    async signout() {
      try {
        await this.$store.dispatch('user/signout')
        this.msgcount = 0
        utils.linkTo('/')
      } catch (e) {
        console.error(e)
      }
    },
    toggleNav() {
      this.navbarActive = !this.navbarActive
    },
    handleScroll() {
      if (window.scrollY > 0) {
        this.$refs.nav.classList.add('scrolled')
      } else {
        this.$refs.nav.classList.remove('scrolled')
      }
    }
  }
}
</script>

<style lang="scss" scoped>
#searchForm {
  .input {
    // box-shadow: inset 0 1px 2px rgba(10,10,10,.1);
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
.navbar {
  opacity: 0.99;
  border-bottom: 1px solid #e7edf3;
  &.scrolled {
    box-shadow: 1px 0px 6px rgba(0, 0, 0, 0.25);
    border-bottom: none;
  }
  .navbar-item {
    font-weight: 600;
    &:hover,
    &.active {
      color: #009a61 !important;
    }
  }
}
</style>
