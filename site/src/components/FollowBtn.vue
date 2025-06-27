<template>
  <div>
    <button
      class="button follow-btn is-link"
      :class="{ 'is-light': followed }"
      @click="follow"
    >
      <i class="iconfont icon-add" />
      <span>{{
        followed
          ? $t("component.followBtn.followed")
          : $t("component.followBtn.follow")
      }}</span>
    </button>
  </div>
</template>

<script setup>
const userStore = useUserStore();
const props = defineProps({
  userId: {
    type: Number,
    required: true,
  },
});
const emits = defineEmits(["onFollowed"]);

const { data: followed } = await useMyFetch(
  `/api/fans/is_followed?userId=${props.userId}`
);

async function follow() {
  if (!userStore.isLogin) {
    useMsgSignIn();
    return;
  }
  try {
    if (followed.value) {
      await useHttpPost(
        "/api/fans/unfollow",
        useJsonToForm({
          userId: props.userId,
        })
      );
      followed.value = false;
      emits("onFollowed", props.userId, false);
    } else {
      await useHttpPost(
        "/api/fans/follow",
        useJsonToForm({
          userId: props.userId,
        })
      );
      followed.value = true;
      emits("onFollowed", props.userId, true);
    }
  } catch (e) {
    useMsgError(e.message || e);
  }
}
</script>

<style lang="scss" scoped>
.follow-btn {
  height: 25px;
  font-size: 12px;
  i {
    font-size: 12px;
    margin-right: 5px;
  }
}
</style>
