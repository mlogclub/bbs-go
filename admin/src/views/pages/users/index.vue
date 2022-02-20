<template>
  <section class="page-container">
    <div ref="toolbar" class="toolbar">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.id" placeholder="编号" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="filters.username" placeholder="用户名" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="filters.nickname" placeholder="昵称" />
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
        <el-table-column type="expand">
          <template slot-scope="scope">
            <div v-if="scope.row.username">
              <span>用户名：</span>
              {{ scope.row.username }}
            </div>
            <div v-if="scope.row.roles && scope.row.roles.length">
              <div>
                <span>角色：</span>
                <el-tag
                  v-for="role in scope.row.roles"
                  :key="role"
                  size="mini"
                  style="margin-right: 3px"
                >
                  {{ role }}
                </el-tag>
              </div>
            </div>
            <div>
              <span>状态：</span>
              {{ scope.row.status === 0 ? "正常" : "删除" }}
            </div>
            <div>
              <span>注册时间：</span>
              {{ scope.row.createTime | formatDate }}
            </div>
            <div>
              <span>更新时间：</span>
              {{ scope.row.updateTime | formatDate }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="id" label="编号" width="100" />
        <el-table-column prop="avatar" label="头像" width="80">
          <template slot-scope="scope">
            <avatar :user="scope.row" />
          </template>
        </el-table-column>
        <el-table-column prop="nickname" label="昵称" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="score" label="积分" />
        <el-table-column prop="forbidden" label="是否禁言">
          <template slot-scope="scope">
            <span v-if="scope.row.forbidden" class="tag is-warning">
              <template v-if="scope.row.forbiddenEndTime === -1">永久禁言</template>
              <template v-else>禁言至：{{ scope.row.forbiddenEndTime | formatDate }}</template>
            </span>
            <span v-else class="tag is-success">正常</span>
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="注册时间">
          <template slot-scope="scope">
            {{ scope.row.createTime | formatDate }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template slot-scope="scope">
            <el-dropdown size="mini" trigger="hover" placement="bottom" @command="handleCommand">
              <el-button type="primary">
                操作<i class="el-icon-arrow-down el-icon--right" />
              </el-button>
              <el-dropdown-menu slot="dropdown">
                <el-dropdown-item :command="{ cmd: 'edit', row: scope.row }">
                  编辑
                </el-dropdown-item>
                <el-dropdown-item
                  v-if="scope.row.forbidden"
                  :command="{ cmd: 'removeForbidden', row: scope.row }"
                >
                  取消禁言
                </el-dropdown-item>
                <el-dropdown-item v-else :command="{ cmd: 'forbidden', row: scope.row }">
                  禁言
                </el-dropdown-item>
                <el-dropdown-item :command="{ cmd: 'scoreLog', row: scope.row }">
                  积分记录
                </el-dropdown-item>
              </el-dropdown-menu>
            </el-dropdown>
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
      <el-form ref="addForm" :model="addForm" :rules="addFormRules" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="addForm.username" />
        </el-form-item>

        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="addForm.nickname" />
        </el-form-item>

        <el-form-item label="邮箱" prop="email">
          <el-input v-model="addForm.email" />
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input v-model="addForm.password" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="addFormVisible = false"> 取消 </el-button>
        <el-button :loading="addLoading" type="primary" @click.native="addSubmit"> 提交 </el-button>
      </div>
    </el-dialog>

    <el-dialog :visible.sync="editFormVisible" :close-on-click-modal="false" title="编辑">
      <el-form ref="editForm" :model="editForm" :rules="editFormRules" label-width="80px">
        <el-input v-model="editForm.id" type="hidden" />
        <el-form-item label="用户名" prop="username">
          <el-input v-model="editForm.username" />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="editForm.nickname" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="editForm.email" />
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
            <el-option v-for="item in editForm.roles" :key="item" :label="item" :value="item" />
          </el-select>
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input v-model="editForm.password" placeholder="不填写表示不更改密码" />
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-select v-model="editForm.status" placeholder="请选择">
            <el-option :key="0" :value="0" label="正常" />
            <el-option :key="1" :value="1" label="删除" />
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

    <el-dialog :visible.sync="forbiddenFormVisible" :close-on-click-modal="false" title="禁言">
      <el-form ref="forbiddenForm" :model="forbiddenForm" label-width="80px">
        <el-input v-model="forbiddenForm.userId" type="hidden" />
        <el-form-item label="禁言时间" prop="reason">
          <el-select v-model="forbiddenForm.days">
            <el-option label="3天" value="3" />
            <el-option label="5天" value="3" />
            <el-option label="7天" value="7" />
            <el-option label="15天" value="15" />
            <el-option label="30天" value="30" />
            <el-option label="永久" value="-1" />
          </el-select>
        </el-form-item>
        <el-form-item label="禁言原因" prop="reason">
          <el-select v-model="forbiddenForm.reason">
            <el-option value="广告" />
            <el-option value="灌水" />
            <el-option value="涉黄" />
            <el-option value="涉政" />
            <el-option value="其他" />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="forbiddenFormVisible = false"> 取消 </el-button>
        <el-button :loading="forbiddenLoading" type="primary" @click.native="forbidden">
          禁言
        </el-button>
      </div>
    </el-dialog>

    <score-log ref="scoreLog" />
  </section>
</template>

<script>
import ScoreLog from "./score-log";
import Avatar from "@/components/Avatar";
import mainHeight from "@/utils/mainHeight";

export default {
  name: "Users",
  components: { ScoreLog, Avatar },
  data() {
    return {
      mainHeight: "300px",
      results: [],
      listLoading: false,
      page: {},
      filters: {
        id: "",
      },
      selectedRows: [],

      addForm: {
        username: "",
        nickname: "",
        avatar: "",
        email: "",
        roles: [],
        password: "",
        status: "",
      },
      addFormVisible: false,
      addFormRules: {},
      addLoading: false,

      editForm: {
        id: "",
        username: "",
        nickname: "",
        avatar: "",
        email: "",
        roles: [],
        password: "",
        status: "",
      },
      editFormVisible: false,
      editFormRules: {},
      editLoading: false,

      forbiddenForm: {
        userId: "",
        days: 0,
        reason: "",
      },
      forbiddenFormVisible: false,
      forbiddenLoading: false,
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
        .form("/api/admin/user/list", params)
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
        .form("/api/admin/user/create", this.addForm)
        .then((data) => {
          me.$message({ message: "提交成功", type: "success" });
          me.addFormVisible = false;
          me.list();
        })
        .catch((rsp) => {
          me.$notify.error({ title: "错误", message: rsp.message });
        });
    },
    handleEdit(row) {
      const me = this;
      this.axios
        .get(`/api/admin/user/${row.id}`)
        .then((data) => {
          me.editForm = Object.assign({}, data);
          me.editForm.password = "";
          me.editFormVisible = true;
        })
        .catch((rsp) => {
          me.$notify.error({ title: "错误", message: rsp.message });
        });
    },
    editSubmit() {
      const params = { ...this.editForm };
      if (params.roles && params.roles.length) {
        params.roles = params.roles.join(",");
      } else {
        params.roles = "";
      }
      const me = this;
      this.axios
        .form("/api/admin/user/update", params)
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
    showForbiddenDialog(row) {
      this.forbiddenForm = {
        userId: row.id,
        days: 7,
        reason: "广告",
      };
      this.forbiddenFormVisible = true;
    },
    async forbidden() {
      this.forbiddenLoading = true;
      try {
        await this.axios.form("/api/admin/user/forbidden", this.forbiddenForm);
        this.forbiddenForm = {};
        this.forbiddenFormVisible = false;
        this.$message.success("禁言成功");
        this.list();
      } catch (e) {
        this.$message.success("禁言失败 " + (e.message || e));
      } finally {
        this.forbiddenLoading = false;
      }
    },
    async removeForbidden(row) {
      try {
        await this.axios.form("/api/admin/user/forbidden", {
          userId: row.id,
          days: 0,
        });
        this.$message.success("取消禁言成功");
        this.list();
      } catch (e) {
        this.$message.success("取消禁言失败 " + (e.message || e));
      }
    },
    showScoreLog(row) {
      this.$refs.scoreLog.showLog(row.id);
    },
    handleCommand(cmd) {
      if (cmd.cmd === "edit") {
        this.handleEdit(cmd.row);
      } else if (cmd.cmd === "removeForbidden") {
        this.removeForbidden(cmd.row);
      } else if (cmd.cmd === "forbidden") {
        this.showForbiddenDialog(cmd.row);
      } else if (cmd.cmd === "scoreLog") {
        this.showScoreLog(cmd.row);
      }
    },
  },
};
</script>

<style scoped></style>
