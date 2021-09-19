<template>
  <div/>
</template>

<script>
export default {
  props: {
    // route object
    item: {
      type: Object,
      required: true,
    },
    isNest: {
      type: Boolean,
      default: false,
    },
    basePath: {
      type: String,
      default: '',
    },
  },
};
</script>

<style scoped>

</style>
