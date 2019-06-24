import axios from 'axios'
import cookieManager from './CookieManager'
import config from './Config'

class HttpClient {
  constructor() {
    this.http = axios.create({
      baseURL: config.host
    })
    this.http.defaults.headers.common['X-Client'] = 'mlog'
    this.http.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded'
    this.http.interceptors.request.use(function (config) {
      let accessToken = cookieManager.getCookie('accessToken')
      if (accessToken) {
        config.headers.common['Authorization'] = 'Bearer ' + accessToken
      }
      return config
    }, function (reason) {
      console.error(reason)
    })
    this.http.interceptors.response.use(function (response) {
      if (response.status != 200) {
        return Promise.reject('请求失败，status=' + response.status)
      }
      if (response.data.success) {
        return response.data.data
      } else {
        if (response.data.errorCode === 1) { // 未登录
          window.vue.$store.dispatch('Login/showLogin')
          return Promise.reject(response.data)
        }
        window.vue.$message({
          showClose: true,
          message: response.data.message,
          type: 'error'
        })
        return Promise.reject(response.data)
      }
    }, function (error) {
      window.vue.$message({
        showClose: true,
        message: error,
        type: 'error'
      })
      return Promise.reject(error)
    })
  }

  get(api) {
    return this.http.get(api)
  }

  post(api, data) {
    return this.http.post(api, toParams(data))
  }
}

// json object to params
function toParams(data) {
  const params = new URLSearchParams()
  if (data) {
    for (let o in data) {
      params.append(o, data[o])
    }
  }
  return params
}

export default new HttpClient()
