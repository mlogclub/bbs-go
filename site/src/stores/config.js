import { defineStore } from "pinia";

export const useConfigStore = defineStore("config", {
  state: () => ({
    config: {},
  }),
  getters: {
    // doubleCounter: state => state.counter * 2,
    // doubleCounterPlusOne(): number {
    //   return this.doubleCounter + 1
    // },
    isEnabledArticle(state) {
      return state.config.modules.article || true;
    },
  },
  actions: {
    async fetchConfig() {
      const { data } = await useMyFetch("/api/config/configs");
      this.config = data.value;
    },
  },
});
