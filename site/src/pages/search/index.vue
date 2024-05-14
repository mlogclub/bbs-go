<template>
  <section class="main">
    <div class="container main-container left-main size-320">
      <div class="left-container">
        <div class="main-content">
          <search-topics-nav />
          <div v-if="searchPage && searchPage.results">
            <search-topic-list :search-page="searchPage" />
            <!-- <pagination
              :page="searchPage.page"
              :url-prefix="'/search?q=' + keyword + '&p='"
            /> -->
          </div>
          <div v-else class="notification is-info empty-results">
            {{ searchLoading ? "加载中..." : "未搜索到内容" }}
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
const searchStore = useSearchStore();
searchStore.initParams({
  keyword: route.query.q || "",
  page: parseInt(route.query.p) || 1,
});

const searchPage = computed(() => {
  return searchStore.searchPage;
});

const searchLoading = computed(() => {
  return searchStore.searchLoading;
});

onMounted(() => {
  searchStore.searchTopic();
});
</script>

<style lang="scss" scoped>
.search-input {
  background-color: var(--bg-color);
  padding: 10px;
  text-align: center;
}
.empty-results {
  margin-top: 10px;
}
</style>
