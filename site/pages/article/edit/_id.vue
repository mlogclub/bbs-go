<template>
  <section class="main">
    <div class="container-wrapper">
      <div class="columns">
        <div class="column is-21">
          <div class="main-body">
            <div class="widget">
              <div class="widget-header">
                <nav class="breadcrumb" aria-label="breadcrumbs">
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
                    <tag-input v-model="postForm.tags" />
                  </div>
                </div>

                <div class="field">
                  <div class="control">
                    <vditor v-model="postForm.content" />
                  </div>
                </div>

                <div class="field is-grouped">
                  <div class="control">
                    <a
                      :class="{ 'is-loading': publishing }"
                      :disabled="publishing"
                      @click="submitCreate"
                      class="button is-success"
                      >提交更改</a
                    >
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="column is-3">
          <div class="main-aside">
            <markdown-help />
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import utils from '~/common/utils'
import TagInput from '~/components/TagInput'
import MarkdownHelp from '~/components/MarkdownHelp'

export default {
  middleware: 'authenticated',
  components: {
    TagInput,
    MarkdownHelp
  },
  async asyncData({ $axios, params }) {
    const [article] = await Promise.all([
      $axios.get('/api/article/edit/' + params.id)
    ])
    return {
      article,
      postForm: {
        title: article.title,
        tags: article.tags,
        content: article.content
      }
    }
  },
  data() {
    return {
      publishing: false, // 当前是否正处于发布中...
      postForm: {
        title: '',
        tags: [],
        content: ''
      }
    }
  },
  computed: {
    currentUser() {
      return this.$store.state.user.current
    }
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
            tags: this.postForm.tags ? this.postForm.tags.join(',') : ''
          }
        )
        this.$toast.success('修改成功', {
          duration: 2000,
          onComplete() {
            utils.linkTo('/article/' + article.articleId)
          }
        })
      } catch (e) {
        console.error(e)
        me.publishing = false
        this.$toast.error('提交失败：' + (e.message || e))
      }
    }
  },
  head() {
    return {
      title: this.$siteTitle('发表话题')
    }
  }
}
</script>

<style lang="scss" scoped></style>
