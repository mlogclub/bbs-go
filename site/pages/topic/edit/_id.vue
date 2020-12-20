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
                  v-model="postForm.content"
                  editor-id="topicEditEditor"
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
  async asyncData({ $axios, params }) {
    const [topic, nodes] = await Promise.all([
      $axios.get('/api/topic/edit/' + params.id),
      $axios.get('/api/topic/nodes'),
    ])
    return {
      topic,
      nodes,
      postForm: {
        nodeId: topic.nodeId,
        title: topic.title,
        tags: topic.tags,
        content: topic.content,
      },
    }
  },
  data() {
    return {
      publishing: false, // 当前是否正处于发布中...
      postForm: {
        nodeId: 0,
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
  mounted() {},
  methods: {
    async submitCreate() {
      const me = this
      if (me.publishing) {
        return
      }
      me.publishing = true

      try {
        const topic = await this.$axios.post(
          '/api/topic/edit/' + this.topic.topicId,
          {
            nodeId: this.postForm.nodeId,
            title: this.postForm.title,
            content: this.postForm.content,
            tags: this.postForm.tags ? this.postForm.tags.join(',') : '',
          }
        )
        this.$msg({
          message: '修改成功',
          onClose() {
            me.$linkTo('/topic/' + topic.topicId)
          },
        })
      } catch (e) {
        console.error(e)
        me.publishing = false
        this.$message.error('提交失败：' + (e.message || e))
      }
    },
  },
  head() {
    return {
      title: this.$siteTitle('修改话题'),
    }
  },
}
</script>

<style lang="scss" scoped></style>
