<template>
  <div class="comment-form">
    <div class="comment-create">
      <div ref="commentEditor" class="comment-input-wrapper">
        <div v-if="quote" class="comment-quote-info">
          回复：
          <label v-text="quote.user.nickname" />
          <i class="iconfont icon-close" alt="取消回复" @click="cancelReply" />
        </div>
        <markdown-editor
          v-if="inputMode === 'markdown'"
          ref="mdEditor"
          v-model="value.content"
          height="200px"
          placeholder="请发表你的观点..."
          @submit="create"
        />
        <simple-editor
          v-else
          ref="simpleEditor"
          v-model="value"
          :max-word-count="500"
          height="150px"
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
  </div>
</template>

<script>
export default {
  props: {
    mode: {
      type: String,
      default: 'markdown',
    },
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
      value: {
        content: '', // 内容
        imageList: [],
      },
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
    inputMode() {
      if (this.$store.state.env.isMobile) {
        // 手机中，强制使用普通文本编辑器
        return 'text'
      }
      return this.mode
    },
    contentType() {
      return this.inputMode === 'markdown' ? 'markdown' : 'text'
    },
  },
  methods: {
    async create() {
      if (!this.value.content) {
        this.$message.error('请输入评论内容')
        return
      }
      if (this.sending) {
        console.log('正在发送中，请不要重复提交...')
        return
      }
      if (this.$refs.simpleEditor && this.$refs.simpleEditor.isOnUpload()) {
        this.$message.warning('正在上传中...请上传完成后提交')
        return
      }
      this.sending = true
      try {
        const data = await this.$axios.post('/api/comment/create', {
          contentType: this.contentType,
          entityType: this.entityType,
          entityId: this.entityId,
          content: this.value.content,
          imageList:
            this.value.imageList && this.value.imageList.length
              ? JSON.stringify(this.value.imageList)
              : '',
          quoteId: this.quote ? this.quote.commentId : '',
        })
        this.$emit('created', data)

        this.value.content = ''
        this.value.imageList = []
        this.quote = null
        if (this.$refs.mdEditor) {
          this.$refs.mdEditor.clear()
        }
        if (this.$refs.simpleEditor) {
          this.$refs.simpleEditor.clear()
        }
      } catch (e) {
        console.error(e)
        this.$message.error(e.message || e)
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
  },
}
</script>

<style scoped lang="scss"></style>
