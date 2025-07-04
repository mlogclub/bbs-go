<template>
  <div class="container">
    <div class="container-header">
      <a-form :model="filters" layout="inline" :size="appStore.table.size">
        <a-form-item>
          <a-input
            v-model="filters.name"
            :placeholder="$t('pages.role.name')"
          />
        </a-form-item>
        <a-form-item>
          <a-input
            v-model="filters.code"
            :placeholder="$t('pages.role.code')"
          />
        </a-form-item>
        <a-form-item>
          <a-select
            v-model="filters.status"
            :placeholder="$t('pages.role.status')"
            allow-clear
            @change="list"
          >
            <a-option :value="0" :label="$t('pages.role.statusNormal')" />
            <a-option :value="1" :label="$t('pages.role.statusDisabled')" />
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" @click="list">
            <template #icon> <icon-search /> </template>
            {{ $t('pages.role.search') }}
          </a-button>
        </a-form-item>
      </a-form>

      <div class="action-btns">
        <a-button type="primary" :size="appStore.table.size" @click="showAdd">
          <template #icon>
            <icon-plus />
          </template>
          {{ $t('pages.role.add') }}
        </a-button>
      </div>
    </div>

    <div class="container-main">
      <a-table
        :loading="loading"
        :data="data.results"
        :size="appStore.table.size"
        :bordered="appStore.table.bordered"
        :pagination="pagination"
        :sticky-header="true"
        style="height: 100%"
        column-resizable
        :draggable="{ type: 'handle', width: 40 }"
        row-key="id"
        @change="handleChange"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
      >
        <template #columns>
          <a-table-column :title="$t('pages.role.id')" data-index="id" />

          <a-table-column :title="$t('pages.role.name')" data-index="name" />

          <a-table-column :title="$t('pages.role.code')" data-index="code" />

          <a-table-column :title="$t('pages.role.type')" data-index="name">
            <template #cell="{ record }">
              <a-tag v-if="record.type === roleTypeSystem" color="blue">{{
                $t('pages.role.typeSystem')
              }}</a-tag>
              <a-tag v-else>{{ $t('pages.role.typeCustom') }}</a-tag>
            </template>
          </a-table-column>

          <a-table-column
            :title="$t('pages.role.remark')"
            data-index="remark"
          />

          <a-table-column :title="$t('pages.role.status')" data-index="status">
            <template #cell="{ record }">
              {{
                record.status === 0
                  ? $t('pages.role.statusNormal')
                  : $t('pages.role.statusDisabled')
              }}
            </template>
          </a-table-column>

          <a-table-column
            :title="$t('pages.role.createTime')"
            data-index="createTime"
          >
            <template #cell="{ record }">
              {{ useFormatDate(record.createTime) }}
            </template>
          </a-table-column>

          <a-table-column :title="$t('pages.role.actions')">
            <template #cell="{ record }">
              <a-button
                type="primary"
                :size="appStore.table.size"
                :disabled="record.type === roleTypeSystem"
                @click="showEdit(record.id)"
                >{{ $t('pages.role.edit') }}</a-button
              >
            </template>
          </a-table-column>
        </template>
      </a-table>
    </div>

    <Edit ref="edit" @ok="list" />
  </div>
</template>

<script setup lang="ts">
  import Edit from './components/Edit.vue';

  const roleTypeSystem = 0;
  const roleTypeCustom = 1;

  const appStore = useAppStore();
  const loading = ref(false);
  const edit = ref();
  const filters = reactive({
    limit: 20,
    page: 1,

    name: undefined,
    code: undefined,
    status: 0,
  });

  const data = reactive({
    page: {
      page: 1,
      limit: 20,
      total: 0,
    },
    results: [],
  });

  const pagination = computed(() => {
    return {
      total: data.page.total,
      current: data.page.page,
      pageSize: data.page.limit,
      showTotal: true,
      showJumper: true,
      showPageSize: true,
      pageSizeOptions: [20, 50, 100, 200, 300, 500],
    };
  });

  onMounted(() => {
    useTableHeight();
  });

  const list = async () => {
    loading.value = true;
    try {
      const ret = await axios.postForm<any>(
        '/api/admin/role/list',
        jsonToFormData(filters)
      );
      data.page = ret.page;
      data.results = ret.results;
    } finally {
      loading.value = false;
    }
  };

  list();

  const showAdd = () => {
    edit.value.show();
  };

  const showEdit = (id: any) => {
    edit.value.showEdit(id);
  };

  const handleChange = async (_data: any[]) => {
    const ids: number[] = [];

    _data.forEach((element) => {
      ids.push(element.id);
    });

    await axios.post('/api/admin/role/update_sort', ids);
    await list();
  };

  const onPageChange = (page: number) => {
    filters.page = page;
    list();
  };

  const onPageSizeChange = (pageSize: number) => {
    filters.limit = pageSize;
    list();
  };
</script>

<style scoped lang="less"></style>
