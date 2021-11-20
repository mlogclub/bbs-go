export const state = () => ({
  // 当前设备是否是手机
  isMobile: false,
  // 移动端是否显示侧边栏
  showMobileSidebar: false,
  // 移动端是否显示nodes选择面板
  showMobileNodes: false,
  // 当前所在节点编号，0：最新、-1：recommend
  currentNodeId: 0,
})

export const mutations = {
  setIsMobile(state, isMobile) {
    state.isMobile = isMobile
  },
  setShowMobileSidebar(state, show) {
    state.showMobileSidebar = show
  },
  setShowMobileNodes(state, show) {
    state.showMobileNodes = show
  },
  setCurrentNodeId(state, nodeId) {
    state.currentNodeId = nodeId
  },
}

export const actions = {}

export const getters = {}
