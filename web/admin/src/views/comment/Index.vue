<template>
  <section>

    <!--
    <el-col :span="24" class="toolbar" style="padding-bottom: 0px;">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.userId" placeholder="名称"></el-input>
        </el-form-item>
        <el-form-item>
          <el-select v-model="filters.status" clearable placeholder="请选择状态" @change="list">
            <el-option label="正常" value="0"></el-option>
            <el-option label="删除" value="1"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" v-on:click="list">查询</el-button>
        </el-form-item>
      </el-form>
    </el-col>
    -->


    <!--
    <el-table :data="results" highlight-current-row border v-loading="listLoading"
              style="width: 100%;" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55"></el-table-column>
      <el-table-column prop="id" label="编号"></el-table-column>
      <el-table-column prop="userId" label="用户编号"></el-table-column>
      <el-table-column prop="entityType" label="entityType"></el-table-column>
      <el-table-column prop="entityId" label="entityId"></el-table-column>
      <el-table-column prop="content" label="内容"></el-table-column>
      <el-table-column prop="quoteId" label="引用"></el-table-column>
      <el-table-column prop="status" label="状态"></el-table-column>
      <el-table-column prop="createTime" label="创建时间">
        <template scope="scope">{{scope.row.createTime | formatDate}}</template>
      </el-table-column>
    </el-table>
    -->

    <ul class="comments">
      <li v-for="item in results" :key="item.id">
        <div class="comment-item">
          <div class="avatar" :style="{backgroundImage:'url(' + item.user.avatar + ')'}">

          </div>
          <div class="content">
            <div class="meta">
              <span class="nickname">{{item.user.nickname}}</span>
              <span class="create-time">{{item.createTime | formatDate}}</span>
            </div>
            <div class="summary" v-html="item.content"></div>
          </div>
        </div>

      </li>
    </ul>

    <el-col :span="24" class="toolbar">
      <el-pagination layout="total, sizes, prev, pager, next, jumper" :page-sizes="[20, 50, 100, 300]"
                     @current-change="handlePageChange"
                     @size-change="handleLimitChange"
                     :current-page="page.page"
                     :page-size="page.limit"
                     :total="page.total"
                     style="float:right;">
      </el-pagination>
    </el-col>

  </section>
</template>

<script>
  import HttpClient from '../../apis/HttpClient'

  export default {
    name: 'List',
    data() {
      return {
        results: [],
        listLoading: false,
        page: {},
        filters: {
          userId: '',
          status: ''
        },
        selectedRows: [],

        addForm: {
          'userId': '',
          'entityType': '',
          'entityId': '',
          'content': '',
          'quoteId': '',
          'status': '',
          'createTime': '',
        },
        addFormVisible: false,
        addFormRules: {},
        addLoading: false,

        editForm: {
          'id': '',
          'userId': '',
          'entityType': '',
          'entityId': '',
          'content': '',
          'quoteId': '',
          'status': '',
          'createTime': '',
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
        let me = this
        me.listLoading = true
        let params = Object.assign(me.filters, {
          page: me.page.page,
          limit: me.page.limit
        })
        HttpClient.post('/api/admin/comment/list', params)
          .then(data => {
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
        let me = this
        HttpClient.post('/api/admin/comment/create', this.addForm)
          .then(data => {
            me.$message({message: '提交成功', type: 'success'})
            me.addFormVisible = false
            me.list()
          })
          .catch(rsp => {
            me.$notify.error({title: '错误', message: rsp.message})
          })
      },
      handleEdit(index, row) {
        let me = this
        HttpClient.get('/api/admin/comment/' + row.id)
          .then(data => {
            me.editForm = Object.assign({}, data)
            me.editFormVisible = true
          })
          .catch(rsp => {
            me.$notify.error({title: '错误', message: rsp.message})
          })
      },
      editSubmit() {
        let me = this
        HttpClient.post('/api/admin/comment/update', me.editForm)
          .then(data => {
            me.list()
            me.editFormVisible = false
          })
          .catch(rsp => {
            me.$notify.error({title: '错误', message: rsp.message})
          })
      },

      handleSelectionChange(val) {
        this.selectedRows = val
      },
    }
  }
</script>

<style scoped lang="scss">
  .comments {
    list-style: none;
    padding: 0px;

    li {
      border-bottom: 1px solid #f2f2f2;
      padding-top: 10px;
      padding-bottom: 10px;

      .comment-item {
        display: flex;

        .avatar {
          min-width: 40px;
          min-height: 40px;
          width: 40px;
          height: 40px;
          border-radius: 50%;
          margin-right: 10px;
          background-repeat: no-repeat;
          background-size: contain;
          background-position: center;
        }

        .content {
          .meta {

            span {
              &:not(:last-child) {
                margin-right: 5px;
              }

              &.nickname {
                color: #1a1a1a;
                font-size: 14px;
                font-weight: bold;
              }

              &.create-time {
                color: #999;
                font-size: 13px;
              }
            }
          }

          .summary {
            font-size: 15px;
            color: #555;
          }
        }
      }
    }
  }
</style>

