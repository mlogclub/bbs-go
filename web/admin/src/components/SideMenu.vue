<template>
  <el-menu router :default-active="$route.path" @open="handleOpen" @close="handleClose"
           @select="handleSelect" :collapse="collapsed" class="side-menu">

    <template v-for="(item, index) in $router.options.routes">
      <el-submenu v-if="!item.hidden && item.children && item.children.length > 0"
                  :index="'root-' + index"
                  :key="item.path">
        <template slot="title">
          <i v-if="item.iconCls" :class="item.iconCls"></i>
          <span slot="title">{{item.name}}</span>
        </template>
        <template v-for="child in item.children">
          <el-menu-item v-if="!child.hidden" :index="child.path" :key="child.path">
            <i v-if="child.iconCls" :class="child.iconCls"></i>
            {{child.name}}
          </el-menu-item>
        </template>
      </el-submenu>

    </template>

  </el-menu>
</template>

<script>
  export default {
    name: 'SideMenu',
    methods: {
      handleOpen() {
      },
      handleClose() {
      },
      handleSelect: function (a, b) {
      },
      showMenu(i, status) {
        this.$refs.menuCollapsed.getElementsByClassName('submenu-hook-' + i)[0].style.display = status ? 'block' : 'none';
      },
    },
    computed: {
      collapsed() {
        return this.$store.state.Default.collapsed;
      }
    }
  };
</script>

<style scoped lang="scss">
  .side-menu:not(.el-menu--collapse) {
    width: 229px;
    min-height: 400px;
  }
</style>
