import Vue from "vue";
import Element from "element-ui";
import "@/styles/element-variables.scss";
import Cookies from "js-cookie";

Vue.use(Element, {
  size: Cookies.get("size") || "medium", // set element-ui default size
});
