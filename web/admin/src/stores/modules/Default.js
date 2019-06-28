const state = {
  collapsed: false,
}
const mutations = {
  collapse(state) {
    state.collapsed = !state.collapsed
  }
}
const actions = {
  collapse(context) {
    context.commit('collapse')
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
