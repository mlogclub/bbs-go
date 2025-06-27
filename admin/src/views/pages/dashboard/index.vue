<template>
  <div class="dashboard">
    <a-card v-if="systemInfo" :title="$t('pages.dashboard.systemInfo')">
      <div class="attr">
        <span class="attr-name">{{ $t('pages.dashboard.os') }}:</span>
        <span class="attr-value">{{ systemInfo.os }}</span>
      </div>
      <div class="attr">
        <span class="attr-name">{{ $t('pages.dashboard.arch') }}:</span>
        <span class="attr-value">{{ systemInfo.arch }}</span>
      </div>
      <div class="attr">
        <span class="attr-name">{{ $t('pages.dashboard.cpu') }}:</span>
        <span class="attr-value">{{ systemInfo.numCpu }}</span>
      </div>
      <div class="attr">
        <span class="attr-name">{{ $t('pages.dashboard.goVersion') }}:</span>
        <span class="attr-value">{{ systemInfo.goVersion }}</span>
      </div>
      <div class="attr">
        <span class="attr-name">{{ $t('pages.dashboard.hostname') }}:</span>
        <span class="attr-value">{{ systemInfo.hostname }}</span>
      </div>
    </a-card>

    <a-card v-if="systemInfo" :title="$t('pages.dashboard.teamInfo')">
      <div class="attr">
        <span class="attr-name"
          >{{ $t('pages.dashboard.officialSite') }}:&nbsp;</span
        >
        <a class="attr-value" href="https://mlog.club" target="_blank"
          >https://mlog.club</a
        >
      </div>
      <div class="attr">
        <span class="attr-name"
          >{{ $t('pages.dashboard.devTeam') }}:&nbsp;</span
        >
        <a class="attr-value" href="https://github.com/mlogclub" target="_blank"
          >mlogclub</a
        >
      </div>
      <div class="attr">
        <span class="attr-name">{{ $t('pages.dashboard.github') }}:&nbsp;</span>
        <a
          class="attr-value"
          href="https://github.com/mlogclub/bbs-go"
          target="_blank"
          >https://github.com/mlogclub/bbs-go</a
        >
      </div>
      <div class="attr">
        <span class="attr-name">{{ $t('pages.dashboard.gitee') }}:&nbsp;</span>
        <a
          class="attr-value"
          href="https://gitee.com/mlogclub/bbs-go"
          target="_blank"
          >https://gitee.com/mlogclub/bbs-go</a
        >
      </div>
    </a-card>
  </div>
</template>

<script setup lang="ts">
  const systemInfo = ref();
  const load = async () => {
    systemInfo.value = await axios.get('/api/admin/common/system_info');
  };
  load();
</script>

<style lang="less" scoped>
  .dashboard {
    margin: 20px;
    display: flex;
    flex-direction: column;
    row-gap: 20px;

    font-size: 14px;
    line-height: 1.8;

    .attr {
      display: flex;
      align-items: center;
      column-gap: 10px;
      font-size: 14px;
      .attr-name {
        color: var(--color-text-1);
      }
      span.attr-value {
        color: var(--color-text-2);
      }
      a.attr-value {
        color: rgb(var(--arcoblue-6));
      }
    }
  }
</style>
