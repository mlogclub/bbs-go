<template>
  <div class="comment-component" id="JComment">
    <div class="comment-header">
      <span>评论</span>
      <span v-if="commentCount > 0">&nbsp;{{ commentCount }}</span>
    </div>

    <template v-if="isLogin">
      <div v-if="isNeedEmailVerify" class="comment-not-login">
        <div class="comment-login-div">
          请先前往
          <nuxt-link style="font-weight: 700" to="/user/profile/account">
            个人中心 &gt; 账号设置 </nuxt-link
          >页面设置邮箱，并完成邮箱认证。
        </div>
      </div>
      <template v-else>
        <comment-input
          ref="input"
          :entity-id="entityId"
          :entity-type="entityType"
          @created="commentCreated"
        />
      </template>
    </template>
    <div v-else class="comment-not-login">
      <div class="comment-login-div">
        请
        <a style="font-weight: 700" @click="useToSignIn()">登录</a>后发表观点
      </div>
    </div>

    <comment-list ref="list" :entity-id="entityId" :entity-type="entityType" />
  </div>
</template>

<script setup>
const props = defineProps({
  entityType: {
    type: String,
    default: "",
    required: true,
  },
  entityId: {
    type: Number,
    default: 0,
    required: true,
  },
  commentCount: {
    type: Number,
    default: 0,
  },
});
const emits = defineEmits(["created"]);
const userStore = useUserStore();
const configStore = useConfigStore();
const user = computed(() => {
  return userStore.user;
});

const isLogin = computed(() => {
  return userStore.user !== null;
});

const config = computed(() => {
  return configStore.config;
});

const input = ref(null);
const list = ref(null);

// 是否需要先邮箱认证
const isNeedEmailVerify = computed(() => {
  return config.createCommentEmailVerified && user && user.emailVerified;
});

const commentCreated = (data) => {
  list.value.append(data);
  emits("created", data);
};
</script>

<style lang="scss" scoped>
.comment-component {
  padding: 16px;
  background-color: var(--bg-color);
  border-radius: var(--border-radius);
  .comment-header {
    display: flex;
    color: var(--text-color);
    font-size: 16px;
    font-weight: 500;
  }

  .comment-not-login {
    margin: 10px 0;
    border: 1px solid var(--border-color);
    border-radius: 0;
    overflow: hidden;
    position: relative;
    padding: 10px;
    box-sizing: border-box;

    .comment-login-div {
      color: var(--text-color4);
      cursor: pointer;
      border-radius: 3px;
      padding: 0 10px;

      a {
        margin-left: 10px;
        margin-right: 10px;
      }
    }
  }
}
</style>
