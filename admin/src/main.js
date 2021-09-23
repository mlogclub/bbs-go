import Vue from "vue";
import App from "@/App";
import store from "@/store";
import router from "@/router";

import "@/plugins/element";
import "@/plugins/axios";
import "@/plugins/filters";
import "@/plugins/lazyload";

import "normalize.css/normalize.css"; // a modern alternative to CSS resets
import "@/styles/index.scss"; // global css
import "@/icons"; // icon
import "@/permission"; // permission control

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount("#app");
