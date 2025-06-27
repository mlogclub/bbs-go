<template>
  <a-modal
    v-model:visible="config.visible"
    :title="config.title"
    :size="appStore.table.size"
    @cancel="handleCancel"
    @before-ok="handleBeforeOk"
  >
    <a-form ref="formRef" :model="form" :rules="rules">
      <a-form-item :label="$t('pages.api.name')" field="name">
        <a-input v-model="form.name" />
      </a-form-item>

      <a-form-item :label="$t('pages.api.method')" field="method">
        <a-select v-model="form.method">
          <a-option value="GET" label="GET" />
          <a-option value="POST" label="POST" />
          <a-option value="DELETE" label="DELETE" />
          <a-option value="PUT" label="PUT" />
          <a-option value="ANY" label="ANY" />
        </a-select>
      </a-form-item>

      <a-form-item :label="$t('pages.api.path')" field="path">
        <a-input v-model="form.path" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
  const { t } = useI18n();
  const emit = defineEmits(['ok']);

  const appStore = useAppStore();
  const formRef = ref();
  const config = reactive({
    visible: false,
    isCreate: false,
    title: '',
  });

  const form = ref({
    name: undefined,
    method: undefined,
    path: undefined,
    createTime: undefined,
    updateTime: undefined,
  });

  const rules = {
    name: [{ required: true, message: t('pages.api.pleaseInputName') }],
    method: [{ required: true, message: t('pages.api.pleaseInputMethod') }],
    path: [{ required: true, message: t('pages.api.pleaseInputPath') }],
  };

  const show = () => {
    formRef.value.resetFields();

    config.isCreate = true;
    config.title = t('pages.api.new');
    config.visible = true;
  };

  const showEdit = async (id: any) => {
    formRef.value.resetFields();

    config.isCreate = false;
    config.title = t('pages.api.editTitle');

    try {
      form.value = await axios.get(`/api/admin/api/${id}`);
    } catch (e: any) {
      useHandleError(e);
    }

    config.visible = true;
  };

  const handleCancel = () => {
    formRef.value.resetFields();
  };
  const handleBeforeOk = async (done: (closed: boolean) => void) => {
    const validateErr = await formRef.value.validate();
    if (validateErr) {
      done(false);
      return;
    }
    try {
      const url = config.isCreate
        ? '/api/admin/api/create'
        : '/api/admin/api/update';
      await axios.postForm<any>(url, jsonToFormData(form.value));
      useNotificationSuccess(t('pages.api.submitSuccess'));
      emit('ok');
      done(true);
    } catch (e: any) {
      useHandleError(e);
      done(false);
    }
  };

  defineExpose({
    show,
    showEdit,
  });
</script>

<style lang="less" scoped></style>
