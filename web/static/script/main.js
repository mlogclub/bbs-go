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

// 跳转到Github登录
function toGithubSignin() {
  var redirectUrl = getQueryParam("redirectUrl") // 优先从url中获取，如果没获取到再取当前网页地址
  if (!redirectUrl) {
    redirectUrl = window.location.pathname
  }
  window.location = '/user/github/login?redirectUrl=' + encodeURI(redirectUrl)
}

// 获取url中的query参数
function getQueryParam(paramName) {
  var query = window.location.search.substring(1);
  var vars = query.split("&");
  for (var i = 0; i < vars.length; i++) {
    var pair = vars[i].split("=");
    if (pair[0] === paramName) {
      return pair[1];
    }
  }
  return "";
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

// 处理toc目录
function handleToc(tocSelector) {
  let tocList = document.querySelectorAll(tocSelector)

  window.addEventListener("scroll", event => {
    let fromTop = window.scrollY;
    let mainNavLinks = document.querySelectorAll(tocSelector + ' a');
    mainNavLinks.forEach((link, index) => {
      let section = document.getElementById(decodeURI(link.hash).substring(1));
      let nextSection = null;
      if (mainNavLinks[index + 1]) {
        nextSection = document.getElementById(decodeURI(mainNavLinks[index + 1].hash).substring(1));
      }
      if (section.offsetTop <= fromTop) {
        if (nextSection) {
          if (nextSection.offsetTop > fromTop) {
            link.classList.add('active');
          } else {
            link.classList.remove('active');
          }
        } else {
          link.classList.add('active');
        }
      } else {
        link.classList.remove('active');
      }
    });
  });

  changeSize()
  window.addEventListener('resize', event => {
    changeSize()
  });

  // 滚动的时候控制toc位置
  window.addEventListener('scroll', event => {
    tocList.forEach((toc, index) => {
      changePos(toc, toc.offsetTop)
    });
  });

  // 更改toc位置
  function changePos(obj, height) {
    let scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
    if (scrollTop < height + 100) { // 这里+100，控制还没滚动到顶部的时候就固定toc
      obj.style.position = 'relative';
    } else {
      obj.style.position = 'fixed';
      obj.style.top = '5px';
    }
  }

  // 设置toc width
  function changeSize() {
    tocList.forEach((toc, index) => {
      toc.style.width = toc.parentNode.clientWidth + 'px'

      let $toc = $(toc)
      let height = (window.innerHeight || document.documentElement.clientHeight || document.body.clientHeight) - 55; // 容器高度
      let contentHeight = $('.content > ul', $toc).height(); // 内容的高度
      if (contentHeight >= height) {
        $('.content', $toc).css('height', height + 'px');
        $('.content', $toc).css('overflow', 'auto')
      }
    });
  }

}
