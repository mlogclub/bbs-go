<template>
  <a-modal
    v-model:visible="config.visible"
    :title="config.title"
    :size="appStore.table.size"
    @cancel="handleCancel"
    @before-ok="handleBeforeOk"
  >
    <a-form ref="formRef" :model="form" :rules="rules">
      <a-form-item label="类型" field="type">
        <a-select v-model="form.type" placeholder="类型">
          <a-option label="词组" value="word" />
          <a-option label="正则表达式" value="regex" />
        </a-select>
      </a-form-item>

      <a-form-item label="违禁词" field="word">
        <a-input v-model="form.word" />
      </a-form-item>

      <a-form-item label="备注" field="remark">
        <a-input v-model="form.remark" />
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
    type: undefined,
    word: undefined,
    remark: undefined,
  });
  const rules = {
    word: [{ required: true, message: '请输入违禁词' }],
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
      form.value = await axios.get(`/api/admin/forbidden-word/${id}`);
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
        ? '/api/admin/forbidden-word/create'
        : '/api/admin/forbidden-word/update';
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
