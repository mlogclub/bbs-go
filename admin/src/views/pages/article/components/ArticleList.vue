<template>
  <div class="articles">
    <template v-if="results && results.length">
      <div v-for="item in results" :key="item.id" class="article">
        <avatar :user="item.user" size="40" />
        <div class="article-right">
          <div class="article-nickname">{{ item.user.nickname }}</div>
          <div class="article-metas">
            <div>
              ID: <span>{{ item.id }}</span>
            </div>
            <div>{{ item.createTime | formatDate }}</div>
          </div>

          <a class="article-title" :href="('/article/' + item.id) | siteUrl" target="_blank">
            {{ item.title }}
          </a>

          <div v-if="item.tags && item.tags.length" class="article-tags">
            <el-tag v-for="tag in item.tags" :key="tag.tagId" type="info" size="mini">
              {{ tag.tagName }}
            </el-tag>
          </div>

          <div class="article-info">
            <el-tag v-if="item.status === 1" type="danger">已删除</el-tag>
          </div>

          <div class="article-summary">{{ item.summary }}</div>
          <div class="article-actions">
            <template v-if="item.status === 0">
              <el-link
                class="action-item"
                icon="el-icon-view"
                :href="('/article/' + item.id) | siteUrl"
                target="_blank"
                >查看详情</el-link
              >
              <el-link class="action-item" icon="el-icon-s-comment" @click="showComments(item.id)"
                >查看评论</el-link
              >
              <el-link class="action-item" icon="el-icon-edit" @click="showUpdateTags(item)"
                >修改标签</el-link
              >
              <el-link
                type="danger"
                icon="el-icon-delete"
                class="action-item"
                @click="deleteSubmit(item)"
                >删除</el-link
              >
            </template>
            <template v-if="item.status === 2">
              <el-link
                type="danger"
                icon="el-icon-delete"
                class="action-item"
                @click="deleteSubmit(item)"
                >删除</el-link
              >
              <el-link
                type="success"
                icon="el-icon-s-check"
                class="action-item"
                @click="auditSubmit(item)"
                >审核通过</el-link
              >
            </template>
          </div>
        </div>
      </div>
    </template>
    <el-empty v-else />

    <el-dialog
      :visible.sync="updateTagsDialogVisible"
      :close-on-click-modal="false"
      title="添加标签"
    >
      <el-form label-width="80px">
        <el-form-item label="标签">
          <el-select
            v-model="updateTagForm.tags"
            style="width: 100%"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="标签"
          />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="updateTagsDialogVisible = false"> 取消 </el-button>
        <el-button type="primary" @click.native="updateTags"> 提交 </el-button>
      </div>
    </el-dialog>

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
  data() {
    return {
      tagOptions: [],
      updateTagsDialogVisible: false,
      updateTagForm: {
        articleId: 0,
        tags: [],
      },
    };
  },
  methods: {
    showComments(articleId) {
      this.$refs.commentsDialog.show("article", articleId);
    },
    deleteSubmit(row) {
      const me = this;
      this.$confirm("确认要删除文章？", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(() => {
          this.axios
            .form("/api/admin/article/delete", { id: row.id })
            .then((data) => {
              me.$message({ message: "删除成功", type: "success" });
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
    auditSubmit(row) {
      const me = this;
      this.$confirm("确认审核通过？", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(() => {
          this.axios
            .form("/api/admin/article/audit", { id: row.id })
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
    async showUpdateTags(article) {
      const tags = [];
      try {
        const tagObjs = await this.axios.get("/api/admin/article/tags?articleId=" + article.id);
        if (tagObjs && tagObjs.length) {
          for (let i = 0; i < tagObjs.length; i++) {
            tags.push(tagObjs[i].tagName);
          }
        }
      } catch (e) {
        this.$message({
          type: "error",
          message: e.message || e,
        });
      }

      this.updateTagForm.articleId = article.id;
      this.updateTagForm.tags = tags;
      this.updateTagsDialogVisible = true;
    },
    async updateTags() {
      try {
        // const nowTags = await this.axios.form("/api/admin/article/tags", {
        //   articleId: this.updateTagForm.articleId,
        //   tags: (this.updateTagForm.tags || []).join(","),
        // });
        // if (this.results && this.results.length) {
        //   for (let i = 0; i < this.results.length; i++) {
        //     if (this.results[i].id === this.updateTagForm.articleId) {
        //       this.results[i].tags = nowTags;
        //     }
        //   }
        // }
        await this.axios.form("/api/admin/article/tags", {
          articleId: this.updateTagForm.articleId,
          tags: (this.updateTagForm.tags || []).join(","),
        });
        this.updateTagsDialogVisible = false;
        this.$emit("change");
      } catch (e) {
        this.$message({
          type: "error",
          message: e.message || e,
        });
      }
    },
  },
};
</script>
<style scoped lang="scss">
.articles {
  width: 100%;
  padding: 0;
  margin: 0;
  overflow-y: auto;

  .article {
    display: flex;
    padding: 20px 10px 10px 10px;
    border-bottom: 1px solid #e9e9e9;

    .article-right {
      width: 100%;
      margin-left: 10px;
      position: relative;

      .article-nickname {
        font-size: 15px;
        font-weight: 500;
        color: #111827;
      }

      .article-title {
        display: block;
        margin-top: 10px;
        color: #555;
        font-size: 16px;
        font-weight: bold;
        cursor: pointer;
        text-decoration: none;
      }

      .article-tags {
        margin-top: 10px;
        .el-tag {
          margin-right: 3px;
        }
      }

      .article-metas {
        margin-top: 10px;
        display: flex;
        font-size: 12px;
        color: #6b7280;

        & > div {
          margin-right: 10px;
        }
      }

      .article-info {
        position: absolute;
        top: 0;
        right: 0;
      }

      .article-summary {
        margin-top: 10px;
        word-break: break-all;
        -webkit-line-clamp: 2;
        overflow: hidden !important;
        text-overflow: ellipsis;
        -webkit-box-orient: vertical;
        display: -webkit-box;
        color: #4a4a4a;
        font-size: 14px;
        font-weight: 400;
        line-height: 1.5;
        margin-bottom: 10px;
      }

      .article-actions {
        margin-top: 10px;
        .action-item {
          margin-right: 10px;
        }
      }
    }
  }
}
</style>
