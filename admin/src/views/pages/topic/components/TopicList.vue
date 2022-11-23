<template>
  <div class="topics">
    <div v-if="results && results.length">
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
                v-else-if="!topic.recommend && topic.status === 0"
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
            <template v-else-if="topic.status === 2">
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
                class="action-item"
                type="danger"
                icon="el-icon-delete"
                @click="deleteSubmit(topic.topicId)"
                >删除</el-link
              >
              <el-link
                type="success"
                icon="el-icon-s-check"
                class="action-item"
                @click="auditSubmit(topic)"
                >审核通过</el-link
              >
            </template>
          </div>
        </div>
      </div>
    </div>
    <el-empty v-else />

    <comments-dialog ref="commentsDialog" />
  </div>
</template>

<script>
import Avatar from "@/components/Avatar";
import CommentsDialog from "../../comments/CommentsDialog";
export default {
  components: { Avatar, CommentsDialog },
  props: {
    results: {
      type: Array,
      default() {
        return [];
      },
    },
  },
  methods: {
    deleteSubmit(topicId) {
      const me = this;
      this.$confirm("是否确认删除该帖子?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(function () {
          me.axios
            .form("/api/admin/topic/delete", { id: topicId })
            .then(function () {
              me.$message({ message: "删除成功", type: "success" });
              me.$emit("change");
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
            this.$emit("change");
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
            this.$emit("change");
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
            this.$emit("change");
          } catch (e) {
            this.$notify.error({ title: "错误", message: e.message });
          }
        }
      } catch (e) {
        this.$message.success("操作已取消");
      }
    },
    auditSubmit(row) {
      const me = this;
      this.$confirm("确认审核通过？", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(() => {
          this.axios
            .form("/api/admin/topic/audit", { id: row.topicId })
            .then((data) => {
              me.$message({ message: "审核成功", type: "success" });
              me.$emit("change");
            })
            .catch((rsp) => {
              me.$notify.error({ title: "错误", message: rsp.message });
            });
        })
        .catch(() => {
          this.$message.success("操作已取消");
        });
    },
    showComments(topicId) {
      this.$refs.commentsDialog.show("topic", topicId);
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
        color: #4e4f53;
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
        margin-top: 10px;
        .action-item {
          margin-right: 15px;
          font-size: 13px;
        }
      }
    }
  }
}
</style>
