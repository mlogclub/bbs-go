<template>
  <a-form :model="config" auto-label-width>
    <a-form-item :label="$t('pages.settings.spam.topicCaptcha')">
      <a-tooltip
        :content="$t('pages.settings.spam.topicCaptchaTooltip')"
        placement="top"
      >
        <a-switch v-model="config.topicCaptcha" />
      </a-tooltip>
    </a-form-item>

    <a-form-item :label="$t('pages.settings.spam.createTopicEmailVerified')">
      <a-tooltip
        :content="$t('pages.settings.spam.createTopicEmailVerifiedTooltip')"
        placement="top"
      >
        <a-switch v-model="config.createTopicEmailVerified" />
      </a-tooltip>
    </a-form-item>

    <a-form-item :label="$t('pages.settings.spam.createArticleEmailVerified')">
      <a-tooltip
        :content="$t('pages.settings.spam.createArticleEmailVerifiedTooltip')"
        placement="top"
      >
        <a-switch v-model="config.createArticleEmailVerified" />
      </a-tooltip>
    </a-form-item>

    <a-form-item :label="$t('pages.settings.spam.createCommentEmailVerified')">
      <a-tooltip
        :content="$t('pages.settings.spam.createCommentEmailVerifiedTooltip')"
        placement="top"
      >
        <a-switch v-model="config.createCommentEmailVerified" />
      </a-tooltip>
    </a-form-item>

    <a-form-item :label="$t('pages.settings.spam.articlePending')">
      <a-tooltip
        :content="$t('pages.settings.spam.articlePendingTooltip')"
        placement="top"
      >
        <a-switch v-model="config.articlePending" />
      </a-tooltip>
    </a-form-item>

    <a-form-item :label="$t('pages.settings.spam.userObserveSeconds')">
      <a-tooltip
        :content="$t('pages.settings.spam.userObserveSecondsTooltip')"
        placement="top"
      >
        <a-input-number
          v-model="config.userObserveSeconds"
          mode="button"
          :min="0"
          :max="720"
        />
      </a-tooltip>
    </a-form-item>

    <a-form-item :label="$t('pages.settings.spam.emailWhitelist')">
      <a-select
        v-model="config.emailWhitelist"
        style="width: 100%"
        multiple
        filterable
        allow-create
        default-first-option
        :placeholder="$t('pages.settings.spam.placeholder.emailWhitelist')"
      />
    </a-form-item>

    <a-form-item>
      <a-button type="primary" :loading="loading" @click="submit">{{
        $t('pages.settings.spam.submit')
      }}</a-button>
    </a-form-item>
  </a-form>
</template>

<script setup lang="ts">
  const { t } = useI18n();

  const loading = ref(false);
  const config = reactive({
    topicCaptcha: undefined,
    createTopicEmailVerified: undefined,
    createArticleEmailVerified: undefined,
    createCommentEmailVerified: undefined,
    articlePending: undefined,
    userObserveSeconds: undefined,
    emailWhitelist: undefined,
  });
  const loadConfig = async () => {
    const ret = await axios.get<any, any>('/api/admin/sys-config/configs');
    config.topicCaptcha = ret.topicCaptcha;
    config.createTopicEmailVerified = ret.createTopicEmailVerified;
    config.createArticleEmailVerified = ret.createArticleEmailVerified;
    config.createCommentEmailVerified = ret.createCommentEmailVerified;
    config.articlePending = ret.articlePending;
    config.userObserveSeconds = ret.userObserveSeconds;
    config.emailWhitelist = ret.emailWhitelist;
  };

  loadConfig();

  const submit = async () => {
    loading.value = true;
    try {
      await axios.post('/api/admin/sys-config/save', config);
      await loadConfig();
      useNotificationSuccess(t('pages.settings.spam.message.submitSuccess'));
    } catch (e) {
      useHandleError(e);
    } finally {
      loading.value = false;
    }
  };
</script>

<style scoped lang="less"></style>
