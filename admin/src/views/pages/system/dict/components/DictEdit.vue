<template>
  <a-modal
    v-model:visible="config.visible"
    :title="config.title"
    :size="appStore.table.size"
    @cancel="handleCancel"
    @before-ok="handleBeforeOk"
  >
    <a-form ref="formRef" :model="form" :rules="rules">
      <a-form-item label="类型" field="typeId">
        <span v-if="dictStore.currentType">
          {{ dictStore.currentType.code }}&nbsp;({{
            dictStore.currentType.name
          }})</span
        >
      </a-form-item>

      <a-form-item label="上级" field="parentId">
        <a-tree-select
          v-model="form.parentId"
          allow-clear
          :data="dicts"
          :field-names="{ key: 'id', title: 'label' }"
          placeholder="请选择上级"
        >
        </a-tree-select>
      </a-form-item>

      <a-form-item label="Name" field="label">
        <a-input v-model="form.name" />
      </a-form-item>

      <a-form-item label="Label" field="label">
        <a-input v-model="form.label" />
      </a-form-item>

      <a-form-item label="Value" field="value">
        <a-input v-model="form.value" />
      </a-form-item>

      <!-- <a-form-item label="sortNo" field="sortNo">
        <a-input v-model="form.sortNo" />
      </a-form-item> -->

      <a-form-item label="状态" field="status">
        <a-select v-model="form.status">
          <a-option :value="0" label="启用" />
          <a-option :value="1" label="禁用" />
        </a-select>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
  const emit = defineEmits(['ok']);

  const dictStore = useDictStore();
  const appStore = useAppStore();
  const formRef = ref();
  const config = reactive({
    visible: false,
    isCreate: false,
    title: '',
  });

  const form = ref({
    typeId: undefined,
    parentId: undefined,
    label: undefined,
    value: undefined,
    status: 0,
  });
  const rules = {
    name: [{ required: true, message: '请输入Name' }],
    label: [{ required: true, message: '请输入Label' }],
    value: [{ required: true, message: '请输入Value' }],
    status: [{ required: true }],
  };

  const dicts = ref([]);
  const loadDicts = async () => {
    dicts.value = await axios.get(
      `/api/admin/dict/dicts?typeId=${dictStore.currentType.id}`
    );
  };

  const show = async () => {
    formRef.value.resetFields();

    config.isCreate = true;
    config.title = '新增';

    await loadDicts();

    config.visible = true;
  };

  const showEdit = async (id) => {
    formRef.value.resetFields();

    config.isCreate = false;
    config.title = '编辑';

    await loadDicts();

    try {
      form.value = await axios.get(`/api/admin/dict/${id}`);
    } catch (e) {
      useHandleError(e);
    }

    config.visible = true;
  };

  const handleCancel = () => {
    formRef.value.resetFields();
  };
  const handleBeforeOk = async (done) => {
    const validateErr = await formRef.value.validate();
    if (validateErr) {
      done(false);
      return;
    }
    try {
      const url = config.isCreate
        ? '/api/admin/dict/create'
        : '/api/admin/dict/update';
      if (config.isCreate) {
        form.value.typeId = dictStore.currentTypeId;
      }
      await axios.postForm(url, jsonToFormData(form.value));
      useNotificationSuccess('提交成功');
      emit('ok');
      done(true);
    } catch (e) {
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
