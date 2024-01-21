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
        <sup v-if="msgcount > 0">{{ msgcount > 9 ? "9+" : msgcount }}</sup>
      </nuxt-link>
    </div>
  </div>
</template>

<script setup>
const msgcount = ref(0);
const messages = ref([]);
const { data } = await useAsyncData(() => useMyFetch("/api/user/msgrecent"));
msgcount.value = data.value.count || 0;
messages.value = data.value.messages || [];
</script>

<style lang="scss" scoped>
.msg-notice {
  .msgicon {
    font-size: 15px;
    color: var(--text-color);

    display: flex;
    align-items: center;
    column-gap: 6px;

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
