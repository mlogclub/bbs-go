<template>
  <div class="container">
    <div class="container-header">
      <a-form :model="filters" layout="inline" :size="appStore.table.size">
        <a-form-item>
          <a-input v-model="filters.id" placeholder="用户ID" />
        </a-form-item>
        <a-form-item>
          <a-input v-model="filters.username" placeholder="用户名" />
        </a-form-item>
        <a-form-item>
          <a-input v-model="filters.nickname" placeholder="昵称" />
        </a-form-item>
        <a-form-item>
          <a-select v-model="filters.type" placeholder="用户类型" allow-clear>
            <a-option :value="0" label="用户" />
            <a-option :value="1" label="员工" />
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" @click="list">
            <template #icon> <icon-search /> </template>
            查询
          </a-button>
        </a-form-item>
      </a-form>

      <div class="action-btns">
        <!-- <a-button type="primary" :size="appStore.table.size" @click="showAdd">
          <template #icon>
            <icon-plus />
          </template>
          新增
        </a-button> -->
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
          <a-table-column title="编号" data-index="id"></a-table-column>
          <a-table-column title="类型" data-index="type">
            <template #cell="{ record }">
              <a-tag v-if="record.type === 1" color="blue">员工</a-tag>
              <a-tag v-else>用户</a-tag>
            </template>
          </a-table-column>
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
          <a-table-column title="操作">
            <template #cell="{ record }">
              <a-button
                type="primary"
                :size="appStore.table.size"
                @click="showEdit(record.id)"
                >编辑</a-button
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

    id: undefined,
    username: undefined,
    nickname: undefined,
    type: undefined,
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
