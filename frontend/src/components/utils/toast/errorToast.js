import { toast } from "react-toastify";
import { options } from "./options";
export const errorToast = (error) => {
  return toast.error(error.response["data"]["error_description"], options);
};
