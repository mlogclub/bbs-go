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

      <a-form-item label="编码" field="code">
        <a-input v-model="form.code" />
      </a-form-item>

      <a-form-item label="状态" field="status">
        <a-select v-model="form.status">
          <a-option :value="0" label="启用" />
          <a-option :value="1" label="禁用" />
        </a-select>
      </a-form-item>

      <a-form-item label="备注" field="remark">
        <a-textarea v-model="form.remark" :max-length="128" />
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

  // const form = ref({
  //   name: undefined,
  //   code: undefined,
  //   status: 0,
  //   remark: undefined,
  // });
  const form = ref<any>({});
  const rules = {
    name: [{ required: true, message: '请输入名称' }],
    code: [{ required: true, message: '请输入编码' }],
    status: [{ required: true }],
  };

  const show = () => {
    form.value = {
      status: 0,
    };
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
