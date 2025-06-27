<template>
  <section v-if="topic" class="main">
    <div v-if="isPending" class="container main-container">
      <div class="notification is-warning" style="width: 100%; margin: 20px 0">
        {{ t('pages.topic.detail.pending') }}
      </div>
    </div>
    <div class="container main-container left-main size-360">
      <div class="left-container">
        <div class="main-content no-padding no-bg">
          <article class="topic-detail">
            <side-action-bar
              class="float-bar"
              entity-type="topic"
              :entity-id="topic.id"
              :liked="topic.liked"
              :like-count="topic.likeCount"
              :comment-count="topic.commentCount"
              :favorited="topic.favorited"
            />
            <div class="topic-header">
              <div class="topic-header-left">
                <my-avatar :user="topic.user" :size="45" />
              </div>
              <div class="topic-header-center">
                <div class="topic-nickname">
                  <nuxt-link :to="`/user/${topic.user.id}`">
                    {{ topic.user.nickname }}
                  </nuxt-link>
                </div>
                <div class="topic-meta">
                  <span class="meta-item">
                    {{ t('pages.topic.detail.publishedAt') }}
                    <time>{{ usePrettyDate(topic.createTime, t) }}</time>
                  </span>
                  <span v-if="topic.ipLocation" class="meta-item">
                    {{ t('pages.topic.detail.ipLocation') }}{{ topic.ipLocation }}
                  </span>
                </div>
              </div>
              <div class="topic-header-right">
                <topic-manage-menu v-model="topic" />
              </div>
            </div>

            <!-- 内容 -->
            <div
              class="topic-content content"
              :class="{
                'topic-tweet': topic.type === 1,
              }"
            >
              <h1 v-if="topic.title" class="topic-title">
                {{ topic.title }}
              </h1>
              <div
                class="topic-content-detail line-numbers"
                v-html="topic.content"
              />
              <ul
                v-if="topic.imageList && topic.imageList.length"
                class="topic-image-list"
              >
                <li v-for="(image, index) in topic.imageList" :key="index">
                  <div class="image-item">
                    <el-image
                      :src="image.preview"
                      :preview-src-list="imageUrls"
                      :initial-index="index"
                    />
                  </div>
                </li>
              </ul>
              <div
                v-if="hideContent && hideContent.exists"
                class="topic-content-detail hide-content"
              >
                <div v-if="hideContent.show" class="widget has-border">
                  <div class="widget-header">
                    <span>
                      <i class="iconfont icon-lock" />
                      <span>&nbsp;{{ t('pages.topic.detail.hideContent') }}</span>
                    </span>
                  </div>
                  <div class="widget-content" v-html="hideContent.content" />
                </div>
                <div v-else class="hide-content-tip">
                  <i class="iconfont icon-lock" />
                  <span>{{ t('pages.topic.detail.hideContentTip') }}</span>
                </div>
              </div>
            </div>

            <!-- 节点、标签 -->
            <div class="topic-tags">
              <nuxt-link
                v-if="topic.node"
                :to="`/topics/node/${topic.node.id}`"
                class="topic-tag"
              >
                {{ topic.node.name }}
              </nuxt-link>
              <nuxt-link
                v-for="tag in topic.tags"
                :key="tag.id"
                :to="`/topics/tag/${tag.id}`"
                class="topic-tag"
              >
                #{{ tag.name }}
              </nuxt-link>
            </div>

            <!-- 点赞用户列表 -->
            <div v-if="likeUsers && likeUsers.length" class="topic-like-users">
              <my-avatar
                v-for="likeUser in likeUsers"
                :key="likeUser.id"
                :user="likeUser"
                :size="24"
                has-border
              />
              <span class="like-count">{{ topic.likeCount }}</span>
            </div>

            <!-- 功能按钮 -->
            <div class="topic-actions">
              <div class="action disabled">
                <i class="action-icon iconfont icon-view" />
                <div class="action-text">
                  <span>{{ t('pages.topic.detail.view') }}</span>
                  <span v-if="topic.viewCount > 0" class="action-text">
                    ({{ topic.viewCount }})
                  </span>
                </div>
              </div>
              <div class="action" @click="like(topic)">
                <i
                  class="action-icon iconfont icon-like"
                  :class="{ 'checked-icon': liked }"
                />
                <div class="action-text">
                  <span>{{ t('pages.topic.detail.like') }}</span>
                  <span v-if="topic.likeCount > 0">
                    ({{ topic.likeCount }})
                  </span>
                </div>
              </div>
              <div class="action" @click="addFavorite(topic.id)">
                <i
                  class="action-icon iconfont icon-favorite"
                  :class="{
                    'icon-has-favorite': topic.favorited,
                    'icon-favorite': !topic.favorited,
                    'checked-icon': topic.favorited,
                  }"
                />
                <div class="action-text">
                  <span>{{ t('pages.topic.detail.favorite') }}</span>
                </div>
              </div>
            </div>
          </article>

          <!-- 评论 -->
          <comment
            :entity-id="topic.id"
            :comment-count="topic.commentCount"
            entity-type="topic"
            @created="commentCreated"
          />
        </div>
      </div>
      <div class="right-container">
        <user-info :user="topic.user" />
      </div>
    </div>
  </section>
