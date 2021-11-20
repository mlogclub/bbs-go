<template>
  <div class="widget">
    <div class="widget-header">
      <div>
        <span>粉丝</span>
        <span class="count">{{ user.fansCount }}</span>
      </div>
      <div class="slot">
        <a @click="showMore">更多</a>
      </div>
    </div>
    <div class="widget-content">
      <div v-if="fansList && fansList.length">
        <user-follow-list :users="fansList" @onFollowed="onFollowed" />
      </div>
      <div v-else class="widget-tips">没有更多内容了</div>
    </div>

    <el-dialog
      title="粉丝"
      :visible.sync="showFansDialog"
      custom-class="my-dialog"
    >
      <div v-loading="fansDialogLoading">
        <load-more
          v-if="fansPage"
          ref="commentsLoadMore"
          v-slot="{ results }"
          :init-data="fansPage"
          :params="{ userId: user.id }"
          url="/api/fans/fans"
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
      fansList: [],
      showFansDialog: false,
      fansDialogLoading: false,
      fansPage: null,
    }
  },
  mounted() {
    this.loadData()
  },
  methods: {
    async loadData() {
      const data = await this.$axios.get(
        '/api/fans/recent/fans?userId=' + this.user.id
      )
      this.fansList = data.results
    },
    async onFollowed(userId, followed) {
      await this.loadData()
    },
    async showMore() {
      this.showFansDialog = true
      this.fansDialogLoading = true
      try {
        this.fansPage = await this.$axios.get('/api/fans/fans', {
          params: {
            userId: this.user.id,
          },
        })
      } catch (e) {
        this.$message.error(e.message || e)
      } finally {
        this.fansDialogLoading = false
      }
    },
  },
}
</script>

<style lang="scss"></style>
