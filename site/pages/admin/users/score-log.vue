<template>
  <el-dialog
    v-if="isShowLog"
    :visible.sync="isShowLog"
    width="80%"
    title="积分记录"
  >
    <el-table
      v-loading="listLoading"
      :data="results"
      highlight-current-row
      border
    >
      <el-table-column prop="id" label="编号"></el-table-column>
      <!-- <el-table-column prop="userId" label="用户编号"></el-table-column> -->
      <el-table-column prop="sourceType" label="来源类型"></el-table-column>
      <el-table-column prop="sourceId" label="来源编号"></el-table-column>
      <el-table-column prop="description" label="描述"></el-table-column>
      <el-table-column prop="type" label="类型">
        <template slot-scope="scope">{{
          scope.row.type === 0 ? '增加' : '减少'
        }}</template>
      </el-table-column>
      <el-table-column prop="score" label="积分"></el-table-column>
      <el-table-column prop="createTime" label="创建时间">
        <template slot-scope="scope">{{
          scope.row.createTime | formatDate
        }}</template>
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
  </el-dialog>
</template>

<script>
export default {
  data() {
    return {
      isShowLog: false,
      userId: 0,
      results: [],
      listLoading: false,
      page: {},
      filters: {},
    }
  },
  mounted() {},
  methods: {
    async showLog(userId) {
      this.userId = userId
      this.isShowLog = true
      await this.list()
    },
    async list() {
      this.listLoading = true
      const params = Object.assign(this.filters, {
        page: this.page.page,
        limit: this.page.limit,
      })
      params.userId = this.userId
      try {
        const data = await this.$axios.post(
          '/api/admin/user-score-log/list',
          params
        )
        this.results = data.results
        this.page = data.page
      } catch (err) {
        this.$notify.error({ title: '错误', message: err.message || err })
      } finally {
        this.listLoading = false
      }
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
