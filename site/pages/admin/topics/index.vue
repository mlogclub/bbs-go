<template>
  <section v-loading="listLoading" class="page-container">
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
            v-model="filters.recommend"
            clearable
            placeholder="是否推荐"
            @change="list"
          >
            <el-option label="推荐" value="1"></el-option>
            <el-option label="未推荐" value="0"></el-option>
          </el-select>
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
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="list">查询</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="page-section topics">
      <div v-for="item in results" :key="item.id" class="topic">
        <div class="topic-header">
          <img :src="item.user.smallAvatar" class="avatar" />
          <div class="topic-right">
            <div class="topic-title">
              <a :href="'/topic/' + item.id" target="_blank">{{
                item.title
              }}</a>
            </div>
            <div class="topic-meta">
              <label>ID: {{ item.id }}</label>
              <label v-if="item.user" class="author">{{
                item.user.nickname
              }}</label>
              <label>{{ item.createTime | formatDate }}</label>
              <label class="node">{{ item.node ? item.node.name : '' }}</label>
              <label v-for="tag in item.tags" :key="tag.tagId" class="tag">{{
                tag.tagName
              }}</label>

              <div class="actions">
                <span v-if="item.status === 1" class="action-item danger"
                  >已删除</span
                >
                <a
                  v-if="item.status === 0"
                  class="action-item btn"
                  @click="deleteSubmit(item)"
                  >删除</a
                >
                <a v-else class="action-item btn" @click="undeleteSubmit(item)"
                  >取消删除</a
                >
                <a
                  v-if="!item.recommend"
                  class="action-item btn"
                  @click="recommend(item.id)"
                  >推荐</a
                >
                <a
                  v-else
                  class="action-item btn"
                  @click="cancelRecommend(item.id)"
                  >取消推荐</a
                >
              </div>
            </div>
          </div>
        </div>

        <div class="summary">
          {{ item.summary }}
        </div>
      </div>
    </div>

    <div class="pagebar">
      <el-pagination
        :page-sizes="[20, 50, 100, 300]"
        :current-page="page.page"
        :page-size="page.limit"
        :total="page.total"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handlePageChange"
        @size-change="handleLimitChange"
      >
      </el-pagination>
    </div>

    <el-dialog
      :visible.sync="addFormVisible"
      :close-on-click-modal="false"
      title="新增"
    >
      <el-form
        ref="addForm"
        :model="addForm"
        :rules="addFormRules"
        label-width="80px"
      >
        <el-form-item label="userId" prop="rule">
          <el-input v-model="addForm.userId"></el-input>
        </el-form-item>

        <el-form-item label="title" prop="rule">
          <el-input v-model="addForm.title"></el-input>
        </el-form-item>

        <el-form-item label="content" prop="rule">
          <el-input v-model="addForm.content"></el-input>
        </el-form-item>

        <el-form-item label="status" prop="rule">
          <el-input v-model="addForm.status"></el-input>
        </el-form-item>

        <el-form-item label="createTime" prop="rule">
          <el-input v-model="addForm.createTime"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="addFormVisible = false">取消</el-button>
        <el-button
          :loading="addLoading"
          type="primary"
          @click.native="addSubmit"
          >提交</el-button
        >
      </div>
    </el-dialog>

    <el-dialog
      :visible.sync="editFormVisible"
      :close-on-click-modal="false"
      title="编辑"
    >
      <el-form
        ref="editForm"
        :model="editForm"
        :rules="editFormRules"
        label-width="80px"
      >
        <el-input v-model="editForm.id" type="hidden"></el-input>

        <el-form-item label="forumId" prop="rule">
          <el-input v-model="editForm.forumId"></el-input>
        </el-form-item>

        <el-form-item label="userId" prop="rule">
          <el-input v-model="editForm.userId"></el-input>
        </el-form-item>

        <el-form-item label="title" prop="rule">
          <el-input v-model="editForm.title"></el-input>
        </el-form-item>

        <el-form-item label="content" prop="rule">
          <el-input v-model="editForm.content"></el-input>
        </el-form-item>

        <el-form-item label="status" prop="rule">
          <el-input v-model="editForm.status"></el-input>
        </el-form-item>

        <el-form-item label="createTime" prop="rule">
          <el-input v-model="editForm.createTime"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="editFormVisible = false">取消</el-button>
        <el-button
          :loading="editLoading"
          type="primary"
          @click.native="editSubmit"
          >提交</el-button
        >
      </div>
    </el-dialog>
  </section>
</template>

