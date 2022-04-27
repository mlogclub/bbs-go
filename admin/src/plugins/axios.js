import axios from "axios";
import VueAxios from "vue-axios";
import { MessageBox, Message } from "element-ui";
import store from "@/store";
import { getToken } from "@/utils/auth";
import qs from "qs";
import Vue from "vue";

const formatFormDataKey = "__formData";

function isFormData(data) {
  return data && data[formatFormDataKey] === "formData";
}

// create an axios instance
const axiosInstance = axios.create({
  baseURL: process.env.VUE_APP_BASE_API, // url = base url + request url
  timeout: 5000, // request timeout,
});

// 设置form请求
axiosInstance.form = function (url, data, config) {
  if (!data) {
    data = {};
  }
  data[formatFormDataKey] = "formData";
  return this.post(url, data, config);
};

// request interceptor
axiosInstance.interceptors.request.use(
  (config) => {
    if (store.getters.token) {
      config.headers["X-User-Token"] = getToken();
    }

    // 如果是form请求
    if (isFormData(config.data)) {
      delete config.data[formatFormDataKey];
      config.data = qs.stringify(config.data); // 转为formdata数据格式
    }

    return config;
  },
  (error) => Promise.reject(error)
);

// response interceptor
axiosInstance.interceptors.response.use(
  (response) => {
    const res = response.data;

    if (res.success !== true) {
      Message({
        message: res.message || "Error",
        type: "error",
        duration: 5 * 1000,
      });

      if (res.errorCode === 1) {
        MessageBox.confirm(
          "You have been logged out, you can cancel to stay on this page, or log in again",
          "Confirm logout",
          {
            confirmButtonText: "Re-Login",
            cancelButtonText: "Cancel",
            type: "warning",
          }
        ).then(() => {
          store.dispatch("user/resetToken").then(() => {
            location.reload();
          });
        });
      }
      return Promise.reject(res);
    }
    return Promise.resolve(res.data);
  },
  (error) => {
    Message({
      message: error.message,
      type: "error",
      duration: 5 * 1000,
    });
    return Promise.reject(error);
  }
);

// export default axiosInstance
Vue.use(VueAxios, axiosInstance);
