<template>
  <section class="page-container">
    <div class="toolbar">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.id" placeholder="编号"></el-input>
        </el-form-item>
        <el-form-item>
          <el-input v-model="filters.username" placeholder="用户名"></el-input>
        </el-form-item>
        <el-form-item>
          <el-input v-model="filters.nickname" placeholder="昵称"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button @click="list" type="primary">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button @click="handleAdd" type="primary">新增</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
      v-loading="listLoading"
      :data="results"
      @selection-change="handleSelectionChange"
      highlight-current-row
      stripe
      style="width: 100%;"
    >
      <el-table-column type="expand">
        <template slot-scope="scope">
          <div
            v-if="scope.row.roles && scope.row.roles.length"
            class="content-form"
          >
            <div class="form-item">
              <div class="field-key">角色：</div>
              <div class="field-value">
                <el-tag
                  v-for="role in scope.row.roles"
                  :key="role"
                  size="mini"
                  style="margin-right:3px;"
                  >{{ role }}
                </el-tag>
              </div>
            </div>
          </div>
          <div class="content-form">
            <div class="form-item">
              <div class="field-key">状态：</div>
              <div class="field-value">
                {{ scope.row.status === 0 ? '正常' : '删除' }}
              </div>
            </div>
          </div>
          <div class="content-form">
            <div class="form-item">
              <div class="field-key">注册时间：</div>
              <div class="field-value">
                {{ scope.row.createTime | formatDate }}
              </div>
            </div>
          </div>
          <div class="content-form">
            <div class="form-item">
              <div class="field-key">更新时间：</div>
              <div class="field-value">
                {{ scope.row.updateTime | formatDate }}
              </div>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="id" label="编号" width="100"></el-table-column>
      <el-table-column prop="avatar" label="头像" width="80">
        <template slot-scope="scope">
          <img
            :src="scope.row.avatar"
            style="max-height: 50px; max-width: 50px; border-radius: 50%;"
          />
        </template>
      </el-table-column>
      <el-table-column prop="username" label="用户名"></el-table-column>
      <el-table-column prop="nickname" label="昵称"></el-table-column>
      <el-table-column prop="email" label="邮箱"></el-table-column>
      <el-table-column prop="score" label="积分"></el-table-column>
      <el-table-column prop="createTime" label="创建时间">
        <template slot-scope="scope">{{
          scope.row.createTime | formatDate
        }}</template>
      </el-table-column>
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
        :current-page="page.page"
        :page-size="page.limit"
        :total="page.total"
        @current-change="handlePageChange"
        @size-change="handleLimitChange"
        layout="total, sizes, prev, pager, next, jumper"
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
        <el-button
          :loading="addLoading"
          @click.native="addSubmit"
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
          <el-select
            v-model="editForm.roles"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="用户角色"
            style="width: 100%"
          >
            <el-option
              v-for="item in editForm.roles"
              :key="item"
              :label="item"
              :value="item"
            ></el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input
            v-model="editForm.password"
            placeholder="不填写表示不更改密码"
          ></el-input>
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
        <el-button
          :loading="editLoading"
          @click.native="editSubmit"
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
      filters: {
        id: ''
      },
      selectedRows: [],

      addForm: {
        username: '',
        nickname: '',
        avatar: '',
        email: '',
        roles: [],
        password: '',
        status: ''
      },
      addFormVisible: false,
      addFormRules: {},
      addLoading: false,

      editForm: {
        id: '',
        username: '',
        nickname: '',
        avatar: '',
        email: '',
        roles: [],
        password: '',
        status: ''
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
      const me = this
      me.listLoading = true
      const params = Object.assign(me.filters, {
        page: me.page.page,
        limit: me.page.limit
      })
      this.$axios
        .post('/api/admin/user/list', params)
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
        description: ''
      }
      this.addFormVisible = true
    },
    addSubmit() {
      const me = this
      this.$axios
        .post('/api/admin/user/create', this.addForm)
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
        .get(`/api/admin/user/${row.id}`)
        .then((data) => {
          me.editForm = Object.assign({}, data)
          me.editForm.password = ''
          me.editFormVisible = true
        })
        .catch((rsp) => {
          me.$notify.error({ title: '错误', message: rsp.message })
        })
    },
    editSubmit() {
      const params = { ...this.editForm }
      if (params.roles && params.roles.length) {
        params.roles = params.roles.join(',')
      } else {
        params.roles = ''
      }
      const me = this
      this.$axios
        .post('/api/admin/user/update', params)
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
