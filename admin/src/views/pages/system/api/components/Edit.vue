<template>
  <a-modal
    v-model:visible="config.visible"
    :title="config.title"
    :size="appStore.table.size"
    @cancel="handleCancel"
    @before-ok="handleBeforeOk"
  >
    <a-form ref="formRef" :model="form" :rules="rules">
      <a-form-item label="名称" field="name">
        <a-input v-model="form.name" />
      </a-form-item>

      <a-form-item label="Method" field="method">
        <a-select v-model="form.method">
          <a-option value="GET" label="GET" />
          <a-option value="POST" label="POST" />
          <a-option value="DELETE" label="DELETE" />
          <a-option value="PUT" label="PUT" />
          <a-option value="ANY" label="ANY" />
        </a-select>
      </a-form-item>

      <a-form-item label="Path" field="path">
        <a-input v-model="form.path" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
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
    name: [{ required: true, message: '请输入名称' }],
    method: [{ required: true, message: '请输入Method' }],
    path: [{ required: true, message: '请输入Path' }],
  };

  const show = () => {
    formRef.value.resetFields();

    config.isCreate = true;
    config.title = '新增';
    config.visible = true;
  };

  const showEdit = async (id: any) => {
    formRef.value.resetFields();

    config.isCreate = false;
    config.title = '编辑';

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
      useNotificationSuccess('提交成功');
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
