<template>
  <div class="right-container">
    <post-btns v-if="isOwner" />

    <div class="widget">
      <div class="widget-header">
        {{ user.nickname }}
      </div>
      <div class="widget-content">
        <img :src="user.smallAvatar" class="img-avatar" />
        <div v-if="user.description" class="description">
          <p>{{ user.description }}</p>
        </div>
        <div class="score">
          <i class="iconfont icon-score" />
          <span>{{ user.score }}</span>
          <a
            v-if="isOwner"
            class="score-log"
            href="/user/scores"
            target="_blank"
            >积分详情&gt;&gt;</a
          >
        </div>
        <ul v-if="isOwner" class="operations">
          <li>
            <i class="iconfont icon-edit" />
            <a href="/user/settings">&nbsp;编辑资料</a>
          </li>
          <li>
            <i class="iconfont icon-message" />
            <a href="/user/messages">&nbsp;消息</a>
          </li>
          <li>
            <i class="iconfont icon-favorites" />
            <a href="/user/favorites">&nbsp;收藏</a>
          </li>
        </ul>
        <ul v-if="isAdmin" class="operations">
          <li v-if="user.forbidden">
            <i class="iconfont icon-forbidden" />
            <a @click="removeForbidden">&nbsp;取消禁言</a>
          </li>
          <li v-else>
            <i class="iconfont icon-forbidden" />
            <a @click="forbidden">&nbsp;禁言</a>
          </li>
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
import PostBtns from '~/components/PostBtns'
import UserHelper from '~/common/UserHelper'
export default {
  components: { PostBtns },
  props: {
    user: {
      type: Object,
      required: true
    }
  },
  data() {
    return {}
  },
  computed: {
    currentUser() {
      return this.$store.state.user.current
    },
    // 是否是主人态
    isOwner() {
      const current = this.$store.state.user.current
      return this.user && current && this.user.id === current.id
    },
    isAdmin() {
      return (
        UserHelper.isOwner(this.currentUser) ||
        UserHelper.isAdmin(this.currentUser)
      )
    }
  },
  methods: {
    async forbidden() {
      try {
        await this.$axios.post('/api/user/forbidden', {
          userId: this.user.id,
          days: 7
        })
        this.user.forbidden = true
        this.$toast.success('禁言成功')
      } catch (e) {
        this.$toast.error('禁言失败')
      }
    },
    async removeForbidden() {
      try {
        await this.$axios.post('/api/user/forbidden', {
          userId: this.user.id,
          days: 0
        })
        this.user.forbidden = false
        this.$toast.success('取消禁言成功')
      } catch (e) {
        this.$toast.error('取消禁言失败')
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.img-avatar {
  margin-top: 5px;
  border: 1px dotted #eeeeee;
  border-radius: 5%;
}

.description {
  font-size: 14px;
  padding: 5px 0 5px 5px;
  // padding: 10px 15px;
  // border: 1px dotted #eeeeee;
  // border-left: 3px solid #eeeeee;
  background-color: #fbfbfb;
}

.score {
  span {
    font-size: 15px;
    font-weight: bold;
    color: #3c3107;
  }

  .score-log {
    color: #3273dc;
    font-size: 75%;
    margin-left: 5px;
    &:hover {
      text-decoration: underline;
    }
  }
}

.operations {
  list-style: none;
  margin-top: 8px;
  margin-left: 0px;

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
