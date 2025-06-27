<template>
  <div class="container">
    <div class="container-header">
      <a-form :model="filters" layout="inline" :size="appStore.table.size">
        <a-form-item>
          <a-select
            v-model="filters.status"
            :placeholder="$t('pages.menu.status')"
            allow-clear
            @change="list"
          >
            <a-option :value="0" :label="$t('pages.menu.statusNormal')" />
            <a-option :value="1" :label="$t('pages.menu.statusDisabled')" />
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" @click="list">
            <template #icon> <icon-search /> </template>
            {{ $t('pages.menu.search') }}
          </a-button>
        </a-form-item>
      </a-form>

      <div class="action-btns">
        <a-button type="primary" :size="appStore.table.size" @click="showAdd">
          <template #icon>
            <icon-plus />
          </template>
          {{ $t('pages.menu.add') }}
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
          <a-table-column :title="$t('pages.menu.title')" data-index="title" />

          <a-table-column :title="$t('pages.menu.type')" data-index="type">
            <template #cell="{ record }">
              <a-tag :color="record.type === 'menu' ? 'green' : 'red'">
                {{
                  record.type === 'menu'
                    ? $t('pages.menu.typeMenu')
                    : $t('pages.menu.typeFunc')
                }}
              </a-tag>
            </template>
          </a-table-column>

          <a-table-column :title="$t('pages.menu.name')" data-index="name" />

          <a-table-column
            :title="$t('pages.menu.component')"
            data-index="component"
          />

          <a-table-column :title="$t('pages.menu.path')" data-index="path" />

          <a-table-column :title="$t('pages.menu.status')" data-index="status">
            <template #cell="{ record }">
              {{
                record.status === 0
                  ? $t('pages.menu.statusNormal')
                  : $t('pages.menu.statusDisabled')
              }}
            </template>
          </a-table-column>

          <a-table-column
            :title="$t('pages.menu.createTime')"
            data-index="createTime"
          >
            <template #cell="{ record }">
              {{ useFormatDate(record.createTime) }}
            </template>
          </a-table-column>

          <a-table-column :title="$t('pages.menu.actions')">
            <template #cell="{ record }">
              <a-button
                type="primary"
                :size="appStore.table.size"
                @click="showEdit(record.id)"
                >{{ $t('pages.menu.edit') }}</a-button
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
