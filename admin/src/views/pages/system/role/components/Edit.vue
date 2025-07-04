<template>
  <a-modal
    v-model:visible="config.visible"
    :title="config.title"
    :size="appStore.table.size"
    @cancel="handleCancel"
    @before-ok="handleBeforeOk"
  >
    <a-form ref="formRef" :model="form" :rules="rules">
      <a-form-item :label="$t('pages.role.name')" field="name">
        <a-input v-model="form.name" />
      </a-form-item>

      <a-form-item :label="$t('pages.role.code')" field="code">
        <a-input v-model="form.code" />
      </a-form-item>

      <a-form-item :label="$t('pages.role.remark')" field="remark">
        <a-input v-model="form.remark" />
      </a-form-item>

      <a-form-item :label="$t('pages.role.status')" field="status">
        <a-select v-model="form.status">
          <a-option :value="0">{{ $t('pages.role.statusNormal') }}</a-option>
          <a-option :value="1">{{ $t('pages.role.statusDisabled') }}</a-option>
        </a-select>
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
    code: undefined,
    sortNo: 0,
    remark: undefined,
    status: 0,
  });
  const rules = {
    name: [{ required: true, message: t('pages.role.pleaseInputName') }],
    code: [{ required: true, message: t('pages.role.pleaseInputCode') }],
    status: [{ required: true, message: t('pages.role.pleaseSelectStatus') }],
  };

  const show = () => {
    formRef.value.resetFields();

    config.isCreate = true;
    config.title = t('pages.role.new');
    config.visible = true;
  };

  const showEdit = async (id: any) => {
    formRef.value.resetFields();

    config.isCreate = false;
    config.title = t('pages.role.editTitle');

    try {
      form.value = await axios.get(`/api/admin/role/${id}`);
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
        ? '/api/admin/role/create'
        : '/api/admin/role/update';
      await axios.postForm<any>(url, jsonToFormData(form.value));
      useNotificationSuccess(t('role.submitSuccess'));
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
