import Vue from 'vue'
import Vditor from '~/components/Vditor.vue'
// import Vditor from 'vditor'

Vue.use({
  install(Vue, options) {
    Vue.component('vditor', Vditor)
  }
})
