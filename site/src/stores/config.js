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
      for (let i = 0; i < state.config.modules.length; i++) {
        if (state.config.modules[i].module === "article") {
          return state.config.modules[i].enabled;
        }
      }
      return true;
    },
  },
  actions: {
    async fetchConfig() {
      this.config = await useMyFetch("/api/config/configs");
    },
  },
});
