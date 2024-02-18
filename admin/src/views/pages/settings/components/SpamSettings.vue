<template>
  <a-form :model="config" auto-label-width>
    <a-form-item label="发帖验证码">
      <a-tooltip content="发帖时是否开启验证码校验" placement="top">
        <a-switch v-model="config.topicCaptcha" />
      </a-tooltip>
    </a-form-item>

    <a-form-item label="邮箱验证后发帖">
      <a-tooltip content="需要验证邮箱后才能发帖" placement="top">
        <a-switch v-model="config.createTopicEmailVerified" />
      </a-tooltip>
    </a-form-item>

    <a-form-item label="邮箱验证后发表文章">
      <a-tooltip content="需要验证邮箱后才能发表文章" placement="top">
        <a-switch v-model="config.createArticleEmailVerified" />
      </a-tooltip>
    </a-form-item>

    <a-form-item label="邮箱验证后评论">
      <a-tooltip content="需要验证邮箱后才能发表评论" placement="top">
        <a-switch v-model="config.createCommentEmailVerified" />
      </a-tooltip>
    </a-form-item>

    <a-form-item label="发表文章审核">
      <a-tooltip content="发布文章后是否开启审核" placement="top">
        <a-switch v-model="config.articlePending" />
      </a-tooltip>
    </a-form-item>

    <a-form-item label="用户观察期(秒)">
      <a-tooltip
        content="观察期内用户无法发表话题、动态等内容，设置为 0 表示无观察期。"
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

    <a-form-item label="邮箱白名单">
      <a-select
        v-model="config.emailWhitelist"
        style="width: 100%"
        multiple
        filterable
        allow-create
        default-first-option
        placeholder="邮箱白名单"
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
    topicCaptcha: undefined,
    createTopicEmailVerified: undefined,
    createArticleEmailVerified: undefined,
    createCommentEmailVerified: undefined,
    articlePending: undefined,
    userObserveSeconds: undefined,
    emailWhitelist: undefined,
  });
  const loadConfig = async () => {
    const ret = await axios.get<any, any>('/api/admin/sys-config/all');
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
      useNotificationSuccess('提交成功');
    } catch (e) {
      useHandleError(e);
    } finally {
      loading.value = false;
    }
  };
</script>

<style scoped lang="less"></style>
