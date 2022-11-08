<template>
  <el-dropdown v-if="hasPermission" trigger="click" @command="handleCommand">
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
    value: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      article: this.value,
    }
  },
  computed: {
    hasPermission() {
      return (
        this.isArticleOwner ||
        UserHelper.isOwner(this.user) ||
        UserHelper.isAdmin(this.user)
      )
    },
    isArticleOwner() {
      if (!this.user || !this.article) {
        return false
      }
      return this.user.id === this.article.user.id
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
      if (!this.article || !this.article.articleId) {
        return
      }
      if (command === 'edit') {
        this.editArticle()
      } else if (command === 'delete') {
        this.deleteArticle()
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
          userId: this.article.user.id,
          days,
        })
        this.$message.success('禁言成功')
      } catch (e) {
        this.$message.error('禁言失败')
      }
    },
    deleteArticle() {
      if (!process.client) {
        return
      }
      const me = this
      this.$confirm('是否确认删除该文章？').then(function () {
        me.$axios
          .post('/api/article/delete/' + me.article.articleId)
          .then(() => {
            me.$msg({
              message: '删除成功',
              onClose() {
                me.$linkTo('/')
              },
            })
          })
          .catch((e) => {
            me.$message.error('删除失败：' + (e.message || e))
          })
      })
    },
    editArticle() {
      this.$linkTo('/article/edit/' + this.article.articleId)
    },
  },
}
</script>

<style lang="scss" scoped>
.el-dropdown-link {
  cursor: pointer;
  color: var(--text-color3);
  font-size: 12px;
}
.el-dropdown-menu__item {
  font-size: 12px;
}
</style>
