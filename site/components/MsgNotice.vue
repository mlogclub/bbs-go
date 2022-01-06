<template>
  <div class="navbar-item dropdown is-hoverable is-right msg-notice">
    <div class="dropdown-trigger">
      <nuxt-link
        :class="{ 'msg-flicker': msgcount > 0 }"
        to="/user/messages"
        class="msgicon"
        title="消息"
      >
        <i class="iconfont icon-message"></i>
        <span>消息</span>
        <sup v-if="msgcount > 0">{{ msgcount > 9 ? '9+' : msgcount }}</sup>
      </nuxt-link>
    </div>
    <!--
    <div v-if="messages && messages.length" class="dropdown-menu">
      <div class="dropdown-content msglist-wrapper">
        <div class="msglist">
          <ul>
            <li v-for="msg in messages" :key="msg.messageId" class="msg-item">
              <nuxt-link to="/user/messages">
                {{ msg.from.id > 0 ? msg.from.nickname : '' }}{{ msg.title }}
              </nuxt-link>
            </li>
          </ul>
        </div>
        <div class="msgfooter">
          <nuxt-link to="/user/messages">消息中心&gt;&gt;</nuxt-link>
        </div>
      </div>
    </div>
    -->
  </div>
</template>

<script>
export default {
  data() {
    return {
      msgcount: 0,
      messages: [],
    }
  },
  computed: {
    user() {
      return this.$store.state.user.current
    },
  },
  mounted() {
    this.getMsgcount()
  },
  methods: {
    async getMsgcount() {
      if (this.user) {
        const ret = await this.$axios.get('/api/user/msgrecent')
        this.msgcount = ret.count
        this.messages = ret.messages
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.msg-notice {
  .msgicon {
    font-size: 16px;
    color: var(--text-color);

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

  .msglist-wrapper {
    padding: 5px 10px;
    .msglist {
      .msg-item {
        padding: 3px 0;
        font-size: 12px;
        line-height: 21px;
        overflow: hidden;
        word-break: break-all;
        -webkit-line-clamp: 1;
        text-overflow: ellipsis;
        -webkit-box-orient: vertical;
        display: -webkit-box;
        &:not(:last-child) {
          border-bottom: 1px solid var(--border-color);
        }
      }
    }
    .msgfooter {
      border-top: 1px solid var(--border-color);
      text-align: right;
      a {
        font-size: 13px;
      }
    }
  }
}
</style>
