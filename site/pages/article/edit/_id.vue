<template>
  <section class="main">
    <div class="container main-container is-white left-main">
      <div class="left-container">
        <div class="widget">
          <div class="widget-header">
            <nav class="breadcrumb">
              <ul>
                <li><a href="/">首页</a></li>
                <li>
                  <a :href="'/user/' + currentUser.id + '?tab=topics'">{{
                    currentUser.nickname
                  }}</a>
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
                  v-model="postForm.content"
                  editor-id="articleEditEditor"
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
                  >提交更改</a
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
import TagInput from '~/components/TagInput'
import MarkdownHelp from '~/components/MarkdownHelp'
import MarkdownEditor from '~/components/MarkdownEditor'

export default {
  middleware: 'authenticated',
  components: {
    TagInput,
    MarkdownHelp,
    MarkdownEditor,
  },
  async asyncData({ $axios, params, error }) {
    try {
      const [article] = await Promise.all([
        $axios.get('/api/article/edit/' + params.id),
      ])
      return {
        article,
        postForm: {
          title: article.title,
          tags: article.tags,
          content: article.content,
        },
      }
    } catch (e) {
      error({
        statusCode: 403,
        message: e.message || '403',
      })
    }
  },
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
  computed: {
    currentUser() {
      return this.$store.state.user.current
    },
  },
  methods: {
    async submitCreate() {
      const me = this
      if (me.publishing) {
        return
      }
      me.publishing = true

      try {
        const article = await this.$axios.post(
          '/api/article/edit/' + this.article.articleId,
          {
            title: this.postForm.title,
            content: this.postForm.content,
            tags: this.postForm.tags ? this.postForm.tags.join(',') : '',
          }
        )
        this.$msg({
          message: '删除成功',
          onClose() {
            me.$linkTo('/article/' + article.articleId)
          },
        })
      } catch (e) {
        me.publishing = false
        this.$message.error('提交失败：' + (e.message || e))
      }
    },
  },
  head() {
    return {
      title: this.$siteTitle('修改文章'),
    }
  },
}
</script>

<style lang="scss" scoped></style>
