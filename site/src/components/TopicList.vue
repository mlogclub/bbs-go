<template>
  <ul class="topic-list">
    <li v-for="topic in topics" :key="topic.id" class="topic-item">
      <div class="topic-avatar" :title="topic.user.nickname">
        <my-avatar :user="topic.user" />
      </div>
      <div class="topic-main-content">
        <div class="topic-top">
          <div class="topic-userinfo">
            <my-avatar
              class="topic-inline-avatar"
              :user="topic.user"
              :size="20"
            />
            <nuxt-link :to="`/user/${topic.user.id}`" target="_blank">
              {{ topic.user.nickname }}
            </nuxt-link>
            <span v-if="showSticky && topic.sticky" class="topic-sticky-icon"
              >置顶</span
            >
          </div>
          <div class="topic-time">
            发布于{{ usePrettyDate(topic.createTime) }}
          </div>
        </div>
        <div class="topic-content" :class="{ 'topic-tweet': topic.type === 1 }">
          <template v-if="topic.type === 0">
            <h1 class="topic-title">
              <nuxt-link :to="`/topic/${topic.id}`" target="_blank">
                {{ topic.title }}
              </nuxt-link>
            </h1>
            <nuxt-link
              :to="`/topic/${topic.id}`"
              class="topic-summary"
              target="_blank"
            >
              {{ topic.summary }}
            </nuxt-link>
          </template>
          <template v-if="topic.type === 1">
            <nuxt-link
              v-if="topic.content"
              :to="`/topic/${topic.id}`"
              class="topic-summary"
              target="_blank"
            >
              {{ topic.content }}
            </nuxt-link>
            <ul
              v-if="topic.imageList && topic.imageList.length"
              class="topic-image-list"
            >
              <li v-for="(image, index) in topic.imageList" :key="index">
                <nuxt-link
                  :to="`/topic/${topic.id}`"
                  class="image-item"
                  target="_blank"
                >
                  <img :src="image.preview" />
                </nuxt-link>
              </li>
            </ul>
          </template>
        </div>
        <div class="topic-bottom">
          <div class="topic-handlers">
            <div
              class="btn"
              :class="{ liked: topic.liked }"
              @click="like(topic)"
            >
              <i class="iconfont icon-like" />{{ topic.liked ? "已赞" : "赞" }}
              <span v-if="topic.likeCount > 0">{{ topic.likeCount }}</span>
            </div>
            <div class="btn" @click="toTopicDetail(topic.id)">
              <i class="iconfont icon-comment" />评论
              <span v-if="topic.commentCount > 0">{{
                topic.commentCount
              }}</span>
            </div>
            <div class="btn" @click="toTopicDetail(topic.id)">
              <i class="iconfont icon-read" />浏览
              <span v-if="topic.viewCount > 0">{{ topic.viewCount }}</span>
            </div>
          </div>
          <div class="topic-tags">
            <nuxt-link
              v-if="topic.node"
              class="topic-tag"
              target="_blank"
              :to="`/topics/node/${topic.node.id}`"
              :alt="topic.node.name"
            >
              {{ topic.node.name }}
            </nuxt-link>
          </div>
        </div>
      </div>
    </li>
  </ul>
</template>
<script>
export default {
  props: {
    topics: {
      type: Array,
      default() {
        return [];
      },
      required: false,
    },
    showAvatar: {
      type: Boolean,
      default: true,
    },
    showSticky: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    async like(topic) {
      try {
        if (topic.liked) {
          await useHttpPostForm("/api/like/unlike", {
            body: {
              entityType: "topic",
              entityId: topic.id,
            },
          });
          topic.liked = false;
          topic.likeCount = topic.likeCount > 0 ? topic.likeCount - 1 : 0;
          useMsgSuccess("已取消点赞");
        } else {
          await useHttpPostForm("/api/like/like", {
            body: {
              entityType: "topic",
              entityId: topic.id,
            },
          });
          topic.liked = true;
          topic.likeCount++;
          useMsgSuccess("点赞成功");
        }
      } catch (e) {
        useCatchError(e);
      }
    },
    toTopicDetail(topicId) {
      useLinkTo(`/topic/${topicId}`);
    },
  },
};
</script>
<style lang="scss" scoped></style>
