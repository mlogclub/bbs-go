<template>
  <!-- <div class="navbar-item dropdown is-hoverable is-right msg-notice"> -->
  <div class="navbar-item is-hoverable is-right msg-notice">
    <div class="dropdown-trigger">
      <a
        :class="{ 'msg-flicker': msgcount > 0 }"
        href="/user/messages"
        class="msgicon"
      >
        <i class="iconfont icon-bell"></i>
        <sup v-if="msgcount > 0">
          {{ msgcount > 9 ? '9+' : msgcount }}
        </sup>
      </a>
    </div>
    <div class="dropdown-menu">
      <div class="dropdown-content msglist-wrapper">
        <div class="msglist">
          我是内容
        </div>
        <div class="msgfooter">
          消息中心
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      msgcount: 0
    }
  },
  computed: {
    user() {
      return this.$store.state.user.current
    }
  },
  mounted() {
    this.getMsgcount()
  },
  methods: {
    async getMsgcount() {
      if (this.user) {
        const ret = await this.$axios.get('/api/user/msgcount')
        this.msgcount = ret.count
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.msg-notice {
  .msgicon {
    font-size: 16px;
    color: #fff;

    &:hover {
      color: red;
    }
  }

  // 闪烁
  .msg-flicker {
    // animation: msgnotice 1s 3;
    animation: msgnotice 1s infinite;
  }

  @keyframes msgnotice {
    50% {
      // color: transparent;
      color: red;
    }
  }
}
</style>
