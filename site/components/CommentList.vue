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
      <ul>
        <li
          v-for="(comment, index) in results"
          :key="comment.commentId"
          class="comment"
          itemprop="comment"
          itemscope
          itemtype="http://schema.org/Comment"
        >
          <adsbygoogle
            v-if="showAd && (index + 1) % 3 === 0 && index !== 0"
            ad-slot="4980294904"
            ad-format="fluid"
            ad-layout-key="-ht-19-1m-3j+mu"
          />
          <div class="comment-avatar">
            <avatar :user="comment.user" size="35" />
          </div>
          <div class="comment-meta">
            <span
              class="comment-nickname"
              itemprop="creator"
              itemscope
              itemtype="http://schema.org/Person"
            >
              <nuxt-link :to="'/user/' + comment.user.id" itemprop="name">
                {{ comment.user.nickname }}
              </nuxt-link>
            </span>
            <span class="comment-time">
              <time
                :datetime="
                  comment.createTime | formatDate('yyyy-MM-ddTHH:mm:ss')
                "
                itemprop="datePublished"
                >{{ comment.createTime | prettyDate }}</time
              >
            </span>
            <span class="comment-reply">
              <a @click="reply(comment)">回复</a>
            </span>
          </div>
          <div v-viewer class="comment-content content">
            <blockquote v-if="comment.quote" class="comment-quote">
              <div class="comment-quote-user">
                <avatar :user="comment.quote.user" size="20" />
                <a class="quote-nickname">{{ comment.quote.user.nickname }}</a>
                <span class="quote-time">
                  {{ comment.quote.createTime | prettyDate }}
                </span>
              </div>
              <div
                v-lazy-container="{ selector: 'img' }"
                itemprop="text"
                v-html="comment.quote.content"
              />
              <div
                v-if="comment.quote.imageList && comment.quote.imageList.length"
                v-lazy-container="{ selector: 'img' }"
                class="comment-image-list"
              >
                <img
                  v-for="(image, imageIndex) in comment.quote.imageList"
                  :key="imageIndex"
                  class="small"
                  :data-src="image.url"
                />
              </div>
            </blockquote>
            <div
              v-lazy-container="{ selector: 'img' }"
              v-html="comment.content"
            ></div>
            <div
              v-if="comment.imageList && comment.imageList.length"
              v-lazy-container="{ selector: 'img' }"
              class="comment-image-list"
            >
              <img
                v-for="(image, imageIndex) in comment.imageList"
                :key="imageIndex"
                :data-src="image.url"
              />
            </div>
          </div>
        </li>
      </ul>
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
    showAd: {
      type: Boolean,
      default: false,
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
      if (!data) return

      this.$refs.commentsLoadMore.unshiftResults(data)
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
  padding: 10px;

  .comment {
    padding: 8px 0;
    overflow: hidden;

    &:not(:last-child) {
      border-bottom: 1px dashed var(--border-color2);
    }

    .comment-avatar {
      float: left;
      padding: 3px;
      margin-right: 10px;
    }

    .comment-meta {
      position: relative;
      height: 36px;

      .comment-nickname {
        position: relative;
        font-size: 14px;
        font-weight: 800;
        margin-right: 5px;
        cursor: pointer;
        // TODO
        color: #1abc9c;
        text-decoration: none;
        display: inline-block;
      }

      .comment-time {
        font-size: 12px;
        color: var(--text-color3);
        line-height: 1;
        display: inline-block;
        position: relative;
      }

      .comment-reply {
        float: right;
        font-size: 12px;
      }
    }

    .comment-content {
      word-wrap: break-word;
      word-break: break-all;
      text-align: justify;
      color: var(--text-color);
      font-size: 14px;
      line-height: 1.6;
      position: relative;
      padding-left: 45px;
      margin-top: -5px;

      img {
        max-width: 50%;
      }

      .comment-image-list {
        margin-top: 10px;
        display: flex;

        & > img {
          margin: 0 8px 8px 0;
          width: 120px;
          height: 120px;
          line-height: 120px;
          max-width: 120px;
          object-fit: cover;
          transition: all 0.5s ease-out 0.1s;

          &.small {
            width: 90px;
            height: 90px;
            line-height: 90px;
            max-width: 90px;
          }

          &:hover {
            transform: matrix(1.04, 0, 0, 1.04, 0, 0);
            backface-visibility: hidden;
          }
        }
      }
    }

    .comment-quote {
      font-size: 12px;
      padding: 10px 10px;
      border-left: 2px solid #5978f3;

      &::after {
        content: '\201D';
        font-size: 60px;
        font-weight: bold;
        color: var(--text-color3);
        position: absolute;
        right: 0;
        top: -18px;
      }

      .comment-quote-user {
        display: flex;
        margin-bottom: 6px;

        .quote-nickname {
          line-height: 20px;
          font-weight: 700;
          margin-left: 5px;
        }

        .quote-time {
          line-height: 20px;
          margin-left: 5px;
          color: var(--text-color3);
        }
      }
    }
  }
}
</style>
