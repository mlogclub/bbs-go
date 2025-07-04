<template>
  <a-modal
    v-model:visible="config.visible"
    :title="config.title"
    :size="appStore.table.size"
    @cancel="handleCancel"
    @before-ok="handleBeforeOk"
  >
    <a-form ref="formRef" :model="form" :rules="rules">
      <a-form-item field="title" :label="$t('pages.link.title')">
        <a-input v-model="form.title" />
      </a-form-item>
      <a-form-item field="url" :label="$t('pages.link.url')">
        <a-input v-model="form.url" />
      </a-form-item>
      <a-form-item field="summary" :label="$t('pages.link.summary')">
        <a-textarea v-model="form.summary" allow-clear />
      </a-form-item>
      <a-form-item field="status" :label="$t('pages.link.status')">
        <a-select v-model="form.status">
          <a-option :value="0">{{ $t('pages.link.statusNormal') }}</a-option>
          <a-option :value="1">{{ $t('pages.link.statusDeleted') }}</a-option>
        </a-select>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
  const emit = defineEmits(['ok']);

  const { t } = useI18n();
  const appStore = useAppStore();
  const formRef = ref();
  const config = reactive({
    visible: false,
    isCreate: false,
    title: '',
  });

  const fields = {
    id: undefined,
    title: undefined,
    url: undefined,
    summary: undefined,
    status: 0,
  };
  const form = ref(fields);
  const rules = {
    title: [{ required: true, message: t('pages.link.pleaseInputTitle') }],
    url: [{ required: true, message: t('pages.link.pleaseInputUrl') }],
  };

  const show = () => {
    formRef.value.resetFields();
    form.value = fields;

    config.isCreate = true;
    config.title = t('pages.link.new');
    config.visible = true;
  };

  const showEdit = async (id: any) => {
    formRef.value.resetFields();

    config.isCreate = false;
    config.title = t('pages.link.editTitle');

    try {
      form.value = await axios.get(`/api/admin/link/${id}`);
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
        ? '/api/admin/link/create'
        : '/api/admin/link/update';
      await axios.postForm<any>(url, jsonToFormData(form.value));
      useNotificationSuccess(t('pages.link.submitSuccess'));
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
