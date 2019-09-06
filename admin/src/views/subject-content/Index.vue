
<template>
  <section>
    <el-col :span="24" class="toolbar" style="padding-bottom: 0px;">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-select v-model="filters.entityType" clearable placeholder="试题类型" @change="list">
            <el-option label="文章" value="article"></el-option>
            <el-option label="话题" value="topic"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-input v-model="filters.entityId" placeholder="实体编号"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" v-on:click="list">查询</el-button>
        </el-form-item>
        <!--
        <el-form-item>
          <el-button type="primary" @click="handleAdd">新增</el-button>
        </el-form-item>
        -->
      </el-form>
    </el-col>

    <el-table
      :data="results"
      highlight-current-row
      stripe
      v-loading="listLoading"
      style="width: 100%;"
      class="results"
      @selection-change="handleSelectionChange"
    >
      <!-- <el-table-column type="selection" width="55"></el-table-column> -->
      <el-table-column prop="id" label="编号"></el-table-column>
      <el-table-column prop="subject" label="专栏">
        <template slot-scope="scope">
          <span v-if="scope.row.subject">{{scope.row.subject.title}}</span>
        </template>
      </el-table-column>
      <el-table-column prop="title" label="标题">
        <template slot-scope="scope">
          <a v-if="scope.row.entityType === 'article'" :href="'https://mlog.club/article/' + scope.row.entityId" target="_blank">{{scope.row.title}}</a>
          <a v-if="scope.row.entityType === 'topic'" :href="'https://mlog.club/topic/' + scope.row.entityId" target="_blank">{{scope.row.title}}</a>
        </template>
      </el-table-column>
      <el-table-column prop="entityType" label="实体">
        <template slot-scope="scope">{{scope.row.entityType}} &nbsp;|&nbsp; {{scope.row.entityId}}</template>
      </el-table-column>
      <!-- <el-table-column prop="deleted" label="是否删除"></el-table-column> -->
      <el-table-column prop="createTime" label="创建时间">
        <template slot-scope="scope">{{scope.row.createTime | formatDate}}</template>
      </el-table-column>
      <!--
      <el-table-column label="操作" width="150">
        <template slot-scope="scope">
          <el-button size="small" @click="handleEdit(scope.$index, scope.row)">编辑</el-button>
        </template>
      </el-table-column>
      -->
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
      <el-form :model="addForm" label-width="80px" :rules="addFormRules" ref="addForm">
        <el-form-item label="专栏编号" prop="rule">
          <el-input v-model="addForm.subjectId"></el-input>
        </el-form-item>

        <el-form-item label="实体类型" prop="rule">
          <el-input v-model="addForm.entityType"></el-input>
        </el-form-item>

        <el-form-item label="实体编号" prop="rule">
          <el-input v-model="addForm.entityId"></el-input>
        </el-form-item>

        <el-form-item label="描述" prop="rule">
          <el-input v-model="addForm.summary"></el-input>
        </el-form-item>

        <el-form-item label="deleted" prop="rule">
          <el-input v-model="addForm.deleted"></el-input>
        </el-form-item>

        <el-form-item label="createTime" prop="rule">
          <el-input v-model="addForm.createTime"></el-input>
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

        <el-form-item label="subjectId" prop="rule">
          <el-input v-model="editForm.subjectId"></el-input>
        </el-form-item>

        <el-form-item label="entityType" prop="rule">
          <el-input v-model="editForm.entityType"></el-input>
        </el-form-item>

        <el-form-item label="entityId" prop="rule">
          <el-input v-model="editForm.entityId"></el-input>
        </el-form-item>

        <el-form-item label="summary" prop="rule">
          <el-input v-model="editForm.summary"></el-input>
        </el-form-item>

        <el-form-item label="deleted" prop="rule">
          <el-input v-model="editForm.deleted"></el-input>
        </el-form-item>

        <el-form-item label="createTime" prop="rule">
          <el-input v-model="editForm.createTime"></el-input>
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

      addForm: {
        subjectId: "",
        entityType: "",
        entityId: "",
        summary: "",
        deleted: "",
        createTime: ""
      },
      addFormVisible: false,
      addFormRules: {},
      addLoading: false,

      editForm: {
        id: "",
        subjectId: "",
        entityType: "",
        entityId: "",
        summary: "",
        deleted: "",
        createTime: ""
      },
      editFormVisible: false,
      editFormRules: {},
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
      HttpClient.post("/api/admin/subject-content/list", params)
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
      HttpClient.post("/api/admin/subject-content/create", this.addForm)
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
      HttpClient.get("/api/admin/subject-content/" + row.id)
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
      HttpClient.post("/api/admin/subject-content/update", me.editForm)
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
    }
  }
};
</script>

<style lang="scss" scoped>
.results {
  a, &:visited {
    text-decoration: none;
    color: #3273dc;
    cursor: pointer;
  }
}
</style>

