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
                  <a href="#" aria-current="page">主题</a>
                </li>
              </ul>
            </nav>
          </div>
          <div class="widget-content">
            <div class="field is-horizontal">
              <div class="field-body">
                <div class="field" style="width: 100%;">
                  <input
                    v-model="postForm.title"
                    class="input"
                    type="text"
                    placeholder="请输入标题"
                  />
                </div>
                <div class="field">
                  <div class="select">
                    <select v-model="postForm.nodeId">
                      <option value="0">选择节点</option>
                      <option
                        v-for="node in nodes"
                        :key="node.nodeId"
                        :value="node.nodeId"
                        >{{ node.name }}
                      </option>
                    </select>
                  </div>
                </div>
              </div>
            </div>

            <div class="field">
              <div class="control">
                <markdown-editor
                  ref="mdEditor"
                  v-model="postForm.content"
                  editor-id="topicCreateEditor"
                  placeholder="可空，将图片复制或拖入编辑器可上传"
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
import utils from '~/common/utils'
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
        this.$toast.error('请输入标题')
        return
      }

      if (!this.postForm.nodeId) {
        this.$toast.error('请选择节点')
        return
      }

      this.publishing = true

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
        this.$toast.success('提交成功', {
          duration: 1000,
          onComplete() {
            utils.linkTo('/topic/' + topic.topicId)
          },
        })
      } catch (e) {
        await this.showCaptcha()
        this.publishing = false
        this.$toast.error('提交失败：' + (e.message || e))
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
          this.$toast.error(e.message || e)
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

<style lang="scss" scoped></style>