</template>

<script setup>
import { useI18n } from 'vue-i18n';

const route = useRoute();
const { t } = useI18n();

const { data: topic } = await useMyFetch(`/api/topic/${route.params.id}`);

const { data: liked } = await useMyFetch("/api/like/liked", {
  params: {
    entityType: "topic",
    entityId: route.params.id,
  },
});

const { data: likeUsers, refresh: refreshLikeUsers } = await useMyFetch(
  `/api/topic/recentlikes/${route.params.id}`
);

const { data: hideContent, refresh: refreshHideContent } = await useMyFetch(
  `/api/topic/hide_content?topicId=${route.params.id}`
);

const imageUrls = computed(() => {
  if (!topic.value.imageList || !topic.value.imageList.length) {
    return [];
  }
  const ret = [];
  for (let i = 0; i < topic.value.imageList.length; i++) {
    ret.push(topic.value.imageList[i].url);
  }
  return ret;
});

useHead({
  title: useTopicSiteTitle(topic.value),
});

const isPending = computed(() => {
  return topic.value?.status === 2;
});

async function like() {
  try {
    if (liked.value) {
      await useHttpPost(
        "/api/like/unlike",
        useJsonToForm({
          entityType: "topic",
          entityId: topic.value.id,
        })
      );
      liked.value = false;
      topic.value.likeCount =
        topic.value.likeCount > 0 ? topic.value.likeCount - 1 : 0;

      useMsgSuccess(t('pages.topic.detail.likeSuccess'));
      await refreshLikeUsers();
    } else {
      await useHttpPost(
        "/api/like/like",
        useJsonToForm({
          entityType: "topic",
          entityId: topic.value.id,
        })
      );
      liked.value = true;
      topic.value.likeCount++;

      useMsgSuccess(t('pages.topic.detail.likeSuccess'));
      await refreshLikeUsers();
    }
  } catch (e) {
    useCatchError(e);
  }
}

async function addFavorite(topicId) {
  try {
    if (topic.value.favorited) {
      await useHttpPost(
        "/api/favorite/delete",
        useJsonToForm({
          entityType: "topic",
          entityId: topicId,
        })
      );
      topic.value.favorited = false;
      useMsgSuccess(t('pages.topic.detail.favoriteSuccess'));
    } else {
      await useHttpPost(
        "/api/favorite/add",
        useJsonToForm({
          entityType: "topic",
          entityId: topicId,
        })
      );
      topic.value.favorited = true;
      useMsgSuccess(t('pages.topic.detail.favoriteSuccess'));
    }
  } catch (e) {
    useCatchError(e);
  }
}

async function commentCreated() {
  refreshHideContent();
}
</script>

