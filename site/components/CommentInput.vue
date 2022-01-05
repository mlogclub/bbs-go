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
        <text-editor
          v-else
          ref="simpleEditor"
          v-model="value"
          @submit="create"
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
      // if (this.$store.state.env.isMobile) {
      //   // 手机中，强制使用普通文本编辑器
      //   return 'text'
      // }
      // return this.mode
      // 强制text模式
      return 'text'
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
        this.$message.success('发布成功')
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

<style scoped lang="scss">
.comment-form {
  background-color: var(--bg-color);
  padding: 10px;
  margin-bottom: 10px;

  .comment-create {
    border-radius: 4px;
    overflow: hidden;
    position: relative;
    padding: 0;
    box-sizing: border-box;

    .comment-quote-info {
      font-size: 13px;
      color: var(--text-color);
      margin-bottom: 3px;
      font-weight: 600;

      i {
        font-size: 12px !important;
        color: var(--text-link-color);
        cursor: pointer;
      }

      i:hover {
        color: red;
      }
    }

    .comment-input-wrapper {
      margin-bottom: 8px;

      .text-input {
        outline: none;
        width: 100%;
        height: 85px;
        font-size: 14px;
        padding: 10px 40px 10px 10px;
        color: var(--text-color);
        line-height: 16px;
        max-width: 100%;
        resize: none;
        border: 1px solid var(--border-color);
        box-sizing: border-box;
        border-radius: var(--jinsom-border-radius);
      }
    }

    .comment-button-wrapper {
      user-select: none;
      display: flex;
      float: right;
      height: 30px;
      line-height: 30px;

      span {
        color: var(--text-color4);
        font-size: 13px;
        margin-right: 5px;
      }

      button {
        font-weight: 500;
      }
    }
  }
}
</style>
