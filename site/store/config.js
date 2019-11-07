export const state = () => ({
  config: {}
})

export const mutations = {
  setConfig(state, config) {
    state.config = config
  }
}

export const actions = {
  // 加载配置
  async loadConfig(context) {
    const ret = await this.$axios.get('/api/config/configs')
    context.commit('setConfig', ret)
    return ret
  }
}

export const getters = {
  siteTitle: function (state) {
    return state.config.siteTitle || ''
  },
  siteDescription: function (state) {
    return state.config.siteDescription || ''
  },
  siteKeywords: function (state) {
    return state.config.siteKeywords || ''
  }
}
