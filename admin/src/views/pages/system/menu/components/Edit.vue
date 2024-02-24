<template>
  <a-modal
    v-model:visible="config.visible"
    :title="config.title"
    :size="appStore.table.size"
    @cancel="handleCancel"
    @before-ok="handleBeforeOk"
  >
    <a-form ref="formRef" :model="form" :rules="rules">
      <a-form-item label="上级菜单" field="parentId">
        <a-tree-select
          v-model="form.parentId"
          allow-clear
          :data="menus"
          placeholder="Please select ..."
        >
        </a-tree-select>
      </a-form-item>

      <a-form-item label="名称" field="title">
        <a-input v-model="form.title" />
      </a-form-item>

      <a-form-item label="编码" field="name">
        <a-input v-model="form.name" />
      </a-form-item>

      <a-form-item label="路径" field="path">
        <a-input v-model="form.path" />
      </a-form-item>

      <a-form-item label="ICON" field="icon">
        <icon-picker v-model="form.icon" />
      </a-form-item>

      <!-- <a-form-item label="排序" field="sortNo">
        <a-input-number v-model="form.sortNo" mode="button" />
      </a-form-item> -->

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
  import IconPicker from './IconPicker.vue';

  const emit = defineEmits(['ok']);

  const appStore = useAppStore();
  const formRef = ref();
  const config = reactive({
    visible: false,
    isCreate: false,
    title: '',
  });

  const form = ref({
    id: undefined,
    parentId: undefined,
    name: undefined,
    title: undefined,
    icon: undefined,
    path: undefined,
    // sortNo: 0,
    status: 0,
  });
  const rules = {
    // parentId: [{ required: true, message: '请选择上级菜单' }],
    title: [{ required: true, message: '请输入标题' }],
    // name: [{ required: true, message: '请输入名称' }],
    // path: [{ required: true, message: '请输入路径' }],
  };

  const menus = ref([]);
  const loadMenus = async () => {
    menus.value = await axios.get('/api/admin/menu/tree');
  };

  const show = async () => {
    form.value.id = undefined;
    formRef.value.resetFields();

    config.isCreate = true;
    config.title = '新增';

    await loadMenus();

    config.visible = true;
  };

  const showEdit = async (id: any) => {
    formRef.value.resetFields();

    config.isCreate = false;
    config.title = '编辑';

    await loadMenus();

    try {
      form.value = await axios.get(`/api/admin/menu/${id}`);
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
        ? '/api/admin/menu/create'
        : '/api/admin/menu/update';
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
