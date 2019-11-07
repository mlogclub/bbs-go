import axios from 'axios';
import qs from 'qs';
import cookies from 'js-cookie'

class HttpClient {
  constructor() {
    this.http = axios.create({
      // 请替换为你自己的接口地址
      baseURL: 'https://mlog.club'
    });
    this.http.defaults.headers.common['X-Client'] = 'bbs-go-admin';
    this.http.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded';
    this.http.interceptors.request.use((config) => {
      const userToken = cookies.get('userToken');
      if (userToken) {
        config.headers.common['X-User-Token'] = userToken
      }
      return config;
    }, (reason) => {
      console.error(reason);
    });
    this.http.interceptors.response.use((response) => {
      if (response.status !== 200) {
        return Promise.reject(response);
      }
      if (response.data.success) {
        return response.data.data;
      }
      if (response.data.errorCode === 1) { // 未登录
        window.vue.$store.dispatch('Login/showLogin');
        return Promise.reject(response.data);
      }
      // window.vue.$message({
      //   showClose: true,
      //   message: response.data.message,
      //   type: 'error',
      // });
      return Promise.reject(response.data);
    }, (error) => {
      window.vue.$message({
        showClose: true,
        message: error,
        type: 'error',
      });
      return Promise.reject(error);
    });
  }

  get(api, params) {
    return this.http.get(api, {
      params: params
    });
  }

  post(api, data) {
    return this.http.post(api, qs.stringify(data));
  }
}

export default new HttpClient();
