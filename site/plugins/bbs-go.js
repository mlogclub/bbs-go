import Vue from 'vue'

Vue.use({
  install(Vue, options) {
    Vue.prototype.$siteTitle = function (subTitle) {
      const siteTitle = this.$store.getters['config/siteTitle'] || ''
      if (subTitle) {
        return subTitle + (siteTitle ? ' | ' + siteTitle : '')
      }
      return siteTitle
    }

    Vue.prototype.$siteDescription = function () {
      return this.$store.getters['config/siteDescription']
    }

    Vue.prototype.$siteKeywords = function () {
      return this.$store.getters['config/siteKeywords']
    }

    Vue.prototype.$topicSiteTitle = function (topic) {
      if (topic.type === 0) {
        return this.$siteTitle(topic.title)
      } else {
        return this.$siteTitle(topic.content)
      }
    }

    /**
     * 是否是移动端
     * @returns {boolean}
     */
    Vue.prototype.$isMobile = function () {
      const sUserAgent = navigator.userAgent.toLowerCase()

      const bIsWeixin = sUserAgent.match(/micromessenger/i)
      if (bIsWeixin && bIsWeixin.index !== -1) {
        return true
      }

      const bIsIphoneOs = sUserAgent.match(/iphone os/i)
      if (bIsIphoneOs && bIsIphoneOs.index !== -1) {
        return true
      }

      const bIsAndroid = sUserAgent.match(/android/i)
      if (bIsAndroid && bIsAndroid.index !== -1) {
        return true
      }

      const bIsUc7 = sUserAgent.match(/rv:1.2.3.4/i)
      if (bIsUc7 && bIsUc7.index !== -1) {
        return true
      }

      const bIsUc = sUserAgent.match(/ucweb/i)
      if (bIsUc && bIsUc.index !== -1) {
        return true
      }

      const bIsMidp = sUserAgent.match(/midp/i)
      if (bIsMidp && bIsMidp.index !== -1) {
        return true
      }
      return false
    }

    /**
     * 链接跳转
     * @param path
     */
    Vue.prototype.$linkTo = function (path) {
      if (window) {
        window.location = path
        // 这里使用$router.push会导致跳转页面之后window.vditor对象undefined，原因未知
        // window.$nuxt.$router.push(path)
      }
    }

    /**
     * 跳转到登录页
     * @param ref
     */
    Vue.prototype.$toSignin = function (ref) {
      if (!ref && process.client) {
        // 如果没配置refUrl，那么取当前地址
        ref = window.location.pathname
      }
      this.$linkTo('/user/signin?ref=' + encodeURIComponent(ref))
    }

    /**
     * 是否是登陆页
     * @param ref
     * @returns {boolean}
     */
    Vue.prototype.$isSigninUrl = function (ref) {
      return ref === '/user/signin'
    }

    /**
     * 弹出错误消息，然后执行登录
     * @param message
     */
    Vue.prototype.$msgSignIn = function () {
      const that = this
      this.$msg({
        type: 'error',
        message: '请先登录',
        onClose() {
          that.$toSignin()
        },
      })
    }

    /**
     * 弹出消息然后执行函数
     * @param type 消息类型，success、error、info...
     * @param message 消息内容
     * @param then 要执行的函数
     */
    Vue.prototype.$msg = function ({
      type = 'success',
      message,
      duration = 800,
      onClose,
    }) {
      this.$message({
        duration: 800,
        type,
        message,
        onClose,
      })
    }
  },
})
