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

      <a-form-item label="描述" field="description">
        <a-input v-model="form.description" />
      </a-form-item>

      <a-form-item label="LOGO" field="logo">
        <image-upload v-model="form.logo" />
      </a-form-item>

      <a-form-item label="排序" field="sortNo">
        <a-input v-model="form.sortNo" />
      </a-form-item>

      <a-form-item field="status" label="状态">
        <a-select v-model="form.status">
          <a-option :value="0">正常</a-option>
          <a-option :value="1">删除</a-option>
        </a-select>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
  import ImageUpload from '@/components/ImageUpload.vue';

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
    description: undefined,
    logo: undefined,
    sortNo: undefined,
    status: undefined,
    createTime: undefined,
  });
  const rules = {};

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
      form.value = await axios.get(`/api/admin/topic-node/${id}`);
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
        ? '/api/admin/topic-node/create'
        : '/api/admin/topic-node/update';
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
