<template>
  <section class="page-container">
    <el-col :span="24" class="toolbar">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input
            v-model="filters.articleId"
            placeholder="文章编号"
          ></el-input>
        </el-form-item>
        <el-form-item>
          <el-input v-model="filters.ruleId" placeholder="规则编号"></el-input>
        </el-form-item>
        <el-form-item>
          <el-input v-model="filters.linkId" placeholder="链接编号"></el-input>
        </el-form-item>
        <el-form-item>
          <el-select
            v-model="filters.status"
            clearable
            placeholder="请选择状态"
            @change="list"
          >
            <el-option label="待审核" :value="0"></el-option>
            <el-option label="审核通过" :value="1"></el-option>
            <el-option label="审核失败" :value="2"></el-option>
            <el-option label="已发布" :value="3"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" v-on:click="list">查询</el-button>
        </el-form-item>
        <!--
        <el-form-item>
          <el-button type="primary" @click="handleAdd">新增</el-button>
        </el-form-item>
        -->
      </el-form>
    </el-col>

    <el-table
      :data="results"
      highlight-current-row
      stripe
      v-loading="listLoading"
      style="width: 100%;"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55"></el-table-column>
      <el-table-column type="expand">
        <template slot-scope="scope">
          <div class="content-form">
            <div class="form-item">
              <div class="field-key">文章编号：</div>
              <div class="field-value">{{ scope.row.articleId }}</div>
            </div>
            <div class="form-item">
              <div class="field-key">规则编号：</div>
              <div class="field-value">{{ scope.row.ruleId }}</div>
            </div>
            <div class="form-item">
              <div class="field-key">描述：</div>
              <div class="field-value">{{ scope.row.summary }}</div>
            </div>
            <div class="form-item">
              <div class="field-key">原链接：</div>
              <div class="field-value">
                <a :href="scope.row.sourceUrl" target="_blank">{{
                  scope.row.sourceUrl
                }}</a>
              </div>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="id" label="编号" width="100"></el-table-column>
      <el-table-column prop="title" label="标题"></el-table-column>
      <!--
      <el-table-column prop="ruleId" label="ruleId"></el-table-column>
      <el-table-column prop="userId" label="用户"></el-table-column>
      <el-table-column prop="linkId" label="linkId"></el-table-column>
      <el-table-column prop="summary" label="summary"></el-table-column>
      <el-table-column prop="content" label="content"></el-table-column>
      <el-table-column prop="contentType" label="contentType"></el-table-column>
      <el-table-column prop="sourceUrl" label="sourceUrl"></el-table-column>
      <el-table-column prop="sourceId" label="sourceId"></el-table-column>
      <el-table-column prop="sourceUrlMd5" label="sourceUrlMd5"></el-table-column>
      <el-table-column prop="sourceTitleMd5" label="sourceTitleMd5"></el-table-column>
      <el-table-column prop="articleId" label="文章编号"></el-table-column>
      -->
      <el-table-column prop="status" label="状态" width="80">
        <template slot-scope="scope">
          <span v-if="scope.row.status === 0">待审核</span>
          <span v-else-if="scope.row.status === 1">审核通过</span>
          <span v-else-if="scope.row.status === 2">审核失败</span>
          <span v-else-if="scope.row.status === 3">已发布</span>
          <span v-else>{{ scope.row.status }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="createTime" label="创建时间">
        <template slot-scope="scope">{{
          scope.row.createTime | formatDate
        }}</template>
      </el-table-column>
      <el-table-column label="操作" width="150">
        <template slot-scope="scope">
          <el-button size="small" @click="handleEdit(scope.$index, scope.row)"
            >编辑</el-button
          >
        </template>
      </el-table-column>
    </el-table>

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

    <el-dialog
      title="新增"
      :visible.sync="addFormVisible"
      :close-on-click-modal="false"
    >
      <el-form :model="addForm" label-width="80px" ref="addForm">
        <el-form-item label="userId">
          <el-input v-model="addForm.userId"></el-input>
        </el-form-item>

        <el-form-item label="ruleId">
          <el-input v-model="addForm.ruleId"></el-input>
        </el-form-item>

        <el-form-item label="linkId">
          <el-input v-model="addForm.linkId"></el-input>
        </el-form-item>

        <el-form-item label="title">
          <el-input v-model="addForm.title"></el-input>
        </el-form-item>

        <el-form-item label="summary">
          <el-input
            v-model="addForm.summary"
            type="textarea"
            autosize
          ></el-input>
        </el-form-item>

        <el-form-item label="content">
          <el-input
            v-model="addForm.content"
            type="textarea"
            autosize
          ></el-input>
        </el-form-item>

        <el-form-item label="contentType">
          <el-input v-model="addForm.contentType"></el-input>
        </el-form-item>

        <el-form-item label="sourceUrl">
          <el-input v-model="addForm.sourceUrl"></el-input>
        </el-form-item>

        <el-form-item label="sourceId">
          <el-input v-model="addForm.sourceId"></el-input>
        </el-form-item>

        <el-form-item label="sourceUrlMd5">
          <el-input v-model="addForm.sourceUrlMd5"></el-input>
        </el-form-item>

        <el-form-item label="sourceTitleMd5">
          <el-input v-model="addForm.sourceTitleMd5"></el-input>
        </el-form-item>

        <el-form-item label="status">
          <el-input v-model="addForm.status"></el-input>
        </el-form-item>

        <el-form-item label="articleId">
          <el-input v-model="addForm.articleId"></el-input>
        </el-form-item>

        <el-form-item label="createTime">
          <el-input v-model="addForm.createTime"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="addFormVisible = false">取消</el-button>
        <el-button
          type="primary"
          @click.native="addSubmit"
          :loading="addLoading"
          >提交</el-button
        >
      </div>
    </el-dialog>

    <el-dialog
      title="编辑"
      :visible.sync="editFormVisible"
      :fullscreen="true"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
    >
      <el-form :model="editForm" label-width="80px" ref="editForm">
        <el-input v-model="editForm.id" type="hidden"></el-input>

        <el-form-item label="标题">
          <el-input v-model="editForm.title"></el-input>
        </el-form-item>

        <el-form-item label="内容">
          <markdown-editor
            v-if="editForm.contentType == 'markdown'"
            v-model="editForm.content"
            :init-value="editForm.content"
            :height="500"
          />

          <html-editor
            v-if="editForm.contentType == 'html'"
            v-model="editForm.content"
          ></html-editor>
        </el-form-item>

        <el-form-item label="摘要">
          <el-input
            v-model="editForm.summary"
            type="textarea"
            autosize
          ></el-input>
        </el-form-item>

        <el-form-item label="状态">
          <el-select v-model="editForm.status" placeholder="请选择状态">
            <el-option label="待审核" :value="0"></el-option>
            <el-option label="审核通过" :value="1"></el-option>
            <el-option label="审核失败" :value="2"></el-option>
            <el-option label="已发布" :value="3"></el-option>
          </el-select>
        </el-form-item>

        <!--
        <el-form-item label="userId">
          <el-input v-model="editForm.userId"></el-input>
        </el-form-item>

        <el-form-item label="ruleId">
          <el-input v-model="editForm.ruleId"></el-input>
        </el-form-item>

        <el-form-item label="linkId">
          <el-input v-model="editForm.linkId"></el-input>
        </el-form-item>

        <el-form-item label="title">
          <el-input v-model="editForm.title"></el-input>
        </el-form-item>

        <el-form-item label="contentType">
          <el-input v-model="editForm.contentType"></el-input>
        </el-form-item>

        <el-form-item label="sourceUrl">
          <el-input v-model="editForm.sourceUrl"></el-input>
        </el-form-item>

        <el-form-item label="sourceId">
          <el-input v-model="editForm.sourceId"></el-input>
        </el-form-item>

        <el-form-item label="sourceUrlMd5">
          <el-input v-model="editForm.sourceUrlMd5"></el-input>
        </el-form-item>

        <el-form-item label="sourceTitleMd5">
          <el-input v-model="editForm.sourceTitleMd5"></el-input>
        </el-form-item>

        <el-form-item label="articleId">
          <el-input v-model="editForm.articleId"></el-input>
        </el-form-item>

        <el-form-item label="createTime">
          <el-input v-model="editForm.createTime"></el-input>
        </el-form-item>
        -->
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="editFormVisible = false">取消</el-button>
        <el-button
          type="primary"
          @click.native="editSubmit"
          :loading="editLoading"
          >提交</el-button
        >
      </div>
    </el-dialog>
  </section>
</template>

<script>
import HttpClient from '../../apis/HttpClient'
import MarkdownEditor from '@/components/MarkdownEditor'
import HtmlEditor from '@/components/HtmlEditor'

export default {
  name: 'List',
  components: { MarkdownEditor, HtmlEditor },
  data() {
    return {
      results: [],
      listLoading: false,
      page: {},
      filters: {},
      selectedRows: [],

      addForm: {
        userId: '',

        ruleId: '',

        linkId: '',

        title: '',

        summary: '',

        content: '',

        contentType: '',

        sourceUrl: '',

        sourceId: '',

        sourceUrlMd5: '',

        sourceTitleMd5: '',

        status: '',

        articleId: '',

        createTime: ''
      },
      addFormVisible: false,
      addLoading: false,

      editForm: {
        id: '',

        userId: '',

        ruleId: '',

        linkId: '',

        title: '',

        summary: '',

        content: '',

        contentType: '',

        sourceUrl: '',

        sourceId: '',

        sourceUrlMd5: '',

        sourceTitleMd5: '',

        status: '',

        articleId: '',

        createTime: ''
      },
      editFormVisible: false,
      editLoading: false
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
      HttpClient.post('/api/admin/collect-article/list', params)
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
    handleAdd() {
      this.addForm = {
        name: '',
        description: ''
      }
      this.addFormVisible = true
    },
    addSubmit() {
      let me = this
      HttpClient.post('/api/admin/collect-article/create', this.addForm)
        .then((data) => {
          me.$message({ message: '提交成功', type: 'success' })
          me.addFormVisible = false
          me.list()
        })
        .catch((rsp) => {
          me.$notify.error({ title: '错误', message: rsp.message })
        })
    },
    handleEdit(index, row) {
      let me = this
      HttpClient.get('/api/admin/collect-article/' + row.id)
        .then((data) => {
          me.editForm = Object.assign({}, data)
          me.editFormVisible = true
        })
        .catch((rsp) => {
          me.$notify.error({ title: '错误', message: rsp.message })
        })
    },
    editSubmit() {
      let me = this
      HttpClient.post('/api/admin/collect-article/update', me.editForm)
        .then((data) => {
          me.list()
          me.editFormVisible = false
        })
        .catch((rsp) => {
          me.$notify.error({ title: '错误', message: rsp.message })
        })
    },

    handleSelectionChange(val) {
      this.selectedRows = val
    }
  }
}
</script>

<style lang="scss" scoped></style>
