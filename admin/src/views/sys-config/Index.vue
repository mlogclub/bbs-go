<template>
  <section class="page-container">
    <el-tabs value="first">
      <el-tab-pane label="通用配置" name="first">
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
                style="width:100%"
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
                style="width:100%"
                multiple
                filterable
                allow-create
                default-first-option
                placeholder="推荐标签"
              ></el-select>
            </el-form-item>

            <el-form-item label="站外链接跳转页面">
              <el-tooltip
                content="在跳转前需手动确认是否前往该站外链接"
                placement="top"
              >
                <el-switch v-model="config.urlRedirect"></el-switch>
              </el-tooltip>
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>
      <el-tab-pane label="导航配置" name="second" class="nav-panel">
        <draggable
          v-model="config.siteNavs"
          draggable=".nav"
          handle=".nav-sort-btn"
          class="navs"
        >
          <div v-for="(nav, index) in config.siteNavs" :key="index" class="nav">
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
    </el-tabs>

    <div style="margin-top: 20px;">
      <el-button type="primary" :loading="loading" @click="save"
        >保存</el-button
      >
    </div>
  </section>
</template>

<script>
import draggable from 'vuedraggable'
import HttpClient from '@/apis/HttpClient'

export default {
  name: 'List',
  components: {
    draggable
  },
  data() {
    return {
      config: {},
      loading: false,
      autocompleteTags: [],
      autocompleteTagLoading: false
    }
  },
  mounted() {
    this.load()
  },
  methods: {
    async load() {
      try {
        this.config = await HttpClient.get('/api/admin/sys-config/all')
      } catch (err) {
        this.$notify.error({ title: '错误', message: err.message })
      }
    },
    async save() {
      this.loading = true
      try {
        await HttpClient.post('/api/admin/sys-config/save', {
          config: JSON.stringify(this.config)
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
        url: ''
      })
    },
    delNav(index) {
      if (!this.config.siteNavs) {
        return
      }
      this.config.siteNavs.splice(index, 1)
    }
  }
}
</script>

<style scoped lang="scss">
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
