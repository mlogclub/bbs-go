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
  },
})
