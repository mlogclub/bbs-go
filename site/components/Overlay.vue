<template>
  <transition name="fade">
    <div
      v-show="{ show }"
      :style="{ zIndex: zIndex }"
      class="bbsgo-overlay"
      @touchmove="onTouchmove"
    ></div>
  </transition>
</template>

<script>
export default {
  props: {
    show: {
      type: Boolean,
      default: false,
    },
    lock: {
      type: Boolean,
      default: true,
    },
    zIndex: {
      type: Number,
      default: 1,
    },
  },
  methods: {
    onTouchmove(event) {
      if (this.lock) {
        this.preventTouchMove(event)
      }
    },
    preventTouchMove(event) {
      this.preventDefault(event, true)
    },
    preventDefault(event, isStopPropagation) {
      if (typeof event.cancelable !== 'boolean' || event.cancelable) {
        event.preventDefault()
      }
      if (isStopPropagation) {
        event.stopPropagation()
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.bbsgo-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.7); // TODO
  z-index: 999;
}
</style>
