<template>
  <div class="container">
    <div class="container-main">
      <div class="type-tool-bar">
        <a-input-search v-model="filters.query" placeholder="搜索" />
        <div class="btn-add">
          <a-button type="primary" @click="showAdd">
            <template #icon>
              <icon-plus />
            </template>
          </a-button>
        </div>
      </div>
      <div v-if="filteredResults && filteredResults.length" class="type-list">
        <div
          v-for="item in filteredResults"
          :key="item.id"
          class="type-item"
          :class="{ active: dictStore.currentTypeId == item.id }"
          @click="dictStore.switchType(item)"
        >
          <div class="type-item-l">
            <div class="type-item-name">{{ item.name }}</div>
            <div class="type-item-code">{{ item.code }}</div>
          </div>
          <div class="type-item-r" @click="showEdit(item.id)">
            <icon-edit />
          </div>
        </div>
      </div>
      <a-empty v-else />
    </div>

    <DictTypeEdit ref="edit" @ok="list" />
  </div>
</template>

<script setup>
  import DictTypeEdit from './DictTypeEdit.vue';

  const dictStore = useDictStore();
  const loading = ref(false);
  const edit = ref();
  const filters = reactive({
    query: '',
  });
  const results = ref([]);

  const filteredResults = computed(() => {
    if (!filters.query) {
      return results.value;
    }
    const query = filters.query.toLowerCase();
    const ret = [];
    results.value.forEach((element) => {
      if (
        element.code.toLowerCase().includes(query) ||
        element.name.toLowerCase().includes(query)
      ) {
        ret.push(element);
      }
    });
    return ret;
  });

  onMounted(() => {
    useTableHeight();
  });

  const list = async () => {
    loading.value = true;
    try {
      const ret = await axios.get('/api/admin/dict-type/list');

      results.value = ret;
      if (!dictStore.currentTypeId && ret && ret.length) {
        const type = ret[0];
        await dictStore.switchType(type);
      }
    } finally {
      loading.value = false;
    }
  };

  list();

  const showAdd = () => {
    edit.value.show();
  };

  const showEdit = (id) => {
    edit.value.showEdit(id);
  };
</script>

<style scoped lang="less">
  .type-tool-bar {
    display: flex;
    align-items: center;
    column-gap: 8px;

    .btn-add {
      width: 30px;
    }
  }
  .type-list {
    margin-top: 8px;
    display: flex;
    flex-direction: column;
    row-gap: 8px;
    .type-item {
      border-radius: 3px;
      cursor: pointer;
      padding: 5px 10px;
      font-size: 14px;
      // background-color: var(--color-neutral-1);
      border: 1px solid var(--color-neutral-2);

      &:hover,
      &.active {
        // background-color: var(--color-neutral-2);
        // color: rgb(var(--arcoblue-5));
        // border: 1px solid rgb(var(--arcoblue-2));
        background-color: var(--color-neutral-2);
        border: 1px solid var(--color-neutral-2);
      }

      display: flex;
      align-items: center;
      justify-content: space-between;
      .type-item-l {
        .type-item-name {
          font-size: 14px;
          font-weight: 500;
          color: var(--color-neutral-10);
        }
        .type-item-code {
          margin-top: 5px;
          font-size: 13px;
          color: var(--color-neutral-6);
        }
      }

      .type-item-r {
        padding: 5px;
        &:hover {
          color: rgb(var(--arcoblue-5));
        }
      }
    }
  }
</style>
