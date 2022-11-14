<template>
  <section class="page-container">
    <div ref="toolbar" class="toolbar">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-select v-model="filters.type" placeholder="类型" clearable @change="list">
            <el-option label="词组" value="word" />
            <el-option label="正则表达式" value="regex" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-input v-model="filters.word" placeholder="违禁词"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="el-icon-search" @click="list">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="el-icon-plus" @click="handleAdd">新增</el-button>
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
        <!-- <el-table-column type="selection" width="55"></el-table-column> -->
        <el-table-column prop="id" label="编号"></el-table-column>

        <el-table-column prop="type" label="类型"></el-table-column>

        <el-table-column prop="word" label="违禁词"></el-table-column>

        <el-table-column prop="remark" label="备注"></el-table-column>

        <el-table-column prop="createTime" label="创建时间">
          <template slot-scope="scope">
            {{ scope.row.createTime | formatDate }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="180">
          <template slot-scope="scope">
            <el-button
              size="small"
              type="primary"
              icon="el-icon-edit"
              @click="handleEdit(scope.$index, scope.row)"
              >编辑</el-button
            >
            <el-button
              size="small"
              type="danger"
              icon="el-icon-delete"
              @click="handleDelete(scope.$index, scope.row)"
              >删除</el-button
            >
          </template>
        </el-table-column>
        <template #empty>
          <el-empty />
        </template>
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
      >
      </el-pagination>
    </div>

    <el-dialog title="新增" :visible.sync="addFormVisible" :close-on-click-modal="false">
      <el-form ref="addForm" :model="addForm" label-width="80px">
        <el-form-item label="类型">
          <el-select v-model="addForm.type" placeholder="类型">
            <el-option label="词组" value="word" />
            <el-option label="正则表达式" value="regex" />
          </el-select>
        </el-form-item>

        <el-form-item label="违禁词">
          <el-input v-model="addForm.word"></el-input>
        </el-form-item>

        <el-form-item label="备注">
          <el-input v-model="addForm.remark"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="addFormVisible = false">取消</el-button>
        <el-button type="primary" :loading="addLoading" @click.native="addSubmit">提交</el-button>
      </div>
    </el-dialog>

    <el-dialog title="编辑" :visible.sync="editFormVisible" :close-on-click-modal="false">
      <el-form ref="editForm" :model="editForm" label-width="80px">
        <el-input v-model="editForm.id" type="hidden"></el-input>

        <el-form-item label="类型">
          <el-select v-model="editForm.type" placeholder="类型">
            <el-option label="词组" value="word" />
            <el-option label="正则表达式" value="regex" />
          </el-select>
        </el-form-item>

        <el-form-item label="违禁词">
          <el-input v-model="editForm.word"></el-input>
        </el-form-item>

        <el-form-item label="备注">
          <el-input v-model="editForm.remark"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="editFormVisible = false">取消</el-button>
        <el-button type="primary" :loading="editLoading" @click.native="editSubmit">提交</el-button>
      </div>
    </el-dialog>
  </section>
</template>

<script>
import mainHeight from "@/utils/mainHeight";
export default {
  data() {
    return {
      mainHeight: "300px",
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
    };
  },
  mounted() {
    mainHeight(this);
    this.list();
  },
  methods: {
    async list() {
      const params = Object.assign(this.filters, {
        page: this.page.page,
        limit: this.page.limit,
      });
      try {
        const data = await this.axios.form("/api/admin/forbidden-word/list", params);
        this.results = data.results;
        this.page = data.page;
      } catch (e) {
        this.$notify.error({ title: "错误", message: e.message || e });
      } finally {
        this.listLoading = false;
      }
    },
    async handlePageChange(val) {
      this.page.page = val;
      await this.list();
    },
    async handleLimitChange(val) {
      this.page.limit = val;
      await this.list();
    },
    handleAdd() {
      this.addForm = {
        type: "word",
      };
      this.addFormVisible = true;
    },
    async addSubmit() {
      try {
        await this.axios.form("/api/admin/forbidden-word/create", this.addForm);
        this.$message({ message: "提交成功", type: "success" });
        this.addFormVisible = false;
        await this.list();
      } catch (e) {
        this.$notify.error({ title: "错误", message: e.message || e });
      }
    },
    async handleEdit(index, row) {
      try {
        const data = await this.axios.get("/api/admin/forbidden-word/" + row.id);
        this.editForm = Object.assign({}, data);
        this.editFormVisible = true;
      } catch (e) {
        this.$notify.error({ title: "错误", message: e.message || e });
      }
    },
    async editSubmit() {
      try {
        await this.axios.form("/api/admin/forbidden-word/update", this.editForm);
        await this.list();
        this.editFormVisible = false;
      } catch (e) {
        this.$notify.error({ title: "错误", message: e.message || e });
      }
    },
    async handleDelete(index, row) {
      try {
        await this.axios.form("/api/admin/forbidden-word/delete", {
          id: row.id,
        });
        await this.list();
        this.$notify.success("删除成功");
      } catch (e) {
        this.$notify.error({ title: "错误", message: e.message || e });
      }
    },
    handleSelectionChange(val) {
      this.selectedRows = val;
    },
  },
};
</script>

<style lang="scss" scoped></style>
