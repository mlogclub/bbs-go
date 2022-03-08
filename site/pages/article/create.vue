<template>
  <section class="main">
    <div class="container">
      <article v-if="isNeedEmailVerify" class="message is-warning">
        <div class="message-header">
          <p>请先验证邮箱</p>
        </div>
        <div class="message-body">
          发表话题前，请先前往
          <strong
            ><nuxt-link
              to="/user/profile/account"
              style="color: var(--text-link-color)"
              >个人中心 &gt; 账号设置</nuxt-link
            ></strong
          >
          页面设置邮箱，并完成邮箱认证。
        </div>
      </article>
      <div v-else class="article-create-form">
        <h1 class="title">发文章</h1>
        <div class="field">
          <div class="control">
            <input
              v-model="postForm.title"
              class="input"
              type="text"
              placeholder="标题"
            />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <markdown-editor
              ref="mdEditor"
              v-model="postForm.content"
              placeholder="请输入内容，将图片复制或拖入编辑器可上传"
            />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <tag-input v-model="postForm.tags" />
          </div>
        </div>

        <div class="field is-grouped">
          <div class="control">
            <a
              :class="{ 'is-loading': publishing }"
              :disabled="publishing"
              class="button is-success"
              @click="submitCreate"
              >发表</a
            >
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  middleware: 'authenticated',
  data() {
    return {
      publishing: false, // 当前是否正处于发布中...
      postForm: {
        title: '',
        tags: [],
        content: '',
      },
    }
  },
  head() {
    return {
      title: this.$siteTitle('发表文章'),
    }
  },
  computed: {
    user() {
      return this.$store.state.user.current
    },
    config() {
      return this.$store.state.config.config
    },
    // 是否需要先邮箱认证
    isNeedEmailVerify() {
      // 发帖必须认证
      return this.config.createArticleEmailVerified && !this.user.emailVerified
    },
  },
  mounted() {},
  methods: {
    async submitCreate() {
      const me = this
      if (me.publishing) {
        return
      }
      me.publishing = true
      try {
        const article = await this.$axios.post('/api/article/create', {
          title: me.postForm.title,
          content: me.postForm.content,
          tags: me.postForm.tags ? me.postForm.tags.join(',') : '',
        })
        this.$refs.mdEditor.clearCache()
        this.$msg({
          message: '提交成功',
          onClose() {
            me.$linkTo('/article/' + article.articleId)
          },
        })
      } catch (e) {
        me.publishing = false
        this.$message.error(e.message || e)
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.article-create-form {
  background-color: var(--bg-color);
  padding: 30px;
}
</style>
