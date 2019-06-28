import Vue from 'vue';
import App from './App.vue';
import router from './router';
import store from './stores';
import format from './plugins/format';
import './plugins/element';
import './plugins/editor';
import './plugins/filters';
import './styles/main.scss'

Vue.config.productionTip = false;
Vue.prototype.format = format

window.vue = new Vue({
  router,
  store,
  render: h => h(App),
}).$mount('#app');
