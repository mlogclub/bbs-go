<template>
  <div class="topics">
    <template v-if="results && results.length">
      <div v-for="topic in results" :key="topic.id" class="topic-item">
        <div class="topic-status">
          <a-tag v-if="topic.recommend" color="green">已推荐</a-tag>
          <a-tag v-if="topic.status === 1" color="red">已删除</a-tag>
        </div>

        <div class="topic-header">
          <a-avatar :size="40">
            <img v-if="topic.user.avatar" :src="topic.user.avatar" />
            <IconUser v-else />
          </a-avatar>
          <div class="topic-head-info">
            <div class="nickname">{{ topic.user.nickname }}</div>
            <div class="topic-metas">
              <div class="meta-item">
                <span>ID:</span>
                <span>{{ topic.id }}</span>
              </div>
              <div class="meta-item">
                <span>时间:</span>
                <span>{{ useFormatDate(topic.createTime) }}</span>
              </div>
              <div class="meta-item">
                <span>查看:</span>
                <span>{{ topic.viewCount }}</span>
              </div>
              <div class="meta-item">
                <span>点赞:</span>
                <span>{{ topic.likeCount }}</span>
              </div>
              <div class="meta-item">
                <span>评论:</span>
                <span>{{ topic.commentCount }}</span>
              </div>
            </div>
          </div>
        </div>
        <div v-if="topic.type === 0 && topic.summary" class="topic-summary">
          {{ topic.summary }}
        </div>
        <div v-if="topic.type === 1 && topic.content" class="topic-summary">
          {{ topic.content }}
        </div>
        <a-image-preview-group
          v-if="topic.imageList && topic.imageList.length"
          infinite
        >
          <div class="topic-image-list">
            <a-image
              v-for="(image, index) in topic.imageList"
              :key="index"
              width="150"
              height="150"
              class="image-item"
              show-loader
              :src="image.url"
              fit="cover"
            />
          </div>
        </a-image-preview-group>
        <div class="topic-footer">
          <div class="topic-tags">
            <a-tag v-if="topic.node" color="green" size="mini">{{
              topic.node.name
            }}</a-tag>
            <template v-if="topic.tags && topic.tags.length">
              <a-tag v-for="tag in topic.tags" :key="tag.id" size="mini"
                >#&nbsp;{{ tag.name }}</a-tag
              >
            </template>
          </div>
          <div class="topic-actions">
            <template v-if="topic.status === 0">
              <a-link
                class="action-item"
                :href="useSiteUrl('/topic/' + topic.id)"
                target="_blank"
                >查看详情</a-link
              >
              <a-link class="action-item" @click="showComments(topic.id)">
                查看评论
              </a-link>

              <a-popconfirm
                v-if="topic.recommend"
                content="是否确定取消推荐？"
                @ok="cancelRecommend(topic.id)"
              >
                <a-button class="action-item" size="mini" type="primary"
                  >取消推荐</a-button
                >
              </a-popconfirm>

              <a-popconfirm
                v-else-if="!topic.recommend && topic.status === 0"
                content="是否确定推荐？"
                @ok="recommend(topic.id)"
              >
                <a-button class="action-item" size="mini" type="primary"
                  >推荐</a-button
                >
              </a-popconfirm>

              <a-popconfirm
                content="是否确定删除？"
                @ok="deleteSubmit(topic.id)"
              >
                <a-button class="action-item" size="mini" type="primary"
                  >删除</a-button
                >
              </a-popconfirm>
            </template>
            <template v-else-if="topic.status === 1">
              <a-popconfirm
                content="是否确定取消删除？"
                @ok="undeleteSubmit(topic.id)"
              >
                <a-button class="action-item" size="mini" type="primary"
                  >取消删除</a-button
                >
              </a-popconfirm>
            </template>
            <template v-else-if="topic.status === 2">
              <a-link
                class="action-item"
                :href="useSiteUrl('/topic/' + topic.id)"
                target="_blank"
                >查看详情</a-link
              >
              <a-link class="action-item" @click="showComments(topic.id)"
                >查看评论</a-link
              >
              <a-popconfirm
                content="是否确定删除？"
                @ok="deleteSubmit(topic.id)"
              >
                <a-button class="action-item" size="mini" type="primary"
                  >删除</a-button
                >
              </a-popconfirm>
              <a-popconfirm
                content="是否确定删除？"
                @ok="auditSubmit(topic.id)"
              >
                <a-button class="action-item" size="mini" type="primary"
                  >审核通过</a-button
                >
              </a-popconfirm>
            </template>
          </div>
        </div>
      </div>
    </template>
    <a-empty v-else />
  </div>
