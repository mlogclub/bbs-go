<template>
  <a-modal
    v-model:visible="config.visible"
    :title="config.title"
    :size="appStore.table.size"
    @cancel="handleCancel"
    @before-ok="handleBeforeOk"
  >
    <a-form ref="formRef" :model="form" :rules="rules">
      <a-form-item :label="$t('pages.dictType.name')" field="name">
        <a-input v-model="form.name" />
      </a-form-item>

      <a-form-item :label="$t('pages.dictType.code')" field="code">
        <a-input v-model="form.code" />
      </a-form-item>

      <a-form-item :label="$t('pages.dictType.status')" field="status">
        <a-select v-model="form.status">
          <a-option :value="0" :label="$t('pages.dictType.enabled')" />
          <a-option :value="1" :label="$t('pages.dictType.disabled')" />
        </a-select>
      </a-form-item>

      <a-form-item :label="$t('pages.dictType.remark')" field="remark">
        <a-textarea v-model="form.remark" :max-length="128" />
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

  const form = ref<any>({});
  const rules = {
    name: [{ required: true, message: t('pages.dictType.pleaseInputName') }],
    code: [{ required: true, message: t('pages.dictType.pleaseInputCode') }],
    status: [{ required: true }],
  };

  const show = () => {
    form.value = {
      status: 0,
    };
    formRef.value.resetFields();

    config.isCreate = true;
    config.title = t('pages.dictType.add');
    config.visible = true;
  };

  const showEdit = async (id: any) => {
    formRef.value.resetFields();

    config.isCreate = false;
    config.title = t('pages.dictType.edit');

    try {
      form.value = await axios.get(`/api/admin/dict-type/${id}`);
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
        ? '/api/admin/dict-type/create'
        : '/api/admin/dict-type/update';
      await axios.postForm<any>(url, jsonToFormData(form.value));
      useNotificationSuccess(t('pages.dictType.submitSuccess'));
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
