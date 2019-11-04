export default async function ({ app, store, $axios }) {
  const userToken = app.$cookies.get('userToken')
  if (userToken) { // 用户登录
    const user = await $axios.get('/api/user/current')
    console.log(user)
    store.dispatch('user/setCurrentUser', user)
  }
}
