<template>
  <el-menu
    class="layout-menu"
    router
  >
    <!--    <el-menu-item index="/home">-->
    <!--      首页-->
    <!--    </el-menu-item>-->
    <!--    <el-menu-item index="/about">-->
    <!--      关于-->
    <!--    </el-menu-item>-->

    <sidebar-menu-item
      v-for="route in routes"
      :key="route.path"
      :item="route"
      :base-path="route.path"
    />
  </el-menu>
</template>

<script>
import path from 'path'
import SidebarMenuItem from '@/components/layout/SidebarMenuItem'

export default {
  components: { SidebarMenuItem },
  computed: {
    routes () {
      return this.$router.options.routes
    }
  },
  methods: {
    hasOneShowingChild (children = [], parent) {
      const showingChildren = children.filter((item) => {
        if (item.hidden) {
          return false
        }
        // Temp set(will be used if only has one showing child)
        this.onlyOneChild = item
        return true
      })

      // When there is only one child router, the child router is displayed by default
      if (showingChildren.length === 1) {
        return true
      }

      // Show parent if there are no child router to display
      if (showingChildren.length === 0) {
        this.onlyOneChild = { ...parent, path: '', noShowingChildren: true }
        return true
      }

      return false
    },
    resolvePath (routePath) {
      // if (isExternal(routePath)) {
      //   return routePath;
      // }
      // if (isExternal(this.basePath)) {
      //   return this.basePath;
      // }
      return path.resolve(this.basePath, routePath)
    }
  }
}
</script>

<style scoped></style>
