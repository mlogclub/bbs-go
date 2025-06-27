<template>
  <BBSDialog
    v-model:visible="visible"
    :title="$t('component.setUsernameDialog.title')"
    @ok="submit"
  >
    <div class="set-username-dialog">
      <div class="notification is-small is-light">
        <ul>
          {{
            $t("component.setUsernameDialog.usernameRule")
          }}
        </ul>
      </div>
      <input
        v-model="username"
        class="input"
        type="text"
        :placeholder="$t('component.setUsernameDialog.usernamePlaceholder')"
      />
    </div>
  </BBSDialog>
</template>

<script setup>
const { t } = useI18n();
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
    await useHttpPost(
      "/api/user/set_username",
      useJsonToForm({
        username: username.value,
      })
    );
    visible.value = false;
    emits("success");
    useMsgSuccess(t("component.setUsernameDialog.success"));
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
