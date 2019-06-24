window.appLogin = null // 登录组件

// 时间格式化
function formatDate(timestamp, fmt) {
  fmt = fmt || 'yyyy-MM-dd HH:mm:ss'
  var date = new Date(timestamp)
  var o = {
    'M+': date.getMonth() + 1,
    'd+': date.getDate(),
    'h+': date.getHours() % 12,
    'H+': date.getHours(),
    'm+': date.getMinutes(),
    's+': date.getSeconds(),
    'q+': Math.floor((date.getMonth() + 3) / 3),
    'S': date.getMilliseconds()
  }
  if (/(y+)/.test(fmt))
    fmt = fmt.replace(RegExp.$1, (date.getFullYear() + '').substr(4 - RegExp.$1.length))
  for (var k in o)
    if (new RegExp('(' + k + ')').test(fmt))
      fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (('00' + o[k]).substr(('' + o[k]).length)))
  return fmt
}

// 添加vue-filter
Vue.filter('formatDate', function (timestamp, fmt) {
  return formatDate(timestamp, fmt)
})

// 跳转到登陆页
function toSignin() {
  var redirectUrl = window.location.pathname
  window.location = '/user/signin?redirectUrl=' + encodeURI(redirectUrl)
}

// Get请求
function httpGet(path, params) {
  params = params || {}
  var def = $.Deferred()
  $.get(path, params, function (response, status, xhr) {
    handleAjaxResponse(def, response, status, xhr)
  })
  return def.promise()
}

// Post请求
function httpPost(path, params) {
  params = params || {}
  var def = $.Deferred()
  $.post(path, params, function (response, status, xhr) {
    handleAjaxResponse(def, response, status, xhr)
  })
  return def.promise()
}

function imageUpload(path, file) {
  var formData = new FormData()
  formData.append('image', file)

  var def = $.Deferred()
  $.ajax({
    url: path,
    type: 'POST',
    cache: false,
    processData: false,
    contentType: false,
    data: formData,
    success: function (response, status, xhr) {
      handleAjaxResponse(def, response, status, xhr)
    },
    error: function (err) {
      def.reject(err)
    }
  })
  return def.promise()
}

// 处理ajax请求返回
function handleAjaxResponse(def, response, status, xhr) {
  if (status === 'success') {
    if (response.success) {
      def.resolve(response.data)
    } else {
      if (response.errorCode === 1) {
        toSignin()
      } else {
        def.reject(response)
      }
    }
  } else {
    def.reject(response, status)
  }
}

// 处理顶部导航菜单
$(document).ready(function () {
  $('.navbar-burger').click(function () {
    $('.navbar-burger').toggleClass('is-active')
    $('.navbar-menu').toggleClass('is-active')
  })
})
