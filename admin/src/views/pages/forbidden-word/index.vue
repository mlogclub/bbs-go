<template>
  <div class="container">
    <div class="container-header">
      <a-form :model="filters" layout="inline" :size="appStore.table.size">
        <a-form-item>
          <a-select
            v-model="filters.type"
            :placeholder="$t('pages.forbiddenWord.type')"
            allow-clear
            @change="list"
          >
            <a-option
              :label="$t('pages.forbiddenWord.typeWord')"
              value="word"
            />
            <a-option
              :label="$t('pages.forbiddenWord.typeRegex')"
              value="regex"
            />
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-input
            v-model="filters.word"
            :placeholder="$t('pages.forbiddenWord.word')"
          ></a-input>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" @click="list">
            <template #icon> <icon-search /> </template>
            {{ $t('pages.forbiddenWord.search') }}
          </a-button>
        </a-form-item>
      </a-form>

      <div class="action-btns">
        <a-button type="primary" :size="appStore.table.size" @click="showAdd">
          <template #icon>
            <icon-plus />
          </template>
          {{ $t('pages.forbiddenWord.add') }}
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
            :title="$t('pages.forbiddenWord.id')"
            data-index="id"
          />

          <a-table-column
            :title="$t('pages.forbiddenWord.type')"
            data-index="type"
          />

          <a-table-column
            :title="$t('pages.forbiddenWord.word')"
            data-index="word"
          />

          <a-table-column
            :title="$t('pages.forbiddenWord.remark')"
            data-index="remark"
          />

          <a-table-column
            :title="$t('pages.forbiddenWord.createTime')"
            data-index="createTime"
          >
            <template #cell="{ record }">
              {{ useFormatDate(record.createTime) }}
            </template>
          </a-table-column>

          <a-table-column
            :title="$t('pages.forbiddenWord.edit')"
            data-index="edit"
          >
            <template #cell="{ record }">
              <a-button
                type="primary"
                :size="appStore.table.size"
                @click="showEdit(record.id)"
                >{{ $t('pages.forbiddenWord.edit') }}</a-button
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

    type: undefined,
    word: undefined,
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
        '/api/admin/forbidden-word/list',
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
