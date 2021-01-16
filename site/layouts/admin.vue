<template>
  <div class="bbs-admin">
    <nav
      class="navbar is-info admin-navbar"
      role="navigation"
      aria-label="main navigation"
    >
      <div class="navbar-brand">
        <div class="navbar-item">
          <img src="~/assets/images/logo.png" />
        </div>

        <a
          :class="{ 'is-active': navbarActive }"
          role="button"
          class="navbar-burger burger"
          aria-label="menu"
          aria-expanded="false"
          data-target="navbarBasic"
          @click="toggleNav"
        >
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
        </a>
      </div>

      <div
        id="navbarBasic"
        :class="{ 'is-active': navbarActive }"
        class="navbar-menu"
      >
        <div class="navbar-start">
          <a class="navbar-item" href="/admin">
            <i class="iconfont icon-dashboard"></i>
            <span>后台</span>
          </a>

          <div class="navbar-item has-dropdown is-hoverable">
            <a class="navbar-link">
              <i class="iconfont icon-topic"></i>
              <span>内容</span>
            </a>
            <div class="navbar-dropdown">
              <a class="navbar-item" href="/admin/topics">
                <i class="iconfont icon-topic"></i>
                <span>话题</span>
              </a>

              <a class="navbar-item" href="/admin/articles">
                <i class="iconfont icon-article"></i>
                <span>文章</span>
              </a>

              <a class="navbar-item" href="/admin/comments">
                <i class="iconfont icon-comments"></i>
                <span>评论</span>
              </a>
            </div>
          </div>

          <div class="navbar-item has-dropdown is-hoverable">
            <a class="navbar-link">
              <i class="iconfont icon-tags"></i>
              <span>节点</span>
            </a>
            <div class="navbar-dropdown">
              <a class="navbar-item" href="/admin/topics/nodes">
                <i class="iconfont icon-tag"></i>
                <span>节点</span>
              </a>
              <a class="navbar-item" href="/admin/tags">
                <i class="iconfont icon-tags"></i>
                <span>标签</span>
              </a>
            </div>
          </div>

          <a class="navbar-item" href="/admin/users">
            <i class="iconfont icon-username"></i>
            <span>用户</span>
          </a>
          <a class="navbar-item" href="/admin/links">
            <i class="iconfont icon-link"></i>
            <span>友链</span>
          </a>
          <a class="navbar-item" href="/admin/settings">
            <i class="iconfont icon-setting"></i>
            <span>设置</span>
          </a>
          <a class="navbar-item" href="/admin/operate-log">
            <i class="iconfont icon-log"></i>
            <span>操作日志</span>
          </a>
        </div>

        <div class="navbar-end">
          <div class="navbar-item">
            <div class="buttons">
              <a class="button is-success" href="/" target="_blank">
                <strong>查看网站</strong>
              </a>
            </div>
          </div>
        </div>
      </div>
    </nav>
    <nuxt />
  </div>
</template>

<script>
// import '~/plugins/element-ui'

export default {
  middleware: ['authenticated'],
  data() {
    return {
      navbarActive: false,
    }
  },
  methods: {
    toggleNav() {
      this.navbarActive = !this.navbarActive
    },
  },
  head() {
    return {
      title: 'BBS-GO管理',
      link: [
        // {
        //   rel: 'stylesheet',
        //   href: '//unpkg.com/element-ui/lib/theme-chalk/index.css',
        // },
      ],
    }
  },
}
</script>

<style lang="scss">
.bbs-admin {
  .admin-navbar {
    .navbar-item {
      font-weight: 700;
      .iconfont {
        margin-right: 3px;
      }
    }
  }
  .page-container {
    /*margin: 10px !important;*/
    /*padding: 10px !important;*/
    /*background: #fff;*/

    .page-section {
      background-color: #fff;
    }

    .toolbar {
      border-bottom: 1px solid #e9ecf0;
      margin: 0 0 5px 0;
      padding: 10px 10px 0 10px;
      background: #fff;

      .el-form-item {
        margin-bottom: 10px;
      }
    }

    .pagebar {
      background: #fff;
      text-align: right;
      margin-top: 5px;
      padding: 3px 0;
    }
  }
}
</style>
