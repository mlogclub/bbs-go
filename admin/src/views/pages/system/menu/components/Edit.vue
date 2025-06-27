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
      <a-form-item :label="$t('pages.menu.parent')" field="parentId">
        <a-tree-select
          v-model="form.parentId"
          allow-clear
          :data="menus"
          :placeholder="$t('pages.menu.pleaseSelectParent')"
        >
        </a-tree-select>
      </a-form-item>

      <a-form-item :label="$t('pages.menu.type')" field="type">
        <a-select
          v-model="form.type"
          :placeholder="$t('pages.menu.pleaseSelectType')"
        >
          <a-option value="menu" :label="$t('pages.menu.typeMenu')" />
          <a-option value="func" :label="$t('pages.menu.typeFunc')" />
        </a-select>
      </a-form-item>

      <a-form-item :label="$t('pages.menu.title')" field="title">
        <a-input v-model="form.title" />
      </a-form-item>

      <a-form-item
        v-if="form.type === 'menu'"
        :label="$t('pages.menu.name')"
        field="name"
      >
        <a-input
          v-model="form.name"
          :placeholder="$t('pages.menu.namePlaceholder')"
        />
      </a-form-item>

      <a-form-item
        v-if="form.type === 'menu'"
        :label="$t('pages.menu.component')"
        field="component"
      >
        <a-input v-model="form.component" />
      </a-form-item>

      <a-form-item
        v-if="form.type === 'menu'"
        :label="$t('pages.menu.path')"
        field="path"
      >
        <a-input v-model="form.path" />
      </a-form-item>

      <a-form-item
        v-if="form.type === 'menu'"
        :label="$t('pages.menu.icon')"
        field="icon"
      >
        <icon-picker v-model="form.icon" />
      </a-form-item>

      <a-form-item :label="$t('pages.menu.status')" field="status">
        <a-select v-model="form.status">
          <a-option :value="0">{{ $t('pages.menu.statusNormal') }}</a-option>
          <a-option :value="1">{{ $t('pages.menu.statusDisabled') }}</a-option>
        </a-select>
      </a-form-item>

      <a-form-item :label="$t('pages.menu.api')">
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
              <span>{{ $t('pages.menu.apiAdd') }}</span>
            </a-button>
          </div>
        </div>
      </a-form-item>
    </a-form>

    <a-modal v-model:visible="apiDialogVisible" :width="950" height="300">
      <div style="margin-bottom: 10px">
        <a-input v-model="apiFilter" :placeholder="$t('pages.menu.apiSearch')">
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
          <a-table-column
            :title="$t('pages.menu.apiId')"
            data-index="id"
            :width="80"
          />
          <a-table-column
            :title="$t('pages.menu.apiMethod')"
            data-index="method"
            :width="80"
          >
            <template #cell="{ record }">
              <a-tag>{{ record.method }}</a-tag>
            </template>
          </a-table-column>
          <a-table-column :title="$t('pages.menu.apiPath')" data-index="path" />
          <a-table-column :title="$t('pages.menu.apiName')" data-index="name" />
        </template>
      </a-table>
    </a-modal>
  </a-drawer>
</template>

<script setup>
  import IconPicker from './IconPicker.vue';

  const emit = defineEmits(['ok']);

  const { t } = useI18n();

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
    type: [{ required: true, message: t('pages.menu.pleaseSelectType') }],
    title: [{ required: true, message: t('pages.menu.pleaseInputTitle') }],
    // name: [{ required: true, message: t('pages.menu.pleaseInputName') }],
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
    config.title = t('pages.menu.new');

    await loadMenus();

    config.visible = true;
  };

  const showEdit = async (id) => {
    formRef.value.resetFields();

    config.isCreate = false;
    config.title = t('pages.menu.editTitle');

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
      (await axios.postForm) < any > (url, jsonToFormData(form.value));
      useNotificationSuccess(t('pages.menu.submitSuccess'));
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
