<template>
  <section class="page-container">
    <!--工具条-->
    <el-col :span="24" class="toolbar">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.id" placeholder="编号"></el-input>
        </el-form-item>
        <el-form-item>
          <el-input v-model="filters.userId" placeholder="用户编号"></el-input>
        </el-form-item>
        <el-form-item>
          <el-input v-model="filters.title" placeholder="标题"></el-input>
        </el-form-item>
        <el-form-item>
          <el-select
            v-model="filters.status"
            clearable
            placeholder="请选择状态"
            @change="list"
          >
            <el-option label="正常" value="0"></el-option>
            <el-option label="删除" value="1"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" v-on:click="list">查询</el-button>
        </el-form-item>
      </el-form>
    </el-col>

    <!--列表-->
    <div class="articles main-content">
      <div class="article" v-for="item in results" :key="item.id">
        <div class="article-header">
          <img class="avatar" :src="item.user.avatar" />
          <div class="article-right">
            <div class="article-title">
              <a @click="toArticle(item)" href="javascript:void(0)">{{
                item.title
              }}</a>
            </div>
            <div class="article-meta">
              <label class="author" v-if="item.user">{{
                item.user.nickname
              }}</label>
              <label>{{ item.createTime | formatDate }}</label>
              <label class="tag" v-for="tag in item.tags" :key="tag.tagId">{{
                tag.tagName
              }}</label>
            </div>
          </div>
        </div>
        <div class="summary">{{ item.summary }}</div>
        <div class="article-footer">
          <span class="danger" v-if="item.status === 1">已删除</span>
          <span class="info">编号：{{ item.id }}</span>
          <a class="btn" @click="deleteSubmit(item)">删除</a>
        </div>
      </div>
    </div>

    <!--工具条-->
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
        title: '',
        status: ''
      },
      tagOptions: []
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
      HttpClient.post('/api/admin/article/list', params)
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
    deleteSubmit(row) {
      const me = this
      HttpClient.post('/api/admin/article/delete', { id: row.id })
        .then((data) => {
          me.$message({ message: '删除成功', type: 'success' })
          me.list()
        })
        .catch((rsp) => {
          me.$notify.error({ title: '错误', message: rsp.message })
        })
    },
    toArticle(row) {
      window.open(`https://mlog.club/article/${row.id}`, '_blank')
    }
  }
}
</script>

<style scoped lang="scss">
.articles {
  display: table;

  .article:not(:last-child) {
    border-bottom: solid 1px rgba(140, 147, 157, 0.14);
  }

  .article {
    padding-top: 10px;
    padding-bottom: 10px;

    .article-header {
      display: flex;

      .avatar {
        max-width: 50px;
        max-height: 50px;
        min-width: 50px;
        min-height: 50px;
        border-radius: 50%;
      }

      .article-right {
        display: block;
        margin-left: 10px;

        .article-title a {
          color: #555;
          font-size: 16px;
          font-weight: bold;
          cursor: pointer;
          text-decoration: none;
        }

        .article-meta {
          margin-top: 8px;

          label:not(:last-child) {
            margin-right: 5px;
          }

          label {
            color: #999;
            font-size: 12px;
          }

          label.tag {
            align-items: center;
            background-color: #f5f5f5;
            border-radius: 4px;
            color: #4a4a4a;
            display: inline-flex;
            justify-content: center;
            line-height: 1.5;
            padding-left: 5px;
            padding-right: 5px;
            white-space: nowrap;
          }
        }
      }
    }

    .summary {
      margin-top: 10px;
      margin-left: 60px;
      word-break: break-all;
      -webkit-line-clamp: 2;
      overflow: hidden !important;
      text-overflow: ellipsis;
      -webkit-box-orient: vertical;
      display: -webkit-box;
      color: #4a4a4a;
      font-size: 12px;
      font-weight: 400;
      line-height: 1.5;
    }

    .article-footer {
      text-align: right;

      span.info {
        font-size: 12px;
        margin-right: 10px;
        background: #eee;
        padding: 2px 5px 2px 5px;
      }

      span.danger {
        font-size: 12px;
        margin-right: 10px;
        background: #eee;
        color: red;
        padding: 2px 5px 2px 5px;
      }

      a.btn {
        font-size: 12px;
        margin-right: 10px;
        color: blue;
        cursor: pointer;
      }
    }
  }
}
</style>
