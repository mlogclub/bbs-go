<template>
  <section>
    <!--工具条-->
    <el-col :span="24" class="toolbar" style="padding-bottom: 0px;">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.title" placeholder="标题"></el-input>
        </el-form-item>
        <el-form-item>
          <el-select v-model="filters.status" clearable placeholder="请选择状态" @change="list">
            <el-option label="正常" value="0"></el-option>
            <el-option label="删除" value="1"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" v-on:click="list">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleAdd">新增</el-button>
        </el-form-item>
      </el-form>
    </el-col>

    <!--列表-->
    <div class="articles">
      <div class="article" v-for="item in results" :key="item.id">

        <div class="article-header">
          <img class="avatar" :src="item.user.avatar"/>
          <div class="article-right">
            <div class="article-title">
              <a @click="toArticle(item)" href="javascript:void(0)">{{item.title}}</a>
            </div>
            <div class="article-meta">
              <label class="author" v-if="item.user">{{item.user.nickname}}</label>
              <label>{{item.createTime | formatDate}}</label>
              <label class="category" v-if="item.category">{{item.category.categoryName}}</label>
              <label class="tag" v-for="tag in item.tags" :key="tag.tagId">{{tag.tagName}}</label>
            </div>
          </div>
        </div>
        <div class="summary">
          {{item.summary}}
        </div>
        <div class="article-footer">
          <span class="danger" v-if="item.status === 1">已删除</span>
          <span class="info">编号：{{item.id}}</span>
          <a class="btn" @click="deleteSubmit(item)">删除</a>
        </div>
      </div>

    </div>

    <!--工具条-->
    <el-col :span="24" class="toolbar">
      <el-pagination layout="total, sizes, prev, pager, next, jumper" :page-sizes="[20, 50, 100, 300]"
                     @current-change="handlePageChange"
                     @size-change="handleLimitChange"
                     :current-page="page.page"
                     :page-size="page.limit"
                     :total="page.total"
                     style="float:right;">
      </el-pagination>
    </el-col>

    <!--新增界面-->
    <el-dialog title="新增" :visible.sync="addFormVisible"
               :fullscreen="true"
               :close-on-click-modal="false"
               :close-on-press-escape="false"
    >
      <el-form :model="addForm" label-width="80px" :rules="addFormRules" ref="addForm">
        <el-input v-model="addForm.contentType" type="hidden"></el-input>
        <el-form-item label="标签" prop="tagId">
          <el-select v-model="addForm.tagId" placeholder="请选择标签">
            <el-option-group v-for="group in tagOptions" :key="group.label" :label="group.label">
              <el-option v-for="item in group.children" :key="item.value" :label="item.label"
                         :value="item.value">
              </el-option>
            </el-option-group>
          </el-select>
        </el-form-item>

        <el-form-item label="采集" prop="title">
          <el-input v-model="collectUrl" auto-complete="off" style="width: 60%"></el-input>
          <el-button type="primary" @click.prevent="doCollect">采集</el-button>
        </el-form-item>

        <el-form-item label="用户编号" prop="title">
          <el-input v-model="addForm.userId" auto-complete="off"></el-input>
        </el-form-item>

        <el-form-item label="标题" prop="title">
          <el-input v-model="addForm.title" auto-complete="off"></el-input>
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <markdown-editor v-model="addForm.content" :init-value="addForm.content"/>
        </el-form-item>
        <!--
        <el-form-item label="摘要" prop="summary">
          <el-input type="textarea" v-model="addForm.summary" autosize></el-input>
        </el-form-item>
        -->
        <el-form-item label="原文地址" prop="sourceUrl">
          <el-input v-model="addForm.sourceUrl"></el-input>
        </el-form-item>
        <!--
        <el-form-item label="状态" prop="status">
          <el-select v-model="addForm.status" placeholder="请选择">
            <el-option :key="0" :value="0" label="正常"></el-option>
            <el-option :key="1" :value="1" label="删除"></el-option>
          </el-select>
        </el-form-item>
        -->
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="addFormVisible = false">取消</el-button>
        <el-button type="primary" @click.native="addSubmit" :loading="addLoading">提交</el-button>
      </div>
    </el-dialog>

    <!--编辑界面-->
    <el-dialog title="编辑" :visible.sync="editFormVisible"
               :fullscreen="true"
               :close-on-click-modal="false"
               :close-on-press-escape="false"
    >
      <el-form :model="editForm" label-width="80px" :rules="editFormRules" ref="editForm">
        <el-input v-model="editForm.id" type="hidden"></el-input>
        <el-input v-model="editForm.contentType" type="hidden"></el-input>
        <el-form-item label="标签" prop="tagId">
          <el-select v-model="editForm.tagId" placeholder="请选择标签">
            <el-option-group v-for="group in tagOptions" :key="group.label" :label="group.label">
              <el-option v-for="item in group.children" :key="item.value" :label="item.label"
                         :value="item.value">
              </el-option>
            </el-option-group>
          </el-select>
        </el-form-item>
        <el-form-item label="标题" prop="title">
          <el-input v-model="editForm.title" auto-complete="off"></el-input>
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <markdown-editor v-if="editForm.contentType == 'markdown'" v-model="editForm.content"
                           :init-value="editForm.content" :height="500"/>

          <html-editor v-if="editForm.contentType == 'html'" v-model="editForm.content"></html-editor>

        </el-form-item>
        <el-form-item label="摘要" prop="summary">
          <el-input type="textarea" v-model="editForm.summary" autosize></el-input>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="editForm.status" placeholder="请选择">
            <el-option :key="0" :value="0" label="正常"></el-option>
            <el-option :key="1" :value="1" label="删除"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="原文地址" prop="sourceUrl">
          <el-input v-model="editForm.sourceUrl"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="editFormVisible = false">取消</el-button>
        <el-button type="primary" @click.native="editSubmit" :loading="editLoading">提交</el-button>
      </div>
    </el-dialog>
  </section>
