<template>
  <div>
    <MyHeader />

    <section v-if="user" class="main">
      <div class="container">
        <user-profile :user="user" />
        <div class="container main-container right-main size-320">
          <user-center-sidebar :user="user" />
          <div class="right-container">
            <slot />
          </div>
        </div>
      </div>
    </section>

    <MyFooter />
  </div>
</template>

<script setup>
const userStore = useUserStore();
const user = computed(() => {
  return userStore.user;
});
</script>

<style lang="scss" scoped></style>
