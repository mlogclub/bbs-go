<template>
  <section v-loading="listLoading" class="page-container">
    <!--工具条-->
    <div class="toolbar">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.id" placeholder="编号"></el-input>
        </el-form-item>
        <el-form-item>
          <el-input v-model="filters.userId" placeholder="用户编号"></el-input>
        </el-form-item>
        <el-form-item>
          <el-input v-model="filters.title" placeholder="标题"></el-input>
        </el-form-item>
        <el-form-item>
          <el-select
            v-model="filters.status"
            clearable
            placeholder="请选择状态"
            @change="list"
          >
            <el-option label="正常" value="0"></el-option>
            <el-option label="删除" value="1"></el-option>
            <el-option label="待审核" value="2"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="list">查询</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!--列表-->
    <div v-if="results && results.length > 0" class="page-section articles">
      <div v-for="item in results" :key="item.id" class="article">
        <div class="article-header">
          <avatar :user="item.user" />
          <div class="article-right">
            <div class="article-title">
              <a :href="'/article/' + item.id" target="_blank">{{
                item.title
              }}</a>
            </div>
            <div class="article-meta">
              <label class="action-item info">ID: {{ item.id }}</label>
              <label v-if="item.user" class="author">{{
                item.user.nickname
              }}</label>
              <label>{{ item.createTime | formatDate }}</label>

              <div class="article-tags">
                <el-tag
                  v-for="tag in item.tags"
                  :key="tag.tagId"
                  type="info"
                  size="mini"
                >
                  {{ tag.tagName }}
                </el-tag>
              </div>
            </div>
          </div>
        </div>
        <div class="summary">{{ item.summary }}</div>
        <div class="actions">
          <a class="action-item btn" @click="showUpdateTags(item)">修改标签</a>
          <span v-if="item.status === 1" class="action-item danger"
            >已删除</span
          >
          <a
            v-if="item.status !== 1"
            class="action-item btn"
            @click="deleteSubmit(item)"
            >删除</a
          >
          <a
            v-if="item.status === 2"
            :href="'/article/edit/' + item.id"
            class="action-item btn"
            >修改</a
          >
          <a
            v-if="item.status === 2"
            class="action-item btn"
            @click="pendingSubmit(item)"
            >审核</a
          >
        </div>
      </div>
    </div>
    <div v-else class="page-section articles">
      <div class="notification is-primary">
        <strong>无数据</strong>
      </div>
    </div>

    <!--工具条-->
    <div v-if="page.total > 0" class="pagebar">
      <el-pagination
        :page-sizes="[20, 50, 100, 300]"
        :current-page="page.page"
        :page-size="page.limit"
        :total="page.total"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handlePageChange"
        @size-change="handleLimitChange"
      ></el-pagination>
    </div>

    <el-dialog
      :visible.sync="updateTagsDialogVisible"
      :close-on-click-modal="false"
      title="添加标签"
    >
      <el-form label-width="80px">
        <el-form-item label="标签">
          <el-select
            v-model="updateTagForm.tags"
            style="width: 100%;"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="标签"
          ></el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="updateTagsDialogVisible = false"
          >取消</el-button
        >
        <el-button type="primary" @click.native="updateTags">提交 </el-button>
      </div>
    </el-dialog>
  </section>
</template>

