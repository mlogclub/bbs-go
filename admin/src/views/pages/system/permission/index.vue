<template>
  <div class="container">
    <div class="container-main">
      <a-card title="角色列表" class="roles-panel" :body-style="cardBodyStyle">
        <div v-for="role in roles" :key="role.id">
          <a>{{ role.name }}</a>
        </div>
      </a-card>
      <a-card title="菜单权限" class="menus-panel" :body-style="cardBodyStyle">
        <div v-for="i in 1000" :key="i">sss{{ i }}</div>
      </a-card>
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
  const cardBodyStyle = {
    overflowY: 'auto',
    height: 'calc(100% - 46px)',
  };

  const data = reactive({
    page: {
      page: 1,
      limit: 20,
      total: 0,
    },
    results: [],
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
</script>

<style lang="scss" scoped>
  .container-main {
    display: flex;
    column-gap: 10px;

    .roles-panel {
      width: 260px;
    }

    .menus-panel {
      flex: 1;
    }
  }
</style>
