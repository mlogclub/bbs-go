import '@/plugins/element'
import '@/plugins/axios'
import '@/plugins/filters'

import Vue from 'vue'
import App from '@/App'
import store from '@/store'
import router from '@/router'

import 'normalize.css/normalize.css' // a modern alternative to CSS resets
import '@/styles/index.scss' // global css
import '@/icons' // icon
import '@/permission' // permission control

Vue.config.productionTip = false

const app = new Vue({
  router,
  store,
  render: (h) => h(App)
}).$mount('#app')

store.$app = app // https://github.com/vuejs/vuex/issues/1399#issuecomment-491553564
