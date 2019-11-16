<template>
  <section class="page-container">
    <el-col :span="24" class="toolbar">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.title" placeholder="标题"></el-input>
        </el-form-item>
        <el-form-item>
          <el-select
            v-model="filters.status"
            placeholder="请选择"
            @change="list"
          >
            <el-option :key="0" :value="0" label="启用"></el-option>
            <el-option :key="1" :value="1" label="禁用"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" v-on:click="list">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleAdd">新增</el-button>
        </el-form-item>
      </el-form>
    </el-col>

    <el-table
      :data="results"
      highlight-current-row
      stripe
      v-loading="listLoading"
      style="width: 100%;"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55"></el-table-column>
      <el-table-column prop="id" label="编号"></el-table-column>
      <el-table-column prop="title" label="标题"></el-table-column>
      <el-table-column prop="description" label="描述"></el-table-column>
      <el-table-column prop="logo" label="Logo"></el-table-column>
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

    <el-col :span="24" class="toolbar">
      <el-pagination
        layout="total, sizes, prev, pager, next, jumper"
        :page-sizes="[20, 50, 100, 300]"
        @current-change="handlePageChange"
        @size-change="handleLimitChange"
        :current-page="page.page"
        :page-size="page.limit"
        :total="page.total"
        style="float:right;"
      ></el-pagination>
    </el-col>

    <el-dialog
      title="新增"
      :visible.sync="addFormVisible"
      :close-on-click-modal="false"
    >
      <el-form
        :model="addForm"
        label-width="80px"
        :rules="addFormRules"
        ref="addForm"
      >
        <el-form-item label="标题">
          <el-input v-model="addForm.title"></el-input>
        </el-form-item>

        <el-form-item label="描述">
          <el-input v-model="addForm.description"></el-input>
        </el-form-item>

        <el-form-item label="Logo">
          <el-input v-model="addForm.logo"></el-input>
        </el-form-item>

        <el-form-item label="状态">
          <el-select v-model="addForm.status" placeholder="请选择">
            <el-option :key="0" :value="0" label="启用"></el-option>
            <el-option :key="1" :value="1" label="禁用"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="addFormVisible = false">取消</el-button>
        <el-button
          type="primary"
          @click.native="addSubmit"
          :loading="addLoading"
          >提交</el-button
        >
      </div>
    </el-dialog>

    <el-dialog
      title="编辑"
      :visible.sync="editFormVisible"
      :close-on-click-modal="false"
    >
      <el-form
        :model="editForm"
        label-width="80px"
        :rules="editFormRules"
        ref="editForm"
      >
        <el-input v-model="editForm.id" type="hidden"></el-input>

        <el-form-item label="标题">
          <el-input v-model="editForm.title"></el-input>
        </el-form-item>

        <el-form-item label="描述">
          <el-input v-model="editForm.description"></el-input>
        </el-form-item>

        <el-form-item label="Logo">
          <el-input v-model="editForm.logo"></el-input>
        </el-form-item>

        <el-form-item label="状态">
          <el-select v-model="editForm.status" placeholder="请选择">
            <el-option :key="0" :value="0" label="启用"></el-option>
            <el-option :key="1" :value="1" label="禁用"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="editFormVisible = false">取消</el-button>
        <el-button
          type="primary"
          @click.native="editSubmit"
          :loading="editLoading"
          >提交</el-button
        >
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
      filters: {},
      selectedRows: [],

      addForm: {
        title: '',
        description: '',
        logo: '',
        status: '',
        createTime: ''
      },
      addFormVisible: false,
      addFormRules: {},
      addLoading: false,

      editForm: {
        id: '',
        title: '',
        description: '',
        logo: '',
        status: '',
        createTime: ''
      },
      editFormVisible: false,
      editFormRules: {},
      editLoading: false
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
      HttpClient.post('/api/admin/subject/list', params)
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
      // this.addForm = {
      //   name: "",
      //   description: ""
      // };
      this.addForm = {}
      this.addFormVisible = true
    },
    addSubmit() {
      let me = this
      HttpClient.post('/api/admin/subject/create', this.addForm)
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
      let me = this
      HttpClient.get('/api/admin/subject/' + row.id)
        .then((data) => {
          me.editForm = Object.assign({}, data)
          me.editFormVisible = true
        })
        .catch((rsp) => {
          me.$notify.error({ title: '错误', message: rsp.message })
        })
    },
    editSubmit() {
      let me = this
      HttpClient.post('/api/admin/subject/update', me.editForm)
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
    }
  }
}
</script>

<style scoped></style>
