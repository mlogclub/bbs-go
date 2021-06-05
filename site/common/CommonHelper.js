class CommonHelper {
  highlightScript = '//cdn.staticfile.org/highlight.js/10.3.2/highlight.min.js'
  highlightLineNumberScript =
    '//cdn.jsdelivr.net/npm/highlightjs-line-numbers.js@2.8.0/dist/highlightjs-line-numbers.min.js'

  highlightCss =
    '//cdn.staticfile.org/highlight.js/10.3.2/styles/dracula.min.css'

  isMobile(ua) {
    return /mobile|android|webos|iphone|blackberry|micromessenger/i.test(ua)
  }

  initHighlight(ctx) {
    if (!process.client) {
      return
    }
    const me = this
    if (window.hljs) {
      window.hljs.initHighlighting()
      if (window.hljs.initLineNumbersOnLoad) {
        window.hljs.initLineNumbersOnLoad()
      } else {
        me.addScript(me.highlightLineNumberScript, function () {
          window.hljs.initLineNumbersOnLoad()
        })
      }
    } else {
      me.addScript(this.highlightScript, function () {
        window.hljs.initHighlighting()
        me.addScript(me.highlightLineNumberScript, function () {
          window.hljs.initLineNumbersOnLoad()
        })
      })
    }
  }

  addScript(url, callback) {
    if (!process.client) {
      console.warn('Add script fail, !process.client, ' + url)
      return
    }
    const script = document.createElement('script')
    script.type = 'text/javascript'
    script.src = url
    script.defer = true
    document.body.appendChild(script)
    script.onload = function () {
      if (callback) callback()
    }
  }
}

export default new CommonHelper()
