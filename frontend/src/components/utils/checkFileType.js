import { toast } from "react-toastify"
import { options } from "./toast/options"

export const checkFileType = (file) => {
  if (
    file?.type !== "image/jpeg" &&
    file?.type !== "image/png" &&
    file?.type !== "image/gif"
  ) {
    toast.error("Wrong file type", options)
    return false
  }
  return true
}
