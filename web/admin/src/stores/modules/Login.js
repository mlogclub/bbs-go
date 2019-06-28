const state = {
  showLogin: false,
  userInfo: null,
}
const mutations = {
  setShowLogin(state, showLogin) {
    state.showLogin = showLogin
  },
  setUserInfo(state, userInfo) {
    state.userInfo = userInfo
  }
}
const actions = {
  showLogin(context) {
    context.commit('setShowLogin', true)
  },
  hideLogin(context) {
    context.commit('setShowLogin', false)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
