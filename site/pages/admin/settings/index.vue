<template>
  <section v-loading="loading" class="page-container">
    <div class="page-section config-panel">
      <el-tabs value="commonConfigTab">
        <el-tab-pane label="通用配置" name="commonConfigTab">
          <div class="config">
            <el-form label-width="160px">
              <el-form-item label="网站名称">
                <el-input
                  v-model="config.siteTitle"
                  type="text"
                  placeholder="网站名称"
                ></el-input>
              </el-form-item>

              <el-form-item label="网站描述">
                <el-input
                  v-model="config.siteDescription"
                  type="textarea"
                  autosize
                  placeholder="网站描述"
                ></el-input>
              </el-form-item>

              <el-form-item label="网站关键字">
                <el-select
                  v-model="config.siteKeywords"
                  style="width: 100%;"
                  multiple
                  filterable
                  allow-create
                  default-first-option
                  placeholder="网站关键字"
                ></el-select>
              </el-form-item>

              <el-form-item label="网站公告">
                <el-input
                  v-model="config.siteNotification"
                  type="textarea"
                  autosize
                  placeholder="网站公告（支持输入HTML）"
                ></el-input>
              </el-form-item>

              <el-form-item label="推荐标签">
                <el-select
                  v-model="config.recommendTags"
                  style="width: 100%;"
                  multiple
                  filterable
                  allow-create
                  default-first-option
                  placeholder="推荐标签"
                ></el-select>
              </el-form-item>

              <el-form-item label="默认节点">
                <el-select
                  v-model="config.defaultNodeId"
                  style="width: 100%;"
                  placeholder="发帖默认节点"
                >
                  <el-option
                    v-for="node in nodes"
                    :key="node.id"
                    :label="node.name"
                    :value="node.id"
                  >
                  </el-option>
                </el-select>
              </el-form-item>

              <template v-if="config.loginMethod">
                <el-form-item label="登录方式">
                  <el-checkbox v-model="config.loginMethod.password"
                    >密码登录</el-checkbox
                  >
                  <el-checkbox v-model="config.loginMethod.qq"
                    >QQ登录</el-checkbox
                  >
                  <el-checkbox v-model="config.loginMethod.github"
                    >Github登录</el-checkbox
                  >
                </el-form-item>
              </template>

              <el-form-item label="站外链接跳转页面">
                <el-tooltip
                  content="在跳转前需手动确认是否前往该站外链接"
                  placement="top"
                >
                  <el-switch v-model="config.urlRedirect"></el-switch>
                </el-tooltip>
              </el-form-item>

              <el-form-item label="发帖验证码">
                <el-tooltip content="发帖时是否开启验证码校验" placement="top">
                  <el-switch v-model="config.topicCaptcha"></el-switch>
                </el-tooltip>
              </el-form-item>

              <el-form-item label="发表文章审核">
                <el-tooltip content="发布文章后是否开启审核" placement="top">
                  <el-switch v-model="config.articlePending"></el-switch>
                </el-tooltip>
              </el-form-item>

              <el-form-item label="用户观察期(秒)">
                <el-tooltip
                  content="观察期内用户无法发表话题、动态等内容，设置为 0 表示无观察期。"
                  placement="top"
                >
                  <el-input-number
                    v-model="config.userObserveSeconds"
                    :min="0"
                    :max="720"
                  ></el-input-number>
                </el-tooltip>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>
        <el-tab-pane label="导航配置" name="navConfigTab" class="nav-panel">
          <draggable
            v-model="config.siteNavs"
            draggable=".nav"
            handle=".nav-sort-btn"
            class="navs"
          >
            <div
              v-for="(nav, index) in config.siteNavs"
              :key="index"
              class="nav"
            >
              <el-row :gutter="20">
                <el-col :span="1">
                  <i class="iconfont icon-sort nav-sort-btn" />
                </el-col>
                <el-col :span="10">
                  <el-input
                    v-model="nav.title"
                    type="text"
                    size="small"
                    placeholder="标题"
                  ></el-input>
                </el-col>
                <el-col :span="11">
                  <el-input
                    v-model="nav.url"
                    type="text"
                    size="small"
                    placeholder="链接"
                  ></el-input>
                </el-col>
                <el-col :span="2">
                  <el-button
                    type="danger"
                    icon="el-icon-delete"
                    circle
                    size="small"
                    @click="delNav(index)"
                  ></el-button>
                </el-col>
              </el-row>
            </div>
          </draggable>
          <div class="add-nav">
            <el-tooltip
              class="item"
              effect="dark"
              content="点击按钮添加导航"
              placement="top"
            >
              <el-button
                type="primary"
                icon="el-icon-plus"
                circle
                @click="addNav"
              ></el-button>
            </el-tooltip>
          </div>
        </el-tab-pane>
        <el-tab-pane
          v-if="config.scoreConfig"
          label="积分配置"
          name="scoreConfigTab"
        >
          <el-form label-width="160px">
            <el-form-item label="发帖积分">
              <el-input-number
                v-model="config.scoreConfig.postTopicScore"
                :min="1"
                type="text"
                placeholder="发帖获得积分"
              ></el-input-number>
            </el-form-item>
            <el-form-item label="跟帖积分">
              <el-input-number
                v-model="config.scoreConfig.postCommentScore"
                :min="1"
                type="text"
                placeholder="跟帖获得积分"
              ></el-input-number>
            </el-form-item>
            <el-form-item label="签到积分">
              <el-input-number
                v-model="config.scoreConfig.checkInScore"
                :min="1"
                type="text"
                placeholder="签到获得积分"
              ></el-input-number>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>

      <div style="text-align: right;">
        <el-button :loading="loading" type="primary" @click="save"
          >保存配置
        </el-button>
      </div>
    </div>
  </section>
