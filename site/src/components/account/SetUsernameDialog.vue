<template>
  <BBSDialog v-model:visible="visible" title="设置用户名" @ok="submit">
    <div class="set-username-dialog">
      <div class="notification is-small is-light">
        <ul>
          用户名只能设置一次，请谨慎操作。用户名必须由5-12位(数字、字母、_、-)组成，且必须以字母开头。
        </ul>
      </div>
      <input
        v-model="username"
        class="input"
        type="text"
        placeholder="请输入用户名"
      />
    </div>
  </BBSDialog>
</template>

<script setup>
const visible = ref(false);
const emits = defineEmits(["success"]);
const username = ref(null);

const show = async () => {
  visible.value = true;

  const user = await useHttpGet("/api/user/current");
  username.value = user.username;
};

const submit = async () => {
  try {
    await useHttpPostForm("/api/user/set_username", {
      body: {
        username: username.value,
      },
    });
    visible.value = false;
    emits("success");
    useMsgSuccess("用户名设置成功");
  } catch (err) {
    useMsgError(err.message || err);
  }
};

defineExpose({
  show,
});
</script>

<style lang="scss" scoped>
.set-username-dialog {
  padding: 30px;
  .notification {
    padding: 10px;
    margin-bottom: 10px;
    font-size: 12px;
  }
}
</style>
