<template>
  <section v-loading="listLoading" class="page-container">
    <!--工具条-->
    <div ref="toolbar" class="toolbar">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.id" placeholder="编号" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="filters.userId" placeholder="用户编号" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="filters.title" placeholder="标题" />
        </el-form-item>
        <el-form-item>
          <el-select v-model="filters.status" clearable placeholder="请选择状态" @change="list">
            <el-option label="正常" :value="0" />
            <el-option label="已删除" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="list"> 查询 </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!--列表-->
    <div ref="mainContent" :style="{ height: mainHeight }" class="page-section">
      <article-list :results="results" @change="list" />
    </div>

    <!--工具条-->
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
  </section>
</template>

<script>
import mainHeight from "@/utils/mainHeight";
import ArticleList from "./components/ArticleList";

export default {
  name: "Article",
  components: { ArticleList },
  data() {
    return {
      mainHeight: "300px",
      results: [],
      listLoading: false,
      page: {},
      filters: {
        status: 0,
      },
    };
  },
  mounted() {
    mainHeight(this);
    this.list();
  },
  methods: {
    async list() {
      this.listLoading = true;
      const params = Object.assign(this.filters, {
        page: this.page.page,
        limit: this.page.limit,
      });
      try {
        const data = await this.axios.form("/api/admin/article/list", params);
        this.results = data.results;
        this.page = data.page;
      } catch (err) {
        // this.$message.error(err.message);
      } finally {
        this.listLoading = false;
      }
    },
    handlePageChange(val) {
      this.page.page = val;
      this.list();
    },
    handleLimitChange(val) {
      this.page.limit = val;
      this.list();
    },
  },
};
</script>

<style scoped lang="scss">
.page-section {
  overflow-y: auto;
}
</style>
