<template>
  <div class="comment-component">
    <div class="comment-form">
      <div v-if="isLogin" class="comment-create">
        <div class="comment-input-wrapper">
          <div v-if="quote" class="comment-quote-info">
            回复：
            <label v-text="quote.user.nickname" />
            <i @click="cancelReply" class="iconfont icon-close" />
          </div>
          <textarea
            ref="commentEditor"
            v-model="content"
            @keydown.ctrl.enter="ctrlEnterCreate"
            @keydown.meta.enter="ctrlEnterCreate"
            @input="autoHeight"
            class="comment-input"
            placeholder="发表你的观点..."
          />
        </div>
        <div class="comment-button-wrapper">
          <div class="comment-help" title="Markdown is supported">
            <a href="https://mlog.club/article/5522" target="_blank">
              <svg
                class="markdown"
                viewBox="0 0 16 16"
                version="1.1"
                width="16"
                height="16"
                aria-hidden="true"
              >
                <path
                  fill-rule="evenodd"
                  d="M14.85 3H1.15C.52 3 0 3.52 0 4.15v7.69C0 12.48.52 13 1.15 13h13.69c.64 0 1.15-.52 1.15-1.15v-7.7C16 3.52 15.48 3 14.85 3zM9 11H7V8L5.5 9.92 4 8v3H2V5h2l1.5 2L7 5h2v6zm2.99.5L9.5 8H11V5h2v3h1.5l-2.51 3.5z"
                />
              </svg>
            </a>
          </div>
          <button @click="create" v-text="btnName" class="button is-light" />
        </div>
      </div>
      <div v-else class="comment-not-login">
        <div class="comment-login-div">
          请
          <a @click="toLogin" style="font-weight: 700;">登录</a>后发表观点
        </div>
      </div>
    </div>

    <div v-if="showAd">
      <!-- 展示广告 -->
      <adsbygoogle ad-slot="1742173616" />
    </div>

    <load-more
      ref="commentsLoadMore"
      v-if="commentsPage"
      v-slot="{ results }"
      :init-data="commentsPage"
      :params="{ entityType: entityType, entityId: entityId }"
      url="/api/comment/list"
    >
      <ul class="comments">
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
            <div
              :style="{ backgroundImage: 'url(' + comment.user.avatar + ')' }"
              class="avatar has-border is-rounded"
            />
          </div>
          <div class="comment-meta">
            <span
              class="comment-nickname"
              itemprop="creator"
              itemscope
              itemtype="http://schema.org/Person"
            >
              <a :href="'/user/' + comment.user.id" itemprop="name">
                {{ comment.user.nickname }}
              </a>
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
          <div class="comment-content content">
            <blockquote v-if="comment.quote" class="comment-quote">
              <div class="comment-quote-user">
                <div
                  :style="{
                    backgroundImage: 'url(' + comment.quote.user.avatar + ')'
                  }"
                  class="quote-avatar avatar has-border is-rounded"
                />
                <a class="quote-nickname">{{ comment.quote.user.nickname }}</a>
                <span class="quote-time">
                  {{ comment.quote.createTime | prettyDate }}
                </span>
              </div>
              <div v-html="comment.quote.content" itemprop="text" />
            </blockquote>
            <p v-html="comment.content" />
          </div>
        </li>
      </ul>
    </load-more>
  </div>
</template>

