<template>
  <section v-loading="listLoading" class="page-container">
    <div class="toolbar">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.id" placeholder="编号"></el-input>
        </el-form-item>
        <el-form-item>
          <el-input v-model="filters.userId" placeholder="用户编号"></el-input>
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
          <el-button type="primary" @click="list">查询</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="tweets">
      <ul>
        <li v-for="tweet in results" :key="tweet.tweetId">
          <div class="tweet">
            <div class="pin-header-row">
              <div class="account-group">
                <div>
                  <a
                    :href="'/user/' + tweet.user.id"
                    :title="tweet.user.nickname"
                  >
                    <img :src="tweet.user.smallAvatar" class="avatar size-45" />
                  </a>
                </div>
                <div class="pin-header-content">
                  <div>
                    <a
                      :href="'/user/' + tweet.user.id"
                      :title="tweet.user.nickname"
                      target="_blank"
                      class="nickname"
                      >{{ tweet.user.nickname }}</a
                    >
                  </div>
                  <div class="meta-box">
                    <div class="position ellipsis">
                      {{ tweet.user.description }}
                    </div>
                    <div class="dot">·</div>
                    <div>ID: {{ tweet.tweetId }}</div>
                    <div class="dot">·</div>
                    <time
                      :datetime="
                        tweet.createTime | formatDate('yyyy-MM-ddTHH:mm:ss')
                      "
                      itemprop="datePublished"
                      >{{
                        tweet.createTime | formatDate('yyyy-MM-dd HH:mm:ss')
                      }}</time
                    >
                  </div>
                </div>
              </div>
            </div>
            <div class="pin-content-row">
              <a
                :href="'/tweet/' + tweet.tweetId"
                target="_blank"
                class="content-box"
                >{{ tweet.content }}</a
              >
            </div>
            <ul
              v-if="tweet.imageList && tweet.imageList.length > 0"
              class="pin-image-row"
            >
              <li
                v-for="(image, index) in tweet.imageList"
                :key="image + index"
              >
                <a
                  :href="'/tweet/' + tweet.tweetId"
                  target="_blank"
                  class="image-item"
                >
                  <img :src="image.preview" />
                </a>
              </li>
            </ul>

            <div class="pin-action-row">
              <div class="action-box">
                <div class="like-action action" @click="like(tweet)">
                  <div class="action-title-box">
                    <i class="iconfont icon-like" />
                    <span class="action-title">{{
                      tweet.likeCount > 0 ? tweet.likeCount : '赞'
                    }}</span>
                  </div>
                </div>
                <a
                  :href="'/tweet/' + tweet.tweetId"
                  target="_blank"
                  class="comment-action action"
                >
                  <div class="action-title-box">
                    <i class="iconfont icon-comments" />
                    <span class="action-title">{{
                      tweet.commentCount > 0 ? tweet.commentCount : '评论'
                    }}</span>
                  </div>
                </a>
                <div
                  v-if="tweet.status === 0"
                  class="like-action action"
                  @click="deleteSubmit(tweet)"
                >
                  <div class="action-title-box">
                    <i class="iconfont icon-delete" />
                    <span class="action-title">删除</span>
                  </div>
                </div>
                <div
                  v-if="tweet.status === 1"
                  class="like-action action"
                  @click="undeleteSubmit(tweet)"
                >
                  <div class="action-title-box">
                    <i class="iconfont icon-delete" />
                    <span class="action-title danger">已删除</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </li>
      </ul>
    </div>

    <div class="pagebar">
      <el-pagination
        :page-sizes="[20, 50, 100, 300]"
        :current-page="page.page"
        :page-size="page.limit"
        :total="page.total"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handlePageChange"
        @size-change="handleLimitChange"
      >
      </el-pagination>
    </div>
  </section>
</template>

<script>
import utils from '~/common/utils'

export default {
  layout: 'admin',
  data() {
    return {
      results: [],
      listLoading: false,
      page: {},
      filters: {
        status: '0',
      },
      selectedRows: [],
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
        limit: me.page.limit,
      })
      this.$axios
        .post('/api/admin/tweet/list', params)
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
    deleteSubmit(tweet) {
      const me = this
      this.$confirm('是否确认删除该动态?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      })
        .then(function () {
          me.$axios
            .post('/api/admin/tweet/delete', { id: tweet.tweetId })
            .then(function () {
              me.$message({ message: '删除成功', type: 'success' })
              me.list()
            })
            .catch(function (err) {
              me.$notify.error({ title: '错误', message: err.message || err })
            })
        })
        .catch(function () {
          me.$message({
            type: 'info',
            message: '已取消删除',
          })
        })
    },
    undeleteSubmit(tweet) {
      const me = this
      me.$axios
        .post('/api/admin/tweet/undelete', { id: tweet.tweetId })
        .then(function () {
          me.$message({ message: '已取消删除', type: 'success' })
          me.list()
        })
        .catch(function (err) {
          me.$notify.error({ title: '错误', message: err.message || err })
        })
    },
    async like(tweet) {
      try {
        await this.$axios.post('/api/tweet/like/' + tweet.tweetId)
        tweet.liked = true
        tweet.likeCount++
      } catch (e) {
        if (e.errorCode === 1) {
          this.$toast.info('请登录后点赞！！！', {
            action: {
              text: '去登录',
              onClick: (e, toastObject) => {
                utils.toSignin()
              },
            },
          })
        } else {
          this.$toast.error(e.message || e)
        }
      }
    },
  },
}
</script>

<style scoped lang="scss">
.action-title.danger {
  color: red !important;
}
</style>
