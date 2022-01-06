<template>
  <div class="comments">
    <load-more
      v-if="commentsPage"
      ref="commentsLoadMore"
      v-slot="{ results }"
      :init-data="commentsPage"
      :params="{ entityType: entityType, entityId: entityId }"
      url="/api/comment/comments"
    >
      <div v-for="comment in results" :key="comment.commentId" class="comment">
        <div class="comment-item-left">
          <avatar :user="comment.user" size="40" round has-border />
        </div>
        <div class="comment-item-main">
          <div class="comment-meta">
            <nuxt-link
              :to="'/user/' + comment.user.id"
              class="comment-nickname"
            >
              {{ comment.user.nickname }}
            </nuxt-link>
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
            <div
              v-if="comment.content"
              class="comment-content content"
              v-html="comment.content"
            ></div>
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
              :class="{ active: reply.commentId === comment.commentId }"
              @click="switchShowReply(comment)"
            >
              <i class="iconfont icon-comment"></i>
              <span>{{
                reply.commentId === comment.commentId ? '取消评论' : '评论'
              }}</span>
            </div>
          </div>
          <div
            v-if="reply.commentId === comment.commentId"
            class="comment-reply-form"
          >
            <text-editor
              :ref="`editor${comment.commentId}`"
              v-model="reply.value"
              :height="100"
              @submit="submitReply(comment)"
            />
          </div>
          <sub-comment-list
            v-if="
              comment.replies &&
              comment.replies.results &&
              comment.replies.results.length
            "
            :comment-id="comment.commentId"
            :data="comment.replies"
            @reply="onReply(comment, $event)"
          />
        </div>
      </div>
    </load-more>
  </div>
</template>

<script>
import SubCommentList from './SubCommentList.vue'
export default {
  components: { SubCommentList },
  props: {
    entityType: {
      type: String,
      default: '',
      required: true,
    },
    entityId: {
      type: Number,
      default: 0,
      required: true,
    },
    commentsPage: {
      type: Object,
      default() {
        return {}
      },
    },
  },
  data() {
    return {
      showReplyCommentId: 0,
      reply: {
        commentId: 0,
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
    isLogin() {
      return this.$store.state.user.current != null
    },
  },
  methods: {
    append(data) {
      if (data) {
        this.$refs.commentsLoadMore.unshiftResults(data)
      }
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

      if (this.reply.commentId === comment.commentId) {
        this.hideReply(comment)
      } else {
        this.reply.commentId = comment.commentId
        setTimeout(() => {
          this.$refs[`editor${comment.commentId}`][0].focus()
        }, 0)
      }
    },
    hideReply(comment) {
      this.reply.commentId = 0
      this.reply.value.content = ''
      this.reply.value.imageList = []
    },
    async submitReply(parent) {
      try {
        const ret = await this.$axios.post('/api/comment/create', {
          entityType: 'comment',
          entityId: parent.commentId,
          content: this.reply.value.content,
          imageList:
            this.reply.value.imageList && this.reply.value.imageList.length
              ? JSON.stringify(this.reply.value.imageList)
              : '',
        })
        this.hideReply()
        this.appendReply(parent, ret)
        this.$message.success('发布成功')
      } catch (e) {
        if (e.errorCode === 1) {
          this.$msgSignIn()
        } else {
          this.$message.error(e.message || e)
        }
      }
    },
    onReply(parent, comment) {
      this.appendReply(parent, comment)
    },
    appendReply(parent, comment) {
      if (parent.replies && parent.replies.results) {
        parent.replies.results.push(comment)
      } else {
        parent.replies = {
          results: [comment],
        }
      }
    },
  },
}
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

        .comment-time {
          font-size: 13px;
          color: var(--text-color3);
        }
      }

      .comment-content-wrapper {
        .comment-content {
          margin-top: 10px;
          margin-bottom: 0;
          color: var(--text-color);
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
