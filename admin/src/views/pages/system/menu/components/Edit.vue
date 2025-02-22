<template>
  <a-drawer
    v-model:visible="config.visible"
    :title="config.title"
    :size="appStore.table.size"
    :width="780"
    @cancel="handleCancel"
    @before-ok="handleBeforeOk"
  >
    <a-form ref="formRef" layout="vertical" :model="form" :rules="rules">
      <a-form-item label="上级" field="parentId">
        <a-tree-select
          v-model="form.parentId"
          allow-clear
          :data="menus"
          placeholder="请选择上级"
        >
        </a-tree-select>
      </a-form-item>

      <a-form-item label="类型" field="type">
        <a-select v-model="form.type" placeholder="请选择类型">
          <a-option value="menu" label="菜单" />
          <a-option value="func" label="功能" />
        </a-select>
      </a-form-item>

      <a-form-item label="标题" field="title">
        <a-input v-model="form.title" />
      </a-form-item>

      <a-form-item v-if="form.type === 'menu'" label="名称" field="name">
        <a-input v-model="form.name" placeholder="组件/功能名称，全局唯一" />
      </a-form-item>

      <a-form-item v-if="form.type === 'menu'" label="组件" field="component">
        <a-input v-model="form.component" />
      </a-form-item>

      <a-form-item v-if="form.type === 'menu'" label="路径" field="path">
        <a-input v-model="form.path" />
      </a-form-item>

      <a-form-item v-if="form.type === 'menu'" label="ICON" field="icon">
        <icon-picker v-model="form.icon" />
      </a-form-item>

      <a-form-item label="状态" field="status">
        <a-select v-model="form.status">
          <a-option :value="0">正常</a-option>
          <a-option :value="1">禁用</a-option>
        </a-select>
      </a-form-item>

      <a-form-item label="接口">
        <div class="api-panel">
          <a-table
            class="api-table"
            size="small"
            :columns="apiColumns"
            :data="form.apis"
            :pagination="false"
          >
            <template #method="{ record }">
              <a-tag>{{ record.method }}</a-tag>
            </template>
            <template #operation="{ record, rowIndex }">
              <a-space>
                <a-button
                  type="primary"
                  status="danger"
                  shape="circle"
                  size="mini"
                >
                  <icon-minus @click="removeApi(record, rowIndex)" />
                </a-button>
              </a-space>
            </template>
          </a-table>
          <div class="api-btns">
            <a-button type="primary" size="mini" @click="showApiDialog">
              <icon-plus />
              <span>添加</span>
            </a-button>
          </div>
        </div>
      </a-form-item>
    </a-form>

    <a-modal v-model:visible="apiDialogVisible" :width="950" height="300">
      <div style="margin-bottom: 10px">
        <a-input v-model="apiFilter" placeholder="搜索接口">
          <template #suffix>
            <icon-search />
          </template>
        </a-input>
      </div>
      <a-table
        :data="apiList"
        :pagination="false"
        show-empty-tree
        column-resizable
        row-key="id"
        :scroll="{
          y: 500,
        }"
        @row-click="selectApi"
      >
        <template #columns>
          <a-table-column title="ID" data-index="id" :width="80" />
          <a-table-column title="Method" data-index="method" :width="80">
            <template #cell="{ record }">
              <a-tag>{{ record.method }}</a-tag>
            </template>
          </a-table-column>
          <a-table-column title="Path" data-index="path" />
          <a-table-column title="Name" data-index="name" />
        </template>
      </a-table>
    </a-modal>
  </a-drawer>
</template>

<script setup>
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
    type: 'menu',
    parentId: undefined,
    name: undefined,
    title: undefined,
    icon: undefined,
    path: undefined,
    status: 0,
    apis: [],
  });
  const rules = {
    type: [{ required: true, message: '请选择类型' }],
    title: [{ required: true, message: '请输入标题' }],
    // name: [{ required: true, message: '请输入名称' }],
  };

  const menus = ref([]);
  const loadMenus = async () => {
    menus.value = await axios.get('/api/admin/menu/tree');
  };

  const show = async () => {
    form.value.id = undefined;
    form.value.apis = [];
    formRef.value.resetFields();

    config.isCreate = true;
    config.title = '新增';

    await loadMenus();

    config.visible = true;
  };

  const showEdit = async (id) => {
    formRef.value.resetFields();

    config.isCreate = false;
    config.title = '编辑';

    await loadMenus();

    try {
      form.value = await axios.get(`/api/admin/menu/${id}`);
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
        ? '/api/admin/menu/create'
        : '/api/admin/menu/update';
      await axios.postForm(url, jsonToFormData(form.value));
      useNotificationSuccess('提交成功');
      emit('ok');
      done(true);
    } catch (e) {
      useHandleError(e);
      done(false);
    }
  };

  const apiColumns = [
    {
      title: 'Method',
      dataIndex: 'method',
      slotName: 'method',
      width: 100,
    },
    {
      title: 'Path',
      dataIndex: 'path',
      slotName: 'path',
    },
    {
      title: 'Name',
      dataIndex: 'name',
      slotName: 'name',
    },
    {
      title: '',
      dataIndex: 'operation',
      slotName: 'operation',
      width: 30,
    },
  ];

  const apiDialogVisible = ref(false);
  const apiFilter = ref('');
  const apis = ref([]);
  const apiList = computed(() => {
    const query = apiFilter.value ? apiFilter.value.trim().toLowerCase() : '';
    if (!query) {
      return apis.value;
    }

    const ret = [];
    apis.value.forEach((api) => {
      if (
        api.path.toLowerCase().includes(query) ||
        api.name.toLowerCase().includes(query)
      ) {
        ret.push(api);
      }
    });
    return ret;
  });
  const removeApi = (record, rowIndex) => {
    form.value.apis?.splice(rowIndex, 1);
  };
  const showApiDialog = async () => {
    apis.value = await axios.get('/api/admin/api/list_all');
    apiDialogVisible.value = true;
  };

  const selectApi = (record) => {
    if (!form.value.apis || !form.value.apis.length) {
      form.value.apis = [record];
    } else {
      form.value.apis?.push(record);
    }
    apiDialogVisible.value = false;
  };

  defineExpose({
    show,
    showEdit,
  });
</script>

<style lang="less" scoped>
  .api-panel {
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    .api-table {
      width: 100%;
    }
    .api-btns {
      margin-top: 10px;
    }
  }
</style>
