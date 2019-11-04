export default async function ({ app, store, $axios }) {
  const userToken = app.$cookies.get('userToken')
  if (userToken) {
    const user = await $axios.get('/api/user/current')
    if (user) {
      store.dispatch('user/setCurrentUser', user)
    }
  }
}
