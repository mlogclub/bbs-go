export const state = () => ({
  current: null,
  userToken: null,
})

export const mutations = {
  setCurrent(state, user) {
    state.current = user
  },
  setUserToken(state, userToken) {
    state.userToken = userToken
  },
}

export const actions = {
  // 登录成功
  loginSuccess(context, { token, user }) {
    const config = context.rootState.config.config
    const cookieMaxAge = 86400 * config.tokenExpireDays
    this.$cookies.set('userToken', token, { maxAge: cookieMaxAge, path: '/' })
    context.commit('setUserToken', token)
    context.commit('setCurrent', user)
  },

  // 获取当前登录用户
  async getCurrentUser(context) {
    const user = await this.$axios.get('/api/user/current')
    context.commit('setCurrent', user)
    return user
  },

  // 登录
  async signin(context, { captchaId, captchaCode, username, password }) {
    const ret = await this.$axios.post('/api/login/signin', {
      captchaId,
      captchaCode,
      username,
      password,
    })
    context.dispatch('loginSuccess', ret)
    return ret.user
  },

  // github登录
  async signinByGithub(context, { code, state }) {
    const ret = await this.$axios.get('/api/github/login/callback', {
      params: {
        code,
        state,
      },
    })
    context.dispatch('loginSuccess', ret)
    return ret.user
  },

  // osc登录
  async signinByOSC(context, { code, state }) {
    const ret = await this.$axios.get('/api/osc/login/callback', {
      params: {
        code,
        state,
      },
    })
    context.dispatch('loginSuccess', ret)
    return ret.user
  },

  // qq登录
  async signinByQQ(context, { code, state }) {
    const ret = await this.$axios.get('/api/qq/login/callback', {
      params: {
        code,
        state,
      },
    })
    context.dispatch('loginSuccess', ret)
    return ret.user
  },

  async signup(
    context,
    { captchaId, captchaCode, nickname, username, email, password, rePassword }
  ) {
    const ret = await this.$axios.post('/api/login/signup', {
      captchaId,
      captchaCode,
      nickname,
      username,
      email,
      password,
      rePassword,
    })
    context.dispatch('loginSuccess', ret)
    return ret.user
  },

  // 退出登录
  async signout(context) {
    const userToken = this.$cookies.get('userToken')
    await this.$axios.get('/api/login/signout', {
      params: {
        userToken,
      },
    })
    context.commit('setUserToken', null)
    context.commit('setCurrent', null)
    this.$cookies.remove('userToken')
  },
}