<style lang="scss" scoped>
.topic-detail {
  margin-bottom: 20px;
  background-color: var(--bg-color);
  border-radius: var(--border-radius);

  .float-bar {
    position: fixed;
    margin-left: -58px;
    top: 300px;

    @media screen and (max-width: 1300px) {
      display: none;
    }
  }

  .topic-header,
  .topic-content,
  .topic-tags,
  .topic-like-users,
  .topic-actions {
    margin: 0 16px 16px 16px;
  }

  .topic-header {
    display: flex;

    .topic-header-left {
      margin: 10px 10px 0 0;
    }

    .topic-header-center {
      margin: 10px 10px 0 0;
      width: 100%;

      .topic-nickname a {
        color: var(--text-color2);
        font-size: 16px;
        font-weight: bold;
        overflow: hidden;
      }

      .topic-meta {
        position: relative;
        font-size: 12px;
        line-height: 24px;
        color: var(--text-color3);
        margin-top: 5px;

        span.meta-item {
          font-size: 12px;

          &:not(:last-child) {
            margin-right: 8px;
          }
        }
      }
    }

    .topic-header-right {
      margin-top: 10px;
      min-width: max-content;
    }
  }

  .topic-content {
    font-size: 15px;
    color: var(--text-color);
    white-space: normal;
    word-break: break-all;
    word-wrap: break-word;
    padding-top: 0 !important;

    .topic-title {
      font-weight: 700;
      font-size: 20px;
      word-wrap: break-word;
      word-break: normal;
      border-bottom: 1px solid var(--border-color4);
      padding-bottom: 10px;
    }

    .topic-content-detail {
      font-size: 16px;
      line-height: 24px;
      word-wrap: break-word;
      -webkit-font-smoothing: antialiased;
    }

    &.topic-tweet {
      .topic-content-detail {
        white-space: pre-line;
      }
    }

    .topic-image-list {
      margin-left: 0;
      margin-top: 10px;

      li {
        cursor: pointer;
        border: 1px dashed var(--border-color2);
        text-align: center;

        display: inline-block;
        vertical-align: middle;
        margin: 0 8px 8px 0;
        background-color: var(--bg-color2);
        background-size: 32px 32px;
        background-position: 50%;
        background-repeat: no-repeat;
        overflow: hidden;
        position: relative;

        .image-item {
          display: block;
          overflow: hidden;
          // transform-style: preserve-3d;

          & > img {
            width: 100%;
            height: 100%;
            object-fit: cover;
            transition: all 0.5s ease-out 0.1s;

            &:hover {
              transform: matrix(1.04, 0, 0, 1.04, 0, 0);
              backface-visibility: hidden;
            }
          }
        }

        /* 只有一个图片时 */
        &:first-child:nth-last-child(1) {
          width: 210px;
          height: 210px;
          line-height: 210px;

          .image-item {
            width: 210px;
            height: 210px;
          }
        }

        /* 只有两个图片时 */
        &:first-child:nth-last-child(2),
        &:first-child:nth-last-child(2) ~ li {
          width: 180px;
          height: 180px;
          line-height: 180px;

          .image-item {
            width: 180px;
            height: 180px;
          }
        }

        /*大于两个图片时*/
        &:first-child:nth-last-child(n + 3),
        &:first-child:nth-last-child(n + 3) ~ li {
          width: 120px;
          height: 120px;
          line-height: 120px;

          .image-item {
            width: 120px;
            height: 120px;
          }
        }
      }
    }

    .hide-content {
      margin: 20px 0;

      .widget-header {
        span {
          font-weight: 500;
        }
      }

      .hide-content-tip {
        border: 1px solid var(--border-hover-color);
        border-radius: 2px;
        padding: 6px 12px;
        font-size: 14px;
        color: #3273dc;
      }
    }
  }

  .topic-tags {
    .topic-tag {
      padding: 2px 8px;
      justify-content: center;
      align-items: center;
      border-radius: 12.5px;
      margin-right: 10px;
      background: var(--bg-color2);
      border: 1px solid var(--border-color);
      color: var(--text-color3);
      font-size: 12px;

      &:hover {
        color: var(--text-link-color);
        background: var(--bg-color);
      }
    }
  }

  .topic-like-users {
    width: 80%;
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    margin-bottom: 10px;

    .avatar-a {
      margin-right: -3px;
    }

    .like-count {
      margin-left: 8px;
      font-size: 14px;
      color: var(--text-color);
    }
  }

  .topic-actions {
    padding: 10px 0;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-top: 1px solid var(--border-color4);

    .action {
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: center;
      background: var(--bg-color);
      cursor: pointer;
      color: var(--text-color3);
      font-size: 14px;

      .checked-icon {
        color: var(--color-red);
      }

      &.disabled {
        cursor: not-allowed;

        &:hover {
          color: var(--text-color3);

          > .action-icon {
            fill: var(--text-color3);
          }
        }
      }

      > .action-icon {
        fill: #8590a6;
      }

      .action-text {
        color: var(--text-color);
        margin-left: 5px;
      }

      &:hover {
        color: var(--text-link-color);

        > .action-icon {
          fill: var(--text-link-color);
        }
      }
    }
  }
}
</style>
