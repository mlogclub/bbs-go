<template>
  <section class="app-main">
    <transition name="fade-transform" mode="out-in">
      <keep-alive :include="cachedViews">
        <router-view :key="key" />
      </keep-alive>
    </transition>
  </section>
</template>

<script>
export default {
  name: "AppMain",
  computed: {
    cachedViews() {
      return this.$store.state.tagsView.cachedViews;
    },
    key() {
      return this.$route.path;
    },
  },
};
</script>

<style lang="scss" scoped>
.app-main {
  /* 50= navbar  50  */
  min-height: calc(100vh - 50px);
  width: 100%;
  height: 100%;
  position: relative;
  // overflow: hidden;
  overflow-y: auto;
  background-color: #f3f6fe;
  // padding: 10px;
  // margin-top: 10px;
  padding: 20px;

  .app-main-container {
    background-color: #fff;
    border-radius: 8px;
  }
}

.fixed-header + .app-main {
  padding-top: 50px;
}

.hasTagsView {
  .app-main {
    /* 84 = navbar + tags-view = 50 + 34 */
    min-height: calc(100vh - 84px);
    // min-height: calc(100vh - 50px - 34px);
  }

  .fixed-header + .app-main {
    // padding-top: 84px;
    padding-top: calc(50px + 20px + 34px);
  }
}
</style>

<style lang="scss">
// fix css style bug in open el-dialog
.el-popup-parent--hidden {
  .fixed-header {
    padding-right: 15px;
  }
}
</style>
