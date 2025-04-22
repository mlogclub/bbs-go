// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: true },
  srcDir: 'src/',

  ssr: false,
  modules: [
    '@pinia/nuxt',
    '@vueuse/nuxt',
    // https://color-mode.nuxtjs.org/#configuration
    '@nuxtjs/color-mode',
    '@element-plus/nuxt',
    ['nuxt-lazy-load', {
      images: true,
      videos: true,
      audios: true,
      iframes: true,
      native: true,
      directiveOnly: false,

      // Default image must be in the public folder
      // defaultImage: '/images/default-image.jpg',

      // To remove class set value to false
      loadingClass: 'isLoading',
      loadedClass: 'isLoaded',
      appendClass: 'lazyLoad',

      observerConfig: {
        // See IntersectionObserver documentation
      },
    }],
  ],

  plugins: [
  ],

  elementPlus: {
    defaultLocale: 'zh-cn',
  },

  colorMode: {
    preference: 'system', // default value of $colorMode.preference
    fallback: 'light', // fallback value if not system preference found
    storageKey: 'bbsgo-color-mode',
    classPrefix: 'theme-',
    classSuffix: '',
  },

  imports: {
    dirs: [
      'apis',
      'stores',
    ],
  },

  app: {
    head: {
      title: 'BBS-GO',
      htmlAttrs: { class: 'theme-light has-navbar-fixed-top' },
    },
  },

  css: [
    '~/assets/css/index.scss',
  ],

  nitro: {
    routeRules: {
      '/api/**': {
        proxy: `${import.meta.env.SERVER_URL}/api/**`,
      },
      '/admin/**': {
        proxy: `${import.meta.env.SERVER_URL}/admin/**`,
      },
    },
  },

  compatibilityDate: '2024-09-15',
})