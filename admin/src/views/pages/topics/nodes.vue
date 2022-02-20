<template>
  <section class="page-container">
    <div ref="toolbar" class="toolbar">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.name" placeholder="名称" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="list"> 查询 </el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleAdd"> 新增 </el-button>
        </el-form-item>
      </el-form>
    </div>

    <div ref="mainContent" :style="{ height: mainHeight }">
      <el-table
        v-loading="listLoading"
        height="100%"
        :data="results"
        highlight-current-row
        stripe
        border
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="编号" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="logo" label="图标">
          <template slot-scope="scope">
            <img v-if="scope.row.logo" :src="scope.row.logo" class="node-logo" />
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="sortNo" label="排序" />
        <el-table-column prop="status" label="状态">
          <template slot-scope="scope">
            {{ scope.row.status === 0 ? "启用" : "禁用" }}
          </template>
        </el-table-column>

        <el-table-column prop="createTime" label="创建时间">
          <template slot-scope="scope">
            {{ scope.row.createTime | formatDate }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="150">
          <template slot-scope="scope">
            <el-button size="small" @click="handleEdit(scope.$index, scope.row)"> 编辑 </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div ref="pagebar" class="pagebar">
      <el-pagination
        :page-sizes="[20, 50, 100, 300]"
        :current-page="page.page"
        :page-size="page.limit"
        :total="page.total"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handlePageChange"
        @size-change="handleLimitChange"
      />
    </div>

    <el-dialog :visible.sync="addFormVisible" :close-on-click-modal="false" title="新增">
      <el-form ref="addForm" :model="addForm" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="addForm.name" />
        </el-form-item>
        <el-form-item label="图标">
          <upload v-model="addForm.logo" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="addForm.description" type="textarea" auto-complete="off" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input v-model="addForm.sortNo" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="addFormVisible = false"> 取消 </el-button>
        <el-button :loading="addLoading" type="primary" @click.native="addSubmit"> 提交 </el-button>
      </div>
    </el-dialog>

    <el-dialog
      :visible.sync="editFormVisible"
      :close-on-click-modal="false"
      :destroy-on-close="true"
      title="编辑"
    >
      <el-form ref="editForm" :model="editForm" label-width="80px">
        <el-input v-model="editForm.id" type="hidden" />
        <el-form-item label="名称">
          <el-input v-model="editForm.name" />
        </el-form-item>
        <el-form-item label="图标">
          <upload v-model="editForm.logo" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="editForm.description" type="textarea" auto-complete="off" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input v-model="editForm.sortNo" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="editForm.status" clearable placeholder="请选择状态">
            <el-option :value="0" label="启用" />
            <el-option :value="1" label="禁用" />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="editFormVisible = false"> 取消 </el-button>
        <el-button :loading="editLoading" type="primary" @click.native="editSubmit">
          提交
        </el-button>
      </div>
    </el-dialog>
  </section>
</template>

<script>
import Upload from "@/components/Upload";
import mainHeight from "@/utils/mainHeight";
export default {
  name: "Nodes",
  components: { Upload },
  data() {
    return {
      mainHeight: "300px",
      results: [],
      listLoading: false,
      page: {},
      filters: {},
      selectedRows: [],

      addForm: {
        name: "",
        logo: "",
        description: "",
        status: "",
        sortNo: "",
        createTime: "",
      },
      addFormVisible: false,
      addLoading: false,

      editForm: {
        id: "",
        name: "",
        logo: "",
        description: "",
        status: "",
        sortNo: "",
        createTime: "",
      },
      editFormVisible: false,
      editLoading: false,
    };
  },
  mounted() {
    mainHeight(this);
    this.list();
  },
  methods: {
    list() {
      const me = this;
      me.listLoading = true;
      const params = Object.assign(me.filters, {
        page: me.page.page,
        limit: me.page.limit,
      });
      this.axios
        .form("/api/admin/topic-node/list", params)
        .then((data) => {
          me.results = data.results;
          me.page = data.page;
        })
        .finally(() => {
          me.listLoading = false;
        });
    },
    handlePageChange(val) {
      this.page.page = val;
      this.list();
    },
    handleLimitChange(val) {
      this.page.limit = val;
      this.list();
    },
    handleAdd() {
      this.addForm = {
        name: "",
        description: "",
      };
      this.addFormVisible = true;
    },
    addSubmit() {
      const me = this;
      this.axios
        .form("/api/admin/topic-node/create", this.addForm)
        .then((data) => {
          me.$message({ message: "提交成功", type: "success" });
          me.addFormVisible = false;
          me.list();
        })
        .catch((rsp) => {
          me.$notify.error({ title: "错误", message: rsp.message });
        });
    },
    handleEdit(index, row) {
      const me = this;
      this.axios
        .get("/api/admin/topic-node/" + row.id)
        .then((data) => {
          me.editForm = Object.assign({}, data);
          me.editFormVisible = true;
        })
        .catch((rsp) => {
          me.$notify.error({ title: "错误", message: rsp.message });
        });
    },
    editSubmit() {
      const me = this;
      this.axios
        .form("/api/admin/topic-node/update", me.editForm)
        .then((data) => {
          me.list();
          me.editFormVisible = false;
        })
        .catch((rsp) => {
          me.$notify.error({ title: "错误", message: rsp.message });
        });
    },

    handleSelectionChange(val) {
      this.selectedRows = val;
    },
  },
};
</script>

<style lang="scss" scoped>
.node-logo {
  width: 80px;
  height: 80px;
  max-height: 80px;
  max-width: 80px;
}
</style>
