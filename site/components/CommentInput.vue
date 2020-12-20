<template>
  <div class="comment-form">
    <div v-if="isLogin" class="comment-create">
      <div ref="commentEditor" class="comment-input-wrapper">
        <div v-if="quote" class="comment-quote-info">
          回复：
          <label v-text="quote.user.nickname" />
          <i class="iconfont icon-close" alt="取消回复" @click="cancelReply" />
        </div>
        <markdown-editor
          ref="mdEditor"
          v-model="content"
          editor-id="createEditor"
          height="200px"
          placeholder="请发表你的观点..."
          @submit="create"
        />
      </div>
      <div class="comment-button-wrapper">
        <span>Ctrl or ⌘ + Enter</span>
        <button
          class="button is-small is-success"
          @click="create"
          v-text="btnName"
        />
      </div>
    </div>
    <div v-else class="comment-not-login">
      <div class="comment-login-div">
        请
        <a style="font-weight: 700;" @click="toLogin">登录</a>后发表观点
      </div>
    </div>
  </div>
</template>

<script>
import MarkdownEditor from '~/components/MarkdownEditor'

export default {
  components: {
    MarkdownEditor,
  },
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
  },
  data() {
    return {
      content: '', // 内容
      sending: false, // 发送中
      quote: null, // 引用的对象
    }
  },
  computed: {
    btnName() {
      return this.sending ? '正在发表...' : '发表'
    },
    user() {
      return this.$store.state.user.current
    },
    isLogin() {
      return this.$store.state.user.current != null
    },
  },
  methods: {
    async create() {
      if (!this.content) {
        this.$message.error('请输入评论内容')
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
          quoteId: this.quote ? this.quote.commentId : '',
        })
        this.$emit('created', data)
        this.content = ''
        this.$refs.mdEditor.clear()
        this.quote = null
      } catch (e) {
        console.error(e)
        this.$message.error('评论失败：' + (e.message || e))
      } finally {
        this.sending = false
      }
    },
    reply(quote) {
      this.quote = quote
      this.$refs.commentEditor.scrollIntoView({
        block: 'start',
        behavior: 'smooth',
      })
    },
    cancelReply() {
      this.quote = null
    },
    toLogin() {
      this.$toSignin()
    },
  },
}
</script>

<style scoped lang="scss"></style>
