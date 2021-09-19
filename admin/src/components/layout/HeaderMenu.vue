<template>
  <!--
  <el-menu
    class="el-menu-demo"
    mode="horizontal"
    background-color="#1e80ff"
    text-color="#fff"
    active-text-color="#fff"
    router
  >
    <el-menu-item index="1">
      处理中心
    </el-menu-item>
  </el-menu>
  -->
  <div class="header-menu">
    <div class="header-menu-logo"><strong>BBS-GO</strong>&nbsp;管理平台</div>
    <div class="header-menu-user">
      <div data-v-d3d800ae="" class="avatar" />
      <el-dropdown>
        <span class="el-dropdown-link">
          用户名<i class="el-icon-arrow-down el-icon--right" />
        </span>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item>退出登录</el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </div>
  </div>
</template>

<script>
export default {}
</script>

<style lang="scss" scoped>
.header-menu {
  display: flex;
  align-items: center;
  background: #fff;
  border-bottom: 1px solid #e4e7f0;
  height: 59px;
  .header-menu-logo {
    height: 33px;
    margin-left: 25px;
    font-size: 24px;
    line-height: 30px;
  }
  .header-menu-user {
    display: flex;
    align-items: center;
    margin: auto 33px auto auto;

    .avatar {
      margin-right: 10px;
      width: 36px;
      height: 36px;
      border-radius: 50%;
      border-color: #e4e7f0;
    }
  }
}
</style>
