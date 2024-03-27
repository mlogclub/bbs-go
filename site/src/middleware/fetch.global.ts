export default defineNuxtRouteMiddleware(async () => {

    const load = async () => {
        const configStore = useConfigStore()
        const userStore = useUserStore()
        await Promise.all([
            configStore.fetchConfig(),
            userStore.fetchCurrent(),
        ])
    }

    const nuxtApp = useNuxtApp()

    // 服务端渲染，或者客服端渲染
    if (import.meta.server || (import.meta.client && !nuxtApp.isHydrating)) {
        await load()
    }

})
