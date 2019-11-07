export default async function ({ store }) {
  await store.dispatch('user/getCurrentUser')
  await store.dispatch('config/loadConfig')
}
