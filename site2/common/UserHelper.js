class UserHelper {
  hasRole(user, role) {
    if (!user || !user.roles || !user.roles.length) {
      return false
    }
    for (let i = 0; i < user.roles.length; i++) {
      if (user.roles[i] === role) {
        return true
      }
    }
    return false
  }

  hasAnyRole(user, ...roles) {
    if (!roles || !roles.length) {
      return false
    }
    for (let i = 0; i < roles.length; i++) {
      if (this.hasRole(user, roles[i])) {
        return true
      }
    }
    return false
  }

  isOwner(user) {
    return this.hasRole(user, 'owner')
  }

  isAdmin(user) {
    return this.hasRole(user, 'admin')
  }
}

export default new UserHelper()