</template>

<script>
  import config from '../../apis/Config'
  import HttpClient from '../../apis/HttpClient'
  import MarkdownEditor from '../../components/MarkdownEditor'
  import HtmlEditor from '../../components/HtmlEditor'

  export default {
    name: 'List',
    components: {MarkdownEditor, HtmlEditor},
    data() {
      return {
        results: [],
        listLoading: false,
        page: {},
        filters: {
          title: '',
          status: ''
        },

        tagOptions: [],

        addForm: {
          contentType: 'markdown',
          userId: '',
          tagId: '',
          title: '',
          summary: '',
          content: '',
          sourceUrl: '',
          status: 0,
        },
        collectUrl: '',
        addFormVisible: false,
        addFormRules: {},
        addLoading: false,

        editForm: {
          id: '',
          contentType: 'markdown',
          tagId: '',
          title: '',
          summary: '',
          content: '',
          sourceUrl: '',
          status: 0
        },
        editFormVisible: false,
        editFormRules: {},
        editLoading: false,
      }
    },
    mounted() {
      this.list()
    },
    methods: {
      list() {
        let me = this
        me.listLoading = true
        let params = Object.assign(me.filters, {
          page: me.page.page,
          limit: me.page.limit
        })
        HttpClient.post('/api/admin/article/list', params)
          .then(data => {
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
      handleAdd() {
        this.addForm = {
          contentType: 'markdown',
          tagId: '',
          title: '',
          summary: '',
          content: '',
          sourceUrl: '',
          status: 0
        }
        this.addFormVisible = true
      },
      addSubmit() {
        let me = this
        HttpClient.post('/api/admin/article/create', this.addForm)
          .then(data => {
            me.$message({message: '提交成功', type: 'success'})
            me.addFormVisible = false
            me.list()
          })
          .catch(rsp => {
            me.$notify.error({title: '错误', message: rsp.message})
          })
      },
      handleEdit(index, row) {
        let me = this
        HttpClient.get('/api/admin/article/' + row.id)
          .then(data => {
            me.editForm = Object.assign({}, data)
            me.editFormVisible = true
          })
          .catch(rsp => {
            me.$notify.error({title: '错误', message: rsp.message})
          })
      },
      editSubmit() {
        let me = this
        HttpClient.post('/api/admin/article/update', me.editForm)
          .then(data => {
            me.list()
            me.editFormVisible = false
          })
          .catch(rsp => {
            me.$notify.error({title: '错误', message: rsp.message})
          })
      },
      deleteSubmit(row) {
        let me = this
        HttpClient.post('/api/admin/article/delete', {id: row.id})
          .then(data => {
            me.$message({message: '删除成功', type: 'success'})
            me.list()
          })
          .catch(rsp => {
            me.$notify.error({title: '错误', message: rsp.message})
          })
      },
      doCollect() {
        let me = this
        if (!me.collectUrl) {
          me.$notify.error({title: '错误', message: '请填写采集地址'})
          return
        }
        this.$confirm('确定采集吗，采集之后将覆盖现有内容?', '提示', {
          type: 'warning'
        }).then(() => {
          HttpClient.post('/api/admin/article/collect', {
            url: me.collectUrl
          }).then(data => {
            me.addForm.title = data.title
            me.addForm.content = data.content
            me.addForm.sourceUrl = me.collectUrl

            me.$message({
              showClose: true,
              message: '采集成功！',
              type: 'success'
            })
          }).catch(rsp => {
            me.$notify.error({title: '错误', message: rsp.message})
          })
        }).catch(() => {
          me.$message({
            showClose: true,
            message: '取消采集！'
          })
        })
      },
      toArticle(row) {
        window.open(config.host + '/article/' + row.id, '_blank')
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

            label.category {
              color: #3273dc;
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
