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
            <el-option label="删除" :value="1" />
            <el-option label="待审核" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="list"> 查询 </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!--列表-->
    <div ref="mainContent" :style="{ height: mainHeight }" class="page-section articles">
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
                type="warning"
                icon="el-icon-s-check"
                class="action-item"
                @click="pendingSubmit(item)"
                >审核</el-link
              >
            </template>
          </div>
        </div>
      </div>
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
  </section>
</template>

<script>
import Avatar from "@/components/Avatar";
import mainHeight from "@/utils/mainHeight";
import CommentsDialog from "../comments/CommentsDialog";

export default {
  name: "Articles",
  components: { Avatar, CommentsDialog },
  data() {
    return {
      mainHeight: "300px",
      results: [],
      listLoading: false,
      page: {},
      filters: {
        status: 0,
      },
      tagOptions: [],
      updateTagsDialogVisible: false,
      updateTagForm: {
        articleId: 0,
        tags: [],
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
              me.list();
            })
            .catch((rsp) => {
              me.$notify.error({ title: "错误", message: rsp.message });
            });
        })
        .catch(() => {
          this.$message.success("操作已取消");
        });
    },
    pendingSubmit(row) {
      const me = this;
      this.$confirm("确认要过审文章？", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(() => {
          this.axios
            .form("/api/admin/article/pending", { id: row.id })
            .then((data) => {
              me.$message({ message: "审核成功", type: "success" });
              me.list();
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
        const nowTags = await this.axios.form("/api/admin/article/tags", {
          articleId: this.updateTagForm.articleId,
          tags: (this.updateTagForm.tags || []).join(","),
        });
        if (this.results && this.results.length) {
          for (let i = 0; i < this.results.length; i++) {
            if (this.results[i].id === this.updateTagForm.articleId) {
              this.results[i].tags = nowTags;
            }
          }
        }
        this.updateTagsDialogVisible = false;
        this.list();
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
