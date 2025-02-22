import { defineStore } from 'pinia';

const useDictStore = defineStore('dict', {
  state: () => ({
    currentType: undefined,

    dicts: [],
    dictsLoading: false,
  }),
  getters: {
    currentTypeId(state) {
      return state.currentType ? state.currentType.id : undefined;
    },
  },
  actions: {
    async switchType(type) {
      this.currentType = type;
      this.loadDicts();
    },
    async loadDicts() {
      if (!this.currentTypeId) {
        return;
      }
      this.dictsLoading = true;
      try {
        this.dicts = await axios.get(
          `/api/admin/dict/list?typeId=${this.currentTypeId}`
        );
      } finally {
        this.dictsLoading = false;
      }
    },
  },
});

export default useDictStore;
