<template>
  <el-container class="layout">
    <el-header class="layout-header">
      <header-menu />
    </el-header>
    <el-container
      class="layout-container"
      :style="{ height: layoutContainerHeight }"
    >
      <el-aside class="layout-aside">
        <side-menu />
      </el-aside>
      <el-main class="layout-main">
        <router-tab :tabs="tabs" />
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import HeaderMenu from '@/components/layout/HeaderMenu'
import SideMenu from '@/components/layout/SideMenu'

export default {
  components: {
    HeaderMenu,
    SideMenu
  },
  data () {
    return {
      tabs: ['/'],
      layoutContainerHeight: '100%',
      routeConfigs: []
    }
  },
  mounted () {
    const me = this
    me.handleLayoutContainerHeight()
    window.onresize = () => {
      me.handleLayoutContainerHeight()
    }
  },
  methods: {
    handleLayoutContainerHeight () {
      this.layoutContainerHeight = `${document.documentElement.offsetHeight - 60}px`
      console.log(document.documentElement.offsetHeight)
    }
  }
}
</script>

<style lang="scss">
.layout {
  .layout-header {
    padding: 0 !important;
  }

  .layout-container {
    .layout-aside {
      width: 220px !important;

      .layout-menu {
        height: 100%;
      }
    }

    .layout-main {
      padding: 0;

      .router-tab {
        height: 100%;

        .router-tab__header {
          .router-tab__nav {
            width: 100%;

            .router-tab__item {
              color: #495060;
              &.is-active {
                background-color: #f3f6fe;
                color: #000;
                border-bottom: none;
                font-weight: 700;
              }
            }
          }
        }

        .router-tab__container {
          background-color: #f3f6fe;

          .router-alive {
            padding: 20px;
          }
        }
      }
    }
  }
}
</style>