</template>

<script>
import draggable from 'vuedraggable'

export default {
  layout: 'admin',
  components: {
    draggable,
  },
  data() {
    return {
      config: {},
      loading: false,
      autocompleteTags: [],
      autocompleteTagLoading: false,
      nodes: [],
    }
  },
  mounted() {
    this.load()
  },
  methods: {
    async load() {
      this.loading = true
      try {
        this.config = await this.$axios.get('/api/admin/sys-config/all')
        this.nodes = await this.$axios.get('/api/admin/topic-node/nodes')
      } catch (err) {
        this.$notify.error({ title: '错误', message: err.message })
      } finally {
        this.loading = false
      }
    },
    async save() {
      this.loading = true
      try {
        await this.$axios.post('/api/admin/sys-config/save', {
          config: JSON.stringify(this.config),
        })
        this.$message({ message: '提交成功', type: 'success' })
        this.load()
      } catch (err) {
        this.$notify.error({ title: '错误', message: err.message })
      } finally {
        this.loading = false
      }
    },
    addNav() {
      if (!this.config.siteNavs) {
        this.config.siteNavs = []
      }
      this.config.siteNavs.push({
        title: '',
        url: '',
      })
    },
    delNav(index) {
      if (!this.config.siteNavs) {
        return
      }
      this.config.siteNavs.splice(index, 1)
    },
  },
}
</script>

<style scoped lang="scss">
.config-panel {
  margin: 20px;
  padding: 10px;
}

.config {
  padding: 10px 0;
}

.nav-panel {
  .navs {
    border: 1px solid #ddd;
    border-radius: 5px;

    .nav {
      padding: 5px 5px;
      margin: 0;

      &:not(:last-child) {
        border-bottom: 1px solid #ddd;
      }

      .nav-sort-btn {
        font-size: 21px;
        font-weight: 700;
        cursor: pointer;
        float: right;
      }
    }
  }

  .add-nav {
    margin-top: 20px;
    text-align: center;
  }
}
</style>
