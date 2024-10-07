<template>
  <section class="main">
    <div class="container main-container">
      <div class="main-content no-padding no-bg topics-wrapper">
        <div class="topics-nav">
          <topics-nav />
        </div>
        <div class="topics-main">
          <div class="topics-main-header">
            <div>全部帖子</div>
            <div></div>
          </div>
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
  title: useSiteTitle(),
});
</script>

<style lang="scss" scoped>
.topics-main-header {
  margin: 0 32px;
  padding: 24px 0;
  position: relative;
  font-style: normal;
  font-weight: 600;
  font-size: 20px;
  line-height: 16px;
  color: var(--text-color);
  border-bottom: 1px solid var(--border-color4);

  @media screen and (max-width: 768px) {
    & {
      margin: 0 12px;
    }
  }
}
</style>
