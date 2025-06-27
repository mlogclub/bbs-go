<template>
  <div class="container">
    <div class="container-header">
      <a-form :model="filters" layout="inline" :size="appStore.table.size">
        <a-form-item>
          <a-input
            v-model="filters.id"
            :placeholder="$t('pages.topic.filter.id')"
          />
        </a-form-item>
        <a-form-item>
          <a-input
            v-model="filters.userId"
            :placeholder="$t('pages.topic.filter.userId')"
          />
        </a-form-item>
        <a-form-item>
          <a-select
            v-model="filters.status"
            :placeholder="$t('pages.topic.filter.status')"
            allow-clear
            @change="list"
          >
            <a-option
              :value="0"
              :label="$t('pages.topic.filter.statusNormal')"
            />
            <a-option
              :value="1"
              :label="$t('pages.topic.filter.statusDeleted')"
            />
            <a-option
              :value="2"
              :label="$t('pages.topic.filter.statusPending')"
            />
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-select
            v-model="filters.recommend"
            :placeholder="$t('pages.topic.filter.recommend')"
            allow-clear
            @change="list"
          >
            <a-option
              :value="1"
              :label="$t('pages.topic.filter.recommendYes')"
            />
            <a-option
              :value="0"
              :label="$t('pages.topic.filter.recommendNo')"
            />
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-input
            v-model="filters.title"
            :placeholder="$t('pages.topic.filter.title')"
          />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" @click="list">
            <template #icon> <icon-search /> </template>
            {{ $t('pages.topic.filter.search') }}
          </a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="container-main">
      <div v-if="data && data.results" class="topic-container">
        <topic-list :results="data.results" @change="list" />
      </div>
      <a-empty v-else />
    </div>
    <div class="container-footer">
      <a-pagination
        style="margin: 10px"
        :total="pagination.total"
        :current="pagination.current"
        :page-size="pagination.pageSize"
        :show-total="pagination.showTotal"
        :show-jumper="pagination.showJumper"
        :show-page-size="pagination.showPageSize"
        @change="onPageChange"
        @page-size-change="onPageSizeChange"
      />
    </div>
  </div>
</template>

<script setup>
  import TopicList from './components/TopicList.vue';

  const appStore = useAppStore();
  const loading = ref(false);
  const filters = reactive({
    limit: 20,
    page: 1,
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
    };
  });

  onMounted(() => {
    useTableHeight();
  });

  const list = async () => {
    loading.value = true;
    try {
      const ret = await axios.postForm(
        '/api/admin/topic/list',
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

<style scoped lang="less">
  .container-footer {
    // padding: 10px;
    display: flex;
    justify-content: end;
  }
</style>
