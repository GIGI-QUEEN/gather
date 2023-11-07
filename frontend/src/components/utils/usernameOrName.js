export const usernameOrName = (user) => {
  let name
  if (user?.username) {
    name = user?.username
  } else {
    name = `${user?.firstname} ${user?.lastname}`
  }
  return name
}
