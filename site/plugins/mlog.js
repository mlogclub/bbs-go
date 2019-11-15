import Vue from 'vue'

Vue.use({
  install(Vue, options) {
    Vue.prototype.$siteTitle = function(subTitle) {
      const siteTitle = this.$store.getters['config/siteTitle'] || ''
      if (subTitle) {
        return subTitle + (siteTitle ? ' | ' + siteTitle : '')
      }
      return siteTitle
    }

    Vue.prototype.$siteDescription = function() {
      return this.$store.getters['config/siteDescription']
    }

    Vue.prototype.$siteKeywords = function() {
      return this.$store.getters['config/siteKeywords']
    }
  }
})
