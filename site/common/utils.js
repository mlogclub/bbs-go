class Utils {
  linkTo(path) {
    window.location = path
    // 这里使用$router.push会导致跳转页面之后window.vditor对象undefined，原因未知
    // window.$nuxt.$router.push(path)
  }

  toSignin(ref) {
    if (!ref && process.client) {
      // 如果没配置refUrl，那么取当前地址
      ref = window.location.pathname
    }
    this.linkTo('/user/signin?ref=' + encodeURIComponent(ref))
  }

  isSigninUrl(ref) {
    return ref === '/user/signin'
  }

  isArray(sources) {
    return Object.prototype.toString.call(sources) === '[object Array]'
  }

  isDate(sources) {
    return (
      {}.toString.call(sources) === '[object Date]' &&
      sources.toString() !== 'Invalid Date' &&
      !isNaN(sources)
    )
  }

  isElement(sources) {
    return !!(sources && sources.nodeName && sources.nodeType === 1)
  }

  isFunction(sources) {
    return Object.prototype.toString.call(sources) === '[object Function]'
  }

  isNumber(sources) {
    return (
      Object.prototype.toString.call(sources) === '[object Number]' &&
      isFinite(sources)
    )
  }

  isObject(sources) {
    return Object.prototype.toString.call(sources) === '[object Object]'
  }

  isString(sources) {
    return Object.prototype.toString.call(sources) === '[object String]'
  }

  isBoolean(sources) {
    return typeof sources === 'boolean'
  }
}
export default new Utils()
