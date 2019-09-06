export const state = () => ({
  configs: {}
})

export const mutations = {
  setConfigs(state, configs) {
    state.configs = configs
  }
}

export const actions = {

}

export const getters = {
  siteTitle: function (state) {
    return state.configs.siteTitle || ''
  },
  siteDescription: function (state) {
    return state.configs.siteDescription || ''
  },
  siteKeywords: function (state) {
    return state.configs.siteKeywords || ''
  }
}
