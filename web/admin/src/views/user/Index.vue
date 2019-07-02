<template>
  <section>

    <el-col :span="24" class="toolbar" style="padding-bottom: 0px;">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.id" placeholder="编号"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" v-on:click="list">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleAdd">新增</el-button>
        </el-form-item>
      </el-form>
    </el-col>


    <el-table :data="results" highlight-current-row border v-loading="listLoading"
              style="width: 100%;" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55"></el-table-column>
      <el-table-column prop="id" label="编号" width="100"></el-table-column>
      <el-table-column prop="avatar" label="头像" width="80">
        <template scope="scope">
          <img :src="scope.row.avatar" style="max-height: 50px; max-width: 50px; border-radius: 50%;"/>
        </template>
      </el-table-column>
      <el-table-column prop="username" label="用户名信息">
        <template scope="scope">
          <div>{{scope.row.nickname}}</div>
          <div>{{scope.row.username}}</div>
          <div>{{scope.row.email}}</div>
          <div v-if="scope.row.roles && scope.row.roles.length">
            <el-tag size="mini" v-for="role in scope.row.roles" :key="role"
                    style="margin-right:3px;">{{role}}
            </el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态">
        <template scope="scope">{{scope.row.status === 0 ? '正常' : '删除'}}</template>
      </el-table-column>
      <el-table-column label="时间" width="200">
        <template scope="scope">
          注册：{{scope.row.createTime | formatDate}}<br/>
          更新：{{scope.row.updateTime | formatDate}}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150">
        <template scope="scope">
          <el-button size="small" @click="handleEdit(scope.$index, scope.row)">编辑</el-button>
        </template>
      </el-table-column>
    </el-table>


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

        <el-form-item label="用户名" prop="username">
          <el-input v-model="addForm.username"></el-input>
        </el-form-item>

        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="addForm.nickname"></el-input>
        </el-form-item>

        <el-form-item label="邮箱" prop="email">
          <el-input v-model="addForm.email"></el-input>
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input v-model="addForm.password"></el-input>
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

        <el-form-item label="用户名" prop="username">
          <el-input v-model="editForm.username"></el-input>
        </el-form-item>

        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="editForm.nickname"></el-input>
        </el-form-item>

        <el-form-item label="邮箱" prop="email">
          <el-input v-model="editForm.email"></el-input>
        </el-form-item>

        <el-form-item label="角色" prop="roles">
          <el-select v-model="editForm.roles" multiple filterable allow-create default-first-option placeholder="用户角色"
                     style="width: 100%">
            <el-option v-for="item in editForm.roles" :key="item" :label="item"
                       :value="item"></el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input v-model="editForm.password" placeholder="不填写标识不更改密码"></el-input>
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-select v-model="editForm.status" placeholder="请选择">
            <el-option :key="0" :value="0" label="正常"></el-option>
            <el-option :key="1" :value="1" label="删除"></el-option>
          </el-select>
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
  import HttpClient from '../../apis/HttpClient'

  export default {
    name: 'List',
    data() {
      return {
        results: [],
        listLoading: false,
        page: {},
        filters: {
          id: ''
        },
        selectedRows: [],

        addForm: {
          'username': '',
          'nickname': '',
          'avatar': '',
          'email': '',
          'roles': [],
          'password': '',
          'status': ''
        },
        addFormVisible: false,
        addFormRules: {},
        addLoading: false,

        editForm: {
          'id': '',
          'username': '',
          'nickname': '',
          'avatar': '',
          'email': '',
          'roles': [],
          'password': '',
          'status': '',
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
        HttpClient.post('/api/admin/user/list', params)
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
        HttpClient.post('/api/admin/user/create', this.addForm)
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
        HttpClient.get('/api/admin/user/' + row.id)
          .then(data => {
            me.editForm = Object.assign({}, data)
            me.editForm.password = ''
            me.editFormVisible = true
          })
          .catch(rsp => {
            me.$notify.error({title: '错误', message: rsp.message})
          })
      },
      editSubmit() {
        let me = this
        HttpClient.post('/api/admin/user/update', me.editForm)
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

<style scoped>

</style>

