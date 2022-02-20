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
      <div v-for="item in results" :key="item.topicId" class="topic-item">
        <div class="topic-avatar">
          <avatar :user="item.user" />
        </div>
        <div class="topic-main">
          <div class="topic-header">
            <a :href="('/user/' + item.user.id) | siteUrl" target="_blank" class="topic-nickname">{{
              item.user.nickname
            }}</a>

            <div class="topic-info">
              <span v-if="item.status === 1" style="color: red; font-weight: bold">已删除</span>
              <span v-if="item.recommend" style="color: red; font-weight: bold">已推荐</span>
            </div>
          </div>

          <div class="topic-metadata">
            <a
              class="topic-metadata-item topic-title-link"
              :href="('/topic/' + item.topicId) | siteUrl"
              target="_blank"
            >
              ID: {{ item.topicId }}
            </a>
            <span class="topic-metadata-item">发布于：{{ item.createTime | formatDate }}</span>
            <span class="topic-metadata-item">查看：{{ item.viewCount }}</span>
            <span class="topic-metadata-item">点赞：{{ item.likeCount }}</span>
            <span class="topic-metadata-item">评论：{{ item.commentCount }}</span>
            <span v-if="item.node" class="topic-metadata-item node">{{ item.node.name }}</span>
            <template v-if="item.tags && item.tags.length">
              <span v-for="tag in item.tags" :key="tag.tagId" class="topic-metadata-item tag"
                >#{{ tag.tagName }}</span
              >
            </template>
          </div>

          <div class="topic-title">
            <a :href="('/topic/' + item.topicId) | siteUrl" target="_blank">{{ item.title }}</a>
          </div>

          <template v-if="item.type === 0">
            <div class="topic-summary">
              {{ item.summary }}
            </div>
          </template>
          <template v-else>
            <div class="topic-summary">
              {{ item.content }}
            </div>
          </template>

          <ul v-if="item.imageList && item.imageList.length" class="topic-image-list">
            <li v-for="(image, index) in item.imageList" :key="index">
              <a :href="('/topic/' + item.topicId) | siteUrl" target="_blank" class="image-item">
                <img v-lazy="image.preview" />
              </a>
            </li>
          </ul>

          <div class="topic-actions">
            <a v-if="item.status === 0" class="action-item btn" @click="deleteSubmit(item.topicId)"
              >删除</a
            >
            <a v-else class="action-item btn" @click="undeleteSubmit(item.topicId)">取消删除</a>

            <a v-if="!item.recommend" class="action-item btn" @click="recommend(item.topicId)"
              >推荐</a
            >
            <a v-else class="action-item btn" @click="cancelRecommend(item.topicId)">取消推荐</a>
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
  },
};
</script>

<style scoped lang="scss">
.topics {
  width: 100%;
  padding: 0;
  margin: 0;
  overflow-y: auto;

  .topic-item:not(:last-child) {
    border-bottom: solid 1px rgba(140, 147, 157, 0.14);
  }

  .topic-item {
    width: 100%;
    padding: 20px 10px;
    display: flex;
    flex: 1;

    .topic-main {
      width: 100%;
      margin-left: 10px;
      font-size: 14px;
      a {
        font-size: 14px;
      }

      .topic-header {
        .topic-nickname {
          font-size: 16px;
          font-weight: 600;
          color: #606266;
        }

        .topic-info {
          margin-left: 10px;
          float: right;
          cursor: pointer;
        }
      }

      .topic-metadata {
        margin: 12px 0;
        color: #8590a6;

        .topic-metadata-item {
          margin-right: 12px;

          &.topic-title-link {
            color: #1482f0;
          }
          &.node {
          }
          &.tag {
          }
        }
      }

      .topic-title {
        margin: 12px 0;
        a {
          font-size: 16px;
          font-weight: 600;
          color: #000;
        }
      }

      .topic-summary {
        margin: 12px 0;
        color: #525252;
      }

      .topic-image-list {
        margin: 12px 0;
        padding: 0;
        li {
          cursor: pointer;
          border: 1px dashed #ddd;
          text-align: center;

          // 图片尺寸
          $image-size: 120px;

          display: inline-block;
          vertical-align: middle;
          width: $image-size;
          height: $image-size;
          line-height: $image-size;
          margin: 0 8px 8px 0;
          background-color: #e8e8e8;
          background-size: 32px 32px;
          background-position: 50%;
          background-repeat: no-repeat;
          overflow: hidden;
          position: relative;

          .image-item {
            display: block;
            width: $image-size;
            height: $image-size;
            overflow: hidden;
            transform-style: preserve-3d;

            & > img {
              width: 100%;
              height: 100%;
              object-fit: cover;
              transition: all 0.5s ease-out 0.1s;

              &:hover {
                transform: matrix(1.04, 0, 0, 1.04, 0, 0);
                backface-visibility: hidden;
              }
            }
          }
        }
      }

      .topic-actions {
        margin-top: 20px;
        .action-item {
          margin-right: 20px;
          color: #1482f0;
          font-weight: 500;
        }
      }
    }
  }
}
</style>
