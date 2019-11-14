<template>
  <div id="comments" ref="comments" class="comments">
    <div v-if="isLogin" class="comment-create">
      <div class="comment-input-wrapper">
        <div v-if="quote" class="comment-quote-info">
          回复：<label v-text="quote.user.nickname" />
          <i class="iconfont icon-close" @click="cancelReply" />
        </div>
        <textarea
          ref="commentEditor"
          v-model="content"
          class="comment-input"
          placeholder="发表你的观点..."
          @keydown.ctrl.enter="ctrlEnterCreate"
          @keydown.meta.enter="ctrlEnterCreate"
          @input="autoHeight"
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
        <button class="button is-light" @click="create" v-text="btnName" />
      </div>
    </div>
    <div v-else class="comment-not-login">
      <div class="comment-login-div">
        请<a style="font-weight: 700;" @click="toLogin">登录</a>后发表观点
      </div>
    </div>

    <div v-if="showAd">
      <ins
        class="adsbygoogle"
        style="display:block"
        data-ad-client="ca-pub-5683711753850351"
        data-ad-slot="1742173616"
        data-ad-format="auto"
        data-full-width-responsive="true"
      />
      <script>
        (adsbygoogle = window.adsbygoogle || []).push({});
      </script>
    </div>

    <load-more
      v-if="commentsPage"
      v-slot="{ results }"
      :init-data="commentsPage"
      :params="{ entityType:entityType, entityId: entityId }"
      url="/api/comment/list"
    >
      <div v-for="comment in results" :key="comment.commentId">
        <div class="comment">
          <div class="comment-avatar">
            <div
              class="avatar has-border is-rounded"
              :style="{backgroundImage:'url(' + comment.user.avatar + ')'}"
            />
          </div>
          <div class="comment-meta">
            <div>
              <span class="comment-nickname"><a :href="'/user/' + comment.user.id">{{ comment.user.nickname }}</a></span>
              <span class="comment-reply"><a @click="reply(comment)">回复</a></span>
            </div>
            <div>
              <small class="comment-time">{{ comment.createTime | prettyDate }}</small>
            </div>
          </div>
          <div v-highlight class="content comment-content">
            <p v-html="comment.content" />
            <blockquote
              v-if="comment.quoteContent"
              class="comment-quote"
              v-html="comment.quoteContent"
            />
          </div>
        </div>
      </div>
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
    showAd: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      commentsPage: null, // 数据
      content: '', // 内容
      currentUser: true, // 当前登录用户
      sending: false, // 发送中
      quote: null // 引用的对象
    }
  },
  computed: {
    btnName: function () {
      return this.sending ? '正在发表...' : '发表 (ctrl+enter)'
    },
    isLogin: function () {
      return this.currentUser != null
    }
  },
  created() {
    this.list()
    this.getCurrentUser()
  },
  methods: {
    async getCurrentUser() {
      try {
        const ret = await this.$axios.get('/api/user/current')
        this.currentUser = ret
      } catch (e) {
        console.error(e)
      }
    },
    async list() {
      this.commentsPage = await this.$axios.get('/api/comment/list', {
        params: {
          entityType: this.entityType,
          entityId: this.entityId
        }
      })
    },
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
        this.results.unshift(data)
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
      this.$refs.comments.scrollIntoView({ block: 'start', behavior: 'smooth' })
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
.comments {
  .comment {
    padding: .5rem 0;
    overflow: hidden;
    border-bottom: 1px dashed #f5f5f5;

    .comment-avatar {
      float: left;
      padding: .125rem;
      margin-right: .7525rem;
    }

    .comment-meta {
      position: relative;
      height: 50px;

      .comment-nickname {
        position: relative;
        font-size: .875rem;
        font-weight: 800;
        margin-right: .875rem;
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
      line-height: 2;
      position: relative;
      padding-left: 62px;
    }

    .comment-quote {
      font-size: 12px;
    }
  }

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
      font-size: .875rem;
      background: transparent;
      resize: vertical;
      -webkit-transition: all .25s ease;
      transition: all .25s ease;
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

  .comment-show-more {
    margin-top: 10px;
    text-align: center;
  }
}
</style>
