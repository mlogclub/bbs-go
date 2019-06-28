import Vue from 'vue';
import Format from './format';

const filters = {
  formatDate: function (timestamp, fmt) {
    return Format.formatDate(timestamp, fmt);
  }
};

Object.keys(filters).forEach(function (key) {
  Vue.filter(key, filters[key]);
});
