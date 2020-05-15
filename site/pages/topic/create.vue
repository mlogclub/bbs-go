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
                <el-form
                  ref="pollForm"
                  :model="pollForm"
                  :gutter="24"
                  label-width="80px"
                  label-position="top"
                >
                  <el-form-item
                    v-show="questionCreated"
                    :label="'投票问题:'"
                    :rules="{
                      required: true,
                      message: '问题不能为空',
                      trigger: 'blur'
                    }"
                    prop="question"
                  >
                    <el-col :span="14">
                      <el-input v-model="pollForm.question"></el-input>
                    </el-col>
                  </el-form-item>
                  <el-form-item
                    v-for="(option, index) in pollForm.options"
                    :key="option.key"
                    :label="'选项' + (index + 1)"
                    :prop="'options.' + index + '.value'"
                    :rules="{
                      required: true,
                      message: '选项不能为空',
                      trigger: 'blur'
                    }"
                  >
                    <el-col :span="14">
                      <el-input v-model="option.value"></el-input>
                    </el-col>
                    <el-col :span="8">
                      <el-button
                        type="danger"
                        plain
                        @click.prevent="removeOption(option)"
                        >删除
                      </el-button>
                    </el-col>
                  </el-form-item>
                  <el-form-item>
                    <el-col :span="24">
                      <el-button @click="addPollOption">新增投票项</el-button>
                      <el-button
                        v-if="pollForm.options.length > 0"
                        type="danger"
                        plain
                        @click="resetForm('pollForm')"
                        >重置投票
                      </el-button>
                      <el-button
                        v-if="
                          pollForm.options.length > 0 || questionCreated == true
                        "
                        type="danger"
                        @click="removePoll"
                        >放弃投票
                      </el-button>
                    </el-col>
                  </el-form-item>
                </el-form>
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
import '~/plugins/element-ui'
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
      },
      pollForm: {
        question: '',
        options: []
      },
      questionCreated: false
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

      if (!this.validatePollForm()) {
        return
      }

      me.publishing = true

      try {
        const me = this
        const pollsInArray = []
        me.pollForm.options.forEach(function(item) {
          pollsInArray.push(item.value)
        })

        const topic = await this.$axios.post('/api/topic/create', {
          nodeId: me.postForm.nodeId,
          title: me.postForm.title,
          content: me.postForm.content,
          tags: me.postForm.tags ? me.postForm.tags.join(',') : '',
          polls: pollsInArray ? pollsInArray.join(',') : '',
          question: me.pollForm.question
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
    },
    validatePollForm() {
      if (this.pollForm.options.length < 1) {
        return false
      }
      let pollFormValidateResult = true
      this.$refs.pollForm.validate((valid) => {
        if (!valid) {
          this.$toast.error('请检查投票内容')
          pollFormValidateResult = false
        }
      })
      return pollFormValidateResult
    },
    resetForm(formName) {
      this.$refs[formName].resetFields()
    },
    removeOption(item) {
      const index = this.pollForm.options.indexOf(item)
      if (index !== -1) {
        this.pollForm.options.splice(index, 1)
      }
    },
    addPollOption() {
      if (!this.questionCreated) {
        this.questionCreated = true
      }
      this.pollForm.options.push({
        value: '',
        key: this.pollForm.options.length
      })
    },
    removePoll() {
      this.pollForm = {
        question: '',
        options: []
      }
      this.questionCreated = false
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
