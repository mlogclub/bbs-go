<template>
  <section class="main">
    <div class="container main-container">
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
  </section>
</template>

<script setup>
const { t } = useI18n();
const route = useRoute();
const tagId = route.params.id;
const { data: tag } = await useMyFetch(`/api/tag/${tagId}`);

useHead({
  title: useSiteTitle(tag.value.name, t("pages.topics.title")),
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
