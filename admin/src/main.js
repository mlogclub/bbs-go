import Vue from 'vue';
import RouterTab from 'vue-router-tab';
import App from './App.vue';
import router from './router';
import store from './store';
import './plugins/element';
import 'vue-router-tab/dist/lib/vue-router-tab.css';

Vue.use(RouterTab);
Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount('#app');
