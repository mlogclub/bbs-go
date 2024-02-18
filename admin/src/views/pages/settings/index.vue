<template>
  <div class="container">
    <div class="container-main">
      <a-tabs position="top" lazy-load animation destroy-on-hide>
        <a-tab-pane key="1" title="通用配置">
          <div class="settings-container">
            <CommonSettings />
          </div>
        </a-tab-pane>
        <a-tab-pane key="2" title="导航配置">
          <div class="settings-container">
            <NavSettings />
          </div>
        </a-tab-pane>
        <a-tab-pane key="3" title="积分配置">
          <div class="settings-container">
            <ScoreSettings />
          </div>
        </a-tab-pane>
        <a-tab-pane key="4" title="反作弊配置">
          <div class="settings-container">
            <SpamSettings />
          </div>
        </a-tab-pane>
      </a-tabs>
    </div>
  </div>
</template>

<script setup>
  import CommonSettings from './components/CommonSettings.vue';
  import NavSettings from './components/NavSettings.vue';
  import ScoreSettings from './components/ScoreSettings.vue';
  import SpamSettings from './components/SpamSettings.vue';
</script>

<style lang="less" scoped>
  .settings-container {
    padding: 20px;
  }
</style>
