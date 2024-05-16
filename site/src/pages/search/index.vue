<template>
  <section class="main">
    <div class="container main-container left-main size-320">
      <div class="left-container">
        <div class="main-content">
          <search-topics-nav />
          <load-more-async
            ref="loadMore"
            v-slot="{ results }"
            url="/api/search/topic"
            :params="params"
          >
            <search-topic-list :results="results" />
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
const loadMore = ref(null);
const params = reactive({
  keyword: route.query.q || "",
  nodeId: route.query.nodeId || 0,
  timeRange: route.query.timeRange,
});

watch(
  () => route.query,
  (newQuery, oldQuery) => {
    params.keyword = newQuery.q || "";
    params.nodeId = newQuery.nodeId || 0;
    params.timeRange = newQuery.timeRange;
    nextTick(() => {
      if (loadMore.value) {
        loadMore.value.refresh();
      }
    });
  },
  { immediate: true }
);
</script>

<style lang="scss" scoped>
.search-input {
  background-color: var(--bg-color);
  padding: 10px;
  text-align: center;
}
</style>
