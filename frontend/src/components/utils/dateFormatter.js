export const shortDate = (date) => {
  return new Date(date * 1000).toLocaleString("en", {
    hour: "2-digit",
    hour12: false,
    minute: "2-digit",
  })
}

export const longDate = (date) => {
  return new Date(date * 1000).toLocaleString("en", {
    day: "numeric",
    month: "short",
    hour: "2-digit",
    hour12: false,
    minute: "2-digit",
  })
}
