/**
 * 用于在切换路由的时候重置页面环境
 */
export default function ({ store, route, $cookies }) {
  store.commit('env/setShowMobileSidebar', false)
  store.commit('env/setShowMobileNodes', false)
}
