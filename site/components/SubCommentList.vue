<template>
  <div class="replies">
    <div
      v-for="comment in replies.results"
      :key="comment.commentId"
      class="comment"
    >
      <div class="comment-item-left">
        <avatar :user="comment.user" size="30" round has-border />
      </div>
      <div class="comment-item-main">
        <div class="comment-meta">
          <div>
            <nuxt-link
              :to="'/user/' + comment.user.id"
              class="comment-nickname"
            >
              {{ comment.user.nickname }}
            </nuxt-link>
            <template v-if="comment.quote">
              <span>回复</span>
              <nuxt-link
                :to="'/user/' + comment.quote.user.id"
                class="comment-nickname"
              >
                {{ comment.quote.user.nickname }}
              </nuxt-link>
            </template>
          </div>
          <time
            class="comment-time"
            :datetime="comment.createTime | formatDate('yyyy-MM-ddTHH:mm:ss')"
            >{{ comment.createTime | prettyDate }}</time
          >
        </div>
        <div
          v-viewer
          v-lazy-container="{ selector: 'img' }"
          class="comment-content-wrapper"
        >
          <div v-if="comment.content" class="comment-content content">
            <div v-html="comment.content"></div>
          </div>
          <div
            v-if="comment.imageList && comment.imageList.length"
            class="comment-image-list"
          >
            <img
              v-for="(image, imageIndex) in comment.imageList"
              :key="imageIndex"
              :data-src="image.url"
            />
          </div>

          <div v-if="comment.quote" class="comment-quote">
            <div
              class="comment-quote-content content"
              v-html="comment.quote.content"
            ></div>
            <div
              v-if="comment.quote.imageList && comment.quote.imageList.length"
              class="comment-quote-image-list"
            >
              <img
                v-for="(image, imageIndex) in comment.imageList"
                :key="imageIndex"
                :data-src="image.url"
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
            <i class="iconfont icon-like"></i>
            <span>{{ comment.liked ? '已赞' : '点赞' }}</span>
            <span v-if="comment.likeCount > 0">{{ comment.likeCount }}</span>
          </div>
          <div
            class="comment-action-item"
            :class="{ active: reply.quoteId === comment.commentId }"
            @click="switchShowReply(comment)"
          >
            <i class="iconfont icon-comment"></i>
            <span>{{
              reply.quoteId === comment.commentId ? '取消评论' : '评论'
            }}</span>
          </div>
        </div>
        <div
          v-if="reply.quoteId === comment.commentId"
          class="comment-reply-form"
        >
          <text-editor
            :ref="`editor${comment.commentId}`"
            v-model="reply.value"
            :height="80"
            @submit="submitReply()"
          />
        </div>
      </div>
    </div>
    <div v-if="replies.hasMore === true" class="comment-more">
      <a @click="loadMore">查看更多回复...</a>
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
      showReplyCommentId: 0,
      reply: {
        quoteId: 0,
        value: {
          content: '',
          imageList: [],
        },
      },
    }
  },
  computed: {
    user() {
      return this.$store.state.user.current
    },
  },
  methods: {
    async loadMore() {
      const ret = await this.$axios.get('/api/comment/replies', {
        params: {
          commentId: this.commentId,
          cursor: this.replies.cursor,
        },
      })
      this.replies.cursor = ret.cursor
      this.replies.hasMore = ret.hasMore
      this.replies.results.push(...ret.results)
    },
    async like(comment) {
      try {
        await this.$axios.post(`/api/comment/like/${comment.commentId}`)
        comment.liked = true
        comment.likeCount = comment.likeCount + 1
        this.$message.success('点赞成功')
      } catch (e) {
        if (e.errorCode === 1) {
          this.$msgSignIn()
        } else {
          this.$message.error(e.message || e)
        }
      }
    },
    switchShowReply(comment) {
      if (!this.user) {
        this.$msgSignIn()
        return
      }

      if (this.reply.quoteId === comment.commentId) {
        this.hideReply(comment)
      } else {
        this.reply.quoteId = comment.commentId
        setTimeout(() => {
          this.$refs[`editor${comment.commentId}`][0].focus()
        }, 0)
      }
    },
    hideReply(comment) {
      this.reply.quoteId = 0
      this.reply.value.content = ''
      this.reply.value.imageList = []
    },
    async submitReply(parent) {
      try {
        const ret = await this.$axios.post('/api/comment/create', {
          entityType: 'comment',
          entityId: this.commentId,
          quoteId: this.reply.quoteId,
          content: this.reply.value.content,
          imageList:
            this.reply.value.imageList && this.reply.value.imageList.length
              ? JSON.stringify(this.reply.value.imageList)
              : '',
        })
        this.hideReply()
        this.$emit('reply', ret)
        this.$message.success('发布成功')
      } catch (e) {
        if (e.errorCode === 1) {
          this.$msgSignIn()
        } else {
          this.$message.error(e.message || e)
        }
      }
    },
  },
}
</script>
<style lang="scss" scoped>
.replies {
  margin-top: 10px;
  padding: 1px 10px;
  font-size: 12px;
  background-color: var(--bg-color2);

  .comment {
    display: flex;
    padding: 8px 0;

    &:not(:last-child) {
      border-bottom: 1px solid var(--border-color);
    }

    .comment-item-main {
      flex: 1 1 auto;
      margin-left: 8px;

      .comment-meta {
        display: flex;
        justify-content: space-between;
        .comment-nickname {
          font-size: 12px;
          font-weight: 600;
          color: var(--text-color);

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
            content: '\201D';
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

        .comment-action-item {
          line-height: 22px;
          font-size: 11px;
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
    }
  }

  .reply {
    display: flex;
  }

  .comment-more {
    margin: 10px;
    font-size: 13px;
    font-weight: 500;
    color: var(--text-link-color);
  }
}
</style>
