<template>
  <div class="container">
    <div class="container-header">
      <a-form :model="filters" layout="inline" :size="appStore.table.size">
        <a-form-item>
          <a-input
            v-model="filters.id"
            :placeholder="$t('pages.user.filter.id')"
          />
        </a-form-item>
        <a-form-item>
          <a-input
            v-model="filters.username"
            :placeholder="$t('pages.user.filter.username')"
          />
        </a-form-item>
        <a-form-item>
          <a-input
            v-model="filters.nickname"
            :placeholder="$t('pages.user.filter.nickname')"
          />
        </a-form-item>
        <a-form-item>
          <a-select
            v-model="filters.type"
            :placeholder="$t('pages.user.filter.type')"
            allow-clear
          >
            <a-option :value="0" :label="$t('pages.user.filter.typeUser')" />
            <a-option :value="1" :label="$t('pages.user.filter.typeStaff')" />
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" @click="list">
            <template #icon> <icon-search /> </template>
            {{ $t('pages.user.filter.search') }}
          </a-button>
        </a-form-item>
      </a-form>

      <div class="action-btns">
        <!-- <a-button type="primary" :size="appStore.table.size" @click="showAdd">
          <template #icon>
            <icon-plus />
          </template>
          {{ $t('pages.user.table.add') }}
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
          <a-table-column
            :title="$t('pages.user.table.id')"
            data-index="id"
          ></a-table-column>
          <a-table-column
            :title="$t('pages.user.table.type')"
            data-index="type"
          >
            <template #cell="{ record }">
              <a-tag v-if="record.type === 1" color="blue">{{
                $t('pages.user.filter.typeStaff')
              }}</a-tag>
              <a-tag v-else>{{ $t('pages.user.filter.typeUser') }}</a-tag>
            </template>
          </a-table-column>
          <a-table-column
            :title="$t('pages.user.table.avatar')"
            data-index="avatar"
          >
            <template #cell="{ record }">
              <a-avatar>
                <img v-if="record.avatar" :src="record.avatar" />
                <span v-else>{{ record.nickname }}</span>
              </a-avatar>
            </template>
          </a-table-column>
          <a-table-column
            :title="$t('pages.user.table.nickname')"
            data-index="nickname"
          ></a-table-column>
          <a-table-column
            :title="$t('pages.user.table.email')"
            data-index="email"
          ></a-table-column>
          <a-table-column
            :title="$t('pages.user.table.score')"
            data-index="score"
          ></a-table-column>
          <a-table-column
            :title="$t('pages.user.table.forbidden')"
            data-index="forbidden"
          >
            <template #cell="{ record }">
              {{
                record.forbidden
                  ? $t('pages.user.table.forbiddenYes')
                  : $t('pages.user.table.forbiddenNo')
              }}
            </template>
          </a-table-column>
          <a-table-column
            :title="$t('pages.user.table.createTime')"
            data-index="createTime"
          >
            <template #cell="{ record }">
              {{ useFormatDate(record.createTime) }}
            </template>
          </a-table-column>
          <a-table-column :title="$t('pages.user.table.action')">
            <template #cell="{ record }">
              <a-button
                type="primary"
                :size="appStore.table.size"
                @click="showEdit(record.id)"
                >{{ $t('pages.user.table.edit') }}</a-button
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
