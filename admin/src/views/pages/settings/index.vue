<template>
  <div class="container">
    <div class="container-main">
      <a-tabs position="top" lazy-load animation destroy-on-hide>
        <a-tab-pane key="1" :title="$t('pages.settings.common.title')">
          <div class="settings-container">
            <CommonSettings />
          </div>
        </a-tab-pane>
        <a-tab-pane key="2" :title="$t('pages.settings.nav.title')">
          <div class="settings-container">
            <NavSettings />
          </div>
        </a-tab-pane>
        <a-tab-pane key="3" :title="$t('pages.settings.score.title')">
          <div class="settings-container">
            <ScoreSettings />
          </div>
        </a-tab-pane>
        <a-tab-pane key="4" :title="$t('pages.settings.spam.title')">
          <div class="settings-container">
            <SpamSettings />
          </div>
        </a-tab-pane>
        <a-tab-pane key="5" :title="$t('pages.settings.upload.title')">
          <div class="settings-container">
            <UploadSettings />
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
  import UploadSettings from './components/UploadSettings.vue';
</script>

<style lang="less" scoped>
  .settings-container {
    padding: 20px;
  }
</style>
