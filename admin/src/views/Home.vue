<template>
  <el-container class="sec">
    <el-header class="sec-header">
      <frame-header />
    </el-header>
    <el-container class="sec-main">
      <el-aside class="sec-aside" :class="{'collapse': collapsed}">
        <side-menu />
      </el-aside>
      <el-main class="sec-main-view">
        <router-tab />
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import FrameHeader from "../components/FrameHeader";
import SideMenu from "../components/SideMenu";

export default {
  components: {
    FrameHeader,
    SideMenu
  },
  data() {
    return {
      form: {
        name: "",
        region: "",
        date1: "",
        date2: "",
        delivery: false,
        type: [],
        resource: "",
        desc: ""
      }
    };
  },
  methods: {},
  computed: {
    collapsed() {
      return this.$store.state.Default.collapsed;
    }
  }
};
</script>

<style scoped lang="scss">
@import "../styles/vars.scss";

.sec {
  margin: 0px;
  padding: 0px;
  height: 100%;

  .sec-header {
    padding: 0px;
    height: 48px !important;
  }

  .sec-main {
    height: 100%;

    // border: solid 3px red;

    .sec-main-view {
      padding: 0px;
      height: 100%;

      padding: 0px;

      /deep/ .router-tab {
        height: 100%;
        // 路由页面
        &-page {
          padding: 0px;
          font-size: 14px;
          line-height: 1.5;
        }
      }
    }
  }

  .sec-aside {
    height: 100%;
    width: $aside-width !important;
    overflow-y: auto;
    overflow-x: hidden;

    transition: width 0.2s;

    &.collapse {
      width: 65px !important;
    }
  }
}
</style>
