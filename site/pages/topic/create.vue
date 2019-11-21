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
                      <a href="#" aria-current="page">主题</a>
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
                      placeholder="请输入标题，如果标题能够表达完整内容，则正文可以为空"
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
                      >发表主题</a
                    >
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="column is-3 ">
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
  async asyncData({ $axios, query }) {
    const currentUser = await $axios.get('/api/user/current')

    // 发帖标签
    const tags = []
    if (query.tagId) {
      try {
        const tag = await $axios.get('/api/tag/' + query.tagId)
        if (tag) {
          tags.push(tag.tagName)
        }
      } catch (e) {
        console.error(e)
      }
    }

    return {
      currentUser,
      postForm: {
        tags
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
  mounted() {},
  methods: {
    async submitCreate() {
      const me = this
      if (me.publishing) {
        return
      }
      me.publishing = true

      try {
        const me = this
        const topic = await this.$axios.post('/api/topic/create', {
          title: me.postForm.title,
          content: me.postForm.content,
          tags: me.postForm.tags ? me.postForm.tags.join(',') : ''
        })
        this.$toast.success('提交成功', {
          duration: 1000,
          onComplete() {
            utils.linkTo('/topic/' + topic.topicId)
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
