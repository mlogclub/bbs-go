<template>
  <div class="container">
    <div class="container-main">
      <div class="left">
        <ul>
          <li v-for="role in roles" :key="role.id">
            {{ role.name }}
          </li>
        </ul>
      </div>
      <div class="right"> 123 </div>
    </div>
  </div>
</template>

<script setup>
  onMounted(() => {
    useTableHeight();

    getRoles();
  });

  const roles = ref([]);
  const menus = ref([]);
  const currentRoleId = ref();

  const getRoles = async () => {
    roles.value = await axios.get('/api/admin/role/all_roles');
  };
  const getMenus = async () => {
    // TODO
  };
</script>

<style lang="scss" scoped></style>
