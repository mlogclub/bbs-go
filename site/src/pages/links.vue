<script setup>
useHead({
  title: "友情链接",
});
const links = await useHttpGet("/api/link/list");
</script>

<template>
  <section class="main">
    <div class="container">
      <div class="main-body">
        <div class="widget">
          <div class="widget-header">友情链接</div>
          <div class="widget-content">
            <ul class="links">
              <li v-for="link in links" :key="link.id" class="link">
                <div class="link-logo">
                  <img v-if="link.logo" :src="link.logo" />
                  <img v-if="!link.logo" src="~/assets/images/net.png" />
                </div>
                <div class="link-content">
                  <a
                    :href="link.url"
                    :title="link.title"
                    class="link-title"
                    target="_blank"
                    >{{ link.title }}</a
                  >
                  <p class="link-summary">
                    {{ link.summary }}
                  </p>
                </div>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<style lang="scss" scoped></style>