</template>

<script setup>
  defineProps({
    results: {
      type: Array,
      default() {
        return [];
      },
    },
  });

  const emits = defineEmits(['change']);

  const deleteSubmit = async (topicId) => {
    try {
      await axios.form(
        '/api/admin/topic/delete',
        jsonToFormData({ id: topicId })
      );
      useNotificationSuccess('删除成功');
      emits('change');
    } catch (e) {
      useHandleError(e);
    }
  };
  const undeleteSubmit = async (topicId) => {
    try {
      await axios.form(
        '/api/admin/topic/undelete',
        jsonToFormData({ id: topicId })
      );
      useNotificationSuccess('取消删除成功');
      emits('change');
    } catch (e) {
      useHandleError(e);
    }
  };
  const recommend = async (topicId) => {
    try {
      await axios.form(
        '/api/admin/topic/recommend',
        jsonToFormData({ id: topicId })
      );
      useNotificationSuccess('推荐成功');
      emits('change');
    } catch (e) {
      useHandleError(e);
    }
  };
  const cancelRecommend = async (topicId) => {
    try {
      await axios.delete(`/api/admin/topic/recommend?id=${topicId}`);
      useNotificationSuccess('取消推荐成功');
      emits('change');
    } catch (e) {
      useHandleError(e);
    }
  };
  const auditSubmit = async (topicId) => {
    try {
      await axios.form(
        '/api/admin/topic/audit',
        jsonToFormData({ id: topicId })
      );
      useNotificationSuccess('审核成功');
      emits('change');
    } catch (e) {
      useHandleError(e);
    }
  };
  const showComments = (topicId) => {
    // this.$refs.commentsDialog.show('topic', topicId);
  };
</script>

<style scoped lang="less">
  .topics {
    .topic-item {
      padding: 10px 20px;
      row-gap: 10px;
      display: flex;
      flex-direction: column;
      border-bottom: 1px solid var(--color-border-1);
      position: relative;
      .topic-status {
        position: absolute;
        right: 0;
        display: flex;
        align-items: center;
        column-gap: 10px;
      }

      .topic-header {
        display: flex;
        align-items: center;
        .topic-head-info {
          margin-left: 10px;
          display: flex;
          flex-direction: column;
          row-gap: 10px;
          .nickname {
            color: var(--color-neutral-8);
            font-size: 14px;
          }
          .topic-metas {
            display: flex;
            align-items: center;
            column-gap: 10px;

            .meta-item {
              color: var(--color-neutral-6);
              font-size: 12px;

              display: flex;
              align-items: center;
              column-gap: 3px;
            }
          }
        }
      }

      .topic-summary {
        color: var(--color-neutral-8);
        font-size: 14px;
      }

      .topic-image-list {
        display: flex;
        row-gap: 10px;
        column-gap: 10px;

        .image-item {
          cursor: pointer;
        }
      }

      .topic-footer {
        display: flex;
        align-items: center;
        justify-content: space-between;
        .topic-tags {
          display: flex;
          align-items: center;
          row-gap: 10px;
          column-gap: 10px;
        }
        .topic-actions {
          display: flex;
          align-items: center;
          row-gap: 10px;
          column-gap: 10px;
        }
      }
    }
  }
</style>
