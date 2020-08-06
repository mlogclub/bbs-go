import Vue from 'vue'
import VueLazyload from 'vue-lazyload'

Vue.use(VueLazyload)

Vue.use(VueLazyload, {
  preLoad: 1.3, // 预加载的宽高比
  attempt: 1, // 尝试加载次数
})
