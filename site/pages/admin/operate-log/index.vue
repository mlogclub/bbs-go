<template>
  <section class="page-container">
    <div class="toolbar">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.userId" placeholder="用户编号"></el-input>
        </el-form-item>
        <el-form-item>
          <el-select
            v-model="filters.opType"
            clearable
            placeholder="操作类型"
            @change="list"
          >
            <el-option label="添加" value="create"></el-option>
            <el-option label="删除" value="delete"></el-option>
            <el-option label="修改" value="update"></el-option>
          </el-select>
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
      stripe
      style="width: 100%"
    >
      <el-table-column type="expand">
        <template slot-scope="scope">
          <div>{{ scope.row.ip }}</div>
          <div>{{ scope.row.userAgent }}</div>
          <div>{{ scope.row.referer }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="id" label="编号" width="100"></el-table-column>
      <el-table-column prop="userId" label="用户编号"></el-table-column>
      <el-table-column prop="opType" label="操作类型"></el-table-column>
      <el-table-column prop="dataType" label="数据类型"></el-table-column>
      <el-table-column prop="dataId" label="数据编号"></el-table-column>
      <el-table-column prop="createTime" label="操作时间">
        <template slot-scope="scope"
          >{{ scope.row.createTime | formatDate }}
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
      ></el-pagination>
    </div>
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
        .post('/api/admin/operate-log/list', params)
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
  },
}
</script>

<style scoped>
.link-logo {
  max-width: 50px;
  max-height: 50px;
}
</style>
