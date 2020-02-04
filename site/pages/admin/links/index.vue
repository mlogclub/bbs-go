<template>
  <section class="page-container">
    <div class="toolbar">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.title" placeholder="标题"></el-input>
        </el-form-item>
        <el-form-item>
          <el-input v-model="filters.url" placeholder="链接"></el-input>
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
          <el-button v-on:click="list" type="primary">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button @click="handleAdd" type="primary">新增</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
      :data="results"
      v-loading="listLoading"
      @selection-change="handleSelectionChange"
      highlight-current-row
      stripe
      style="width: 100%;"
    >
      <el-table-column type="expand">
        <template slot-scope="scope">
          <div class="content-form">
            <div v-if="scope.row.category" class="form-item">
              <div class="field-key">分类：</div>
              <div class="field-value">{{ scope.row.category }}</div>
            </div>
          </div>
          <div class="content-form">
            <div v-if="scope.row.summary" class="form-item">
              <div class="field-key">描述：</div>
              <div class="field-value">{{ scope.row.summary }}</div>
            </div>
          </div>
          <div class="content-form">
            <div class="form-item">
              <div class="field-key">创建时间：</div>
              <div class="field-value">
                {{ scope.row.createTime | formatDate }}
              </div>
            </div>
          </div>
          <div class="content-form">
            <div v-if="scope.row.remark" class="form-item">
              <div class="field-key">备注：</div>
              <div class="field-value">{{ scope.row.remark }}</div>
            </div>
          </div>
        </template>
      </el-table-column>
      <!-- <el-table-column type="selection" width="55"></el-table-column> -->
      <el-table-column prop="id" label="编号" width="100"></el-table-column>
      <!-- <el-table-column prop="category" label="分类"></el-table-column> -->
      <el-table-column prop="url" label="链接">
        <template slot-scope="scope">
          <a :href="scope.row.url" target="_blank">{{ scope.row.url }}</a>
        </template>
      </el-table-column>
      <el-table-column prop="title" label="标题"></el-table-column>
      <!-- <el-table-column prop="summary" label="描述"></el-table-column> -->
      <!-- <el-table-column prop="logo" label="Logo"></el-table-column> -->
      <el-table-column prop="status" label="状态" width="50">
        <template slot-scope="scope">{{
          scope.row.status === 0
            ? '正常'
            : scope.row.status === 1
            ? '删除'
            : '待审核'
        }}</template>
      </el-table-column>
      <el-table-column prop="score" label="分数" width="80">
        <template slot-scope="scope">{{ scope.row.score || 0 }}</template>
      </el-table-column>
      <!--
      <el-table-column prop="createTime" label="创建时间">
        <template slot-scope="scope">{{scope.row.createTime | formatDate}}</template>
      </el-table-column>
      -->
      <el-table-column label="操作" width="150">
        <template slot-scope="scope">
          <el-button @click="handleEdit(scope.$index, scope.row)" size="small"
            >编辑</el-button
          >
        </template>
      </el-table-column>
    </el-table>

    <div class="pagebar">
      <el-pagination
        :page-sizes="[20, 50, 100, 300]"
        @current-change="handlePageChange"
        @size-change="handleLimitChange"
        :current-page="page.page"
        :page-size="page.limit"
        :total="page.total"
        layout="total, sizes, prev, pager, next, jumper"
      ></el-pagination>
    </div>

    <el-dialog
      :visible.sync="addFormVisible"
      :close-on-click-modal="false"
      title="新增"
    >
      <el-form ref="addForm" :model="addForm" label-width="80px">
        <el-form-item label="链接">
          <el-input v-model="addForm.url" style="width: 80%;"></el-input>&nbsp;
          <el-button @click="detect" type="primary">Detect</el-button>
        </el-form-item>

        <el-form-item label="标题">
          <el-input v-model="addForm.title"></el-input>
        </el-form-item>

        <el-form-item label="描述">
          <el-input
            v-model="addForm.summary"
            :autosize="{ minRows: 2, maxRows: 4 }"
          ></el-input>
        </el-form-item>

        <el-form-item label="Logo">
          <el-input v-model="addForm.logo"></el-input>
        </el-form-item>

        <el-form-item label="分类">
          <el-input v-model="addForm.category"></el-input>
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-select v-model="addForm.status" placeholder="请选择">
            <el-option :key="0" :value="0" label="正常"></el-option>
            <el-option :key="1" :value="1" label="删除"></el-option>
            <el-option :key="2" :value="2" label="待审核"></el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="分数">
          <el-input-number
            v-model="addForm.score"
            :min="1"
            :max="100"
            label="分数越高越优质"
          ></el-input-number>
        </el-form-item>

        <el-form-item label="备注">
          <el-input
            v-model="addForm.remark"
            :autosize="{ minRows: 2, maxRows: 4 }"
          ></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="addFormVisible = false">取消</el-button>
        <el-button
          @click.native="addSubmit"
          :loading="addLoading"
          type="primary"
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

        <el-form-item label="链接">
          <el-input v-model="editForm.url"></el-input>
        </el-form-item>

        <el-form-item label="标题">
          <el-input v-model="editForm.title"></el-input>
        </el-form-item>

        <el-form-item label="描述">
          <el-input
            v-model="editForm.summary"
            :autosize="{ minRows: 2, maxRows: 4 }"
          ></el-input>
        </el-form-item>

        <el-form-item label="Logo">
          <el-input v-model="editForm.logo"></el-input>
        </el-form-item>

        <el-form-item label="分类">
          <el-input v-model="editForm.category"></el-input>
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-select v-model="editForm.status" placeholder="请选择">
            <el-option :key="0" :value="0" label="正常"></el-option>
            <el-option :key="1" :value="1" label="删除"></el-option>
            <el-option :key="2" :value="2" label="待审核"></el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="分数">
          <el-input-number
            v-model="editForm.score"
            :min="1"
            :max="100"
            label="分数越高越优质"
          ></el-input-number>
        </el-form-item>

        <el-form-item label="备注">
          <el-input
            v-model="editForm.remark"
            :autosize="{ minRows: 2, maxRows: 4 }"
          ></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="editFormVisible = false">取消</el-button>
        <el-button
          @click.native="editSubmit"
          :loading="editLoading"
          type="primary"
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

      addForm: {},
      addFormVisible: false,
      addFormRules: {},
      addLoading: false,

      editForm: {},
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
      const me = this
      me.listLoading = true
      const params = Object.assign(me.filters, {
        page: me.page.page,
        limit: me.page.limit
      })
      this.$axios
        .post('/api/admin/link/list', params)
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
      this.addForm = {}
      this.addFormVisible = true
    },
    addSubmit() {
      const me = this
      this.$axios
        .post('/api/admin/link/create', this.addForm)
        .then((data) => {
          me.$message({ message: '提交成功', type: 'success' })
          me.addFormVisible = false
          me.list()
        })
        .catch((rsp) => {
          me.$notify.error({ title: '错误', message: rsp.message })
        })
    },
    async detect() {
      if (!this.addForm.url) {
        return
      }
      try {
        const flag = await this.$confirm(
          '确定采集吗，采集之后将覆盖现有内容?',
          '提示',
          { type: 'warning' }
        )
        if (flag) {
          const data = await this.$axios.get('/api/admin/link/detect', {
            url: this.addForm.url
          })
          if (data) {
            this.addForm.title = data.title
            this.addForm.summary = data.description
          }
        }
      } catch (e) {
        this.$notify.error({ title: '错误', message: e.message || e })
      }
    },
    handleEdit(index, row) {
      const me = this
      this.$axios
        .get(`/api/admin/link/${row.id}`)
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
        .post('/api/admin/link/update', me.editForm)
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

<style scoped>
.demo-table-expand {
  font-size: 0;
}
.demo-table-expand label {
  width: 90px;
  color: #99a9bf;
}
.demo-table-expand .el-form-item {
  margin-right: 0;
  margin-bottom: 0;
  width: 50%;
}
</style>