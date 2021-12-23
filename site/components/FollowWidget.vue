<template>
  <div class="widget">
    <div class="widget-header">
      <div>
        <span>关注</span>
        <span class="count">{{ user.followCount }}</span>
      </div>
      <div class="slot">
        <a @click="showMore">更多</a>
      </div>
    </div>
    <div class="widget-content">
      <div v-if="followList && followList.length">
        <user-follow-list :users="followList" @onFollowed="onFollowed" />
      </div>
      <div v-else class="widget-tips">没有更多内容了</div>
    </div>

    <el-dialog
      title="关注"
      :visible.sync="showFollowDialog"
      custom-class="my-dialog"
    >
      <div v-loading="followDialogLoading">
        <load-more
          v-if="followPage"
          ref="commentsLoadMore"
          v-slot="{ results }"
          :init-data="followPage"
          :params="{ userId: user.id }"
          url="/api/fans/follows"
        >
          <user-follow-list :users="results" />
        </load-more>
        <div v-else>没数据</div>
      </div>
    </el-dialog>
  </div>
</template>

<script>
export default {
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      followList: [],
      showFollowDialog: false,
      followDialogLoading: false,
      followPage: null,
    }
  },
  mounted() {
    this.loadData()
  },
  methods: {
    async loadData() {
      const data = await this.$axios.get(
        '/api/fans/recent/follow?userId=' + this.user.id
      )
      this.followList = data.results
    },
    async onFollowed(userId, followed) {
      await this.loadData()
    },
    async showMore() {
      this.showFollowDialog = true
      this.followDialogLoading = true
      try {
        this.followPage = await this.$axios.get('/api/fans/follows', {
          params: {
            userId: this.user.id,
          },
        })
      } catch (e) {
        this.$message.error(e.message || e)
      } finally {
        this.followDialogLoading = false
      }
    },
  },
}
</script>

<style lang="scss"></style>
