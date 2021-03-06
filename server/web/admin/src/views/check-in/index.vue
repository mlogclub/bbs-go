
<template>
    <section class="page-container">
        
        <el-col :span="24" class="toolbar">
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
            
			<el-table-column prop="userId" label="userId"></el-table-column>
			
			<el-table-column prop="latestDayName" label="latestDayName"></el-table-column>
			
			<el-table-column prop="consecutiveDays" label="consecutiveDays"></el-table-column>
			
			<el-table-column prop="createTime" label="createTime"></el-table-column>
			
			<el-table-column prop="updateTime" label="updateTime"></el-table-column>
			
            <el-table-column label="操作" width="150">
                <template slot-scope="scope">
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
            <el-form :model="addForm" label-width="80px" ref="addForm">
                
				<el-form-item label="userId">
					<el-input v-model="addForm.userId"></el-input>
				</el-form-item>
                
				<el-form-item label="latestDayName">
					<el-input v-model="addForm.latestDayName"></el-input>
				</el-form-item>
                
				<el-form-item label="consecutiveDays">
					<el-input v-model="addForm.consecutiveDays"></el-input>
				</el-form-item>
                
				<el-form-item label="createTime">
					<el-input v-model="addForm.createTime"></el-input>
				</el-form-item>
                
				<el-form-item label="updateTime">
					<el-input v-model="addForm.updateTime"></el-input>
				</el-form-item>
                
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click.native="addFormVisible = false">取消</el-button>
                <el-button type="primary" @click.native="addSubmit" :loading="addLoading">提交</el-button>
            </div>
        </el-dialog>

        
        <el-dialog title="编辑" :visible.sync="editFormVisible" :close-on-click-modal="false">
            <el-form :model="editForm" label-width="80px" ref="editForm">
                <el-input v-model="editForm.id" type="hidden"></el-input>
                
				<el-form-item label="userId">
					<el-input v-model="editForm.userId"></el-input>
				</el-form-item>
                
				<el-form-item label="latestDayName">
					<el-input v-model="editForm.latestDayName"></el-input>
				</el-form-item>
                
				<el-form-item label="consecutiveDays">
					<el-input v-model="editForm.consecutiveDays"></el-input>
				</el-form-item>
                
				<el-form-item label="createTime">
					<el-input v-model="editForm.createTime"></el-input>
				</el-form-item>
                
				<el-form-item label="updateTime">
					<el-input v-model="editForm.updateTime"></el-input>
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
  export default {
    data() {
      return {
        results: [],
        listLoading: false,
        page: {},
        filters: {},
        selectedRows: [],

        addForm: {},
        addFormVisible: false,
        addLoading: false,

        editForm: {},
        editFormVisible: false,
        editLoading: false,
      }
    },
    mounted() {
      this.list();
    },
    methods: {
		async list() {
		  const params = Object.assign(this.filters, {
			page: this.page.page,
			limit: this.page.limit,
		  })
		  try {
			const data = await this.$axios.post('/api/admin/check-in/list', params)
			this.results = data.results
			this.page = data.page
		  } catch (e) {
			this.$notify.error({ title: '错误', message: e || e.message })
		  } finally {
			this.listLoading = false
		  }
		},
		async handlePageChange(val) {
		  this.page.page = val
		  await this.list()
		},
		async handleLimitChange(val) {
		  this.page.limit = val
		  await this.list()
		},
		handleAdd() {
		  this.addForm = {
			name: '',
			description: '',
		  }
		  this.addFormVisible = true
		},
		async addSubmit() {
		  try {
			await this.$axios.post('/api/admin/check-in/create', this.addForm)
			this.$message({ message: '提交成功', type: 'success' })
			this.addFormVisible = false
			await this.list()
		  } catch (e) {
			this.$notify.error({ title: '错误', message: e || e.message })
		  }
		},
		async handleEdit(index, row) {
		  try {
			const data = await this.$axios.get('/api/admin/check-in/' + row.id)
			this.editForm = Object.assign({}, data)
			this.editFormVisible = true
		  } catch (e) {
			this.$notify.error({ title: '错误', message: e || e.message })
		  }
		},
		async editSubmit() {
		  try {
			await this.$axios.post('/api/admin/check-in/update', this.editForm)
			await this.list()
			this.editFormVisible = false
		  } catch (e) {
			this.$notify.error({ title: '错误', message: e || e.message })
		  }
		},
	
		handleSelectionChange(val) {
		  this.selectedRows = val
		},
    }
  }
</script>

<style lang="scss" scoped>

</style>

