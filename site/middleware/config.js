export default async function (context) {
  if (process.server) {
    const configs = await context.$axios.get('/api/config/configs')
    context.store.commit('config/setConfigs', configs)
    context.app.head.title = configs.siteTitle
  }
}
