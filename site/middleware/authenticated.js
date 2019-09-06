export default async function (context) {
  const signInUrl = getSignInUrl(context)
  const userToken = getUserToken(context)
  if (!userToken) {
    context.redirect(signInUrl)
  } else {
    const user = await checkLogin(context)
    if (!user) {
      context.redirect(signInUrl)
    }
  }
}

// 检查登录
async function checkLogin(context) {
  try {
    return await context.$axios.get('/api/user/current')
  } catch (e) {
    console.error(e)
    return null
  }
}

// 获取UserToken
function getUserToken(context) {
  return context.app.$cookies.get('userToken')
}

// 获取登录跳转地址
function getSignInUrl(context) {
  let ref // 来源地址
  if (process.server) { // 服务端
    ref = context.req.originalUrl
  } else if (process.client) { // 客户端
    ref = context.route.path
  }
  let signinUrl = '/user/signin'
  if (ref) {
    signinUrl += '?ref=' + encodeURIComponent(ref)
  }
  return signinUrl
}
