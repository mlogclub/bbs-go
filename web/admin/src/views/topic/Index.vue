<template>
  <section>

    <el-col :span="24" class="toolbar" style="padding-bottom: 0px;">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.title" placeholder="标题"></el-input>
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
        <!--
        <el-form-item>
          <el-button type="primary" @click="handleAdd">新增</el-button>
        </el-form-item>
        -->
      </el-form>
    </el-col>


    <div class="topics">
      <div class="topic" v-for="item in results" :key="item.id">

        <div class="topic-header">
          <img class="avatar" :src="item.user.avatar"/>
          <div class="topic-right">
            <div class="topic-title">
              <a @click="toTopic(item)" href="javascript:void(0)">{{item.title}}</a>
            </div>
            <div class="topic-meta">
              <label class="author" v-if="item.user">{{item.user.nickname}}</label>
              <label>{{item.createTime | formatDate}}</label>
              <label class="tag" v-for="tag in item.tags" :key="tag.tagId">{{tag.tagName}}</label>
            </div>
          </div>
        </div>

        <div class="summary">
          {{item.summary}}
        </div>

        <div class="topic-footer">
          <span class="danger" v-if="item.status === 1">已删除</span>
          <span class="info">编号：{{item.id}}</span>
          <a class="btn" @click="deleteSubmit(item)">删除</a>
        </div>
      </div>
    </div>


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


    <el-dialog title="新增" :visible.sync="addFormVisible" :close-on-click-modal="false">
      <el-form :model="addForm" label-width="80px" :rules="addFormRules" ref="addForm">

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
        <el-button type="primary" @click.native="addSubmit" :loading="addLoading">提交</el-button>
      </div>
    </el-dialog>


    <el-dialog title="编辑" :visible.sync="editFormVisible" :close-on-click-modal="false">
      <el-form :model="editForm" label-width="80px" :rules="editFormRules" ref="editForm">
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
        <el-button type="primary" @click.native="editSubmit" :loading="editLoading">提交</el-button>
      </div>
    </el-dialog>
  </section>
</template>

<script>
  import config from '../../apis/Config'
  import HttpClient from '../../apis/HttpClient'

  export default {
    name: "List",
    data() {
      return {
        results: [],
        listLoading: false,
        page: {},
        filters: {},
        selectedRows: [],

        addForm: {
          'forumId': '',
          'userId': '',
          'title': '',
          'content': '',
          'status': '',
          'createTime': '',
        },
        addFormVisible: false,
        addFormRules: {},
        addLoading: false,

        editForm: {
          'id': '',
          'forumId': '',
          'userId': '',
          'title': '',
          'content': '',
          'status': '',
          'createTime': '',
        },
        editFormVisible: false,
        editFormRules: {},
        editLoading: false,
      }
    },
    mounted() {
      this.list();
    },
    methods: {
      list() {
        let me = this
        me.listLoading = true
        let params = Object.assign(me.filters, {
          page: me.page.page,
          limit: me.page.limit
        })
        HttpClient.post('/api/admin/topic/list', params)
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
        HttpClient.post('/api/admin/topic/create', this.addForm)
          .then(data => {
            me.$message({message: '提交成功', type: 'success'});
            me.addFormVisible = false
            me.list()
          })
          .catch(rsp => {
            me.$notify.error({title: '错误', message: rsp.message})
          })
      },
      handleEdit(index, row) {
        let me = this
        HttpClient.get('/api/admin/topic/' + row.id)
          .then(data => {
            me.editForm = Object.assign({}, data);
            me.editFormVisible = true
          })
          .catch(rsp => {
            me.$notify.error({title: '错误', message: rsp.message})
          })
      },
      editSubmit() {
        let me = this
        HttpClient.post('/api/admin/topic/update', me.editForm)
          .then(data => {
            me.list()
            me.editFormVisible = false
          })
          .catch(rsp => {
            me.$notify.error({title: '错误', message: rsp.message})
          })
      },
      deleteSubmit(row) {
        let me = this
        HttpClient.post('/api/admin/topic/delete', {id: row.id})
          .then(data => {
            me.$message({message: '删除成功', type: 'success'})
            me.list()
          })
          .catch(rsp => {
            me.$notify.error({title: '错误', message: rsp.message})
          })
      },
      handleSelectionChange(val) {
        this.selectedRows = val
      },

      toTopic(row) {
        window.open(config.host + '/topic/' + row.id, '_blank')
      }
    }
  }
</script>

<style scoped lang="scss">
  .topics {
    display: table;

    .topic:not(:last-child) {
      border-bottom: solid 1px rgba(140, 147, 157, 0.14);
    }

    .topic {

      padding-top: 10px;
      padding-bottom: 10px;

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
            margin-top: 8px;

            label:not(:last-child) {
              margin-right: 5px;
            }

            label {
              color: #999;
              font-size: 12px;
            }

            label.author {
              /*color: #dc2323;*/
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
          }
        }
      }

      .summary {
        margin-top: 10px;
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

      .topic-footer {
        text-align: right;

        span.info {
          font-size: 12px;
          margin-right: 10px;
          background: #eee;
          padding: 2px 5px 2px 5px;
        }

        span.danger {
          font-size: 12px;
          margin-right: 10px;
          background: #eee;
          color: red;
          padding: 2px 5px 2px 5px;
        }

        a.btn {
          font-size: 12px;
          margin-right: 10px;
          color: blue;
          cursor: pointer;
        }
      }

    }

  }
</style>