<script>
import utils from '~/common/utils'
import LoadMore from '~/components/LoadMore'
export default {
  name: 'Comment',
  components: {
    LoadMore
  },
  props: {
    entityType: {
      type: String,
      default: '',
      required: true
    },
    entityId: {
      type: Number,
      default: 0,
      required: true
    },
    commentsPage: {
      type: Object,
      default() {
        return {}
      }
    },
    showAd: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      content: '', // 内容
      sending: false, // 发送中
      quote: null // 引用的对象
    }
  },
  computed: {
    btnName() {
      return this.sending ? '正在发表...' : '发表 (Ctrl/Cmd + Enter)'
    },
    user() {
      return this.$store.state.user.current
    },
    isLogin() {
      return this.$store.state.user.current != null
    }
  },
  methods: {
    ctrlEnterCreate(event) {
      event.stopPropagation()
      event.preventDefault()
      this.create()
    },
    async create() {
      if (!this.content) {
        this.$toast.error('请输入评论内容')
        return
      }
      if (this.sending) {
        console.log('正在发送中，请不要重复提交...')
        return
      }
      this.sending = true
      try {
        const data = await this.$axios.post('/api/comment/create', {
          entityType: this.entityType,
          entityId: this.entityId,
          content: this.content,
          quoteId: this.quote ? this.quote.commentId : ''
        })
        this.$refs.commentsLoadMore.unshiftResults(data)
        this.content = ''
        this.quote = null
      } catch (e) {
        console.error(e)
        this.$toast.error('评论失败：' + (e.message || e))
      } finally {
        this.sending = false
      }
    },
    reply(quote) {
      if (!this.isLogin) {
        utils.toSignin()
      }
      this.quote = quote
      this.$refs.commentEditor.scrollIntoView({
        block: 'start',
        behavior: 'smooth'
      })
    },
    cancelReply() {
      this.quote = null
    },
    autoHeight() {
      const elem = this.$refs.commentEditor
      elem.style.height = 'auto'
      elem.scrollTop = 0 // 防抖动
      elem.style.height = elem.scrollHeight + 'px'
    },
    toLogin() {
      utils.toSignin()
    }
  }
}
</script>

<style lang="scss" scoped>
.comment-component {
  .comment-form {
    .comment-create {
      border: 1px solid #f0f0f0;
      border-radius: 4px;
      margin-bottom: 10px;
      overflow: hidden;
      position: relative;
      padding: 10px;
      box-sizing: border-box;

      .comment-quote-info {
        font-size: 12px;
        color: #000;

        i {
          font-size: 12px !important;
          color: blue;
          cursor: pointer;
        }

        i:hover {
          color: red;
        }
      }

      .comment-input {
        width: 100%;
        min-height: 8.75rem;
        font-size: 0.875rem;
        background: transparent;
        resize: vertical;
        -webkit-transition: all 0.25s ease;
        transition: all 0.25s ease;
        border: none;
        outline: none;
        padding: 10px 5px;
        max-width: 100%;
        margin-top: 0;
        margin-bottom: 0;
        overflow: hidden;
      }

      .comment-button-wrapper {
        .comment-help {
          float: left;
          margin-top: 5px;
        }

        button {
          float: right;
        }
      }
    }

    .comment-not-login {
      border: 1px solid #f0f0f0;
      border-radius: 0px;
      margin-bottom: 10px;
      overflow: hidden;
      position: relative;
      padding: 10px;
      box-sizing: border-box;

      .comment-login-div {
        color: #d5d5d5;
        cursor: pointer;
        border-radius: 3px;
        padding: 0 10px;

        a {
          margin-left: 10px;
          margin-right: 10px;
        }
      }
    }
  }

  .comments {
    .comment {
      padding: 8px 0;
      overflow: hidden;

      &:not(:last-child) {
        border-bottom: 1px dashed #d1d1d1;
      }

      .comment-avatar {
        float: left;
        padding: 3px;
        margin-right: 5px;

        .avatar {
          min-width: 36px;
          min-height: 36px;
          width: 36px;
          height: 36px;
        }
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
          color: #1abc9c;
          text-decoration: none;
          display: inline-block;
        }

        .comment-time {
          font-size: 12px;
          color: #999999;
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
        color: #4a4a4a;
        font-size: 14px;
        line-height: 1.6;
        position: relative;
        padding-left: 45px;
        margin-top: -5px;
      }

      .comment-quote {
        font-size: 12px;
        padding: 10px 10px;
        border-left: 2px solid #5978f3;

        &::after {
          content: '\201D';
          font-family: Georgia, serif;
          font-size: 60px;
          font-weight: bold;
          color: #aaa;
          position: absolute;
          right: 0px;
          top: -18px;
        }

        .comment-quote-user {
          display: flex;
          .quote-avatar {
            min-width: 20px;
            min-height: 20px;
            width: 20px;
            height: 20px;
          }

          .quote-nickname {
            line-height: 20px;
            font-weight: 700;
            margin-left: 5px;
          }

          .quote-time {
            line-height: 20px;
            margin-left: 5px;
            color: rgba(134, 142, 150, 0.8);
          }
        }
      }
    }
  }
}
</style>
