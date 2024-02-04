<template>
  <section class="main">
    <div class="container main-container left-main size-320">
      <div class="left-container">
        <div class="main-content no-padding no-bg topics-wrapper">
          <div class="topics-nav">
            <topics-nav />
          </div>
          <div class="topics-main">
            <load-more-async v-slot="{ results }" url="/api/topic/topics">
              <topic-list :topics="results" show-sticky />
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
useHead({
  title: useSiteTitle("话题"),
});
</script>

<style scoped></style>
