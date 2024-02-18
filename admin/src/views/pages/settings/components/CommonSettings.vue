<template>
  <a-form :model="config" auto-label-width>
    <a-form-item label="网站名称">
      <a-input v-model="config.siteTitle" type="text" placeholder="网站名称" />
    </a-form-item>

    <a-form-item label="网站描述">
      <a-textarea
        v-model="config.siteDescription"
        :auto-size="{
          minRows: 2,
          maxRows: 5,
        }"
        placeholder="网站描述"
      />
    </a-form-item>

    <a-form-item label="网站关键字">
      <a-select
        v-model="config.siteKeywords"
        multiple
        filterable
        allow-create
        default-first-option
        placeholder="网站关键字"
      />
    </a-form-item>

    <a-form-item label="网站公告">
      <a-textarea
        v-model="config.siteNotification"
        :auto-size="{
          minRows: 2,
          maxRows: 5,
        }"
        placeholder="网站公告（支持输入HTML）"
      />
    </a-form-item>

    <a-form-item label="推荐标签">
      <a-select
        v-model="config.recommendTags"
        multiple
        filterable
        allow-create
        default-first-option
        placeholder="推荐标签"
      />
    </a-form-item>

    <a-form-item label="默认节点">
      <a-select v-model="config.defaultNodeId" placeholder="发帖默认节点">
        <a-option
          v-for="node in nodes"
          :key="node.id"
          :label="node.name"
          :value="node.id"
        />
      </a-select>
    </a-form-item>

    <a-form-item label="功能模块">
      <a-checkbox v-model="config.modules.tweet" border>动态</a-checkbox>
      <a-checkbox v-model="config.modules.topic" border>帖子</a-checkbox>
      <a-checkbox v-model="config.modules.article" border>文章</a-checkbox>
    </a-form-item>

    <a-form-item label="站外链接跳转页面">
      <a-tooltip content="在跳转前需手动确认是否前往该站外链接" placement="top">
        <a-switch v-model="config.urlRedirect" />
      </a-tooltip>
    </a-form-item>

    <a-form-item label="启用评论可见内容">
      <a-tooltip content="发帖时支持设置评论后可见内容" placement="top">
        <a-switch v-model="config.enableHideContent" />
      </a-tooltip>
    </a-form-item>

    <a-form-item>
      <a-button type="primary" :loading="loading" @click="submit"
        >提交</a-button
      >
    </a-form-item>
  </a-form>
</template>

<script setup lang="ts">
  import { NodeDTO } from '@/composables/types';

  const loading = ref(false);
  const config = reactive({
    siteTitle: '',
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
    const ret = await axios.get<any, any>('/api/admin/sys-config/all');
    config.siteTitle = ret.siteTitle;
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
      useNotificationSuccess('提交成功');
    } catch (e) {
      useHandleError(e);
    } finally {
      loading.value = false;
    }
  };
</script>

<style lang="scss" scoped></style>
