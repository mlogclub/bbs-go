<template>
  <div id="app">
    <transition name="fade" mode="out-in">
      <router-view></router-view>
    </transition>

    <login />
  </div>
</template>

<script>
import Login from "./components/Login";

export default {
  name: "app",
  components: { Login }
};
</script>

<style lang="scss">
body {
  margin: 0px;
  padding: 0px;
  font-family: Helvetica Neue, Helvetica, PingFang SC, Hiragino Sans GB,
    Microsoft YaHei, SimSun, sans-serif;
  font-size: 14px;
  -webkit-font-smoothing: antialiased;
}

a {
  color: #3273dc;
  text-decoration: none;

  &:visited {
    color: #3273dc;
  }
}

#app {
  position: absolute;
  top: 0px;
  bottom: 0px;
  width: 100%;
}

.el-submenu [class^="fa"] {
  vertical-align: baseline;
  margin-right: 10px;
}

.el-menu-item [class^="fa"] {
  vertical-align: baseline;
  margin-right: 10px;
}

.main-content {
  padding: 0 10px;
}

.fade-enter-active,
.fade-leave-active {
  transition: all 0.2s ease;
}

.fade-enter,
.fade-leave-active {
  opacity: 0;
}

.content-form {
  .form-item {
    display: flex;
    .field-key {
      width: 100px;
      text-align: right;
      font-weight: bold;
    }
  }
}

/*滚动条整体样式*/
::-webkit-scrollbar {
  width: 6px;
  // height: 1px;
}
/*滚动条滑块*/
::-webkit-scrollbar-thumb {
  // box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.2);
  border-radius: 10px;
  background: rgba(144, 147, 153, 0.5);
  height: 10px;
}
/*滚动条轨道*/
::-webkit-scrollbar-track {
  // box-shadow: inset 0 0 1px rgba(0, 0, 0, 0);
  border-radius: 0px;
  background: rgba(144, 147, 153, 0.2);
}
</style>
