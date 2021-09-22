import Vue from "vue";
import * as filters from "@/filters"; // global filters

// register global utility filters
Object.keys(filters).forEach((key) => {
  Vue.filter(key, filters[key]);
});
