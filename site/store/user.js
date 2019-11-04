export const state = () => ({
  current: null,
  userToken: null
})

export const mutations = {
  setCurrent(state, user) {
    state.current = user
  },
  setUserToken(state, userToken) {
    state.userToken = userToken
  }
}

export const actions = {
  // 登录成功
  loginSuccess(context, { token, user }) {
    this.$cookies.set('userToken', token, { maxAge: 86400 * 7, path: '/' })
    context.commit('setUserToken', token)
    context.commit('setCurrent', user)
  },

  // 设置当前登录用户
  setCurrentUser(context, user) {
    context.commit('setCurrent', user)
  },

  // 登录
  async signin(context, { username, password }) {
    const ret = await this.$axios.post('/api/login/signin', {
      username: username,
      password: password
    })
    context.dispatch('loginSuccess', ret)
    return ret.user
  },

  // github登录
  async signinByGithub(context, { code, state }) {
    const ret = await this.$axios.get('/api/login/github/callback', {
      params: {
        code: code,
        state: state
      }
    })
    context.dispatch('loginSuccess', ret)
    return ret.user
  },

  // qq登录
  async signinByQQ(context, { code, state }) {
    const ret = await this.$axios.get('/api/login/qq/callback', {
      params: {
        code: code,
        state: state
      }
    })
    context.dispatch('loginSuccess', ret)
    return ret.user
  },

  async signup(context, { nickname, username, password, rePassword }) {
    const ret = await this.$axios.post('/api/login/signup', {
      nickname: nickname,
      username: username,
      password: password,
      rePassword: rePassword
    })
    context.dispatch('loginSuccess', ret)
    return ret.user
  },

  // 退出登录
  async signout(context) {
    const userToken = this.$cookies.get('userToken')
    await this.$axios.get('/api/login/signout', {
      params: {
        userToken: userToken
      }
    })
    context.commit('setUserToken', null)
    context.commit('setCurrent', null)
    this.$cookies.remove('userToken')
  }
}
