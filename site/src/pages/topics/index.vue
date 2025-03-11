<template>
  <section class="main">
    <div class="container main-container">
      <div class="main-content no-padding no-bg topics-wrapper">
        <div class="topics-nav">
          <topics-nav />
        </div>
        <div class="topics-main">
          <!-- <div class="topics-main-header">
            <div>全部帖子</div>
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
useHead({
  title: useSiteTitle("话题"),
});
</script>

<style scoped></style>
