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
          <el-select v-model="filters.entityType" clearable placeholder="评论对象">
            <el-option label="话题" value="topic" />
            <el-option label="文章" value="article" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-input v-model="filters.entityId" placeholder="对象编号" />
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

    <div ref="mainContent" :style="{ height: mainHeight }">
      <div v-if="results && results.length > 0" class="page-section comments-div">
        <ul class="comments">
          <li v-for="item in results" :key="item.id">
            <div class="comment-item">
              <avatar :user="item.user" />
              <div class="content">
                <div class="meta">
                  <span class="nickname">
                    <a :href="('/user/' + item.user.id) | siteUrl" target="_blank">{{
                      item.user.nickname
                    }}</a>
                  </span>

                  <span>ID: {{ item.id }}</span>

                  <span class="create-time">@{{ item.createTime | formatDate }}</span>
                  <span v-if="item.entityType === 'article'">
                    <a :href="('/article/' + item.entityId) | siteUrl" target="_blank"
                      >文章：{{ item.entityId }}</a
                    >
                  </span>

                  <span v-if="item.entityType === 'topic'">
                    <a :href="('/topic/' + item.entityId) | siteUrl" target="_blank"
                      >话题：{{ item.entityId }}</a
                    >
                  </span>

                  <div class="tools">
                    <span v-if="item.status === 1" class="item info">已删除</span>
                    <a v-else class="item" @click="handleDelete(item)">删除</a>
                  </div>
                </div>
                <div class="summary" v-html="item.content" />
              </div>
            </div>
          </li>
        </ul>
      </div>
      <div v-else class="page-section comments-div">
        <div class="notification is-primary">
          <strong>无数据 或 输入相应参数进行查询</strong>
        </div>
      </div>
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
</template>

<script>
import Avatar from "@/components/Avatar";
import mainHeight from "@/utils/mainHeight";

export default {
  name: "Comments",
  components: { Avatar },
  data() {
    return {
      mainHeight: "300px",
      results: [],
      listLoading: false,
      page: {},
      filters: {},
      selectedRows: [],
    };
  },
  mounted() {
    mainHeight(this);
  },
  methods: {
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
        const me = this;
        // 因为这个时候界面上的pagebar显示出来了，所以需要重新计算一下高度
        setTimeout(function () {
          mainHeight(me);
          me.listLoading = false;
        }, 100);
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
    handleDelete(row) {
      const me = this;
      this.axios
        .form(`/api/admin/comment/delete/${row.id}`)
        .then((data) => {
          me.$message.success("删除成功");
          me.list();
        })
        .catch((rsp) => {
          me.$notify.error({ title: "错误", message: rsp.message });
        });
    },
  },
};
</script>

<style scoped lang="scss">
.comments-div {
  //padding: 10px 20px;
  height: 100%;
  overflow-y: auto;

  .notification {
    margin: 10px;
    text-align: center;
  }

  .comments {
    width: 100%;
    list-style: none;

    li {
      width: 100%;
      padding: 10px;

      &:not(:last-child) {
        border-bottom: 1px solid #f2f2f2;
      }

      .comment-item {
        width: 100%;
        display: flex;

        .content {
          width: 100%;
          margin-left: 10px;

          .meta {
            span {
              &:not(:last-child) {
                margin-right: 5px;
              }

              font-size: 13px;
              color: #999;
              font-weight: bold;

              &.nickname {
                color: #1a1a1a;
                font-size: 14px;
                font-weight: bold;
              }

              &.create-time {
                color: #999;
                font-size: 13px;
                font-weight: normal;
              }
            }

            .tools {
              float: right;
              font-size: 13px;

              .item {
                color: blue;
                cursor: pointer;

                &:not(:last-child) {
                  margin-right: 10px;
                }

                &.info {
                  color: red;
                }
              }
            }
          }

          .summary {
            font-size: 15px;
            color: #555;
          }
        }
      }
    }
  }
}
</style>
