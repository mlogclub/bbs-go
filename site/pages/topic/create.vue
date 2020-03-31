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
                <div class="field" style="width:100%;">
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
                        >{{ node.name }}</option
                      >
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
    MarkdownEditor
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
        nodeId: currentNode ? currentNode.nodeId : 0
      }
    }
  },
  data() {
    return {
      publishing: false, // 当前是否正处于发布中...
      postForm: {
        nodeId: 0,
        title: '',
        tags: [],
        content: ''
      }
    }
  },
  computed: {
    user() {
      return this.$store.state.user.current
    }
  },
  mounted() {},
  methods: {
    async submitCreate() {
      const me = this
      if (me.publishing) {
        return
      }

      if (!me.postForm.title) {
        this.$toast.error('请输入标题')
        return
      }

      if (!me.postForm.nodeId) {
        this.$toast.error('请选择节点')
        return
      }

      me.publishing = true

      try {
        const me = this
        const topic = await this.$axios.post('/api/topic/create', {
          nodeId: me.postForm.nodeId,
          title: me.postForm.title,
          content: me.postForm.content,
          tags: me.postForm.tags ? me.postForm.tags.join(',') : ''
        })
        this.$refs.mdEditor.clearCache()
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
