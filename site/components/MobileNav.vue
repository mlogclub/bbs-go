<template>
  <div>
    <nav class="mobile-nav">
      <div class="nav-left">
        <div class="sidebar-btn" @click="switchSidebar">
          <i class="iconfont icon-menu" />
        </div>
      </div>
      <div class="nav-center">
        <div class="menu-item">
          <nuxt-link to="/">首页</nuxt-link>
        </div>
        <div
          class="menu-item"
          :class="{ active: isShowNodes }"
          @click="switchNodes"
        >
          话题
          <i v-if="isShowNodes" class="iconfont icon-drop-up" />
          <i v-else class="iconfont icon-drop-down" />
        </div>
      </div>
      <div class="nav-right">
        <create-topic-btn class="create-topic-btn">
          <i class="iconfont icon-plus" style="color: var(--text-color)" />
        </create-topic-btn>
      </div>
    </nav>
    <mobile-sidebar @click.native="close" />
    <mobile-nodes @click.native="close" />
    <overlay
      v-show="isShowOverlay"
      :z-index="overlayZIndex"
      @click.native="close"
    />
  </div>
</template>

<script>
export default {
  data() {
    return {
      sidebarActive: false,
      nodesActive: false,
    }
  },
  computed: {
    isShowOverlay() {
      return (
        this.$store.state.env.showMobileSidebar ||
        this.$store.state.env.showMobileNodes
      )
    },
    overlayZIndex() {
      return this.isShowSidebar ? 40 : 20
    },
    isShowSidebar() {
      return this.$store.state.env.showMobileSidebar
    },
    isShowNodes() {
      return this.$store.state.env.showMobileNodes
    },
  },
  methods: {
    switchSidebar() {
      this.$store.commit('env/setShowMobileSidebar', !this.isShowSidebar)
      if (this.isShowSidebar) {
        this.$store.commit('env/setShowMobileNodes', false)
      }
    },
    switchNodes() {
      this.$store.commit('env/setShowMobileNodes', !this.isShowNodes)
    },
    close() {
      this.$store.commit('env/setShowMobileNodes', false)
      this.$store.commit('env/setShowMobileSidebar', false)
    },
  },
}
</script>
<style lang="scss" scoped></style>
