import cookies from 'js-cookie'
import httpClient from '@/apis/HttpClient'

const state = {
  showLogin: false,
  userToken: '',
  userInfo: null
}
const mutations = {
  setShowLogin(state, showLogin) {
    state.showLogin = showLogin
  },
  setUserInfo(state, userInfo) {
    state.userInfo = userInfo
  },
  loginSuccess(state, ret) {
    state.userToken = ret.token
    state.userInfo = ret.user
    state.showLogin = false
  }
}
const actions = {
  showLogin(context) {
    context.commit('setShowLogin', true)
  },
  hideLogin(context) {
    context.commit('setShowLogin', false)
  },
  setUserInfo(context, userInfo) {
    context.commit('setUserInfo', userInfo)
  },
  async doLogin(context, params) {
    try {
      const ret = await httpClient.post('/api/login/signin', params)
      cookies.set('userToken', ret.token, { expires: 7 })
      cookies.set('userInfo', ret.user, { expires: 7 })
      context.commit('loginSuccess', ret)
      this._vm.$message.success('登录成功')
    } catch (e) {
      this._vm.$message.error(`登录失败：${e.message || e}`)
    }
  },
  async doLogout(context) {
    try {
      await httpClient.get('/api/login/signout')
      cookies.remove('userToken')
      cookies.remove('userInfo')
      context.commit('setShowLogin', true)
      this._vm.$message.success('退出登录成功')
    } catch (e) {
      this._vm.$message.error(`退出登录失败：${e.message || e}`)
    }
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
