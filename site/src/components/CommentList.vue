<template>
  <div class="comments">
    <load-more-async
      ref="loadMore"
      v-slot="{ results }"
      :params="{ entityType, entityId }"
      url="/api/comment/comments"
    >
      <div v-for="comment in results" :key="comment.id" class="comment">
        <div class="comment-item-left">
          <my-avatar :user="comment.user" :size="40" has-border />
        </div>
        <div class="comment-item-main">
          <div class="comment-meta">
            <nuxt-link
              :to="`/user/${comment.user.id}`"
              class="comment-nickname"
            >
              {{ comment.user.nickname }}
            </nuxt-link>
            <div class="comment-meta-right">
              <time class="comment-time">{{
                usePrettyDate(comment.createTime)
              }}</time>
              <span v-if="comment.ipLocation" class="comment-ip-area"
                >IP属地{{ comment.ipLocation }}</span
              >
            </div>
          </div>
          <div class="comment-content-wrapper">
            <template v-if="comment.content">
              <div
                v-if="comment.contentType === 'text'"
                class="comment-content content"
                v-text="comment.content"
              />
              <div
                v-else
                class="comment-content content"
                v-html="comment.content"
              />
            </template>
            <div
              v-if="comment.imageList && comment.imageList.length"
              class="comment-image-list"
            >
              <img
                v-for="(image, imageIndex) in comment.imageList"
                :key="imageIndex"
                :src="image.url"
              />
            </div>
          </div>
          <div class="comment-actions">
            <div
              class="comment-action-item"
              :class="{ active: comment.liked }"
              @click="like(comment)"
            >
              <i class="iconfont icon-like" />
              <span>{{ comment.liked ? "已赞" : "点赞" }}</span>
              <span v-if="comment.likeCount > 0">{{ comment.likeCount }}</span>
            </div>
            <div
              class="comment-action-item"
              :class="{ active: reply.commentId === comment.id }"
              @click="switchShowReply(comment)"
            >
              <i class="iconfont icon-comment" />
              <span>{{
                reply.commentId === comment.id ? "取消评论" : "评论"
              }}</span>
            </div>
          </div>
          <div v-if="reply.commentId === comment.id" class="comment-reply-form">
            <text-editor
              :ref="`editor${comment.id}`"
              v-model="reply.value"
              :height="100"
              @submit="submitReply(comment)"
            />
          </div>
          <CommentSubList
            v-if="
              comment.replies &&
              comment.replies.results &&
              comment.replies.results.length
            "
            :comment-id="comment.id"
            :data="comment.replies"
            @reply="onReply(comment, $event)"
          />
        </div>
      </div>
    </load-more-async>
  </div>
</template>

<script setup>
const props = defineProps({
  entityType: {
    type: String,
    default: "",
    required: true,
  },
  entityId: {
    type: Number,
    default: 0,
    required: true,
  },
});
const reply = reactive({
  commentId: 0,
  value: {
    content: "",
    imageList: [],
  },
});

const userStore = useUserStore();
const loadMore = ref(null);

const append = (data) => {
  if (loadMore.value) {
    // console.log(loadMore.value);
    // console.log(loadMore.value.unshiftResults);
    // loadMore.value.unshiftResults(data);
    loadMore.value.refresh();
  }
};

const like = async (comment) => {
  try {
    if (comment.liked) {
      await useHttpPostForm("/api/like/unlike", {
        body: {
          entityType: "comment",
          entityId: comment.id,
        },
      });
      comment.liked = false;
      comment.likeCount = comment.likeCount > 0 ? comment.likeCount - 1 : 0;
      useMsgSuccess("已取消点赞");
    } else {
      await useHttpPostForm("/api/like/like", {
        body: {
          entityType: "comment",
          entityId: comment.id,
        },
      });
      comment.liked = true;
      comment.likeCount = comment.likeCount + 1;
      useMsgSuccess("点赞成功");
    }
  } catch (e) {
    useCatchError(e);
  }
};

const switchShowReply = (comment) => {
  if (!userStore.user) {
    useMsgSignIn();
    return;
  }

  if (reply.commentId === comment.id) {
    hideReply(comment);
  } else {
    reply.commentId = comment.id;
    // // TODO
    // setTimeout(() => {
    //   this.$refs[`editor${comment.id}`][0].focus();
    // }, 0);
  }
};

const hideReply = (comment) => {
  reply.commentId = 0;
  reply.value.content = "";
  reply.value.imageList = [];
};

const submitReply = async (parent) => {
  try {
    const ret = await useHttpPostForm("/api/comment/create", {
      body: {
        entityType: "comment",
        entityId: parent.id,
        content: reply.value.content,
        imageList:
          reply.value.imageList && reply.value.imageList.length
            ? JSON.stringify(reply.value.imageList)
            : "",
      },
    });
    hideReply();
    appendReply(parent, ret);
    useMsgSuccess("发布成功");
  } catch (e) {
    useCatchError(e);
  }
};

const onReply = (parent, comment) => {
  appendReply(parent, comment);
};

const appendReply = (parent, comment) => {
  if (parent.replies && parent.replies.results) {
    parent.replies.results.push(comment);
  } else {
    parent.replies = {
      results: [comment],
    };
  }
};

defineExpose({
  append,
});
</script>

<style scoped lang="scss">
.comments {
  padding: 10px;
  font-size: 14px;

  .comment {
    display: flex;
    padding: 10px 0;

    &:not(:last-child) {
      border-bottom: 1px solid var(--border-color);
    }

    .comment-item-main {
      flex: 1 1 auto;
      margin-left: 16px;

      .comment-meta {
        display: flex;
        justify-content: space-between;
        .comment-nickname {
          font-size: 14px;
          font-weight: 600;
          color: var(--text-color);

          &:hover {
            color: var(--text-link-color);
          }
        }

        .comment-meta-right {
          .comment-time {
            font-size: 13px;
            color: var(--text-color3);
          }
          .comment-ip-area {
            font-size: 13px;
            color: var(--text-color3);
            margin-left: 10px;
          }
        }
      }

      .comment-content-wrapper {
        .comment-content {
          margin-top: 10px;
          margin-bottom: 0;
          color: var(--text-color);
          white-space: pre-wrap;
        }
        .comment-image-list {
          margin-top: 10px;

          img {
            width: 72px;
            height: 72px;
            line-height: 72px;
            cursor: pointer;
            &:not(:last-child) {
              margin-right: 8px;
            }

            object-fit: cover;
            transition: all 0.5s ease-out 0.1s;

            &:hover {
              transform: matrix(1.04, 0, 0, 1.04, 0, 0);
              backface-visibility: hidden;
            }
          }
        }
      }

      .comment-actions {
        margin-top: 10px;
        display: flex;
        align-items: center;

        .comment-action-item {
          line-height: 22px;
          font-size: 13px;
          cursor: pointer;
          color: var(--text-color3);
          user-select: none;

          &:hover {
            color: var(--text-link-color);
          }

          &.active {
            color: var(--text-link-color);
            font-weight: 500;
          }

          &:not(:last-child) {
            margin-right: 16px;
          }
        }
      }

      .comment-reply-form {
        margin-top: 10px;
      }

      .comment-replies {
        margin-top: 10px;
        // padding: 10px;
        background-color: var(--bg-color2);
      }
    }
  }

  .reply {
    display: flex;
  }
}
</style>
