<template>
  <div>
    <my-nav />

    <section class="main">
      <div class="container">
        <user-profile :user="user" />
        <div class="container main-container right-main size-320">
          <user-center-sidebar :user="user" />
          <div class="right-container">
            <nuxt />
          </div>
        </div>
      </div>
    </section>

    <my-footer />
  </div>
</template>

<script>
export default {
  middleware: 'authenticated',
  computed: {
    user() {
      return this.$store.state.user.current
    },
  },
  methods: {},
}
</script>

<style lang="scss" scoped></style>
