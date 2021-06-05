const highlightScript =
  '//cdn.staticfile.org/highlight.js/10.3.2/highlight.min.js'
const highlightLineNumberScript =
  '//cdn.jsdelivr.net/npm/highlightjs-line-numbers.js@2.8.0/dist/highlightjs-line-numbers.min.js'
const highlightCss =
  '//cdn.staticfile.org/highlight.js/10.3.2/styles/dracula.min.css'

class CommonHelper {
  isMobile(ua) {
    return /mobile|android|webos|iphone|blackberry|micromessenger/i.test(ua)
  }

  initHighlight() {
    if (!process.client) {
      return
    }

    const me = this
    if (window.hljs) {
      window.hljs.initHighlighting()
      window.hljs.initLineNumbersOnLoad()
    } else {
      me.addScript(highlightScript, function () {
        me.addScript(highlightLineNumberScript, function () {
          window.hljs.initHighlighting()
          window.hljs.initLineNumbersOnLoad()
        })
      })
    }
  }

  getHighlightCss() {
    return highlightCss
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
