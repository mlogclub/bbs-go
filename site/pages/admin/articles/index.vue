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
            @change="list"
            clearable
            placeholder="请选择状态"
          >
            <el-option label="正常" value="0"></el-option>
            <el-option label="删除" value="1"></el-option>
            <el-option label="待审核" value="2"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button @click="list" type="primary">查询</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!--列表-->
    <div class="page-section articles">
      <div v-for="item in results" :key="item.id" class="article">
        <div class="article-header">
          <img :src="item.user.smallAvatar" class="avatar" />
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
              <label v-for="tag in item.tags" :key="tag.tagId" class="tag">{{
                tag.tagName
              }}</label>

              <div class="actions">
                <span v-if="item.status === 1" class="action-item danger"
                  >已删除</span
                >
                <a
                  v-if="item.status !== 1"
                  @click="deleteSubmit(item)"
                  class="action-item btn"
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
                  @click="PendingSubmit(item)"
                  class="action-item btn"
                  >审核</a
                >
              </div>
            </div>
          </div>
        </div>
        <div class="summary">{{ item.summary }}</div>
      </div>
    </div>

    <!--工具条-->
    <div class="pagebar">
      <el-pagination
        :page-sizes="[20, 50, 100, 300]"
        :current-page="page.page"
        :page-size="page.limit"
        :total="page.total"
        @current-change="handlePageChange"
        @size-change="handleLimitChange"
        layout="total, sizes, prev, pager, next, jumper"
      ></el-pagination>
    </div>
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
        title: '',
        status: ''
      },
      tagOptions: []
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
        limit: me.page.limit
      })
      this.$axios
        .post('/api/admin/article/list', params)
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
    deleteSubmit(row) {
      const me = this
      this.$confirm('确认要删除文章？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
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
            message: '已取消删除'
          })
        })
    },
    PendingSubmit(row) {
      const me = this
      this.$confirm('确认要过审文章？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
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
            message: '已取消审核'
          })
        })
    }
  }
}
</script>

<style scoped lang="scss">
.articles {
  display: table;
  width: 100%;

  .article:not(:last-child) {
    border-bottom: solid 1px rgba(140, 147, 157, 0.14);
  }

  .article {
    width: 100%;
    padding: 10px;

    .article-header {
      display: flex;

      .avatar {
        max-width: 50px;
        max-height: 50px;
        min-width: 50px;
        min-height: 50px;
        border-radius: 50%;
      }

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
