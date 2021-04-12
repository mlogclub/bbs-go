<template>
  <el-dropdown v-if="hasPermission" @command="handleCommand">
    <span class="el-dropdown-link">
      管理<i class="el-icon-arrow-down el-icon--right"></i>
    </span>
    <el-dropdown-menu slot="dropdown">
      <el-dropdown-item v-if="isOwner && topic.type === 0" command="edit"
        >修改</el-dropdown-item
      >
      <el-dropdown-item v-if="isOwner || isAdmin" command="recommend">{{
        topic.recommend ? '取消推荐' : '推荐'
      }}</el-dropdown-item>
      <el-dropdown-item v-if="isTopicOwner" command="delete"
        >删除</el-dropdown-item
      >
      <el-dropdown-item v-if="isOwner || isAdmin" command="forbidden7Days"
        >禁言7天</el-dropdown-item
      >
      <el-dropdown-item v-if="isOwner" command="forbiddenForever"
        >永久禁言</el-dropdown-item
      >
    </el-dropdown-menu>
  </el-dropdown>
</template>

<script>
import UserHelper from '~/common/UserHelper'

export default {
  props: {
    topic: {
      type: Object,
      required: true,
    },
  },
  computed: {
    hasPermission() {
      return (
        this.isTopicOwner ||
        UserHelper.isOwner(this.user) ||
        UserHelper.isAdmin(this.user)
      )
    },
    isTopicOwner() {
      if (!this.user || !this.topic) {
        return false
      }
      return this.user.id === this.topic.user.id
    },
    isOwner() {
      return UserHelper.isOwner(this.user)
    },
    isAdmin() {
      return UserHelper.isAdmin(this.user)
    },
    user() {
      return this.$store.state.user.current
    },
  },
  methods: {
    async handleCommand(command) {
      if (!this.topic || !this.topic.topicId) {
        return
      }
      if (command === 'edit') {
        this.editTopic(this.topic.topicId)
      } else if (command === 'delete') {
        this.deleteTopic(this.topic.topicId)
      } else if (command === 'recommend') {
        this.switchRecommend(this.topic)
      } else if (command === 'forbidden7Days') {
        await this.forbidden(7)
      } else if (command === 'forbiddenForever') {
        await this.forbidden(-1)
      } else {
        console.log('click on item ' + command)
      }
    },
    async forbidden(days) {
      try {
        await this.$axios.post('/api/user/forbidden', {
          userId: this.topic.user.id,
          days,
        })
        this.$message.success('禁言成功')
      } catch (e) {
        this.$message.error('禁言失败')
      }
    },
    deleteTopic(topicId) {
      if (!process.client) {
        return
      }
      const me = this
      this.$confirm('是否确认删除该帖子？').then(function () {
        me.$axios
          .post('/api/topic/delete/' + topicId)
          .then(() => {
            me.$msg({
              message: '删除成功',
              onClose() {
                me.$linkTo('/topics')
              },
            })
          })
          .catch((e) => {
            me.$message.error('删除失败：' + (e.message || e))
          })
      })
    },
    editTopic(topicId) {
      this.$linkTo('/topic/edit/' + topicId)
    },
    switchRecommend(topic) {
      const me = this
      const action = topic.recommend ? '取消推荐' : '推荐'
      this.$confirm(`是否确认${action}该帖子？`).then(function () {
        me.$axios
          .post('/api/topic/recommend/' + topic.topicId, {
            recommend: !topic.recommend,
          })
          .then(() => {
            me.$msg({
              message: `${action}成功`,
            })
          })
          .catch((e) => {
            me.$message.error(`${action}失败：` + (e.message || e))
          })
      })
    },
  },
}
</script>

<style lang="scss" scoped>
.el-dropdown-link {
  cursor: pointer;
  color: #909399;
  font-size: 12px;
}
.el-dropdown-menu__item {
  font-size: 12px;
}
</style>
