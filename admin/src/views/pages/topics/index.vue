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
            <el-option label="正常" value="0" />
            <el-option label="删除" value="1" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="list"> 查询 </el-button>
        </el-form-item>
      </el-form>
    </div>

    <div ref="mainContent" :style="{ height: mainHeight }" class="page-section topics">
      <div v-for="topic in results" :key="topic.topicId" class="topic-item">
        <div class="topic-left">
          <avatar :user="topic.user" size="30" />
        </div>
        <div class="topic-main">
          <div class="nickname">{{ topic.user.nickname }}</div>
          <div class="create-time">{{ topic.createTime | formatDate }}</div>
          <div v-if="topic.type === 0 && topic.summary" class="summary">
            {{ topic.summary }}
          </div>
          <div v-if="topic.type === 1 && topic.content" class="summary">
            {{ topic.content }}
          </div>
          <ul v-if="topic.imageList && topic.imageList.length" class="image-list">
            <li v-for="(image, index) in topic.imageList" :key="index">
              <el-image
                class="image-item"
                lazy
                :src="image.url"
                fit="cover"
                :preview-src-list="imagePreviewList(topic.imageList)"
              />
            </li>
          </ul>
          <div class="actions">
            <el-link class="action-item">默认链接</el-link>
            <el-link class="action-item" type="primary">主要链接</el-link>
            <el-link class="action-item" type="success">成功链接</el-link>
            <el-link class="action-item" type="warning">警告链接</el-link>
            <el-link class="action-item" type="danger">危险链接</el-link>
            <el-link class="action-item" type="info">信息链接</el-link>
          </div>
        </div>
      </div>
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
import Avatar from "@/components/Avatar";
import mainHeight from "@/utils/mainHeight";

export default {
  name: "Topics",
  components: { Avatar },
  data() {
    return {
      mainHeight: "300px",
      results: [],
      listLoading: false,
      page: {},
      filters: {
        status: "0",
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
    deleteSubmit(topicId) {
      const me = this;
      this.$confirm("是否确认删除该话题?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(function () {
          me.axios
            .form("/api/admin/topic/delete", { id: topicId })
            .then(function () {
              me.$message({ message: "删除成功", type: "success" });
              me.list();
            })
            .catch(function (err) {
              me.$notify.error({ title: "错误", message: err.message || err });
            });
        })
        .catch(function () {
          me.$message({
            type: "info",
            message: "已取消删除",
          });
        });
    },
    async undeleteSubmit(topicId) {
      try {
        await this.axios.form("/api/admin/topic/undelete", { id: topicId });
        this.list();
        this.$message({ message: "取消删除成功", type: "success" });
      } catch (err) {
        this.$notify.error({ title: "错误", message: err.message || err });
      }
    },
    async recommend(id) {
      try {
        await this.axios.form("/api/admin/topic/recommend", {
          id,
        });
        this.$message({ message: "推荐成功", type: "success" });
        this.list();
      } catch (e) {
        this.$notify.error({ title: "错误", message: e.message });
      }
    },
    async cancelRecommend(id) {
      try {
        await this.axios.delete("/api/admin/topic/recommend", {
          params: {
            id,
          },
        });
        this.$message({ message: "取消推荐成功", type: "success" });
        this.list();
      } catch (e) {
        this.$notify.error({ title: "错误", message: e.message });
      }
    },
    handleSelectionChange(val) {
      this.selectedRows = val;
    },
    imagePreviewList(imageList) {
      var ret = [];
      for (let i = 0; i < imageList.length; i++) {
        const ele = imageList[i];
        ret.push(ele.url);
      }
      return ret;
    },
  },
};
</script>

<style scoped lang="scss">
.topics {
  width: 100%;
  padding: 0;
  margin: 0;
  overflow-y: auto;

  .topic-item {
    display: flex;
    padding: 20px 10px 10px 10px;
    border-bottom: 1px solid #e9e9e9;
    .topic-left {
    }
    .topic-main {
      margin-left: 10px;
      .nickname {
        font-size: 13px;
        font-weight: 500;
        color: #111827;
      }
      .create-time {
        margin-top: 3px;
        font-size: 12px;
        color: #6b7280;
      }
      .summary {
        margin-top: 10px;
        font-size: 14px;
      }
      .image-list {
        display: flex;
        list-style: none;
        padding: 0;
        margin: 10px 0 0 0;

        .image-item {
          width: 150px;
          height: 150px;
          margin: 0 10px 10px 0;
        }
      }

      .actions {
        margin-top: 20px;
        .action-item {
          margin-right: 10px;
        }
      }
    }
  }
}
</style>
