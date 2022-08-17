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
          <avatar :user="topic.user" size="40" />
        </div>
        <div class="topic-main">
          <div class="topic-status">
            <el-tag v-if="topic.recommend" type="success">已推荐</el-tag>
            <el-tag v-if="topic.status === 1" type="danger">已删除</el-tag>
          </div>
          <div class="topic-nickname">{{ topic.user.nickname }}</div>
          <div class="topic-mates">
            <span> ID: {{ topic.topicId }} </span>
            <span> 时间: {{ topic.createTime | formatDate }} </span>
            <span> 查看: {{ topic.viewCount }} </span>
            <span> 点赞: {{ topic.likeCount }} </span>
            <span> 评论: {{ topic.commentCount }} </span>
          </div>
          <div v-if="topic.type === 0 && topic.summary" class="topic-summary">
            {{ topic.summary }}
          </div>
          <div v-if="topic.type === 1 && topic.content" class="topic-summary">
            {{ topic.content }}
          </div>
          <ul v-if="topic.imageList && topic.imageList.length" class="topic-image-list">
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
          <div class="topic-tags">
            <el-tag type="success" size="mini">{{ topic.node.name }}</el-tag>
            <template v-if="topic.tags && topic.tags.length">
              <el-tag v-for="tag in topic.tags" :key="tag.tagId" type="info" size="mini"
                >#&nbsp;{{ tag.tagName }}</el-tag
              >
            </template>
          </div>
          <div class="actions">
            <template v-if="topic.status === 0">
              <el-link
                class="action-item"
                icon="el-icon-view"
                :href="('/topic/' + topic.topicId) | siteUrl"
                target="_blank"
                >查看详情</el-link
              >
              <el-link
                class="action-item"
                icon="el-icon-s-comment"
                @click="showComments(topic.topicId)"
                >查看评论</el-link
              >
              <el-link
                v-if="topic.recommend"
                class="action-item"
                icon="el-icon-s-flag"
                @click="cancelRecommend(topic.topicId)"
                >取消推荐</el-link
              >
              <el-link
                v-else
                class="action-item"
                icon="el-icon-s-flag"
                @click="recommend(topic.topicId)"
                >推荐</el-link
              >
              <el-link
                class="action-item"
                type="danger"
                icon="el-icon-delete"
                @click="deleteSubmit(topic.topicId)"
                >删除</el-link
              >
            </template>
            <template v-else-if="topic.status === 1">
              <el-link
                class="action-item"
                type="info"
                icon="el-icon-delete"
                @click="undeleteSubmit(topic.topicId)"
                >取消删除</el-link
              >
            </template>
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

    <comments-dialog ref="commentsDialog" />
  </section>
</template>

<script>
import Avatar from "@/components/Avatar";
import mainHeight from "@/utils/mainHeight";
import CommentsDialog from "../comments/CommentsDialog";

export default {
  name: "Topics",
  components: { Avatar, CommentsDialog },
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
    showComments(topicId) {
      this.$refs.commentsDialog.show("topic", topicId);
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
        const flag = await this.$confirm("是否确认取消删除?", "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        });
        if (flag) {
          try {
            await this.axios.form("/api/admin/topic/undelete", { id: topicId });
            this.list();
            this.$message.success("操作成功");
          } catch (err) {
            this.$notify.error({ title: "错误", message: err.message || err });
          }
        }
      } catch (e) {
        this.$message.success("操作已取消");
      }
    },
    async recommend(id) {
      try {
        const flag = await this.$confirm("是否确认推荐?", "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        });
        if (flag) {
          try {
            await this.axios.form("/api/admin/topic/recommend", {
              id,
            });
            this.$message.success("操作成功");
            this.list();
          } catch (e) {
            this.$notify.error({ title: "错误", message: e.message });
          }
        }
      } catch (e) {
        this.$message.success("操作已取消");
      }
    },
    async cancelRecommend(id) {
      try {
        const flag = await this.$confirm("是否取消推荐?", "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        });
        if (flag) {
          try {
            await this.axios.delete("/api/admin/topic/recommend", {
              params: {
                id,
              },
            });
            this.$message.success("操作成功");
            this.list();
          } catch (e) {
            this.$notify.error({ title: "错误", message: e.message });
          }
        }
      } catch (e) {
        this.$message.success("操作已取消");
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
    .topic-main {
      position: relative;
      width: 100%;
      margin-left: 10px;
      .topic-status {
        position: absolute;
        right: 20px;

        .el-tag {
          margin-left: 10px;
        }
      }
      .topic-nickname {
        font-size: 15px;
        font-weight: 500;
        color: #111827;
      }
      .topic-mates {
        margin-top: 10px;
        font-size: 13px;
        color: #6b7280;

        i {
          font-size: 12px;
        }

        span {
          margin-right: 15px;
        }
      }
      .topic-summary {
        margin-top: 10px;
        font-size: 14px;
      }
      .topic-image-list {
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

      .topic-tags {
        margin-top: 10px;

        .el-tag {
          margin-right: 10px;
        }
      }

      .actions {
        margin-top: 20px;
        .action-item {
          margin-right: 15px;
          font-size: 13px;
        }
      }
    }
  }
}
</style>
