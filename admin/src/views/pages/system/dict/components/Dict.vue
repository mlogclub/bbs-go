<template>
  <div class="container">
    <div class="container-header">
      <a-form :model="filters" layout="inline" :size="appStore.table.size">
        <a-form-item>
          <a-input-search v-model="filters.query" placeholder="搜索" />
        </a-form-item>
        <!-- <a-form-item>
          <a-button type="primary" html-type="submit" @click="list">
            <template #icon> <icon-search /> </template>
            查询
          </a-button>
        </a-form-item> -->
      </a-form>

      <div class="action-btns">
        <a-button
          type="primary"
          :size="appStore.table.size"
          :disabled="!dictStore.currentType"
          @click="showAdd"
        >
          <template #icon>
            <icon-plus />
          </template>
          新增
        </a-button>
      </div>
    </div>

    <div class="container-main">
      <a-table
        :loading="loading"
        :data="results"
        :size="appStore.table.size"
        :bordered="appStore.table.bordered"
        :pagination="false"
        :sticky-header="true"
        style="height: 100%"
        show-empty-tree
        column-resizable
        row-key="id"
        :draggable="{ type: 'handle', width: 40 }"
        @change="handleChange"
      >
        <template #columns>
          <a-table-column title="编号" data-index="id" />

          <!-- <a-table-column title="类型" data-index="typeId" /> -->

          <a-table-column title="Name" data-index="name" />

          <a-table-column title="Label" data-index="label" />

          <a-table-column title="Value" data-index="value" />

          <a-table-column title="状态" data-index="status">
            <template #cell="{ record }">
              {{ record.status === 0 ? '启用' : '禁用' }}
            </template>
          </a-table-column>

          <a-table-column title="更新时间" data-index="updateTime">
            <template #cell="{ record }">
              {{ useFormatDate(record.updateTime) }}
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

    <DictEdit ref="edit" @ok="dictStore.loadDicts()" />
  </div>
</template>

<script setup>
  import DictEdit from './DictEdit.vue';

  const dictStore = useDictStore();
  const appStore = useAppStore();
  const loading = ref(false);
  const edit = ref();

  const filters = reactive({
    query: '',
  });

  const results = computed(() => {
    if (!filters.query) {
      return dictStore.dicts;
    }
    const query = filters.query.toLowerCase();
    const ret = [];
    dictStore.dicts.forEach((element) => {
      if (
        element.label.toLowerCase().includes(query) ||
        element.value.toLowerCase().includes(query)
      ) {
        ret.push(element);
      }
    });
    return ret;
  });

  onMounted(() => {
    useTableHeight();
  });

  const showAdd = () => {
    edit.value.show();
  };

  const showEdit = (id) => {
    edit.value.showEdit(id);
  };

  const handleChange = async (_data) => {
    const ids = [];

    getSortedIds(_data);

    function getSortedIds(elements) {
      elements.forEach((element) => {
        ids.push(element.id);
        // 有children，children中的元素同样参与排序
        if (element.children && element.children.length) {
          getSortedIds(element.children);
        }
      });
    }

    await axios.post('/api/admin/dict/update_sort', ids);
    await dictStore.loadDicts();
  };
</script>

<style scoped lang="less">
  .dict-wrapper {
    display: flex;
    .type-container {
      width: 600px;
    }
    .dict-container {
      flex: 1;
    }
  }
</style>
