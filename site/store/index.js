import CommonHelper from '~/common/CommonHelper'
export const actions = {
  /**
   * see https://zh.nuxtjs.org/guide/vuex-store/#nuxtserverinit-%E6%96%B9%E6%B3%95
   *
   * @param commit
   * @param dispatch
   * @param req
   * @param app
   * @returns {Promise<void>}
   */
  async nuxtServerInit({ commit, dispatch }, { req, app }) {
    const config = await dispatch('config/loadConfig')
    app.head.title = config.siteTitle

    if (req && req.headers) {
      const ua = req.headers['user-agent']
      const isMobile = CommonHelper.isMobile(ua)
      commit('env/setIsMobile', isMobile)
    }

    await dispatch('user/getCurrentUser')
  },
}

export const getters = {}
