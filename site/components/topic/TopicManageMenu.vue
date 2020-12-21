<template>
  <el-dropdown @command="handleCommand">
    <span class="el-dropdown-link">
      管理<i class="el-icon-arrow-down el-icon--right"></i>
    </span>
    <el-dropdown-menu slot="dropdown">
      <el-dropdown-item v-if="hasPermission" command="edit"
        >修改</el-dropdown-item
      >
      <el-dropdown-item v-if="hasPermission" command="delete"
        >删除</el-dropdown-item
      >
      <!--<el-dropdown-item command="report">举报</el-dropdown-item>-->
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
        this.isOwner ||
        UserHelper.isOwner(this.user) ||
        UserHelper.isAdmin(this.user)
      )
    },
    isOwner() {
      if (!this.user || !this.topic) {
        return false
      }
      return this.user.id === this.topic.user.id
    },
    user() {
      return this.$store.state.user.current
    },
  },
  methods: {
    handleCommand(command) {
      if (!this.topic || !this.topic.topicId) {
        return
      }
      if (command === 'edit') {
        this.editTopic(this.topic.topicId)
      } else if (command === 'delete') {
        this.deleteTopic(this.topic.topicId)
      } else {
        console.log('click on item ' + command)
      }
    },
    deleteTopic(topicId) {
      if (!process.client) {
        return
      }
      const me = this
      this.$confirm('是否确认删除该话题？').then(function () {
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
