import qs from 'qs'

export default function ({ $axios, app }) {
  $axios.onRequest((config) => {
    config.headers.common['X-Client'] = 'bbs-go-site'
    config.headers.post['Content-Type'] = 'application/x-www-form-urlencoded'
    const userToken = app.$cookies.get('userToken')
    if (userToken) {
      config.headers.common['X-User-Token'] = userToken
    }
    config.transformRequest = [
      function (data) {
        if (process.client && data instanceof FormData) {
          // 如果是FormData就不转换
          return data
        }
        data = qs.stringify(data)
        return data
      },
    ]
  })

  $axios.onResponse((response) => {
    if (response.status !== 200) {
      return Promise.reject(response)
    }
    const jsonResult = response.data
    if (jsonResult.success) {
      return Promise.resolve(jsonResult.data)
    } else {
      return Promise.reject(jsonResult)
    }
  })
}
