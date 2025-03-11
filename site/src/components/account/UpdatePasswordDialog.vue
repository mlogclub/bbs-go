<template>
  <BBSDialog v-model:visible="visible" title="设置密码" @ok="submit">
    <div class="update-password-dialog">
      <div class="field">
        <div class="control has-icons-left">
          <input
            v-model="form.oldPassword"
            class="input"
            type="password"
            placeholder="请输入当前密码"
            @keydown.enter="submit"
          />
          <span class="icon is-left">
            <i class="iconfont icon-captcha" />
          </span>
        </div>
      </div>
      <div class="field">
        <div class="control has-icons-left">
          <input
            v-model="form.password"
            class="input"
            type="password"
            placeholder="请输入新密码"
            @keydown.enter="submit"
          />
          <span class="icon is-left">
            <i class="iconfont icon-password" />
          </span>
        </div>
      </div>
      <div class="field">
        <div class="control has-icons-left">
          <input
            v-model="form.rePassword"
            class="input"
            type="password"
            placeholder="请再次输入新密码"
            @keydown.enter="submit"
          />
          <span class="icon is-left">
            <i class="iconfont icon-password" />
          </span>
        </div>
      </div>
    </div>
  </BBSDialog>
</template>

<script setup>
const visible = ref(false);
const emits = defineEmits(["success"]);
const form = ref({
  oldPassword: "",
  password: "",
  rePassword: "",
});

const show = async () => {
  visible.value = true;

  const user = await useHttpGet("/api/user/current");
  email.value = user.email;
};

const submit = async () => {
  try {
    await useHttpPostForm("/api/user/update_password", {
      body: form.value,
    });
    visible.value = false;
    emits("success");
    useMsgSuccess("修改密码成功");
  } catch (err) {
    useMsgError(err.message || err);
  }
};

defineExpose({
  show,
});
</script>

<style lang="scss" scoped>
.update-password-dialog {
  padding: 30px 30px;

  .field:not(:last-child) {
    margin-bottom: 30px;
  }
}
</style>
