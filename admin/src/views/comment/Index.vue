<template>
  <section class="page-container">
    <div class="main-content">
      <ul class="comments">
        <li v-for="item in results" :key="item.id">
          <div class="comment-item">
            <div
              class="avatar"
              :style="{ backgroundImage: 'url(' + item.user.avatar + ')' }"
            ></div>
            <div class="content">
              <div class="meta">
                <span class="nickname"
                  ><a
                    target="_blank"
                    :href="'https://mlog.club/user/' + item.user.id"
                    >{{ item.user.nickname }}</a
                  ></span
                >
                <span class="create-time"
                  >@{{ item.createTime | formatDate }}</span
                >
                <span v-if="item.entityType === 'article'">
                  <a
                    target="_blank"
                    :href="'https://mlog.club/article/' + item.entityId"
                    >文章：{{ item.entityId }}</a
                  >
                </span>

                <span v-if="item.entityType === 'topic'">
                  <a
                    target="_blank"
                    :href="'https://mlog.club/topic/' + item.entityId"
                    >文章：{{ item.entityId }}</a
                  >
                </span>
              </div>
              <div class="summary" v-html="item.content"></div>
              <div class="tools">
                <span class="item info" v-if="item.status === 1">已删除</span>
                <a class="item" @click="handleDelete(item)">删除</a>
              </div>
            </div>
          </div>
        </li>
      </ul>
    </div>

    <el-col :span="24" class="toolbar">
      <el-pagination
        layout="total, sizes, prev, pager, next, jumper"
        :page-sizes="[20, 50, 100, 300]"
        @current-change="handlePageChange"
        @size-change="handleLimitChange"
        :current-page="page.page"
        :page-size="page.limit"
        :total="page.total"
        style="float:right;"
      ></el-pagination>
    </el-col>
  </section>
</template>

<script>
import HttpClient from '@/apis/HttpClient'

export default {
  name: 'List',
  data() {
    return {
      results: [],
      listLoading: false,
      page: {},
      filters: {
        userId: '',
        status: ''
      },
      selectedRows: [],

      addForm: {
        userId: '',
        entityType: '',
        entityId: '',
        content: '',
        quoteId: '',
        status: '',
        createTime: ''
      },
      addFormVisible: false,
      addFormRules: {},
      addLoading: false,

      editForm: {
        id: '',
        userId: '',
        entityType: '',
        entityId: '',
        content: '',
        quoteId: '',
        status: '',
        createTime: ''
      },
      editFormVisible: false,
      editFormRules: {},
      editLoading: false
    }
  },
  mounted() {
    this.list()
  },
  methods: {
    list() {
      const me = this
      me.listLoading = true
      const params = Object.assign(me.filters, {
        page: me.page.page,
        limit: me.page.limit
      })
      HttpClient.post('/api/admin/comment/list', params)
        .then((data) => {
          me.results = data.results
          me.page = data.page
        })
        .finally(() => {
          me.listLoading = false
        })
    },
    handlePageChange(val) {
      this.page.page = val
      this.list()
    },
    handleLimitChange(val) {
      this.page.limit = val
      this.list()
    },
    handleSelectionChange(val) {
      this.selectedRows = val
    },
    handleDelete(row) {
      const me = this
      HttpClient.post(`/api/admin/comment/delete/${row.id}`)
        .then((data) => {
          me.$message.success('删除成功')
          me.list()
        })
        .catch((rsp) => {
          me.$notify.error({ title: '错误', message: rsp.message })
        })
    }
  }
}
</script>

<style scoped lang="scss">
.comments {
  list-style: none;
  padding: 0px;

  li {
    border-bottom: 1px solid #f2f2f2;
    padding-top: 10px;
    padding-bottom: 10px;

    .comment-item {
      display: flex;

      .avatar {
        min-width: 40px;
        min-height: 40px;
        width: 40px;
        height: 40px;
        border-radius: 50%;
        margin-right: 10px;
        background-repeat: no-repeat;
        background-size: contain;
        background-position: center;
      }

      .content {
        width: 100%;
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
        }

        .summary {
          font-size: 15px;
          color: #555;
        }

        .tools {
          float: right;
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
    }
  }
}
</style>
