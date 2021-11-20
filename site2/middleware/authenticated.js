import UserHelper from '~/common/UserHelper'

export default function (context) {
  const user = context.store.state.user.current
  if (!user) {
    toSignIn(context)
    return
  }
  if (isAdminUrl(context)) {
    if (!UserHelper.isOwner(user) && !UserHelper.isAdmin(user)) {
      context.error({
        statusCode: 403,
        message: '403 forbidden',
      })
    }
  }
}

// 当前访问URL是否是管理后台
function isAdminUrl(context) {
  return context.route.path.indexOf('/admin') === 0
}

// 前往登录地址
function toSignIn(context) {
  const signInUrl = getSignInUrl(context)
  context.redirect(signInUrl)
}

// 获取登录跳转地址
function getSignInUrl(context) {
  let ref // 来源地址
  if (process.server) {
    // 服务端
    ref = context.req.originalUrl
  } else if (process.client) {
    // 客户端
    ref = context.route.path
  }
  let signinUrl = '/user/signin'
  if (ref) {
    signinUrl += '?ref=' + encodeURIComponent(ref)
  }
  return signinUrl
}
