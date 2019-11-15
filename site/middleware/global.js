export default async function({ store, app }) {
  await store.dispatch('user/getCurrentUser')
  const config = await store.dispatch('config/loadConfig')
  app.head.title = config.siteTitle
}
