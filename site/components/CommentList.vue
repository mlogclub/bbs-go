<template>
  <div class="comments">
    <load-more
      v-if="commentsPage"
      ref="commentsLoadMore"
      v-slot="{ results }"
      :init-data="commentsPage"
      :params="{ entityType: entityType, entityId: entityId }"
      url="/api/comment/list"
    >
      <div v-for="comment in results" :key="comment.commentId" class="comment">
        <div class="comment-item-left">
          <avatar :user="comment.user" size="40" round has-border />
        </div>
        <div class="comment-item-right">
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
            <div class="comment-action-item">
              <i class="iconfont icon-like"></i>
              <span>点赞</span>
            </div>
            <div class="comment-action-item">
              <i class="iconfont icon-comment"></i>
              <span>评论</span>
            </div>
          </div>
        </div>
      </div>
    </load-more>
  </div>
</template>

<script>
export default {
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
    reply(quote) {
      if (!this.isLogin) {
        this.$toSignin()
      }
      this.$emit('reply', quote)
    },
    cancelReply() {
      this.quote = null
    },
  },
}
</script>

<style scoped lang="scss">
.comments {
  margin: 20px;
  .comment {
    display: flex;
    padding: 16px 0;
    .comment-item-left {
    }
    .comment-item-right {
      flex: 1 1 auto;
      margin-left: 16px;

      .comment-meta {
        display: flex;
        justify-content: space-between;
        .comment-nickname {
          font-size: 15px;
          font-weight: 600;
          color: var(--text-color);

          &:hover {
            color: var(--text-link-color);
          }
        }

        .comment-time {
          font-size: 14px;
          color: var(--text-color3);
        }
      }

      .comment-content-wrapper {
        .comment-content {
          margin-top: 10px;
          margin-bottom: 0;
          color: var(--text-color3);
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
          font-size: 14px;
          cursor: pointer;
          color: var(--text-color3);

          &:hover {
            color: var(--text-link-color);
          }

          &:not(:last-child) {
            margin-right: 16px;
          }
        }
      }
    }
  }
}
</style>
