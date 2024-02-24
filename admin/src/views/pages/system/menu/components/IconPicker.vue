<template>
  <div class="icon-picker">
    <div class="icon-picker-input">
      <a-input
        :model-value="modelValue"
        class="icon-picker-input-input"
        placeholder="请选择图标"
        @input="onInput"
      >
        <template #prefix>
          <component :is="modelValue" />
        </template>
      </a-input>
      <a-button class="icon-picker-input-btn" type="primary" @click="showIcons"
        >选择图标</a-button
      >
    </div>

    <a-modal v-model:visible="visible" :width="650">
      <template #title> 选择图标 </template>
      <div class="icons">
        <div
          v-for="(icon, i) in icons"
          :key="i"
          class="icon-item"
          @click="selectIcon(icon)"
        >
          <a-popover>
            <component :is="icon" :size="30" />
            <template #content>
              {{ icon }}
            </template>
          </a-popover>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
  const icons = [
    'icon-dashboard',
    'icon-apps',
    'icon-home',
    'icon-stop',
    'icon-link',
    'icon-import',
    'icon-launch',
    'icon-list',
    'icon-message',
    'icon-poweroff',
    'icon-scan',
    'icon-save',
    'icon-settings',
    'icon-share-alt',
    'icon-archive',
    'icon-bar-chart',
    'icon-book',
    'icon-bookmark',
    'icon-branch',
    'icon-bug',
    'icon-bulb',
    'icon-calendar-clock',
    'icon-calendar',
    'icon-camera',
    'icon-cloud',
    'icon-command',
    'icon-common',
    'icon-compass',
    'icon-computer',
    'icon-copyright',
    'icon-desktop',
    'icon-dice',
    'icon-drag-dot-vertical',
    'icon-drag-dot',
    'icon-drive-file',
    'icon-ear',
    'icon-email',
    'icon-empty',
    'icon-experiment',
    'icon-file-audio',
    'icon-file-image',
    'icon-file-pdf',
    'icon-file-video',
    'icon-file',
    'icon-fire',
    'icon-folder-add',
    'icon-folder-delete',
    'icon-folder',
    'icon-gift',
    'icon-idcard',
    'icon-image-close',
    'icon-image',
    'icon-interaction',
    'icon-language',
    'icon-layers',
    'icon-layout',
    'icon-location',
    'icon-lock',
    'icon-loop',
    'icon-man',
    'icon-menu',
    'icon-mind-mapping',
    'icon-mobile',
    'icon-moon',
    'icon-mosaic',
    'icon-nav',
    'icon-notification-close',
    'icon-notification',
    'icon-palette',
    'icon-pen',
    'icon-phone',
    'icon-printer',
    'icon-public',
    'icon-pushpin',
    'icon-qrcode',
    'icon-relation',
    'icon-robot-add',
    'icon-robot',
    'icon-safe',
    'icon-schedule',
    'icon-shake',
    'icon-skin',
    'icon-stamp',
    'icon-storage',
    'icon-subscribe-add',
    'icon-subscribe',
    'icon-subscribed',
    'icon-sun',
    'icon-tag',
    'icon-tags',
    'icon-thunderbolt',
    'icon-tool',
    'icon-trophy',
    'icon-unlock',
    'icon-user-add',
    'icon-user-group',
    'icon-user',
    'icon-video-camera',
    'icon-wifi',
    'icon-woman',
    'icon-live-broadcast',
  ];

  const visible = ref(false);

  const emits = defineEmits(['update:modelValue']);

  const props = defineProps({
    modelValue: {
      type: String,
      default: '',
    },
  });
  //   const value = ref(props.modelValue);

  const showIcons = () => {
    visible.value = true;
  };

  const selectIcon = (icon) => {
    visible.value = false;

    emits('update:modelValue', icon);
  };

  const onInput = (value) => {
    emits('update:modelValue', value);
  };
</script>

<style lang="less" scoped>
  .icon-picker {
    width: 100%;
    .icon-picker-input {
      display: flex;

      .icon-picker-input-btn {
        margin-left: 10px;
      }
    }
  }

  .icons {
    display: flex;
    flex-wrap: wrap;
    border-left: 1px solid var(--color-neutral-3);
    border-top: 1px solid var(--color-neutral-3);

    .icon-item {
      flex: 1;
      cursor: pointer;
      padding: 10px;
      border-right: 1px solid var(--color-neutral-3);
      border-bottom: 1px solid var(--color-neutral-3);
    }
  }
</style>
