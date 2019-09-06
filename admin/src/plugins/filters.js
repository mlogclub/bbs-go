import Vue from 'vue';
import Format from './format';

const filters = {
  formatDate(timestamp, fmt) {
    return Format.formatDate(timestamp, fmt);
  },
};

Object.keys(filters).forEach((key) => {
  Vue.filter(key, filters[key]);
});
