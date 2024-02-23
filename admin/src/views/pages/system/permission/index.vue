<template>
  <div class="container">
    <div class="container-main">
      <a-card class="roles-panel">
        <div v-for="role in roles" :key="role.id">
          <a>{{ role.name }}</a>
        </div>
      </a-card>
      <div>
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
            <a-table-column title="编号" data-index="id"></a-table-column>
            <a-table-column title="头像" data-index="avatar">
              <template #cell="{ record }">
                <a-avatar>
                  <img v-if="record.avatar" :src="record.avatar" />
                  <span v-else>{{ record.nickname }}</span>
                </a-avatar>
              </template>
            </a-table-column>
            <a-table-column title="昵称" data-index="nickname"></a-table-column>
            <a-table-column title="邮箱" data-index="email"></a-table-column>
            <a-table-column title="积分" data-index="score"></a-table-column>
            <a-table-column title="是否禁言" data-index="forbidden">
              <template #cell="{ record }">
                {{ record.forbidden ? '禁言' : '-' }}
              </template>
            </a-table-column>
            <a-table-column title="注册时间" data-index="createTime">
              <template #cell="{ record }">
                {{ useFormatDate(record.createTime) }}
              </template>
            </a-table-column>
          </template>
        </a-table>
      </div>
    </div>
  </div>
</template>

<script setup>
  const appStore = useAppStore();
  const loading = ref(false);
  const filters = reactive({
    limit: 20,
    page: 1,

    username: '',
    nickname: '',
  });
  const roles = ref([]);
  const menus = ref([]);
  const currentRoleId = ref();

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

    getRoles();
  });

  const getRoles = async () => {
    roles.value = await axios.get('/api/admin/role/all_roles');
  };
  const getMenus = async () => {
    // TODO
  };

  const list = async () => {
    loading.value = true;
    try {
      const ret = await axios.postForm(
        '/api/admin/user/list',
        jsonToFormData(filters)
      );
      data.page = ret.page;
      data.results = ret.results;
    } finally {
      loading.value = false;
    }
  };

  list();

  const onPageChange = (page) => {
    filters.page = page;
    list();
  };

  const onPageSizeChange = (pageSize) => {
    filters.limit = pageSize;
    list();
  };
</script>

<style lang="scss" scoped>
  .container-main {
    display: flex;
    column-gap: 10px;

    .roles-panel {
      min-width: 220px;
      padding: 10px;
    }
  }
</style>
