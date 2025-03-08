<template>
  <section class="main">
    <div class="container">
      <div class="widget">
        <div class="widget-header">友情链接</div>
        <div class="widget-content">
          <ul v-if="links && links.length" class="links">
            <li v-for="link in links" :key="link.id" class="link">
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
            </li>
          </ul>
          <my-empty v-else />
        </div>
      </div>
    </div>
  </section>
</template>
<script setup>
useHead({
  title: "友情链接",
});
const links = await useHttpGet("/api/link/list");
</script>

<style lang="scss" scoped>
.links {
  padding: 10px 15px;
  .link {
    display: block;
    height: 62px;
    padding-top: 5px;

    &:not(:last-child) {
      border-bottom: 1px solid var(--border-color4);
    }

    .link-title {
      font-size: 15px;
      font-weight: 500;
      color: var(--text-link-color);

      overflow: hidden;
      word-break: break-all;
      -webkit-line-clamp: 1;
      text-overflow: ellipsis;
      -webkit-box-orient: vertical;
      display: -webkit-box;
    }

    .link-summary {
      font-size: 13px;
      margin-top: 3px;
      color: var(--text-color3);

      overflow: hidden;
      word-break: break-all;
      -webkit-line-clamp: 1;
      text-overflow: ellipsis;
      -webkit-box-orient: vertical;
      display: -webkit-box;
    }
  }
}
</style>
