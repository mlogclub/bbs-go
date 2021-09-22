import Vue from "vue";
import Element from "element-ui";
import "element-ui/lib/theme-chalk/index.css";
import Cookies from "js-cookie";

Vue.use(Element, {
  size: Cookies.get("size") || "medium", // set element-ui default size
});
