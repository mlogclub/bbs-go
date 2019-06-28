
<template>
    <section>
        
        <el-col :span="24" class="toolbar" style="padding-bottom: 0px;">
            <el-form :inline="true" :model="filters">
                <el-form-item>
                    <el-input v-model="filters.name" placeholder="名称"></el-input>
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
            <el-table-column prop="id" label="编号"></el-table-column>
            
                <el-table-column prop="expiredAt" label="expiredAt"></el-table-column>
            
                <el-table-column prop="code" label="code"></el-table-column>
            
                <el-table-column prop="accessToken" label="accessToken"></el-table-column>
            
                <el-table-column prop="refreshToken" label="refreshToken"></el-table-column>
            
                <el-table-column prop="data" label="data"></el-table-column>
            
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
                
                    <el-form-item label="expiredAt" prop="rule">
                        <el-input v-model="addForm.expiredAt"></el-input>
                    </el-form-item>
                
                    <el-form-item label="code" prop="rule">
                        <el-input v-model="addForm.code"></el-input>
                    </el-form-item>
                
                    <el-form-item label="accessToken" prop="rule">
                        <el-input v-model="addForm.accessToken"></el-input>
                    </el-form-item>
                
                    <el-form-item label="refreshToken" prop="rule">
                        <el-input v-model="addForm.refreshToken"></el-input>
                    </el-form-item>
                
                    <el-form-item label="data" prop="rule">
                        <el-input v-model="addForm.data"></el-input>
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
                
                    <el-form-item label="expiredAt" prop="rule">
                        <el-input v-model="editForm.expiredAt"></el-input>
                    </el-form-item>
                
                    <el-form-item label="code" prop="rule">
                        <el-input v-model="editForm.code"></el-input>
                    </el-form-item>
                
                    <el-form-item label="accessToken" prop="rule">
                        <el-input v-model="editForm.accessToken"></el-input>
                    </el-form-item>
                
                    <el-form-item label="refreshToken" prop="rule">
                        <el-input v-model="editForm.refreshToken"></el-input>
                    </el-form-item>
                
                    <el-form-item label="data" prop="rule">
                        <el-input v-model="editForm.data"></el-input>
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
    name: "List",
    data() {
      return {
        results: [],
        listLoading: false,
        page: {},
        filters: {},
        selectedRows: [],

        addForm: {
          
            'expiredAt': '',
          
            'code': '',
          
            'accessToken': '',
          
            'refreshToken': '',
          
            'data': '',
          
        },
        addFormVisible: false,
        addFormRules: {},
        addLoading: false,

        editForm: {
          'id': '',
          
            'expiredAt': '',
          
            'code': '',
          
            'accessToken': '',
          
            'refreshToken': '',
          
            'data': '',
          
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
        HttpClient.post('/api/admin/oauth-token/list', params)
          .then(data => {
            me.results = data.results
            me.page = data.page
          })
          .finally(() => {
            me.listLoading = false
          })
      },
      handlePageChange (val) {
        this.page.page = val
        this.list()
      },
      handleLimitChange (val) {
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
        HttpClient.post('/api/admin/oauth-token/create', this.addForm)
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
        HttpClient.get('/api/admin/oauth-token/' + row.id)
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
        HttpClient.post('/api/admin/oauth-token/update', me.editForm)
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

