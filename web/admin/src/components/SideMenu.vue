<template>
  <el-menu router :default-active="$route.path" @open="handleOpen" @close="handleClose"
           @select="handleSelect" :collapse="collapsed" class="side-menu">

    <template v-for="(item, index) in $router.options.routes" v-if="!item.hidden">

      <el-submenu :index="'root-' +index" v-if="!item.leaf && item.children.length > 0">
        <template slot="title">
          <i v-if="item.iconCls" :class="item.iconCls"></i>
          <span slot="title">{{item.name}}</span>
        </template>
        <el-menu-item v-for="child in item.children" :index="child.path" v-if="!child.hidden">
          <i v-if="child.iconCls" :class="child.iconCls"></i>
          {{child.name}}
        </el-menu-item>
      </el-submenu>

      <el-menu-item v-if="item.leaf && item.children.length === 1" :index="item.children[0].path">
        <i v-if="item.children[0].iconCls" :class="item.children[0].iconCls"></i>
        {{item.children[0].name}}
      </el-menu-item>

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
