export function userHasRole(user, role) {
  if (!user || !user.roles || !user.roles.length) {
    return false;
  }
  for (let i = 0; i < user.roles.length; i++) {
    if (user.roles[i] === role) {
      return true;
    }
  }
  return false;
}

export function userHasAnyRole(user, ...roles) {
  if (!roles || !roles.length) {
    return false;
  }
  for (let i = 0; i < roles.length; i++) {
    if (userHasRole(user, roles[i])) {
      return true;
    }
  }
  return false;
}

export function userIsOwner(user) {
  return userHasRole(user, "owner");
}

export function userIsAdmin(user) {
  return userHasRole(user, "admin");
}
