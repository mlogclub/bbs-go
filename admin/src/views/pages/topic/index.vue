<template>
  <section v-loading="listLoading" class="page-container">
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
          <el-select v-model="filters.recommend" clearable placeholder="是否推荐" @change="list">
            <el-option label="推荐" value="1" />
            <el-option label="未推荐" value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="filters.status" clearable placeholder="请选择状态" @change="list">
            <el-option label="正常" :value="0" />
            <el-option label="删除" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="list"> 查询 </el-button>
        </el-form-item>
      </el-form>
    </div>

    <div ref="mainContent" :style="{ height: mainHeight }" class="page-section">
      <topic-list :results="results" @change="list" />
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
  </section>
</template>

<script>
import mainHeight from "@/utils/mainHeight";
import TopicList from "./components/TopicList";

export default {
  name: "Topic",
  components: { TopicList },
  data() {
    return {
      mainHeight: "300px",
      results: [],
      listLoading: false,
      page: {},
      filters: {
        status: 0,
      },
      selectedRows: [],
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
        .form("/api/admin/topic/list", params)
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
    handleSelectionChange(val) {
      this.selectedRows = val;
    },
  },
};
</script>
<style scoped lang="scss">
.page-section {
  overflow-y: auto;
}
</style>
