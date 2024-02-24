<template>
  <div class="container">
    <div class="container-header">
      <div style="width: max-content">
        <a-alert type="warning">设置角色对应权限</a-alert>
      </div>

      <a-button type="primary" @click="saveRoleMenus">保存</a-button>
    </div>
    <div class="container-main">
      <a-card title="角色列表" class="roles-panel" :body-style="cardBodyStyle">
        <div class="role-item-list">
          <div
            v-for="role in roles"
            :key="role.id"
            class="role-item"
            :class="{ active: role.id === currentRoleId }"
            @click="changeRole(role)"
          >
            <span>{{ role.name }}</span>
            <icon-right />
          </div>
        </div>
      </a-card>
      <a-card title="菜单权限" class="menus-panel" :body-style="cardBodyStyle">
        <a-spin :loading="loading" dot style="width: 100%">
          <a-tree
            v-if="menus && menus.length"
            v-model:checked-keys="checkedMenuIds"
            :data="menus"
            :checkable="true"
            :default-expand-all="true"
            size="large"
          />
        </a-spin>
      </a-card>
    </div>
  </div>
</template>

<script setup>
  const loading = ref(false);
  const roles = ref([]);
  const menus = ref([]);
  const currentRoleId = ref();
  const checkedMenuIds = ref([]);
  const cardBodyStyle = {
    overflowY: 'auto',
    height: 'calc(100% - 46px)',
  };

  onMounted(() => {
    useTableHeight();

    init();
  });

  const init = async () => {
    await Promise.all([getRoles(), getMenus()]);

    if (!currentRoleId.value) {
      if (roles.value && roles.value.length) {
        currentRoleId.value = roles.value[0].id;
      }
    }

    await getRoleMenusIds();
  };

  const changeRole = async (role) => {
    currentRoleId.value = role.id;
    await getRoleMenusIds();
  };

  const saveRoleMenus = async () => {
    try {
      await axios.postForm(
        '/api/admin/role/save_role_menus',
        jsonToFormData({
          roleId: currentRoleId.value,
          menuIds: checkedMenuIds.value ? checkedMenuIds.value.join(',') : '',
        })
      );
      useNotificationSuccess('保存成功');
    } catch (e) {
      useHandleError(e);
    }
    await getRoleMenusIds();
  };

  const getRoles = async () => {
    roles.value = await axios.get('/api/admin/role/all_roles');
  };

  const getMenus = async () => {
    try {
      loading.value = true;
      menus.value = await axios.get(`/api/admin/menu/tree`);
    } finally {
      loading.value = false;
    }
  };

  const getRoleMenusIds = async () => {
    checkedMenuIds.value = await axios.get(
      `/api/admin/role/role_menu_ids?roleId=${currentRoleId.value}`
    );
  };
</script>

<style lang="scss" scoped>
  .container-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .container-main {
    display: flex;
    column-gap: 10px;

    .roles-panel {
      width: 220px;

      .role-item-list {
        display: flex;
        flex-direction: column;
        row-gap: 6px;

        .role-item {
          border: 1px solid var(--color-neutral-3);
          border-radius: 4px;
          padding: 10px;
          cursor: pointer;
          font-size: 14px;
          font-weight: 500;
          background-color: var(--color-neutral-2);
          display: flex;
          align-items: center;
          justify-content: space-between;

          span {
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
          }

          &.active {
            // color: rgb(var(--arcoblue-6));

            color: rgb(var(--link-1));
            background-color: rgb(var(--arcoblue-6));
          }
        }
      }
    }

    .menus-panel {
      flex: 1;
    }
  }
</style>
