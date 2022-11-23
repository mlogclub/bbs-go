<template>
  <el-dialog
    :visible.sync="showDialog"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    width="80%"
    title="查看评论"
  >
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

      <div ref="mainContent" class="main-content" :style="{ height: mainHeight }">
        <div v-if="results && results.length" class="page-section comments">
          <div v-for="item in results" :key="item.id" class="comment-item">
            <avatar :user="item.user" size="40" />
            <div class="comment-main">
              <div class="comment-nickname">
                <a :href="('/user/' + item.user.id) | siteUrl" target="_blank">{{
                  item.user.nickname
                }}</a>
              </div>
              <div class="comment-meta">
                <div>
                  ID: <span>{{ item.id }}</span>
                </div>
                <div class="create-time">
                  时间: <span>{{ item.createTime | formatDate }}</span>
                </div>
              </div>
              <div class="comment-summary" v-html="item.content" />
              <ul v-if="item.imageList && item.imageList.length" class="comment-image-list">
                <li v-for="(image, index) in item.imageList" :key="index">
                  <el-image
                    class="image-item"
                    lazy
                    :src="image.url"
                    fit="cover"
                    :preview-src-list="imagePreviewList(item.imageList)"
                  />
                </li>
              </ul>
              <div class="comment-actions">
                <el-link
                  v-if="item.status === 0"
                  class="action-item"
                  type="danger"
                  icon="el-icon-delete"
                  @click="handleDelete(item)"
                  >删除</el-link
                >
                <el-link
                  v-if="item.status === 1"
                  class="action-item"
                  type="danger"
                  icon="el-icon-delete"
                  disabled
                  >已删除</el-link
                >
              </div>
            </div>
          </div>
        </div>
        <el-empty v-else />
      </div>

      <div v-if="page.total > 0" ref="pagebar" class="pagebar">
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
  </el-dialog>
</template>

<script>
import Avatar from "@/components/Avatar";

export default {
  components: { Avatar },
  data() {
    return {
      showDialog: false,
      mainHeight: "500px",
      results: [],
      listLoading: false,
      page: {},
      filters: {
        status: 0,
      },
      selectedRows: [],
    };
  },
  mounted() {},
  methods: {
    async show(entityType, entityId) {
      this.showDialog = true;
      this.filters.entityType = entityType;
      this.filters.entityId = entityId;
      await this.list();
    },
    async list() {
      this.listLoading = true;
      const params = Object.assign(this.filters, {
        page: this.page.page,
        limit: this.page.limit,
      });
      try {
        let data = await this.axios.form("/api/admin/comment/list", params);
        data = data || {};
        this.results = data.results || [];
        this.page = data.page || {};
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
    handleSelectionChange(val) {
      this.selectedRows = val;
    },
    async handleDelete(row) {
      try {
        const flag = await this.$confirm("是否确认取消删除?", "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        });
        if (flag) {
          try {
            await this.axios.form(`/api/admin/comment/delete/${row.id}`);
            this.list();
            this.$message.success("删除成功");
          } catch (err) {
            this.$notify.error({ title: "错误", message: err.message || err });
          }
        }
      } catch (e) {
        this.$message.success("操作已取消");
      }
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
.main-content {
  overflow-y: auto;
}
.comments {
  width: 100%;
  overflow-y: auto;

  .comment-item {
    width: 100%;
    display: flex;
    padding: 10px;

    &:not(:last-child) {
      border-bottom: 1px solid #f2f2f2;
    }

    .comment-main {
      width: 100%;
      margin-left: 15px;

      .comment-nickname {
        color: #000;
        font-size: 14px;
        font-weight: 500;
      }

      .comment-meta {
        margin-top: 10px;
        display: flex;
        font-size: 12px;
        color: #555;

        & > div {
          margin-right: 10px;
        }
      }

      .comment-summary {
        margin-top: 10px;
        font-size: 15px;
        color: #555;
      }

      .comment-image-list {
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

      .comment-actions {
        margin-top: 10px;
      }
    }
  }
}
</style>
