class CookieManager {
  setCookie(name, value, expireDays) {
    let expireDate = new Date()
    expireDate.setDate(expireDate.getDate() + expireDays)
    document.cookie = name + '=' + escape(value) + ((expireDays == null) ? '' : ';expires=' + expireDate.toGMTString())
  }

  getCookie(name) {
    if (document.cookie.length > 0) {
      let startIndex = document.cookie.indexOf(name + '=')
      if (startIndex != -1) {
        startIndex = startIndex + name.length + 1
        let endIndex = document.cookie.indexOf(';', startIndex)
        if (endIndex == -1) endIndex = document.cookie.length
        return unescape(document.cookie.substring(startIndex, endIndex))
      }
    }
    return ''
  }
}

export default new CookieManager()
