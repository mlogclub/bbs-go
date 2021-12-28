<template>
  <div class="left-container">
    <my-counts :user="localUser" />
    <my-profile :user="localUser" />

    <fans-widget :user="localUser" />
    <follow-widget :user="localUser" />

    <div v-if="isAdmin" class="widget">
      <div class="widget-header">操作</div>
      <div class="widget-content">
        <ul class="operations">
          <li v-if="localUser.forbidden">
            <i class="iconfont icon-forbidden" />
            <a @click="removeForbidden">&nbsp;取消禁言</a>
          </li>
          <template v-else>
            <li>
              <i class="iconfont icon-forbidden" />
              <a @click="forbidden(7)">&nbsp;禁言7天</a>
            </li>
            <li>
              <i v-if="isSiteOwner" class="iconfont icon-forbidden" />
              <a @click="forbidden(-1)">&nbsp;永久禁言</a>
            </li>
          </template>
        </ul>
      </div>
    </div>

    <div class="ad">
      <!-- 展示广告 -->
      <adsbygoogle ad-slot="1742173616" />
    </div>
  </div>
</template>

<script>
import MyCounts from './MyCounts.vue'
import UserHelper from '~/common/UserHelper'
export default {
  components: { MyCounts },
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      localUser: Object.assign({}, this.user),
    }
  },
  computed: {
    currentUser() {
      return this.$store.state.user.current
    },
    // 是否是主人态
    isOwner() {
      const current = this.$store.state.user.current
      return this.localUser && current && this.localUser.id === current.id
    },
    isSiteOwner() {
      return UserHelper.isOwner(this.currentUser)
    },
    isAdmin() {
      return (
        UserHelper.isOwner(this.currentUser) ||
        UserHelper.isAdmin(this.currentUser)
      )
    },
  },
  methods: {
    forbidden(days) {
      const me = this
      const msg = days > 0 ? '是否禁言该用户？' : '是否永久禁言该用户？'
      this.$confirm(msg, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      })
        .then(() => {
          me.doForbidden(days)
        })
        .catch(() => {})
    },
    async doForbidden(days) {
      try {
        await this.$axios.post('/api/user/forbidden', {
          userId: this.localUser.id,
          days,
        })
        this.localUser.forbidden = true
        this.$message.success('禁言成功')
      } catch (e) {
        this.$message.error('禁言失败')
      }
    },
    async removeForbidden() {
      try {
        await this.$axios.post('/api/user/forbidden', {
          userId: this.localUser.id,
          days: 0,
        })
        this.localUser.forbidden = false
        this.$message.success('取消禁言成功')
      } catch (e) {
        this.$message.error('取消禁言失败')
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.img-avatar {
  margin-top: 5px;
  border: 1px dotted var(--border-color);
  border-radius: 5%;
}

.operations {
  list-style: none;

  li {
    padding-left: 3px;

    font-size: 13px;
    &:hover {
      cursor: pointer;
      background-color: #fcf8e3;
      color: #8a6d3b;
      font-weight: bold;
    }

    a {
      color: var(--text-link-color);
    }
  }
}
</style>
