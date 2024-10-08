<template>
  <section class="main">
    <div class="container main-container left-main size-320">
      <div class="left-container">
        <div class="articles-main">
          <load-more-async v-slot="{ results }" url="/api/article/articles">
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
useHead({
  title: useSiteTitle("文章"),
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
