<template>
  <div class="container">
    <div class="container-header">
      <a-form :model="filters" layout="inline" :size="appStore.table.size">
        <a-form-item>
          <a-input
            v-model="filters.name"
            :placeholder="$t('pages.topicNode.filter.name')"
          />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" @click="list">
            <template #icon> <icon-search /> </template>
            {{ $t('pages.topicNode.filter.search') }}
          </a-button>
        </a-form-item>
      </a-form>

      <div class="action-btns">
        <a-button type="primary" :size="appStore.table.size" @click="showAdd">
          <template #icon>
            <icon-plus />
          </template>
          {{ $t('pages.topicNode.table.add') }}
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
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
      >
        <template #columns>
          <a-table-column
            :title="$t('pages.topicNode.table.id')"
            data-index="id"
          />

          <a-table-column
            :title="$t('pages.topicNode.table.name')"
            data-index="name"
          />

          <a-table-column
            :title="$t('pages.topicNode.table.description')"
            data-index="description"
          />

          <a-table-column
            :title="$t('pages.topicNode.table.logo')"
            data-index="logo"
          >
            <template #cell="{ record }">
              <a-image
                v-if="record.logo"
                width="60"
                height="60"
                fit="cover"
                :src="record.logo"
              />
            </template>
          </a-table-column>

          <a-table-column
            :title="$t('pages.topicNode.table.sortNo')"
            data-index="sortNo"
          />

          <a-table-column
            :title="$t('pages.topicNode.table.status')"
            data-index="status"
          >
            <template #cell="{ record }">
              {{
                record.status === 0
                  ? $t('pages.topicNode.table.statusNormal')
                  : $t('pages.topicNode.table.statusDeleted')
              }}
            </template>
          </a-table-column>

          <a-table-column
            :title="$t('pages.topicNode.table.createTime')"
            data-index="createTime"
          >
            <template #cell="{ record }">
              {{ useFormatDate(record.createTime) }}
            </template>
          </a-table-column>

          <a-table-column :title="$t('pages.topicNode.table.action')">
            <template #cell="{ record }">
              <a-button
                type="primary"
                :size="appStore.table.size"
                @click="showEdit(record.id)"
                >{{ $t('pages.topicNode.table.edit') }}</a-button
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

  const appStore = useAppStore();
  const loading = ref(false);
  const edit = ref();
  const filters = reactive({
    limit: 20,
    page: 1,

    name: undefined,
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
        '/api/admin/topic-node/list',
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
