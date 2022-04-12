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
      <div v-else class="topic-create-form">
        <h1 class="title">{{ postForm.type === 0 ? '发帖子' : '发动态' }}</h1>

        <div class="field">
          <div class="control">
            <div
              v-for="node in nodes"
              :key="node.nodeId"
              class="tag"
              :class="{ selected: postForm.nodeId === node.nodeId }"
              @click="postForm.nodeId = node.nodeId"
            >
              <span>{{ node.name }}</span>
            </div>
          </div>
        </div>

        <div v-if="postForm.type === 0" class="field">
          <div class="control">
            <input
              v-model="postForm.title"
              class="input topic-title"
              type="text"
              placeholder="请输入帖子标题"
            />
          </div>
        </div>

        <div v-if="postForm.type === 0" class="field">
          <div class="control">
            <markdown-editor
              ref="mdEditor"
              v-model="postForm.content"
              placeholder="请输入你要发表的内容..."
            />
          </div>
        </div>

        <div v-if="postForm.type === 0 && isEnableHideContent" class="field">
          <div class="control">
            <markdown-editor
              ref="mdEditor"
              v-model="postForm.hideContent"
              height="200px"
              placeholder="隐藏内容，评论后可见"
            />
          </div>
        </div>

        <div v-if="postForm.type === 1" class="field">
          <div class="control">
            <simple-editor
              ref="simpleEditor"
              @input="onSimpleEditorInput"
              @submit="submitCreate"
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
              style="max-width: 150px; margin-right: 20px"
            />
            <span class="icon is-small is-left">
              <i class="iconfont icon-captcha" />
            </span>
          </div>
          <div class="field">
            <a @click="showCaptcha">
              <img :src="captchaUrl" style="height: 40px" />
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
              >{{ postForm.type === 1 ? '发表动态' : '发表帖子' }}</a
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

    const type = parseInt(query.type || 0) || 0

    return {
      nodes,
      postForm: {
        type,
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
        type: 0,
        nodeId: 0,
        title: '',
        tags: [],
        content: '',
        hideContent: '',
        imageList: [],
      },
    }
  },
  head() {
    return {
      title: this.$siteTitle(this.postForm.type === 1 ? '发动态' : '发帖子'),
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
      return this.config.createTopicEmailVerified && !this.user.emailVerified
    },
    isEnableHideContent() {
      return this.config.enableHideContent
    },
  },
  watchQuery: ['type', 'nodeId'],
  mounted() {
    this.showCaptcha()
  },
  methods: {
    async submitCreate() {
      if (this.publishing) {
        return
      }

      if (!this.postForm.nodeId) {
        this.$message.error('请选择节点')
        return
      }

      this.publishing = true

      if (this.$refs.simpleEditor && this.$refs.simpleEditor.isOnUpload()) {
        this.$message.warning('正在上传中...请上传完成后提交')
        return
      }

      const me = this
      try {
        const topic = await this.$axios.post('/api/topic/create', {
          captchaId: this.captchaId,
          captchaCode: this.captchaCode,
          type: this.postForm.type,
          nodeId: this.postForm.nodeId,
          title: this.postForm.title,
          content: this.postForm.content,
          hideContent: this.postForm.hideContent,
          imageList:
            this.postForm.imageList && this.postForm.imageList.length
              ? JSON.stringify(this.postForm.imageList)
              : '',
          tags: this.postForm.tags ? this.postForm.tags.join(',') : '',
        })
        if (this.$refs.mdEditor) {
          this.$refs.mdEditor.clearCache()
        }
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
    onSimpleEditorInput(value) {
      this.postForm.content = value.content
      this.postForm.imageList = value.imageList
    },
  },
}
</script>

<style lang="scss" scoped></style>
