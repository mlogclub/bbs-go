<template>
  <a-modal
    v-model:visible="config.visible"
    :title="config.title"
    :size="appStore.table.size"
    @cancel="handleCancel"
    @before-ok="handleBeforeOk"
  >
    <a-form ref="formRef" :model="form" :rules="rules">
      <a-form-item label="角色名称" field="name">
        <a-input v-model="form.name" />
      </a-form-item>

      <a-form-item label="角色编码" field="code">
        <a-input v-model="form.code" />
      </a-form-item>

      <!-- <a-form-item label="排序" field="sortNo">
        <a-input-number v-model="form.sortNo" mode="button" />
      </a-form-item> -->

      <a-form-item label="备注" field="remark">
        <a-input v-model="form.remark" />
      </a-form-item>

      <a-form-item label="状态" field="status">
        <a-select v-model="form.status">
          <a-option :value="0">正常</a-option>
          <a-option :value="1">禁用</a-option>
        </a-select>
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
    code: undefined,
    sortNo: 0,
    remark: undefined,
    status: 0,
  });
  const rules = {
    name: [{ required: true, message: '请输入角色名称' }],
    code: [{ required: true, message: '请输入角色编码' }],
    status: [{ required: true, message: '请选择状态' }],
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
