<script setup>
definePageMeta({
  layout: "topics",
});
useHead({
  title: useSiteTitle(),
});
</script>

<template>
  <div class="topics-main">
    <load-more-async v-slot="{ results }" url="/api/topic/topics">
      <topic-list :topics="results" show-sticky />
    </load-more-async>
  </div>
</template>

<style scoped></style>
