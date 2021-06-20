<template>
  <section class="main">
    <div v-if="isNeedEmailVerify" class="container">
      <article class="message is-warning">
        <div class="message-header">
          <p>请先验证邮箱</p>
        </div>
        <div class="message-body">
          发表文章前，请先前往
          <strong
            ><nuxt-link to="/user/settings" style="color: #1878f3"
              >个人中心 &gt; 编辑资料</nuxt-link
            ></strong
          >
          页面设置邮箱，并完成邮箱认证。
        </div>
      </article>
    </div>
    <div v-else class="container main-container is-white left-main">
      <div class="left-container">
        <div class="widget">
          <div class="widget-header">
            <nav class="breadcrumb">
              <ul>
                <li><nuxt-link to="/">首页</nuxt-link></li>
                <li>
                  <nuxt-link :to="'/user/' + user.id + '?tab=articles'">{{
                    user.nickname
                  }}</nuxt-link>
                </li>
                <li class="is-active">
                  <a href="#" aria-current="page">文章</a>
                </li>
              </ul>
            </nav>
          </div>
          <div class="widget-content">
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
                  editor-id="articleCreateEditor"
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
      </div>
      <div class="right-container">
        <markdown-help />
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

<style lang="scss" scoped></style>
