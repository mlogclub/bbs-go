<template>
  <div class="topics-main">
    <load-more-async v-slot="{ results }" url="/api/topic/topics">
      <topic-list :topics="results" />
    </load-more-async>
  </div>
</template>

<script setup>
definePageMeta({
  layout: "topics",
});
useHead({
  title: useSiteTitle("话题"),
});
</script>

<style scoped></style>
