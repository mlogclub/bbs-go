<template>
  <el-col :span="24" class="header">
    <el-col :span="10" class="logo" :class="collapsed?'logo-collapse-width':'logo-width'">
      {{collapsed?'':sysName}}
    </el-col>
    <el-col :span="10">
      <div class="tools" @click.prevent="collapse">
        <i class="iconfont" :class="collapsed ? 'icon-right' : 'icon-left'"></i>
      </div>
    </el-col>
    <el-col :span="4" class="userinfo">
      <!--
      <el-dropdown trigger="hover">
          <span class="el-dropdown-link userinfo-inner">
            <img :src="avatar"/> {{nickname}}
          </span>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item>设置</el-dropdown-item>
          <el-dropdown-item divided @click.native="logout">退出登录</el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
      -->
    </el-col>
  </el-col>
</template>

<script>
  export default {
    name: 'Header',
    data() {
      return {
        sysName: 'M-LOG CLUB',
      }
    },
    methods: {
      // 折叠导航栏
      collapse: function () {
        this.$store.dispatch('Default/collapse')
      },
      // 退出登录
      logout: function () {
        var _this = this
        this.$confirm('确认退出吗?', '提示', {
          type: 'warning'
        }).then(() => {
          sessionStorage.removeItem('user')
          _this.$router.push('/login')
        }).catch(() => {

        })
      },
    },
    mounted() {
      var user = sessionStorage.getItem('user')
      if (user) {
        user = JSON.parse(user)
      }
    },
    computed: {
      collapsed() {
        return this.$store.state.Default.collapsed
      },
      nickname() {
        var userInfo = this.$store.state.Login.userInfo || {}
        return userInfo.nickname || ''
      },
      avatar() {
        var userInfo = this.$store.state.Login.userInfo || {}
        return userInfo.avatar || ''
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

    .logo {
      height: 60px;
      font-size: 22px;
      padding-left: 20px;
      padding-right: 20px;
      border-color: rgba(238, 241, 146, 0.3);
      border-right-width: 1px;
      border-right-style: solid;

      img {
        width: 40px;
        float: left;
        margin: 10px 10px 10px 18px;
      }

      .txt {
        color: #fff;
      }
    }

    .logo-width {
      width: 230px;
    }

    .logo-collapse-width {
      width: 65px
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
