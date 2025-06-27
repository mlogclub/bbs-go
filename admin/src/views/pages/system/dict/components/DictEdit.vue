<template>
  <a-modal
    v-model:visible="config.visible"
    :title="config.title"
    :size="appStore.table.size"
    @cancel="handleCancel"
    @before-ok="handleBeforeOk"
  >
    <a-form ref="formRef" :model="form" :rules="rules">
      <a-form-item :label="$t('pages.dict.type')" field="typeId">
        <span v-if="dictStore.currentType">
          {{ dictStore.currentType.code }}&nbsp;({{
            dictStore.currentType.name
          }})</span
        >
      </a-form-item>

      <a-form-item :label="$t('pages.dict.parent')" field="parentId">
        <a-tree-select
          v-model="form.parentId"
          allow-clear
          :data="dicts"
          :field-names="{ key: 'id', title: 'label' }"
          :placeholder="$t('pages.dict.pleaseSelectParent')"
        >
        </a-tree-select>
      </a-form-item>

      <a-form-item :label="$t('pages.dict.name')" field="name">
        <a-input v-model="form.name" />
      </a-form-item>

      <a-form-item :label="$t('pages.dict.label')" field="label">
        <a-input v-model="form.label" />
      </a-form-item>

      <a-form-item :label="$t('pages.dict.value')" field="value">
        <a-input v-model="form.value" />
      </a-form-item>

      <a-form-item :label="$t('pages.dict.status')" field="status">
        <a-select v-model="form.status">
          <a-option :value="0" :label="$t('pages.dict.enabled')" />
          <a-option :value="1" :label="$t('pages.dict.disabled')" />
        </a-select>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
  const { t } = useI18n();

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
    name: [{ required: true, message: t('pages.dict.pleaseInputName') }],
    label: [{ required: true, message: t('pages.dict.pleaseInputLabel') }],
    value: [{ required: true, message: t('pages.dict.pleaseInputValue') }],
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
    config.title = t('pages.dict.add');

    await loadDicts();

    config.visible = true;
  };

  const showEdit = async (id) => {
    formRef.value.resetFields();

    config.isCreate = false;
    config.title = t('pages.dict.edit');

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
      useNotificationSuccess(t('pages.dict.submitSuccess'));
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
