<template>
  <section class="page-container">
    <div class="toolbar">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.userId" placeholder="用户编号"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="list">查询</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
      v-loading="listLoading"
      :data="results"
      highlight-current-row
      border
      style="width: 100%;"
    >
      <el-table-column prop="id" label="编号"></el-table-column>
      <el-table-column prop="userId" label="用户">
        <template slot-scope="scope">
          <user-info :user="scope.row.user" />
        </template>
      </el-table-column>
      <el-table-column prop="score" label="积分"></el-table-column>
      <el-table-column prop="createTime" label="创建时间">
        <template slot-scope="scope">{{
          scope.row.createTime | formatDate
        }}</template>
      </el-table-column>
      <el-table-column prop="updateTime" label="更新时间">
        <template slot-scope="scope">{{
          scope.row.updateTime | formatDate
        }}</template>
      </el-table-column>
      <el-table-column label="操作" width="150">
        <template slot-scope="scope">
          <el-button
            type="success"
            size="small"
            @click="showLog(scope.$index, scope.row)"
            >积分记录</el-button
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

    <score-log ref="scoreLog" />
  </section>
</template>

<script>
import ScoreLog from './score-log'
import UserInfo from '~/pages/admin/components/UserInfo'
export default {
  layout: 'admin',
  components: { ScoreLog, UserInfo },
  data() {
    return {
      results: [],
      scoreLogs: [],
      listLoading: false,
      page: {},
      filters: {},
      isShowLog: false,
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
        .post('/api/admin/user-score/list', params)
        .then((data) => {
          me.results = data.results
          me.page = data.page
        })
        .finally(() => {
          me.listLoading = false
        })
    },
    showLog(index, row) {
      this.$refs.scoreLog.showLog(row.userId)
    },
    handlePageChange(val) {
      this.page.page = val
      this.list()
    },
    handleLimitChange(val) {
      this.page.limit = val
      this.list()
    },
  },
}
</script>

<style lang="scss" scoped></style>
