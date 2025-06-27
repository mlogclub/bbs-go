<template>
  <a-form :model="config" auto-label-width>
    <a-form-item :label="$t('pages.settings.common.siteTitle')">
      <a-input
        v-model="config.siteTitle"
        type="text"
        :placeholder="$t('pages.settings.common.placeholder.siteTitle')"
      />
    </a-form-item>

    <a-form-item :label="$t('pages.settings.common.siteLogo')">
      <image-upload v-model="config.siteLogo" />
    </a-form-item>

    <a-form-item :label="$t('pages.settings.common.siteDescription')">
      <a-textarea
        v-model="config.siteDescription"
        :auto-size="{
          minRows: 2,
          maxRows: 5,
        }"
        :placeholder="$t('pages.settings.common.placeholder.siteDescription')"
      />
    </a-form-item>

    <a-form-item :label="$t('pages.settings.common.siteKeywords')">
      <a-select
        v-model="config.siteKeywords"
        multiple
        filterable
        allow-create
        default-first-option
        :placeholder="$t('pages.settings.common.placeholder.siteKeywords')"
      />
    </a-form-item>

    <a-form-item :label="$t('pages.settings.common.siteNotification')">
      <a-textarea
        v-model="config.siteNotification"
        :auto-size="{
          minRows: 2,
          maxRows: 5,
        }"
        :placeholder="$t('pages.settings.common.placeholder.siteNotification')"
      />
    </a-form-item>

    <a-form-item :label="$t('pages.settings.common.recommendTags')">
      <a-select
        v-model="config.recommendTags"
        multiple
        filterable
        allow-create
        default-first-option
        :placeholder="$t('pages.settings.common.placeholder.recommendTags')"
      />
    </a-form-item>

    <a-form-item :label="$t('pages.settings.common.defaultNodeId')">
      <a-select
        v-model="config.defaultNodeId"
        :placeholder="$t('pages.settings.common.placeholder.defaultNodeId')"
      >
        <a-option
          v-for="node in nodes"
          :key="node.id"
          :label="node.name"
          :value="node.id"
        />
      </a-select>
    </a-form-item>

    <a-form-item :label="$t('pages.settings.common.modules')">
      <a-checkbox v-model="config.modules.tweet" border>{{
        $t('pages.settings.common.tweet')
      }}</a-checkbox>
      <a-checkbox v-model="config.modules.topic" border>{{
        $t('pages.settings.common.topic')
      }}</a-checkbox>
      <a-checkbox v-model="config.modules.article" border>{{
        $t('pages.settings.common.article')
      }}</a-checkbox>
    </a-form-item>

    <a-form-item :label="$t('pages.settings.common.urlRedirect')">
      <a-tooltip
        :content="$t('pages.settings.common.urlRedirectTooltip')"
        placement="top"
      >
        <a-switch v-model="config.urlRedirect" />
      </a-tooltip>
    </a-form-item>

    <a-form-item :label="$t('pages.settings.common.enableHideContent')">
      <a-tooltip
        :content="$t('pages.settings.common.enableHideContentTooltip')"
        placement="top"
      >
        <a-switch v-model="config.enableHideContent" />
      </a-tooltip>
    </a-form-item>

    <a-form-item>
      <a-button type="primary" :loading="loading" @click="submit">{{
        $t('pages.settings.common.submit')
      }}</a-button>
    </a-form-item>
  </a-form>
</template>

<script setup lang="ts">
  import { NodeDTO } from '@/composables/types';
  import ImageUpload from '@/components/ImageUpload.vue';

  const { t } = useI18n();

  const loading = ref(false);
  const config = reactive({
    siteTitle: '',
    siteLogo: '',
    siteDescription: '',
    siteKeywords: [],
    siteNotification: '',
    recommendTags: [],
    defaultNodeId: undefined,
    urlRedirect: false,
    enableHideContent: false,
    modules: {
      tweet: false,
      topic: false,
      article: false,
    },
  });
  const nodes = ref<NodeDTO[]>([]);

  const loadConfig = async () => {
    const ret = await axios.get<any, any>('/api/admin/sys-config/configs');
    config.siteTitle = ret.siteTitle;
    config.siteLogo = ret.siteLogo;
    config.siteDescription = ret.siteDescription;
    config.siteKeywords = ret.siteKeywords;
    config.siteNotification = ret.siteNotification;
    config.recommendTags = ret.recommendTags;
    config.defaultNodeId = ret.defaultNodeId;
    config.urlRedirect = ret.urlRedirect;
    config.enableHideContent = ret.enableHideContent;
    config.modules = ret.modules;
    nodes.value = await axios.get<any, NodeDTO[]>(
      '/api/admin/topic-node/nodes'
    );
  };

  loadConfig();

  const submit = async () => {
    loading.value = true;
    try {
      await axios.post('/api/admin/sys-config/save', config);
      await loadConfig();
      useNotificationSuccess(t('pages.settings.common.message.submitSuccess'));
    } catch (e) {
      useHandleError(e);
    } finally {
      loading.value = false;
    }
  };
</script>

<style lang="scss" scoped></style>
