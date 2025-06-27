export default defineNuxtRouteMiddleware(async (to) => {

    const configStore = useConfigStore()
    const userStore = useUserStore()
    const { $i18n } = useNuxtApp()

    const load = async () => {
        await Promise.all([
            configStore.fetchConfig(),
            userStore.fetchCurrent(),
        ])
    }

    // const nuxtApp = useNuxtApp()
    // // 服务端渲染，或者客服端渲染
    // if (import.meta.server || (import.meta.client && !nuxtApp.isHydrating)) {
    //     await load()
    // }

    const isInstallPage = () => {
        return to.path.startsWith('/install')
    }

    await load()
    const config: any = configStore.config

    // i18n
    if (config.language && config.language !== $i18n.locale.value) {
        console.log('language change: ', $i18n.locale.value, '->', config.language)
        await $i18n.setLocale(config.language)
    }
    if (!isInstallPage() && !config.installed) {
        return navigateTo('/install')
    }
    if (isInstallPage() && config.installed) {
        return navigateTo('/')
    }
})
