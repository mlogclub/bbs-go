<template>
  <el-menu
    router
    :default-active="$route.path"
    @open="handleOpen"
    @close="handleClose"
    @select="handleSelect"
    :collapse="collapsed"
    class="side-menu"
  >
    <template v-for="(item, index) in $router.options.routes">
      <el-submenu
        v-if="!item.hidden && item.children && item.children.length > 0"
        :index="'root-' + index"
        :key="item.path"
      >
        <template slot="title">
          <i v-if="item.meta && item.meta.icon" :class="item.meta.icon"></i>
          <span slot="title">{{item.meta ? (item.meta.title || '无标题') : '无标题'}}</span>
        </template>
        <template v-for="child in item.children">
          <el-menu-item v-if="!child.hidden" :index="child.path" :key="child.path">
            <i v-if="child.meta && child.meta.icon" :class="child.meta.icon"></i>
            {{child.meta ? (child.meta.title || '无标题') : '无标题'}}
          </el-menu-item>
        </template>
      </el-submenu>
    </template>
  </el-menu>
</template>

<script>
export default {
  name: "SideMenu",
  methods: {
    handleOpen() {},
    handleClose() {},
    handleSelect(a, b) {},
    showMenu(i, status) {
      this.$refs.menuCollapsed.getElementsByClassName(
        `submenu-hook-${i}`
      )[0].style.display = status ? "block" : "none";
    }
  },
  computed: {
    collapsed() {
      return this.$store.state.Default.collapsed;
    }
  }
};
</script>

<style scoped lang="scss">
@import "../styles/vars.scss";

.side-menu {
  height: 100%;
  &:not(.el-menu--collapse) {
    width: $aside-width-1;
  }

  &.el-menu--collapse {
    width: 64px;
  }
}
</style>
