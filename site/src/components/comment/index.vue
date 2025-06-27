<template>
  <div class="comment-component" id="JComment">
    <div class="comment-header">
      <span>{{ $t("component.comment.title") }}</span>
      <span v-if="commentCount > 0">&nbsp;{{ commentCount }}</span>
    </div>

    <template v-if="isLogin">
      <div v-if="isNeedEmailVerify" class="comment-not-login">
        <div class="comment-login-div">
          {{ $t("component.comment.emailVerifyPrompt") }}
          <nuxt-link to="/user/profile/account">
            {{ $t("component.comment.accountSettingsLink") }} </nuxt-link
          >{{ $t("component.comment.emailVerifyAction") }}
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
        <a @click="useToSignIn()">{{ $t("component.comment.loginLink") }}</a>
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

const isLogin = computed(() => {
  return userStore.isLogin;
});

const input = ref(null);
const list = ref(null);

// 是否需要先邮箱认证
const isNeedEmailVerify = computed(() => {
  return (
    configStore.config.createCommentEmailVerified &&
    !userStore.user.emailVerified
  );
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
    border-radius: 3px;
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
        // color: var(--text-color2);
        margin-left: 10px;
        margin-right: 10px;
      }
    }
  }
}
</style>
