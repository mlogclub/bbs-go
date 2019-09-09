import Vue from 'vue'
import Vditor from '~/components/Vditor.vue'

Vue.use({
  install(Vue, options) {
    Vue.component('vditor', Vditor)
  }
})