<script>
export default {
  layout: 'admin',
  data() {
    return {
      results: [],
      listLoading: false,
      page: {},
      filters: {
        status: '0',
      },
      selectedRows: [],
      addForm: {
        forumId: '',
        userId: '',
        title: '',
        content: '',
        status: '',
        createTime: '',
      },
      addFormVisible: false,
      addFormRules: {},
      addLoading: false,

      editForm: {
        id: '',
        forumId: '',
        userId: '',
        title: '',
        content: '',
        status: '',
        createTime: '',
      },
      editFormVisible: false,
      editFormRules: {},
      editLoading: false,
    }
  },
  mounted() {
    this.list()
  },
  methods: {
    list() {
      const me = this
      me.listLoading = true
      const params = Object.assign(me.filters, {
        page: me.page.page,
        limit: me.page.limit,
      })
      this.$axios
        .post('/api/admin/topic/list', params)
        .then((data) => {
          me.results = data.results
          me.page = data.page
        })
        .finally(() => {
          me.listLoading = false
        })
    },
    handlePageChange(val) {
      this.page.page = val
      this.list()
    },
    handleLimitChange(val) {
      this.page.limit = val
      this.list()
    },
    handleAdd() {
      this.addForm = {
        name: '',
        description: '',
      }
      this.addFormVisible = true
    },
    addSubmit() {
      const me = this
      this.$axios
        .post('/api/admin/topic/create', this.addForm)
        .then((data) => {
          me.$message({ message: '提交成功', type: 'success' })
          me.addFormVisible = false
          me.list()
        })
        .catch((rsp) => {
          me.$notify.error({ title: '错误', message: rsp.message })
        })
    },
    handleEdit(index, row) {
      const me = this
      this.$axios
        .get(`/api/admin/topic/${row.id}`)
        .then((data) => {
          me.editForm = Object.assign({}, data)
          me.editFormVisible = true
        })
        .catch((rsp) => {
          me.$notify.error({ title: '错误', message: rsp.message })
        })
    },
    editSubmit() {
      const me = this
      this.$axios
        .post('/api/admin/topic/update', me.editForm)
        .then((data) => {
          me.list()
          me.editFormVisible = false
        })
        .catch((rsp) => {
          me.$notify.error({ title: '错误', message: rsp.message })
        })
    },
    deleteSubmit(row) {
      const me = this
      this.$confirm('是否确认删除该话题?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      })
        .then(function () {
          me.$axios
            .post('/api/admin/topic/delete', { id: row.id })
            .then(function () {
              me.$message({ message: '删除成功', type: 'success' })
              me.list()
            })
            .catch(function (err) {
              me.$notify.error({ title: '错误', message: err.message || err })
            })
        })
        .catch(function () {
          me.$message({
            type: 'info',
            message: '已取消删除',
          })
        })
    },
    async undeleteSubmit(row) {
      try {
        await this.$axios.post('/api/admin/topic/undelete', { id: row.id })
        this.list()
        this.$message({ message: '取消删除成功', type: 'success' })
      } catch (err) {
        this.$notify.error({ title: '错误', message: err.message || err })
      }
    },
    async recommend(id) {
      try {
        await this.$axios.post('/api/admin/topic/recommend', {
          id,
        })
        this.$message({ message: '推荐成功', type: 'success' })
        this.list()
      } catch (e) {
        this.$notify.error({ title: '错误', message: e.message })
      }
    },
    async cancelRecommend(id) {
      try {
        await this.$axios.delete('/api/admin/topic/recommend', {
          params: {
            id,
          },
        })
        this.$message({ message: '取消推荐成功', type: 'success' })
        this.list()
      } catch (e) {
        this.$notify.error({ title: '错误', message: e.message })
      }
    },
    handleSelectionChange(val) {
      this.selectedRows = val
    },
  },
}
</script>

<style scoped lang="scss">
.topics {
  width: 100%;

  .topic:not(:last-child) {
    border-bottom: solid 1px rgba(140, 147, 157, 0.14);
  }

  .topic {
    width: 100%;
    padding: 10px;

    .topic-header {
      display: flex;

      .avatar {
        max-width: 50px;
        max-height: 50px;
        min-width: 50px;
        min-height: 50px;
        border-radius: 50%;
      }

      .topic-right {
        display: block;
        margin-left: 10px;

        .topic-title a {
          color: #555;
          font-size: 16px;
          font-weight: bold;
          cursor: pointer;
          text-decoration: none;
        }

        .topic-meta {
          display: flex;
          font-size: 12px;

          label:not(:last-child) {
            margin-right: 8px;
          }

          label {
            color: #999;
          }

          label.author {
            font-weight: bold;
          }

          label.node {
            color: #dc2323;
            font-weight: bold;
          }

          label.tag {
            align-items: center;
            background-color: #f5f5f5;
            border-radius: 4px;
            color: #4a4a4a;
            display: inline-flex;
            justify-content: center;
            line-height: 1.5;
            padding-left: 5px;
            padding-right: 5px;
            white-space: nowrap;
          }

          .actions {
            margin-left: 20px;
            text-align: right;

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
  }
}
</style>
