<template>
  <div>
    <a-table
      :columns="columns"
      :data="data"
      :pagination="false"
      :draggable="{ type: 'handle', width: 40 }"
      @change="handleChange"
    >
      <template #name="{ record }">
        <a-input v-model="record.title" />
      </template>
      <template #url="{ record }">
        <a-input v-model="record.url" />
      </template>
      <template #operation="{ record, rowIndex }">
        <a-space>
          <a-button
            type="primary"
            status="danger"
            shape="circle"
            @click="handleDelete(record, rowIndex)"
          >
            <icon-minus />
          </a-button>
          <a-button
            type="primary"
            status="success"
            shape="circle"
            @click="handleAdd(record, rowIndex)"
          >
            <icon-plus />
          </a-button>
        </a-space>
      </template>
    </a-table>
    <div style="margin-top: 20px">
      <a-button type="primary" :loading="loading" @click="submit">{{
        t('pages.settings.nav.submit')
      }}</a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { NavDTO } from '@/composables/types';
  import { TableData } from '@arco-design/web-vue';

  const { t } = useI18n();

  const loading = ref(false);
  const data = ref<NavDTO[]>([]);
  const columns = [
    {
      title: t('pages.settings.nav.tableTitle'),
      dataIndex: 'title',
      slotName: 'name',
    },
    {
      title: t('pages.settings.nav.tableUrl'),
      dataIndex: 'url',
      slotName: 'url',
    },
    {
      title: '',
      dataIndex: 'operation',
      slotName: 'operation',
    },
  ];

  const loadConfig = async () => {
    const ret = await axios.get<any, any>('/api/admin/sys-config/configs');
    data.value = ret.siteNavs as NavDTO[];
  };

  loadConfig();

  const handleChange = (
    newData: TableData[]
    // extra: TableChangeExtra,
    // currentData: TableData[]
  ) => {
    data.value = newData as NavDTO[];
  };

  const handleDelete = (record: TableData, rowIndex: number) => {
    data.value.splice(rowIndex, 1);
  };

  const handleAdd = (record: TableData, rowIndex: number) => {
    data.value.splice(rowIndex + 1, 0, {} as NavDTO);
  };

  const submit = async () => {
    loading.value = true;
    try {
      if (data.value.some((item) => isAnyBlank(item.title, item.url))) {
        useNotificationError(t('pages.settings.nav.message.validation'));
        return;
      }
      await axios.post('/api/admin/sys-config/save', {
        siteNavs: data.value,
      });
      await loadConfig();
      useNotificationSuccess(t('pages.settings.nav.message.submitSuccess'));
    } catch (e) {
      useHandleError(e);
    } finally {
      loading.value = false;
    }
  };
</script>

<style lang="scss" scoped></style>
