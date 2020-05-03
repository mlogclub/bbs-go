<template>
  <section v-loading="listLoading" class="page-container">
    <div class="toolbar">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.id" placeholder="编号"></el-input>
        </el-form-item>
        <el-form-item>
          <el-input v-model="filters.userId" placeholder="用户编号"></el-input>
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
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button @click="list" type="primary">查询</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="main-content">
      <div v-for="item in results" :key="item.id">
        {{ item.content }}
      </div>
    </div>

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
        status: '0'
      },
      selectedRows: []
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
        .post('/api/admin/tweet/list', params)
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
    async deleteSubmit(row) {
      await this.$confirm('是否确认删除该话题?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })

      try {
        await this.$axios.post('/api/admin/tweet/delete', { id: row.id })
        this.$message({ message: '删除成功', type: 'success' })
        this.list()
      } catch (err) {
        this.$notify.error({ title: '错误', message: err.message || err })
      }
    },
    handleSelectionChange(val) {
      this.selectedRows = val
    }
  }
}
</script>

<style scoped lang="scss"></style>
