<template>
  <section class="page-container">
    <div class="toolbar">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.name" placeholder="名称"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="list">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleAdd">新增</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
      v-loading="listLoading"
      :data="results"
      highlight-current-row
      border
      style="width: 100%;"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55"></el-table-column>
      <el-table-column prop="id" label="编号"></el-table-column>
      <el-table-column prop="name" label="名称"></el-table-column>
      <el-table-column prop="description" label="描述"></el-table-column>
      <el-table-column prop="sortNo" label="排序"></el-table-column>
      <el-table-column prop="status" label="状态">
        <template slot-scope="scope">{{
          scope.row.status === 0 ? '启用' : '禁用'
        }}</template>
      </el-table-column>

      <el-table-column prop="createTime" label="创建时间">
        <template slot-scope="scope">{{
          scope.row.createTime | formatDate
        }}</template>
      </el-table-column>

      <el-table-column label="操作" width="150">
        <template slot-scope="scope">
          <el-button size="small" @click="handleEdit(scope.$index, scope.row)"
            >编辑</el-button
          >
        </template>
      </el-table-column>
    </el-table>

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
      <el-form ref="addForm" :model="addForm" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="addForm.name"></el-input>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="addForm.description"
            type="textarea"
            auto-complete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="排序">
          <el-input v-model="addForm.sortNo"></el-input>
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
      <el-form ref="editForm" :model="editForm" label-width="80px">
        <el-input v-model="editForm.id" type="hidden"></el-input>
        <el-form-item label="名称">
          <el-input v-model="editForm.name"></el-input>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="editForm.description"
            type="textarea"
            auto-complete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="排序">
          <el-input v-model="editForm.sortNo"></el-input>
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="editForm.status"
            clearable
            placeholder="请选择状态"
          >
            <el-option :value="0" label="启用"></el-option>
            <el-option :value="1" label="禁用"></el-option>
          </el-select>
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
      filters: {},
      selectedRows: [],

      addForm: {
        name: '',
        description: '',
        status: '',
        sortNo: '',
        createTime: '',
      },
      addFormVisible: false,
      addLoading: false,

      editForm: {
        id: '',
        name: '',
        description: '',
        status: '',
        sortNo: '',
        createTime: '',
      },
      editFormVisible: false,
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
        .post('/api/admin/topic-node/list', params)
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
        .post('/api/admin/topic-node/create', this.addForm)
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
        .get('/api/admin/topic-node/' + row.id)
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
        .post('/api/admin/topic-node/update', me.editForm)
        .then((data) => {
          me.list()
          me.editFormVisible = false
        })
        .catch((rsp) => {
          me.$notify.error({ title: '错误', message: rsp.message })
        })
    },

    handleSelectionChange(val) {
      this.selectedRows = val
    },
  },
}
</script>

<style lang="scss" scoped></style>
