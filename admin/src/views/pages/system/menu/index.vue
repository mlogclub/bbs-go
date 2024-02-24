<template>
  <div class="container">
    <div class="container-header">
      <a-form :model="filters" layout="inline" :size="appStore.table.size">
        <a-form-item>
          <a-select
            v-model="filters.status"
            placeholder="状态"
            allow-clear
            @change="list"
          >
            <a-option :value="0" label="正常" />
            <a-option :value="1" label="禁用" />
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
        <a-button type="primary" :size="appStore.table.size" @click="showAdd">
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
        :data="data.results"
        :size="appStore.table.size"
        :bordered="appStore.table.bordered"
        :pagination="false"
        show-empty-tree
        column-resizable
        :draggable="{ type: 'handle', width: 40 }"
        row-key="id"
        @change="handleChange"
      >
        <template #columns>
          <!-- <a-table-column title="编号" data-index="id" /> -->

          <a-table-column title="名称" data-index="title" />

          <a-table-column title="编码" data-index="name" />

          <a-table-column title="ICON" data-index="icon">
            <template #cell="{ record }">
              <component :is="record.icon" v-if="record.icon" :size="18" />
            </template>
          </a-table-column>

          <a-table-column title="路径" data-index="path" />

          <!-- <a-table-column title="排序" data-index="sortNo" /> -->

          <a-table-column title="状态" data-index="status">
            <template #cell="{ record }">
              {{ record.status === 0 ? '正常' : '禁用' }}
            </template>
          </a-table-column>

          <a-table-column title="创建时间" data-index="createTime">
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
    status: 0,
  });

  const data = reactive({
    results: [],
  });

  const list = async () => {
    loading.value = true;
    try {
      data.results = await axios.postForm<any>(
        '/api/admin/menu/list',
        jsonToFormData(filters)
      );
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

    await axios.post('/api/admin/menu/update_sort', ids);
    await list();
  };
</script>

<style scoped lang="less"></style>
