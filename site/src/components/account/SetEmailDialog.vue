<template>
  <BBSDialog
    v-model:visible="visible"
    :title="$t('component.setEmailDialog.title')"
    @ok="submit"
  >
    <div style="padding: 60px 30px">
      <input
        v-model="email"
        class="input"
        type="text"
        :placeholder="$t('component.setEmailDialog.emailPlaceholder')"
      />
    </div>
  </BBSDialog>
</template>

<script setup>
const { t } = useI18n();
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
    await useHttpPost(
      "/api/user/set_email",
      useJsonToForm({
        email: email.value,
      })
    );
    visible.value = false;
    emits("success");
    useMsgSuccess(t("component.setEmailDialog.success"));
  } catch (err) {
    useMsgError(err.message || err);
  }
};

defineExpose({
  show,
});
</script>

<style lang="scss" scoped></style>
