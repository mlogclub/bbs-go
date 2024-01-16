export default defineNuxtRouteMiddleware(async () => {
  const configStore = useConfigStore()
  const userStore = useUserStore()
  await Promise.all([
    configStore.fetchConfig(),
    userStore.fetchCurrent(),
  ])
})
