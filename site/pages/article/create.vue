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
                  <a :href="'/user/' + user.id + '?tab=articles'">{{
                    user.nickname
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
    user() {
      return this.$store.state.user.current
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
        this.$message.error('提交失败：' + (e.message || e))
      }
    },
  },
  head() {
    return {
      title: this.$siteTitle('发表文章'),
    }
  },
}
</script>

<style lang="scss" scoped></style>
