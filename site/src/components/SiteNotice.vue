<template>
  <div v-if="siteNotification" class="widget">
    <div class="widget-header">公告</div>
    <div class="widget-content content notice-content">
      <p v-html="siteNotification"></p>
    </div>
  </div>
</template>

<script setup>
const siteNotification = computed(() => {
  const configStore = useConfigStore();
  return configStore.config.siteNotification;
});
</script>

<style lang="scss" scoped>
.notice-content {
  font-size: 80%;
}
</style>
