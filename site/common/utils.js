class Utils {
  linkTo(path) {
    // debugger
    // window.location = path
    window.$nuxt.$router.push(path)
  }

  toSignin(ref) {
    if (!ref && process.client) {
      // 如果没配置refUrl，那么取当前地址
      ref = window.location.pathname
    }
    this.linkTo('/user/signin?ref=' + encodeURIComponent(ref))
  }

  handleToc(tocDom) {
    if (!window || !window.document || !tocDom) {
      return
    }
    const tocSelector = '.toc'
    window.addEventListener('scroll', (event) => {
      const fromTop = window.scrollY
      const mainNavLinks = document.querySelectorAll(tocSelector + ' a')
      mainNavLinks.forEach((link, index) => {
        const section = document.getElementById(
          decodeURI(link.hash).substring(1)
        )
        if (!section) {
          return
        }
        let nextSection = null
        if (mainNavLinks[index + 1]) {
          nextSection = document.getElementById(
            decodeURI(mainNavLinks[index + 1].hash).substring(1)
          )
        }
        if (section.offsetTop <= fromTop) {
          if (nextSection) {
            if (nextSection.offsetTop > fromTop) {
              link.classList.add('active')
            } else {
              link.classList.remove('active')
            }
          } else {
            link.classList.add('active')
          }
        } else {
          link.classList.remove('active')
        }
      })
    })

    // 滚动的时候控制toc位置
    const oldTop = tocDom.offsetTop
    window.addEventListener('scroll', (event) => {
      // 更改toc位置
      const scrollTop = Math.max(
        document.body.scrollTop || document.documentElement.scrollTop
      )
      if (scrollTop < oldTop) {
        tocDom.style.position = 'relative'
        tocDom.style.top = 'unset'
      } else {
        tocDom.style.position = 'fixed'
        tocDom.style.top = '52px'
      }
    })
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
