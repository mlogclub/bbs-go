<template>
  <section class="main">
    <div class="container main-container left-main size-320">
      <div class="left-container">
        <div class="articles-main">
          <load-more-async
            v-slot="{ results }"
            url="/api/article/tag/articles"
            :params="{ tagId: tag.id }"
          >
            <article-list :articles="results" />
          </load-more-async>
        </div>
      </div>
      <div class="right-container">
        <check-in />
        <site-notice />
        <score-rank />
        <friend-links />
      </div>
    </div>
  </section>
</template>

<script setup>
const route = useRoute();
const { data: tag } = await useMyFetch(`/api/tag/${route.params.id}`);

useHead({
  title: useSiteTitle(tag.value.name, "文章"),
  meta: [
    {
      hid: "description",
      name: "description",
      content: useSiteDescription(),
    },
    { hid: "keywords", name: "keywords", content: useSiteKeywords() },
  ],
});
</script>

<style lang="scss" scoped>
.articles-main {
  background-color: var(--bg-color);
  border-radius: var(--border-radius);
}
</style>
