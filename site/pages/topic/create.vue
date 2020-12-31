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
                  <a :href="'/user/' + user.id + '?tab=topics'">{{
                    user.nickname
                  }}</a>
                </li>
                <li class="is-active">
                  <a href="#" aria-current="page">发帖</a>
                </li>
              </ul>
            </nav>
          </div>
          <div class="widget-content topic-create-form" style="padding: 20px;">
            <div class="field">
              <div class="control">
                <span
                  v-for="node in nodes"
                  :key="node.nodeId"
                  class="tag"
                  :class="{ selected: postForm.nodeId === node.nodeId }"
                  @click="postForm.nodeId = node.nodeId"
                >
                  {{ node.name }}
                </span>
              </div>
            </div>

            <div class="field" style="width: 100%;">
              <div class="control">
                <input
                  v-model="postForm.title"
                  class="input topic-title"
                  type="text"
                  placeholder="请输入帖子标题"
                />
              </div>
            </div>

            <div class="field">
              <div class="control">
                <markdown-editor
                  ref="mdEditor"
                  v-model="postForm.content"
                  editor-id="topicCreateEditor"
                  placeholder="请输入你要发表的内容..."
                />
              </div>
            </div>

            <div class="field">
              <div class="control">
                <tag-input v-model="postForm.tags" />
              </div>
            </div>

            <div v-if="captchaUrl" class="field is-horizontal">
              <div class="field control has-icons-left">
                <input
                  v-model="captchaCode"
                  class="input"
                  type="text"
                  placeholder="验证码"
                  style="max-width: 150px; margin-right: 20px;"
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-captcha" />
                </span>
              </div>
              <div class="field">
                <a @click="showCaptcha">
                  <img :src="captchaUrl" style="height: 40px;" />
                </a>
              </div>
            </div>

            <div class="field is-grouped">
              <div class="control">
                <a
                  :class="{ 'is-loading': publishing }"
                  :disabled="publishing"
                  class="button is-success"
                  @click="submitCreate"
                  >发表话题</a
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
  async asyncData({ $axios, query, store }) {
    // 节点
    const nodes = await $axios.get('/api/topic/nodes')

    // 发帖标签
    const config = store.state.config.config || {}
    const nodeId = query.nodeId || config.defaultNodeId
    let currentNode = null
    if (nodeId) {
      try {
        currentNode = await $axios.get('/api/topic/node?nodeId=' + nodeId)
      } catch (e) {
        console.error(e)
      }
    }

    return {
      nodes,
      postForm: {
        nodeId: currentNode ? currentNode.nodeId : 0,
      },
    }
  },
  data() {
    return {
      publishing: false, // 当前是否正处于发布中...
      captchaId: '',
      captchaUrl: '',
      captchaCode: '',
      postForm: {
        nodeId: 0,
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
    config() {
      return this.$store.state.config.config
    },
  },
  mounted() {
    this.showCaptcha()
  },
  methods: {
    async submitCreate() {
      if (this.publishing) {
        return
      }

      if (!this.postForm.title) {
        this.$message.error('请输入标题')
        return
      }

      if (!this.postForm.nodeId) {
        this.$message.error('请选择节点')
        return
      }

      this.publishing = true

      const me = this
      try {
        const topic = await this.$axios.post('/api/topic/create', {
          captchaId: this.captchaId,
          captchaCode: this.captchaCode,
          nodeId: this.postForm.nodeId,
          title: this.postForm.title,
          content: this.postForm.content,
          tags: this.postForm.tags ? this.postForm.tags.join(',') : '',
        })
        this.$refs.mdEditor.clearCache()
        this.$msg({
          message: '提交成功',
          onClose() {
            me.$linkTo('/topic/' + topic.topicId)
          },
        })
      } catch (e) {
        await this.showCaptcha()
        this.publishing = false
        this.$message.error(e.message || e)
      }
    },
    async showCaptcha() {
      if (this.config.topicCaptcha) {
        try {
          const ret = await this.$axios.get('/api/captcha/request', {
            params: {
              captchaId: this.captchaId || '',
            },
          })
          this.captchaId = ret.captchaId
          this.captchaUrl = ret.captchaUrl
        } catch (e) {
          this.$message.error(e.message || e)
        }
      }
    },
  },
  head() {
    return {
      title: this.$siteTitle('发表话题'),
    }
  },
}
</script>

<style lang="scss" scoped>
.topic-create-form {
  .tag {
    margin-right: 10px;
    cursor: pointer;
    font-size: 0.85rem;
    font-weight: 700;
    color: #777;

    &.selected {
      color: #fff;
      background: #1878f3;
    }
  }

  .topic-title {
    background-color: rgb(247, 247, 247);
    border: 1px solid hsla(0, 0%, 89.4%, 0.6);
  }
}
</style>
