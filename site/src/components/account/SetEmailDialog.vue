<template>
  <BBSDialog v-model:visible="visible" title="设置邮箱" @ok="submit">
    <div style="padding: 60px 30px">
      <input
        v-model="email"
        class="input"
        type="text"
        placeholder="请输入邮箱"
      />
    </div>
  </BBSDialog>
</template>

<script setup>
const visible = ref(false);
const emits = defineEmits(["success"]);
const email = ref(null);

const show = async () => {
  visible.value = true;

  const user = await useHttpGet("/api/user/current");
  email.value = user.email;
};

const submit = async () => {
  try {
    await useHttpPostForm("/api/user/set_email", {
      body: {
        email: email.value,
      },
    });
    visible.value = false;
    emits("success");
    useMsgSuccess("邮箱设置成功");
  } catch (err) {
    useMsgError(err.message || err);
  }
};

defineExpose({
  show,
});
</script>

<style lang="scss" scoped></style>