<script>
import Avatar from '~/components/Avatar'
export default {
  layout: 'admin',
  components: { Avatar },
  data() {
    return {
      results: [],
      listLoading: false,
      page: {},
      filters: {},
      tagOptions: [],
      updateTagsDialogVisible: false,
      updateTagForm: {
        articleId: 0,
        tags: [],
      },
    }
  },
  mounted() {
    this.recent()
  },
  methods: {
    async list() {
      this.listLoading = true
      const params = Object.assign(this.filters, {
        page: this.page.page,
        limit: this.page.limit,
      })
      try {
        const data = await this.$axios.post('/api/admin/article/list', params)
        this.results = data.results
        this.page = data.page
      } catch (err) {
        this.$message.error(err.message)
      } finally {
        this.listLoading = false
      }
    },
    async recent() {
      this.listLoading = true
      try {
        this.results = await this.$axios.get('/api/admin/article/recent')
      } catch (err) {
        this.$message.error(err.message)
      } finally {
        this.listLoading = false
      }
    },
    handlePageChange(val) {
      this.page.page = val
      this.list()
    },
    handleLimitChange(val) {
      this.page.limit = val
      this.list()
    },
    deleteSubmit(row) {
      const me = this
      this.$confirm('确认要删除文章？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      })
        .then(() => {
          this.$axios
            .post('/api/admin/article/delete', { id: row.id })
            .then((data) => {
              me.$message({ message: '删除成功', type: 'success' })
              me.list()
            })
            .catch((rsp) => {
              me.$notify.error({ title: '错误', message: rsp.message })
            })
        })
        .catch(() => {
          this.$message({
            type: 'info',
            message: '已取消删除',
          })
        })
    },
    pendingSubmit(row) {
      const me = this
      this.$confirm('确认要过审文章？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      })
        .then(() => {
          this.$axios
            .post('/api/admin/article/pending', { id: row.id })
            .then((data) => {
              me.$message({ message: '审核成功', type: 'success' })
              me.list()
            })
            .catch((rsp) => {
              me.$notify.error({ title: '错误', message: rsp.message })
            })
        })
        .catch(() => {
          this.$message({
            type: 'info',
            message: '已取消审核',
          })
        })
    },
    async showUpdateTags(article) {
      const tags = []
      try {
        const tagObjs = await this.$axios.get(
          '/api/admin/article/tags?articleId=' + article.id
        )
        if (tagObjs && tagObjs.length) {
          for (let i = 0; i < tagObjs.length; i++) {
            tags.push(tagObjs[i].tagName)
          }
        }
      } catch (e) {
        this.$message({
          type: 'error',
          message: e.message || e,
        })
      }

      this.updateTagForm.articleId = article.id
      this.updateTagForm.tags = tags
      this.updateTagsDialogVisible = true
    },
    async updateTags() {
      try {
        const nowTags = await this.$axios.put('/api/admin/article/tags', {
          articleId: this.updateTagForm.articleId,
          tags: (this.updateTagForm.tags || []).join(','),
        })
        if (this.results && this.results.length) {
          for (let i = 0; i < this.results.length; i++) {
            if (this.results[i].id === this.updateTagForm.articleId) {
              this.results[i].tags = nowTags
            }
          }
        }
        this.updateTagsDialogVisible = false
      } catch (e) {
        this.$message({
          type: 'error',
          message: e.message || e,
        })
      }
    },
  },
}
</script>

<style scoped lang="scss">
.articles {
  display: table;
  width: 100%;

  .notification {
    margin: 20px;
    text-align: center;
  }

  .article:not(:last-child) {
    border-bottom: solid 1px rgba(140, 147, 157, 0.14);
  }

  .article {
    width: 100%;
    padding: 10px;

    .article-header {
      display: flex;

      .article-right {
        display: block;
        margin-left: 10px;

        .article-title a {
          color: #555;
          font-size: 16px;
          font-weight: bold;
          cursor: pointer;
          text-decoration: none;
        }

        .article-meta {
          display: flex;
          font-size: 12px;

          label:not(:last-child) {
            margin-right: 8px;
          }

          label {
            color: #999;
            font-size: 12px;
          }

          .article-tags {
            .el-tag + .el-tag {
              margin-left: 5px;
            }

            .button-new-tag {
              margin-left: 10px;
              height: 32px;
              line-height: 30px;
              padding-top: 0;
              padding-bottom: 0;
            }

            .input-new-tag {
              width: 90px;
              margin-left: 10px;
              vertical-align: bottom;
            }
          }
        }
      }
    }

    .summary {
      margin-left: 60px;
      word-break: break-all;
      -webkit-line-clamp: 2;
      overflow: hidden !important;
      text-overflow: ellipsis;
      -webkit-box-orient: vertical;
      display: -webkit-box;
      color: #4a4a4a;
      font-size: 12px;
      font-weight: 400;
      line-height: 1.5;
    }

    .actions {
      text-align: right;
      font-size: 12px;
      font-weight: 400;

      .action-item {
        margin-right: 9px;
      }

      span.danger {
        background: #eee;
        color: red;
        padding: 2px 5px 2px 5px;
      }

      a.btn {
        color: blue;
        cursor: pointer;
      }
    }
  }
}
</style>
