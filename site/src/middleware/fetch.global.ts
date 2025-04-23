export default defineNuxtRouteMiddleware(async (to) => {

    const configStore = useConfigStore()
    const userStore = useUserStore()

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

    if (!isInstallPage() && !config.installed) {
        return navigateTo('/install')
    }
    if (isInstallPage() && config.installed) {
        return navigateTo('/')
    }
})
