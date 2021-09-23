<template>
  <div class="fullscreen-icon" @click="click">
    <i v-if="!isFullscreen" class="iconfont icon-fullscreen" />
    <i v-else class="iconfont icon-exit-fullscreen" />
  </div>
</template>

<script>
import screenfull from "screenfull";

export default {
  name: "Screenfull",
  data() {
    return {
      isFullscreen: false,
    };
  },
  mounted() {
    this.init();
  },
  beforeDestroy() {
    this.destroy();
  },
  methods: {
    click() {
      if (!screenfull.enabled) {
        this.$message({
          message: "you browser can not work",
          type: "warning",
        });
        return false;
      }
      screenfull.toggle();
    },
    change() {
      this.isFullscreen = screenfull.isFullscreen;
    },
    init() {
      if (screenfull.enabled) {
        screenfull.on("change", this.change);
      }
    },
    destroy() {
      if (screenfull.enabled) {
        screenfull.off("change", this.change);
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.fullscreen-icon {
  cursor: pointer;
  i {
    color: #5a5e66;
    font-size: 24px;
  }
}
</style>
