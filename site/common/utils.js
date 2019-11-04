class Utils {
  linkTo(path) {
    // debugger
    // window.location = path
    window.$nuxt.$router.push(path)
  }

  toSignin(ref) {
    if (!ref && process.client) { // 如果没配置refUrl，那么取当前地址
      ref = window.location.pathname
    }
    this.linkTo('/user/signin?ref=' + encodeURIComponent(ref))
  }

  handleToc() {
    if (!window || !window.document) {
      return
    }
    const tocSelector = '.toc'
    const tocList = document.querySelectorAll(tocSelector)
    window.addEventListener('scroll', (event) => {
      const fromTop = window.scrollY
      const mainNavLinks = document.querySelectorAll(tocSelector + ' a')
      mainNavLinks.forEach((link, index) => {
        const section = document.getElementById(decodeURI(link.hash).substring(1))
        let nextSection = null
        if (mainNavLinks[index + 1]) {
          nextSection = document.getElementById(decodeURI(mainNavLinks[index + 1].hash).substring(1))
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

    changeSize()
    window.addEventListener('resize', (event) => {
      changeSize()
    })

    // 滚动的时候控制toc位置
    window.addEventListener('scroll', (event) => {
      tocList.forEach((toc, index) => {
        changePos(toc, toc.offsetTop)
      })
    })

    // 更改toc位置
    function changePos(obj, height) {
      const scrollTop = document.documentElement.scrollTop || document.body.scrollTop
      if (scrollTop < height + 100) { // 这里+100，控制还没滚动到顶部的时候就固定toc
        obj.style.position = 'relative'
      } else {
        obj.style.position = 'fixed'
        obj.style.top = '5px'
      }
    }

    // 设置toc width
    function changeSize() {
      tocList.forEach((toc) => {
        toc.style.width = toc.parentNode.clientWidth + 'px'
        const height = (window.innerHeight || document.documentElement.clientHeight || document.body.clientHeight) - 55 // 容器高度
        const contentHeight = toc.querySelectorAll('.content > ul')[0].clientHeight
        if (contentHeight >= height) {
          toc.querySelectorAll('.content').forEach((content) => {
            content.style.height = height + 'px'
            content.style.overflow = 'auto'
          })
        }
      })
    }
  }

  isArray(sources) {
    return Object.prototype.toString.call(sources) === '[object Array]'
  }

  isDate(sources) {
    return {}.toString.call(sources) === '[object Date]' && sources.toString() !== 'Invalid Date' && !isNaN(sources)
  }

  isElement(sources) {
    return !!(sources && sources.nodeName && sources.nodeType === 1)
  }

  isFunction(sources) {
    return Object.prototype.toString.call(sources) === '[object Function]'
  }

  isNumber(sources) {
    return Object.prototype.toString.call(sources) === '[object Number]' && isFinite(sources)
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
