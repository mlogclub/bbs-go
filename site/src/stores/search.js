import { defineStore } from "pinia";

export const useSearchStore = defineStore("search", {
  state: () => ({
    keyword: "",
    nodeId: 0,
    timeRange: 0,
    page: 1,
    searchPage: null,
    loading: false,
  }),
  getters: {},
  actions: {
    initParams({ keyword, nodeId, page }) {
      this.keyword = keyword || "";
      this.nodeId = nodeId || 0;
      this.page = page || 1;
    },
    changeNodeId(nodeId) {
      this.nodeId = nodeId || 0;
      this.searchTopic();
    },
    changeTimeRange(timeRange) {
      this.timeRange == timeRange || 0;
      this.searchTopic();
    },
    async searchTopic() {
      this.loading = true;
      try {
        this.searchPage = await useHttpGet("/api/search/topic", {
          params: {
            keyword: this.keyword,
            nodeId: this.nodeId,
            timeRange: this.timeRange,
            page: this.page,
          },
        });
      } finally {
        this.loading = false;
      }
    },
  },
});
