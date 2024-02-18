<template>
  <a-form :model="config" auto-label-width>
    <a-form-item label="发帖积分">
      <a-input-number
        v-model="config.postTopicScore"
        :min="1"
        mode="button"
        placeholder="发帖获得积分"
      />
    </a-form-item>
    <a-form-item label="跟帖积分">
      <a-input-number
        v-model="config.postCommentScore"
        :min="1"
        mode="button"
        placeholder="跟帖获得积分"
      />
    </a-form-item>
    <a-form-item label="签到积分">
      <a-input-number
        v-model="config.checkInScore"
        :min="1"
        mode="button"
        placeholder="签到获得积分"
      />
    </a-form-item>

    <a-form-item>
      <a-button type="primary" :loading="loading" @click="submit"
        >提交</a-button
      >
    </a-form-item>
  </a-form>
</template>

<script setup lang="ts">
  const loading = ref(false);
  const config = reactive({
    postTopicScore: undefined,
    postCommentScore: undefined,
    checkInScore: undefined,
  });
  const loadConfig = async () => {
    const ret = await axios.get<any, any>('/api/admin/sys-config/all');
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
      useNotificationSuccess('提交成功');
    } catch (e) {
      useHandleError(e);
    } finally {
      loading.value = false;
    }
  };
</script>

<style scoped lang="less"></style>
