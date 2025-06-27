<template>
  <section class="main">
    <div class="container main-container">
      <div class="main-content no-padding no-bg topics-wrapper">
        <div class="topics-nav">
          <topics-nav />
        </div>
        <div class="topics-main">
          <!-- <div class="topics-main-header">
            <div>ALL</div>
          </div> -->
          <load-more-async v-slot="{ results }" url="/api/topic/topics">
            <topic-list :topics="results" show-sticky />
          </load-more-async>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
const { t } = useI18n();
useHead({
  title: useSiteTitle(t("pages.topics.title")),
});
</script>

<style scoped></style>
