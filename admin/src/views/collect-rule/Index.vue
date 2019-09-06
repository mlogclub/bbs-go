
<template>
  <section>
    <el-col :span="24" class="toolbar" style="padding-bottom: 0px;">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.title" placeholder="名称"></el-input>
        </el-form-item>
        <el-form-item>
          <el-select v-model="filters.status" clearable placeholder="请选择状态" @change="list">
            <el-option label="启用" value="0"></el-option>
            <el-option label="禁用" value="1"></el-option>
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
      <el-table-column type="expand">
        <template slot-scope="scope">
          <p>{{ scope.row.description }}</p>
          <pre>{{ scope.row.rule }}</pre>
        </template>
      </el-table-column>
      <el-table-column prop="id" label="编号"></el-table-column>
      <el-table-column prop="title" label="名称"></el-table-column>
      <!-- <el-table-column prop="description" label="描述"></el-table-column> -->
      <el-table-column prop="status" label="状态">
        <template slot-scope="scope">{{scope.row.status === 0 ? '启用' : '禁用'}}</template>
      </el-table-column>
      <el-table-column prop="createTime" label="创建时间">
        <template slot-scope="scope">{{scope.row.createTime | formatDate}}</template>
      </el-table-column>
      <el-table-column prop="updateTime" label="更新时间">
        <template slot-scope="scope">{{scope.row.updateTime | formatDate}}</template>
      </el-table-column>
      <el-table-column label="操作" width="150">
        <template slot-scope="scope">
          <el-button size="small" @click="handleEdit(scope.$index, scope.row)">编辑</el-button>
          <el-button size="small" type="primary" @click="handleRun(scope.row.id)">启动</el-button>
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

    <el-dialog title="新增" :visible.sync="addFormVisible" :close-on-click-modal="false">
      <el-form :model="addForm" label-width="80px" ref="addForm">
        <el-form-item label="名称">
          <el-input v-model="addForm.title"></el-input>
        </el-form-item>

        <el-form-item label="描述">
          <el-input v-model="addForm.description" type="textarea" autosize></el-input>
        </el-form-item>

        <el-form-item label="规则">
          <el-input v-model="addForm.rule" type="textarea" autosize></el-input>
        </el-form-item>

        <!--
        <el-form-item label="状态">
          <el-select v-model="addForm.status" clearable placeholder="请选择状态">
            <el-option label="启用" :value="0"></el-option>
            <el-option label="禁用" :value="1"></el-option>
          </el-select>
        </el-form-item>
        -->
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="addFormVisible = false">取消</el-button>
        <el-button type="primary" @click.native="addSubmit" :loading="addLoading">提交</el-button>
      </div>
    </el-dialog>

    <el-dialog title="编辑" :visible.sync="editFormVisible" :close-on-click-modal="false">
      <el-form :model="editForm" label-width="80px" ref="editForm">
        <el-input v-model="editForm.id" type="hidden"></el-input>

        <el-form-item label="名称">
          <el-input v-model="editForm.title"></el-input>
        </el-form-item>

        <el-form-item label="描述">
          <el-input v-model="editForm.description" type="textarea" autosize></el-input>
        </el-form-item>

        <el-form-item label="规则">
          <el-input v-model="editForm.rule" type="textarea" autosize></el-input>
        </el-form-item>

        <el-form-item label="状态">
          <el-select v-model="editForm.status" clearable placeholder="请选择状态">
            <el-option label="启用" :value="0"></el-option>
            <el-option label="禁用" :value="1"></el-option>
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
import HttpClient from "../../apis/HttpClient";

export default {
  name: "List",
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
      editLoading: false
    };
  },
  mounted() {
    this.list();
  },
  methods: {
    list() {
      let me = this;
      me.listLoading = true;
      let params = Object.assign(me.filters, {
        page: me.page.page,
        limit: me.page.limit
      });
      HttpClient.post("/api/admin/collect-rule/list", params)
        .then(data => {
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
        description: ""
      };
      this.addFormVisible = true;
    },
    addSubmit() {
      let me = this;
      HttpClient.post("/api/admin/collect-rule/create", this.addForm)
        .then(data => {
          me.$message({ message: "提交成功", type: "success" });
          me.addFormVisible = false;
          me.list();
        })
        .catch(rsp => {
          me.$notify.error({ title: "错误", message: rsp.message });
        });
    },
    handleEdit(index, row) {
      let me = this;
      HttpClient.get("/api/admin/collect-rule/" + row.id)
        .then(data => {
          me.editForm = Object.assign({}, data);
          me.editFormVisible = true;
        })
        .catch(rsp => {
          me.$notify.error({ title: "错误", message: rsp.message });
        });
    },
    editSubmit() {
      let me = this;
      HttpClient.post("/api/admin/collect-rule/update", me.editForm)
        .then(data => {
          me.list();
          me.editFormVisible = false;
        })
        .catch(rsp => {
          me.$notify.error({ title: "错误", message: rsp.message });
        });
    },

    handleSelectionChange(val) {
      this.selectedRows = val;
    },

    async handleRun(id) {
      try {
        await HttpClient.get("/api/admin/collect-rule/run", {
          id: id
        });
        this.$message({ message: "启动成功", type: "success" });
      } catch (e) {
          this.$notify.error({ title: "错误", message: e.message || e });
      }
    }
  }
};
</script>

<style lang="scss" scoped>
pre {
  background: #23241f;
  margin: 0;
  padding: 2px 10px;
  overflow: auto;
  font-size: 13px;
  // color: #4d4d4c;
  color: green;
  line-height: 1.5;
  overflow-x: none;
}
</style>

