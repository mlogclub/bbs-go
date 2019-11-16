<template>
  <el-col :span="24" class="header">
    <el-col :span="10" class="logo" :class="{ collapsed: collapsed }">
      <template v-if="collapsed">
        <div class="logo-img" />
      </template>
      <template v-else>
        <div class="logo-name">
          <img src="../assets/logo.png" />
          <span>{{ sysName }}</span>
        </div>
      </template>
    </el-col>
    <el-col :span="10">
      <div class="tools" @click.prevent="collapse">
        <i class="iconfont" :class="collapsed ? 'icon-right' : 'icon-left'"></i>
      </div>
    </el-col>
    <el-col :span="4" class="userinfo" v-if="userInfo">
      <el-dropdown trigger="hover">
        <span class="el-dropdown-link userinfo-inner">
          <img :src="userInfo.avatar" />
          {{ userInfo.nickname }}
        </span>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item @click.native="logout">退出登录</el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </el-col>
  </el-col>
</template>

<script>
import cookies from 'js-cookie'

export default {
  name: 'Header',
  data() {
    return {
      sysName: 'BBS-GO'
    }
  },
  methods: {
    // 折叠导航栏
    collapse() {
      this.$store.dispatch('Default/collapse')
    },
    // 退出登录
    logout() {
      this.$message({ message: '暂未实现', type: 'success' })
    }
  },
  mounted() {
    let user = sessionStorage.getItem('user')
    if (user) {
      user = JSON.parse(user)
    }
  },
  computed: {
    collapsed() {
      return this.$store.state.Default.collapsed
    },
    userInfo() {
      let { userInfo } = this.$store.state.Login
      if (!userInfo) {
        const userInfoStr = cookies.get('userInfo')
        if (userInfoStr) {
          try {
            userInfo = JSON.parse(userInfoStr)
            this.$store.dispatch('Login/setUserInfo', userInfo)
          } catch (e) {
            console.error(e)
          }
        }
      }
      return userInfo
    }
  }
}
</script>

<style scoped lang="scss">
@import '../styles/vars.scss';

.header {
  height: 50px;
  line-height: 50px;
  background: $color-primary;
  color: #fff;
  border-bottom: 2px solid #fff;

  .logo {
    height: 48px;
    line-height: 48px;
    font-size: 22px;
    background: #0085e8 !important;
    width: $aside-width;
    cursor: pointer;

    &.collapsed {
      width: 64px;
    }

    .logo-img {
      width: 100%;
      height: 100%;
      background: url(../assets/logo.png);
      background-repeat: no-repeat;
      background-size: contain;
      background-position: center;
    }

    .logo-name {
      margin-left: 20px;
      img {
        width: 48px;
        height: 48px;
        margin: 0px;
      }

      span {
        margin-left: 10px;
      }
    }

    img {
      width: 40px;
      float: left;
      margin: 10px 10px 10px 18px;
    }

    .txt {
      color: #fff;
    }
  }

  .userinfo {
    text-align: right;
    padding-right: 35px;
    float: right;

    .userinfo-inner {
      cursor: pointer;
      color: #fff;

      img {
        width: 40px;
        height: 40px;
        border-radius: 20px;
        margin: 5px 0px 5px 10px;
        float: right;
      }
    }
  }

  .tools {
    padding: 0px 23px;
    width: 14px;
    height: 60px;
    line-height: 60px;
    cursor: pointer;
  }
}
</style>
