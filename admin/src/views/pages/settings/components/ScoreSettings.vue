<template>
  <a-form :model="config" auto-label-width>
    <a-form-item :label="$t('pages.settings.score.postTopicScore')">
      <a-input-number
        v-model="config.postTopicScore"
        :min="1"
        mode="button"
        :placeholder="$t('pages.settings.score.placeholder.postTopicScore')"
      />
    </a-form-item>
    <a-form-item :label="$t('pages.settings.score.postCommentScore')">
      <a-input-number
        v-model="config.postCommentScore"
        :min="1"
        mode="button"
        :placeholder="$t('pages.settings.score.placeholder.postCommentScore')"
      />
    </a-form-item>
    <a-form-item :label="$t('pages.settings.score.checkInScore')">
      <a-input-number
        v-model="config.checkInScore"
        :min="1"
        mode="button"
        :placeholder="$t('pages.settings.score.placeholder.checkInScore')"
      />
    </a-form-item>

    <a-form-item>
      <a-button type="primary" :loading="loading" @click="submit">{{
        $t('pages.settings.score.submit')
      }}</a-button>
    </a-form-item>
  </a-form>
</template>

<script setup lang="ts">
  const { t } = useI18n();

  const loading = ref(false);
  const config = reactive({
    postTopicScore: undefined,
    postCommentScore: undefined,
    checkInScore: undefined,
  });
  const loadConfig = async () => {
    const ret = await axios.get<any, any>('/api/admin/sys-config/configs');
    config.postTopicScore = ret.scoreConfig.postTopicScore;
    config.postCommentScore = ret.scoreConfig.postCommentScore;
    config.checkInScore = ret.scoreConfig.checkInScore;
  };

  loadConfig();

  const submit = async () => {
    loading.value = true;
    try {
      await axios.post('/api/admin/sys-config/save', {
        scoreConfig: config,
      });
      await loadConfig();
      useNotificationSuccess(t('pages.settings.score.message.submitSuccess'));
    } catch (e) {
      useHandleError(e);
    } finally {
      loading.value = false;
    }
  };
</script>

<style scoped lang="less"></style>
