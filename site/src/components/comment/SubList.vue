<template>
  <div class="replies">
    <div v-for="comment in replies.results" :key="comment.id" class="comment">
      <div class="comment-item-left">
        <my-avatar :user="comment.user" :size="24" has-border />
      </div>
      <div class="comment-item-main">
        <div class="comment-meta">
          <div>
            <nuxt-link
              :to="`/user/${comment.user.id}`"
              class="comment-nickname"
            >
              {{ comment.user.nickname }}
            </nuxt-link>
            <template v-if="comment.quote">
              &nbsp;<span>回复</span>&nbsp;
              <nuxt-link
                :to="`/user/${comment.quote.user.id}`"
                class="comment-nickname"
              >
                {{ comment.quote.user.nickname }}
              </nuxt-link>
            </template>
          </div>
          <time class="comment-time">{{
            usePrettyDate(comment.createTime)
          }}</time>
        </div>
        <div class="comment-content-wrapper">
          <div v-if="comment.content" class="comment-content content">
            <div v-text="comment.content" />
          </div>
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

          <div v-if="comment.quote" class="comment-quote">
            <div
              class="comment-quote-content content"
              v-html="comment.quote.content"
            />
            <div
              v-if="comment.quote.imageList && comment.quote.imageList.length"
              class="comment-quote-image-list"
            >
              <img
                v-for="(image, imageIndex) in comment.imageList"
                :key="imageIndex"
                :src="image.url"
              />
            </div>
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
            :class="{ active: reply.quoteId === comment.id }"
            @click="switchShowReply(comment)"
          >
            <i class="iconfont icon-comment" />
            <span>{{
              reply.quoteId === comment.id ? "取消评论" : "评论"
            }}</span>
          </div>
        </div>
        <div v-if="reply.quoteId === comment.id" class="comment-reply-form">
          <text-editor
            :ref="`editor${comment.id}`"
            v-model:content="reply.value.content"
            v-model:imageList="reply.value.imageList"
            :height="80"
            @submit="submitReply(comment)"
          />
        </div>
      </div>
    </div>
    <div v-if="replies.hasMore === true" class="comment-more">
      <a @click="loadMore">
        <span>查看更多回复</span>
        <i class="iconfont icon-right" />
      </a>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    commentId: {
      type: Number,
      required: true,
    },
    data: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      replies: this.data,
      reply: {
        quoteId: 0,
        value: {
          content: "",
          imageList: [],
        },
      },
    };
  },
  computed: {
    user() {
      const userStore = useUserStore();
      return userStore.user;
    },
  },
  methods: {
    async loadMore() {
      const ret = await useHttpGet("/api/comment/replies", {
        params: {
          commentId: this.commentId,
          cursor: this.replies.cursor,
        },
      });
      this.replies.cursor = ret.cursor;
      this.replies.hasMore = ret.hasMore;
      this.replies.results.push(...ret.results);
    },
    async like(comment) {
      try {
        if (comment.liked) {
          await useHttpPost(
            "/api/like/unlike",
            useJsonToForm({
              entityType: "comment",
              entityId: comment.id,
            })
          );
          comment.liked = false;
          comment.likeCount = comment.likeCount > 0 ? comment.likeCount - 1 : 0;
          useMsgSuccess("已取消点赞");
        } else {
          await useHttpPost(
            "/api/like/like",
            useJsonToForm({
              entityType: "comment",
              entityId: comment.id,
            })
          );
          comment.liked = true;
          comment.likeCount = comment.likeCount + 1;
          useMsgSuccess("点赞成功");
        }
      } catch (e) {
        useCatchError(e);
      }
    },
    switchShowReply(comment) {
      if (!this.user) {
        useMsgSignIn();
        return;
      }

      if (this.reply.quoteId === comment.id) {
        this.hideReply(comment);
      } else {
        this.reply.quoteId = comment.id;
        setTimeout(() => {
          const refs = this.$refs[`editor${comment.id}`];
          if (refs && refs.length > 0) {
            refs[0].focus();
          }
        }, 100);
      }
    },
    hideReply(comment) {
      this.reply.quoteId = 0;
      this.reply.value.content = "";
      this.reply.value.imageList = [];
    },
    async submitReply(parent) {
      try {
        const ret = await useHttpPost(
          "/api/comment/create",
          useJsonToForm({
            entityType: "comment",
            entityId: this.commentId,
            quoteId: this.reply.quoteId,
            content: this.reply.value.content,
            imageList:
              this.reply.value.imageList && this.reply.value.imageList.length
                ? JSON.stringify(this.reply.value.imageList)
                : "",
          })
        );
        this.hideReply(parent);
        this.$emit("reply", ret);
        useMsgSuccess("发布成功");
      } catch (e) {
        useCatchError(e);
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.replies {
  margin-top: 10px;
  font-size: 12px;
  // padding: 1px 10px;
  // background-color: var(--bg-color2);

  .comment {
    display: flex;
    padding: 8px 0;

    // &:not(:last-child) {
    //   border-bottom: 1px solid var(--border-color);
    // }

    .comment-item-main {
      flex: 1 1 auto;
      margin-left: 6px;

      .comment-meta {
        display: flex;
        justify-content: space-between;
        .comment-nickname {
          font-size: 12px;
          color: var(--text-color2);

          &:hover {
            color: var(--text-link-color);
          }
        }

        .comment-time {
          font-size: 11px;
          color: var(--text-color3);
        }
      }

      .comment-content-wrapper {
        .comment-content {
          margin-top: 5px;
          margin-bottom: 0;
          color: var(--text-color2);
          white-space: pre-wrap;
        }
        .comment-image-list {
          margin-top: 5px;

          img {
            width: 62px;
            height: 62px;
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

        .comment-quote {
          position: relative;
          background-color: var(--bg-color3);
          border: 1px solid var(--border-color2);
          color: var(--text-color3);
          padding: 0 12px;
          margin: 5px 0;
          box-sizing: border-box;
          border-radius: 4px;

          &::after {
            position: absolute;
            content: "\201D";
            font-family: Georgia, serif;
            font-size: 36px;
            font-weight: bold;
            color: var(--text-color3);
            right: 2px;
            top: -8px;
          }

          .comment-quote-content {
            margin: 5px 0;
            color: var(--text-color3);
          }

          .comment-quote-image-list {
            margin-top: 5px;

            img {
              width: 50px;
              height: 50px;
              line-height: 50px;
              cursor: pointer;
              &:not(:last-child) {
                margin-right: 4px;
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
      }

      .comment-actions {
        margin-top: 5px;
        display: flex;
        align-items: center;
        column-gap: 10px;

        .comment-action-item {
          font-size: 11px;
          cursor: pointer;
          color: var(--text-color3);
          user-select: none;
          display: flex;
          align-items: center;
          column-gap: 2px;

          &:hover {
            color: var(--text-link-color);
          }

          &.active {
            color: var(--text-link-color);
            font-weight: 500;
          }

          i {
            font-size: 11px;
          }
        }
      }

      .comment-reply-form {
        margin-top: 10px;
      }
    }
  }

  .reply {
    display: flex;
  }

  .comment-more {
    margin: 10px 10px 10px 30px;
    a {
      font-size: 13px;
      display: flex;
      align-items: center;
      color: var(--text-color2);
      &:hover {
        color: var(--text-link-color);
      }

      i {
        font-size: 13px;
        rotate: 90deg;
      }
    }
  }
}
</style>
