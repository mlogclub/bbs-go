<template>
  <el-col :span="24" class="header">
    <el-col :span="10" class="logo" :class="collapsed?'logo-collapse-width':'logo-width'">{{collapsed?'':sysName}}</el-col>
    <el-col :span="10">
      <div class="tools" @click.prevent="collapse">
        <i class="iconfont" :class="collapsed ? 'icon-right' : 'icon-left'"></i>
      </div>
    </el-col>
    <el-col :span="4" class="userinfo" v-if="userInfo">
      <el-dropdown trigger="hover">
        <span class="el-dropdown-link userinfo-inner">
          <img :src="userInfo.avatar" />
          {{userInfo.nickname}}
        </span>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item @click.native="logout">退出登录</el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </el-col>
  </el-col>
</template>

<script>
import cookies from "js-cookie";
export default {
  name: "Header",
  data() {
    return {
      sysName: "BBS-GO"
    };
  },
  methods: {
    // 折叠导航栏
    collapse() {
      this.$store.dispatch("Default/collapse");
    },
    // 退出登录
    logout() {
      this.$message({ message: "暂未实现", type: "success" });
    }
  },
  mounted() {
    let user = sessionStorage.getItem("user");
    if (user) {
      user = JSON.parse(user);
    }
  },
  computed: {
    collapsed() {
      return this.$store.state.Default.collapsed;
    },
    userInfo() {
      let userInfo = this.$store.state.Login.userInfo;
      if (!userInfo) {
        const userInfoStr = cookies.get("userInfo");
        if (userInfoStr) {
          try {
            userInfo = JSON.parse(userInfoStr);
            this.$store.dispatch("Login/setUserInfo", userInfo);
          } catch (e) {
            console.error(e);
          }
        }
      }
      return userInfo;
    }
  }
};
</script>

<style scoped lang="scss">
@import "../styles/vars.scss";

.header {
  height: 50px;
  line-height: 50px;
  background: $color-primary;
  color: #fff;
  border-bottom: 2px solid #fff;

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
    height: 48px;
    line-height: 48px;
    font-size: 22px;
    padding-left: 20px;
    background: #0085e8 !important;

    // border-color: rgba(238, 241, 146, 0.3);
    // border-color: #e6e6e6;
    // border-right-width: 1px;
    // border-right-style: solid;

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
    width: $aside-width;
  }

  .logo-collapse-width {
    width: 65px;
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
