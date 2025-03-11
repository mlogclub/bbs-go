<template>
  <ul class="topic-list">
    <li v-for="topic in topics" :key="topic.id" class="topic-item">
      <div class="topic-header">
        <div class="topic-header-main">
          <my-avatar :user="topic.user" />
          <div class="topic-userinfo">
            <a :href="`/user/${topic.user.id}`" class="topic-nickname">
              {{ topic.user.nickname }}
            </a>
            <div class="topic-time">
              发布于{{ usePrettyDate(topic.createTime) }}
            </div>
          </div>
        </div>
        <div class="topic-header-right">
          <span v-if="showSticky && topic.sticky" class="topic-sticky-icon"
            >置顶</span
          >
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
        <div class="topic-tags">
          <nuxt-link
            v-if="topic.node"
            class="topic-tag"
            target="_blank"
            :to="`/topics/node/${topic.node.id}`"
            :alt="topic.node.name"
          >
            <img v-if="topic.node.logo" :src="topic.node.logo" />
            <span>{{ topic.node.name }}</span>
          </nuxt-link>
        </div>

        <div class="topic-actions">
          <div
            class="btn EASE"
            :class="{ liked: topic.liked }"
            @click="like(topic)"
          >
            <i class="iconfont icon-like" />
            <span v-if="topic.likeCount > 0">{{ topic.likeCount }}</span>
          </div>
          <div class="btn EASE" @click="toTopicDetail(topic.id)">
            <i class="iconfont icon-comment" />
            <span v-if="topic.commentCount > 0">{{ topic.commentCount }}</span>
          </div>
          <!-- <div class="btn EASE" @click="toTopicDetail(topic.id)">
            <i class="iconfont icon-view" />
            <span v-if="topic.viewCount > 0">{{ topic.viewCount }}</span>
          </div> -->
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
<style lang="scss" scoped>
.topic-list {
  .topic-item {
    padding: 16px 32px;
    position: relative;
    overflow: hidden;
    border-radius: 3px;

    &:not(:last-child):after {
      position: absolute;
      content: "";
      bottom: 0;
      left: 32px;
      right: 32px;
      height: 1px;
      background: var(--border-color4);
    }

    .topic-header {
      display: flex;
      align-items: center;
      justify-content: space-between;

      .topic-header-main {
        display: flex;
        align-items: center;

        .topic-userinfo {
          margin-left: 10px;
          .topic-nickname {
            font-weight: 500;
            font-size: 14px;
            color: var(--text-color);
          }

          .topic-time {
            margin-top: 3px;
            font-size: 13px;
            color: var(--text-color3);
          }
        }
      }

      .topic-header-right {
        .topic-sticky-icon {
          font-size: 13px;
          line-height: 13px;
          color: #ff7827;
          background: #ffe7d9;
          border-radius: 2px;
          padding: 3px 6px;
          white-space: nowrap;
        }
      }
    }

    .topic-content {
      margin-top: 6px;
      .topic-title {
        display: inline-block;
        margin-bottom: 6px;
        word-wrap: break-word;
        word-break: break-all;
        width: 100%;

        a {
          font-size: 16px;
          font-weight: 600;
          color: var(--text-color);

          &:hover {
            //color: #3273dc;
            text-decoration: underline;
          }
        }
      }

      .topic-summary {
        font-size: 14px;
        margin-bottom: 6px;
        width: 100%;
        text-decoration: none;
        color: var(--text-color3);
        word-wrap: break-word;

        overflow: hidden;
        display: -webkit-box;
        -webkit-box-orient: vertical;
        -webkit-line-clamp: 3;
        text-align: justify;
        word-break: break-all;
        text-overflow: ellipsis;
      }

      &.topic-tweet {
        .topic-summary {
          color: var(--text-color);
          white-space: pre-line;
        }
      }

      .topic-image-list {
        margin-top: 10px;

        li {
          cursor: pointer;
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
            transform-style: preserve-3d;

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
    }

    .topic-bottom {
      display: flex;
      align-items: center;
      justify-content: space-between;

      .topic-tags {
        display: flex;

        .topic-tag {
          display: flex;
          justify-content: center;
          align-items: center;
          padding: 4px 10px;
          border-radius: 18px;
          background: var(--bg-color6);
          color: var(--text-color3);
          font-size: 13px;

          &:hover {
            color: var(--text-color3-hover);
            background: var(--bg-color6-hover);
          }

          img {
            display: block;
            width: 20px;
            height: 20px;
            margin: 0 4px 0 0;
            border-radius: 50%;
            object-fit: cover;
          }
        }
      }

      .topic-actions {
        display: flex;
        align-items: center;
        margin-top: 6px;
        font-size: 12px;
        user-select: none;

        .btn {
          color: var(--text-color3);
          cursor: pointer;
          display: flex;
          align-items: center;

          &:not(:last-child) {
            margin-right: 20px;
          }

          &:hover {
            color: var(--text-link-color);
          }

          i {
            margin-right: 3px;
            font-size: 16px;
            position: relative;
          }

          span {
            line-height: 24px;
            font-size: 15px;
          }

          &.liked {
            color: var(--color-red) !important;
          }
        }
      }
    }
  }
}

@media screen and (max-width: 768px) {
  .topic-list {
    .topic-item {
      padding: 12px 12px;

      &:after {
        left: 12px !important;
        right: 12px !important;
      }
    }
  }
}
</style>
