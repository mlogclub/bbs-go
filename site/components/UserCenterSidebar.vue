<template>
  <div class="left-container">
    <div class="widget">
      <div class="widget-header">个人成就</div>
      <div class="widget-content extra-info">
        <ul>
          <li>
            <span>积分</span><br />
            <nuxt-link to="/user/scores">
              <b>{{ localUser.score }}</b>
            </nuxt-link>
          </li>
          <li>
            <span>文章</span><br />
            <b>{{ localUser.topicCount }}</b>
          </li>
          <li>
            <span>评论</span><br />
            <b>{{ localUser.commentCount }}</b>
          </li>
          <li>
            <span>注册排名</span><br />
            <b>{{ localUser.id }}</b>
          </li>
        </ul>
      </div>
    </div>

    <div v-if="isOwner || isAdmin" class="widget">
      <div class="widget-header">操作</div>
      <div class="widget-content">
        <ul class="operations">
          <template v-if="isOwner">
            <li>
              <i class="iconfont icon-message" />
              <nuxt-link to="/user/messages">&nbsp;消息</nuxt-link>
            </li>
            <li>
              <i class="iconfont icon-favorites" />
              <nuxt-link to="/user/favorites">&nbsp;收藏</nuxt-link>
            </li>
            <li>
              <i class="iconfont icon-edit" />
              <nuxt-link to="/user/settings">&nbsp;编辑资料</nuxt-link>
            </li>
          </template>
          <template v-if="isAdmin">
            <hr />
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
import UserHelper from '~/common/UserHelper'
export default {
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
.extra-info {
  ul {
    display: flex;
    li {
      width: 100%;
      text-align: center;
      span {
        font-size: 13px;
        font-weight: 400;
        color: #868e96;
      }
      a,
      b {
        color: #000;
      }
    }
  }
}

.img-avatar {
  margin-top: 5px;
  border: 1px dotted #eeeeee;
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
      color: #3273dc;
    }
  }
}
</style>
