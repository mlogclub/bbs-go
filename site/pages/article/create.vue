<template>
  <section class="main">
    <div class="container">
      <div class="columns">
        <div class="column is-21">
          <div class="widget">
            <div class="header">
              <nav class="breadcrumb" aria-label="breadcrumbs">
                <ul>
                  <li> <a href="/">首页</a> </li>
                  <li>
                    <a :href="'/user/' + currentUser.id + '?tab=articles'">{{ currentUser.nickname }}</a>
                  </li>
                  <li class="is-active">
                    <a href="#" aria-current="page">文章</a>
                  </li>
                </ul>
              </nav>
            </div>
            <div class="content">
              <div class="field">
                <div class="control">
                  <input
                    v-model="postForm.title"
                    class="input"
                    type="text"
                    placeholder="标题"
                  >
                </div>
              </div>

              <div class="field">
                <div class="control">
                  <tag-input v-model="postForm.tags" />
                </div>
              </div>

              <div class="field">
                <div class="control">
                  <editor v-model="postForm.content" />
                </div>
              </div>

              <div class="field is-grouped">
                <div class="control">
                  <a class="button is-success" :class="{'is-loading': publishing}" :disabled="publishing" @click="submitCreate">发表</a>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="column is-3">
          <div class="widget">
            <div class="header">
              Markdown 语法参考
            </div>
            <div class="content">
              <ol>
                <li> <tt>### 单行的标题</tt> </li>
                <li> <tt>**粗体**</tt> </li>
                <li> <tt>`console.log('行内代码')`</tt> </li>
                <li> <tt>```js\n code \n```</tt> 标记代码块 </li>
                <li> <tt>[内容](链接)</tt> </li>
                <li> <tt>![文字说明](图片链接)</tt> </li>
              </ol>
              <span>
                <a href="https://mlog.club/article/5522" target="_blank">Markdown 文档</a>
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import utils from '~/common/utils'
import Editor from '~/components/Editor'
import TagInput from '~/components/TagInput'
export default {
  middleware: 'authenticated',
  components: {
    Editor, TagInput
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
  head() {
    return {
      title: this.$siteTitle('发表话题')
    }
  },
  async asyncData({ $axios }) {
    const [currentUser] = await Promise.all([
      $axios.get('/api/user/current')
    ])
    return {
      currentUser: currentUser
    }
  },
  mounted() {
  },
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
          tags: me.postForm.tags ? me.postForm.tags.join(',') : ''
        })
        this.$toast.success('提交成功', {
          duration: 1000,
          onComplete: function () {
            utils.linkTo('/article/' + article.articleId)
          }
        })
      } catch (e) {
        me.publishing = false
        this.$toast.error('提交失败：' + (e.message || e))
      }
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
