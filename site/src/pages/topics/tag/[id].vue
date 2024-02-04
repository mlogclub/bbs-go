<template>
  <section class="main">
    <div class="container main-container left-main size-320">
      <div class="left-container">
        <div class="main-content no-padding no-bg topics-wrapper">
          <div class="topics-nav">
            <topics-nav />
          </div>
          <div class="topics-main">
            <load-more-async
              v-slot="{ results }"
              url="/api/topic/tag/topics"
              :params="{ tagId: tagId }"
            >
              <topic-list :topics="results" />
            </load-more-async>
          </div>
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
const tagId = route.params.id;
const { data: tag } = await useAsyncData(() => useMyFetch(`/api/tag/${tagId}`));

useHead({
  title: useSiteTitle(tag.value.name, "话题"),
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
<style lang="scss" scoped></style>
