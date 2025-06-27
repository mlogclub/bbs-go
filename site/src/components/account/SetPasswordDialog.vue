<template>
  <BBSDialog
    v-model:visible="visible"
    :title="$t('component.setPasswordDialog.title')"
    @ok="submit"
  >
    <div class="set-password-dialog">
      <div class="field">
        <div class="control has-icons-left">
          <input
            v-model="form.password"
            class="input"
            type="password"
            :placeholder="$t('component.setPasswordDialog.passwordPlaceholder')"
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
            :placeholder="
              $t('component.setPasswordDialog.rePasswordPlaceholder')
            "
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
const { t } = useI18n();
const visible = ref(false);
const emits = defineEmits(["success"]);
const form = ref({
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
    await useHttpPost("/api/user/set_password", useJsonToForm(form.value));
    visible.value = false;
    emits("success");
    useMsgSuccess(t("component.setPasswordDialog.success"));
  } catch (err) {
    useMsgError(err.message || err);
  }
};

defineExpose({
  show,
});
</script>

<style lang="scss" scoped>
.set-password-dialog {
  padding: 30px 30px;

  .field:not(:last-child) {
    margin-bottom: 30px;
  }
}
</style>
